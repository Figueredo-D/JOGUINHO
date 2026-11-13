package entities

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
	"joguinho/internal/progression"
	"joguinho/internal/sprites"
)

// ClassType define o arquétipo do herói.
type ClassType string

const (
	ClassWarrior    ClassType = "Guerreiro"
	ClassArcher     ClassType = "Arqueiro"
	ClassMage       ClassType = "Mago"
	ClassShotgunner ClassType = "Atirador Pesado"
	ClassSheriff    ClassType = "Xerife"
)

// ActiveEffect representa um buff ou debuff ativo no jogador.
type ActiveEffect struct {
	Name      string
	IsDebuff  bool
	Remaining int
	Total     int
	Color     color.RGBA
}

// AfterImage representa uma silhueta fantasma deixada durante o Dash.
type AfterImage struct {
	Cx, Cy     float64
	FacingLeft bool
	Sprite     *ebiten.Image
	Life       int
	MaxLife    int
}

// Player representa o herói controlado pelo jogador com suporte a upgrades meta e habilidades ativas.
type Player struct {
	ClassType    ClassType
	X, Y         float64
	Width        float64
	Height       float64
	BaseSpeed    float64
	HP           float64
	MaxHP        float64
	Mana         float64
	MaxMana      float64
	ManaRegen    float64
	AttackCD     int
	BaseAttackCD int
	InvincibleCD    int
	HitFlashCD      int
	FacingAngle     float64
	FacingLeft      bool
	FrameCounter    int
	IsMoving        bool
	WalkTimer       int
	AttackAnimTimer int
	JustStepped     bool

	// Atributos derivados de upgrades meta
	DamageMultiplier float64
	SlashRange       float64
	ArrowSpeed       float64
	SpellMultiplier  float64

	// Habilidade Especial Ativa (Tecla E ou Botão Direito)
	SpecialCD    int
	SpecialMaxCD int

	// Habilidade Guerreiro: Postura de Escudo
	Shield ShieldState

	// Habilidade Arqueiro: Chuva de Flechas
	ArrowRainRadius float64
	ArrowRainDamage float64

	// Habilidade Mago: Pilar de Fogo
	FirePillarRadius   float64
	FirePillarDuration int
	FirePillarDamage   float64
	FirePillarManaCost float64

	// Atirador Pesado (Espingarda Calibre 12)
	PelletCount      int
	PelletSpread     float64
	PelletDamage     float64
	PelletKnockback  float64
	ShockwaveRadius  float64
	ShockwaveDamage  float64
	MuzzleSmokeTimer int

	// Xerife (Tambor de 6 Tiros & Deadeye)
	Ammo          int
	MaxAmmo       int
	ReloadTimer   int
	MaxReloadTime int
	CritChance    float64
	CritMult      float64
	DeadeyeQueue  int
	DeadeyeTimer  int

	// Buffs e Debuffs (Efeitos de Drops)
	SpeedBoostTimer    int
	AgilityTimer       int
	CursedChaliceTimer int
	PoisonTimer        int
	FatigueTimer       int

	// Dash Universal (Espaço / Shift)
	DashTimer   int
	DashCD      int
	DashMaxCD   int
	DashVx      float64
	DashVy      float64
	AfterImages []AfterImage

	// Relíquias Passivas Acumuláveis na Run
	Relics       map[RelicID]int
	TotalRelics  int
	HitCounter   int // Contador de hits certeiros para Faísca Eletrostática
	LastMoveDx   float64
	LastMoveDy   float64

	// Sprites procedurais cacheados
	spriteWarrior     *ebiten.Image
	spriteArcherRest  *ebiten.Image
	spriteArcherShoot *ebiten.Image
	spriteMage        *ebiten.Image
	spriteShotgunner  *ebiten.Image
	spriteSheriff     *ebiten.Image
}

// NewPlayer inicializa o jogador aplicando as melhorias meta do save e configurando a habilidade especial.
func NewPlayer(class ClassType, startX, startY float64) *Player {
	save := progression.GetSave()

	p := &Player{
		ClassType:         class,
		X:                 startX,
		Y:                 startY,
		Width:             32,
		Height:            32,
		DamageMultiplier:  1.0,
		SpellMultiplier:   1.0,
		ArrowSpeed:        9.5,
		SlashRange:        72.0,
		DashMaxCD:         110, // ~1.8s de recarga para esquiva
		Relics:            make(map[RelicID]int),
		spriteWarrior:     sprites.GetWarriorSprite(),
		spriteArcherRest:  sprites.GetArcherSprite(false),
		spriteArcherShoot: sprites.GetArcherSprite(true),
		spriteMage:        sprites.GetMageSprite(),
		spriteShotgunner:  sprites.GetShotgunnerSprite(),
		spriteSheriff:     sprites.GetSheriffSprite(),
	}

	switch class {
	case ClassWarrior:
		bonusHP := float64(save.Warrior.MaxHPLevel) * 25.0
		p.MaxHP = 200 + bonusHP
		p.HP = p.MaxHP
		p.BaseSpeed = 3.2
		p.BaseAttackCD = 22

		p.DamageMultiplier = 1.0 + float64(save.Warrior.DamageLevel)*0.12
		p.SlashRange = 72.0 + float64(save.Warrior.DamageLevel)*6.0

		shieldCD := 420 - save.Warrior.ShieldLevel*45
		if shieldCD < 180 {
			shieldCD = 180
		}
		shieldDur := 75 + save.Warrior.ShieldLevel*15
		p.SpecialMaxCD = shieldCD
		p.Shield = ShieldState{
			Active:        false,
			Timer:         0,
			MaxDuration:   shieldDur,
			Cooldown:      0,
			MaxCooldown:   shieldCD,
			ArcSpan:       math.Pi * 0.72,
			Radius:        42.0,
			ReflectDamage: 2.2 + float64(save.Warrior.ShieldLevel)*0.35,
		}

	case ClassArcher:
		bonusSpeed := float64(save.Archer.SpeedLevel) * 0.35
		p.MaxHP = 75
		p.HP = p.MaxHP
		p.BaseSpeed = 4.0 + bonusSpeed
		p.BaseAttackCD = 13 - int(float64(save.Archer.AttackSpeedLvl)*0.8)
		if p.BaseAttackCD < 7 {
			p.BaseAttackCD = 7
		}
		p.ArrowSpeed = 9.5 + float64(save.Archer.AttackSpeedLvl)*1.2

		p.SpecialMaxCD = 460
		p.ArrowRainRadius = 70.0 * (1.0 + float64(save.Archer.ArrowRainLevel)*0.15)
		p.ArrowRainDamage = 13.0 * (1.0 + float64(save.Archer.ArrowRainLevel)*0.20)

	case ClassMage:
		bonusMana := float64(save.Mage.ManaLevel) * 20.0
		bonusRegen := float64(save.Mage.ManaLevel) * 3.5
		p.MaxHP = 110
		p.HP = p.MaxHP
		p.MaxMana = 100 + bonusMana
		p.Mana = p.MaxMana
		p.ManaRegen = 16 + bonusRegen
		p.BaseSpeed = 3.0
		p.BaseAttackCD = 20
		p.SpellMultiplier = 1.0 + float64(save.Mage.SpellPowerLevel)*0.15

		p.SpecialMaxCD = 360
		p.FirePillarRadius = 55.0 * (1.0 + float64(save.Mage.FirePillarLevel)*0.15)
		p.FirePillarDuration = 90 + save.Mage.FirePillarLevel*20
		p.FirePillarDamage = 24.0 * (1.0 + float64(save.Mage.SpellPowerLevel)*0.15)
		p.FirePillarManaCost = 30.0

	case ClassShotgunner:
		p.MaxHP = 140
		p.HP = p.MaxHP
		p.BaseSpeed = 2.95
		p.BaseAttackCD = 42 // Cadência lenta de espingarda

		// Upgrades: +1 balote e dispersão mais fechada, +15% dano, +alcance/dano de choque
		p.PelletCount = 6 + save.Shotgunner.SpreadPelletsLevel
		p.PelletSpread = 0.52 - float64(save.Shotgunner.SpreadPelletsLevel)*0.04
		if p.PelletSpread < 0.25 {
			p.PelletSpread = 0.25
		}
		p.PelletDamage = 18.0 * (1.0 + float64(save.Shotgunner.ShrapnelDmgLevel)*0.15)
		p.PelletKnockback = 5.5 + float64(save.Shotgunner.ShrapnelDmgLevel)*0.8

		p.SpecialMaxCD = 480 // ~8s Disparo Sísmico
		p.ShockwaveRadius = 125.0 + float64(save.Shotgunner.ShockwaveLevel)*20.0
		p.ShockwaveDamage = 80.0 * (1.0 + float64(save.Shotgunner.ShockwaveLevel)*0.20)

	case ClassSheriff:
		p.MaxHP = 95
		p.HP = p.MaxHP
		p.BaseSpeed = 3.7
		p.BaseAttackCD = 15

		// Upgrades: recarga mais rápida, dano crítico, balas extras no tambor
		p.MaxAmmo = 6 + save.Sheriff.ExtraAmmoLevel
		p.Ammo = p.MaxAmmo
		p.MaxReloadTime = 72 - save.Sheriff.ReloadSpeedLevel*9 // 1.2s base até ~0.45s
		if p.MaxReloadTime < 25 {
			p.MaxReloadTime = 25
		}
		p.CritChance = 0.20 + float64(save.Sheriff.CritDamageLevel)*0.08
		p.CritMult = 1.8 + float64(save.Sheriff.CritDamageLevel)*0.22

		p.SpecialMaxCD = 600 // ~10s Deadeye
	}

	return p
}

// Update processa movimentação, recarga do Xerife, colisão de bordas, efeitos contínuos e escudo.
func (p *Player) Update(screenWidth, screenHeight float64) {
	p.FrameCounter++

	if p.AttackCD > 0 {
		p.AttackCD--
	}
	if p.SpecialCD > 0 {
		p.SpecialCD--
	}
	if p.DashCD > 0 {
		p.DashCD--
	}
	if p.InvincibleCD > 0 {
		p.InvincibleCD--
	}
	if p.HitFlashCD > 0 {
		p.HitFlashCD--
	}
	if p.MuzzleSmokeTimer > 0 {
		p.MuzzleSmokeTimer--
	}
	if p.AttackAnimTimer > 0 {
		p.AttackAnimTimer--
	}
	p.JustStepped = false

	// Atualiza decaimento das pós-imagens do Dash
	for i := len(p.AfterImages) - 1; i >= 0; i-- {
		p.AfterImages[i].Life--
		if p.AfterImages[i].Life <= 0 {
			p.AfterImages = append(p.AfterImages[:i], p.AfterImages[i+1:]...)
		}
	}

	// Regeneração da Relíquia Coração de Ferro (1.5 HP/segundo por acúmulo)
	if stacks := p.Relics[RelicIronHeart]; stacks > 0 && p.HP < p.MaxHP {
		p.HP += (1.5 * float64(stacks)) / 60.0
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
	}

	// Recarga do tambor do Xerife
	if p.ClassType == ClassSheriff && p.ReloadTimer > 0 {
		p.ReloadTimer--
		if p.ReloadTimer <= 0 {
			p.Ammo = p.MaxAmmo
		}
	}

	// Postura de Escudo do Guerreiro
	if p.ClassType == ClassWarrior && p.Shield.Timer > 0 {
		p.Shield.Timer--
		p.Shield.Angle = p.FacingAngle
		if p.Shield.Timer <= 0 {
			p.Shield.Active = false
		}
	}

	// Atualiza Duração de Efeitos (Buffs / Debuffs)
	if p.SpeedBoostTimer > 0 {
		p.SpeedBoostTimer--
	}
	if p.AgilityTimer > 0 {
		p.AgilityTimer--
	}
	if p.FatigueTimer > 0 {
		p.FatigueTimer--
	}

	if p.CursedChaliceTimer > 0 {
		p.CursedChaliceTimer--
		p.HP -= 4.0 / 60.0
		if p.HP < 1 {
			p.HP = 1
		}
	}

	if p.PoisonTimer > 0 {
		p.PoisonTimer--
		p.HP -= 3.0 / 60.0
		if p.HP < 0 {
			p.HP = 0
		}
	}

	// Regeneração de Mana
	if p.MaxMana > 0 && p.Mana < p.MaxMana {
		p.Mana += p.ManaRegen / 60.0
		if p.Mana > p.MaxMana {
			p.Mana = p.MaxMana
		}
	}

	// Cálculo da Velocidade Efetiva
	speed := p.BaseSpeed
	if p.SpeedBoostTimer > 0 {
		speed *= 1.45
	}
	if p.PoisonTimer > 0 {
		speed *= 0.65
	}

	// Movimentação WASD ou Setas
	var dx, dy float64
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dy -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dy += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dx -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dx += 1
	}

	if dx != 0 && dy != 0 {
		factor := 1.0 / math.Sqrt2
		dx *= factor
		dy *= factor
	}

	p.LastMoveDx = dx
	p.LastMoveDy = dy
	p.IsMoving = dx != 0 || dy != 0

	if p.DashTimer > 0 {
		// Durante o Dash, usa o impulso rápido e deixa pós-imagens translúcidas
		p.DashTimer--
		p.X += p.DashVx
		p.Y += p.DashVy
		if p.DashTimer%2 == 0 {
			cx, cy := p.Center()
			p.AfterImages = append(p.AfterImages, AfterImage{
				Cx:         cx,
				Cy:         cy,
				FacingLeft: p.FacingLeft,
				Sprite:     p.currentSprite(),
				Life:       14,
				MaxLife:    14,
			})
		}
	} else {
		if p.IsMoving {
			p.WalkTimer++
			if p.WalkTimer%13 == 0 {
				p.JustStepped = true
			}
		} else {
			p.WalkTimer = 0
		}

		p.X += dx * speed
		p.Y += dy * speed
	}

	// Restringe aos limites da arena
	hudHeight := 52.0
	if p.X < 14 {
		p.X = 14
	}
	if p.X+p.Width > screenWidth-14 {
		p.X = screenWidth - 14 - p.Width
	}
	if p.Y < hudHeight+12 {
		p.Y = hudHeight + 12
	}
	if p.Y+p.Height > screenHeight-14 {
		p.Y = screenHeight - 14 - p.Height
	}

	mx, my := ebiten.CursorPosition()
	cx, cy := p.Center()
	p.FacingAngle = math.Atan2(float64(my)-cy, float64(mx)-cx)
	p.FacingLeft = float64(mx) < cx
}

// Center retorna o centro do jogador.
func (p *Player) Center() (float64, float64) {
	return p.X + p.Width/2, p.Y + p.Height/2
}

// Bounds retorna a caixa delimitadora AABB.
func (p *Player) Bounds() collision.Rect {
	return collision.NewRect(p.X, p.Y, p.Width, p.Height)
}

// CanAttack verifica se o jogador pode atacar (Xerife não pode durante recarga ou sem balas).
func (p *Player) CanAttack() bool {
	if p.AttackCD > 0 {
		return false
	}
	if p.ClassType == ClassSheriff {
		if p.ReloadTimer > 0 || p.Ammo <= 0 {
			return false
		}
	}
	return true
}

// StartReload inicia a recarga manual ou automática do tambor do Xerife.
func (p *Player) StartReload() {
	if p.ClassType != ClassSheriff {
		return
	}
	if p.ReloadTimer > 0 || p.Ammo >= p.MaxAmmo {
		return
	}
	p.ReloadTimer = p.MaxReloadTime
}

// CanUseSpecial verifica se a habilidade especial está disponível.
func (p *Player) CanUseSpecial() bool {
	if p.SpecialCD > 0 {
		return false
	}
	if p.ClassType == ClassMage && p.Mana < p.FirePillarManaCost {
		return false
	}
	return true
}

// ActivateShield ativa a postura de escudo do Guerreiro.
func (p *Player) ActivateShield() {
	p.Shield.Active = true
	p.Shield.Timer = p.Shield.MaxDuration
	p.Shield.Angle = p.FacingAngle
	p.SpecialCD = p.SpecialMaxCD
}

// CheckShieldReflect verifica e reflete projéteis inimigos.
func (p *Player) CheckShieldReflect(proj *Projectile) bool {
	if !p.Shield.Active || !proj.IsEnemy || !proj.Active {
		return false
	}

	cx, cy := p.Center()
	dx := proj.X - cx
	dy := proj.Y - cy
	dist := math.Hypot(dx, dy)

	if dist > p.Shield.Radius+proj.Radius+10.0 {
		return false
	}

	angleToProj := math.Atan2(dy, dx)
	diff := math.Abs(normalizeAngle(angleToProj - p.Shield.Angle))

	if diff <= p.Shield.ArcSpan/2 {
		proj.IsEnemy = false
		proj.IsHoming = false // Para de perseguir o jogador ao ser refletido
		proj.Vx = -proj.Vx * 1.45
		proj.Vy = -proj.Vy * 1.45
		proj.Damage *= p.Shield.ReflectDamage
		proj.Color = color.RGBA{R: 255, G: 235, B: 120, A: 255}
		return true
	}

	return false
}

// GetEffectiveAttackCD calcula o cooldown considerando Agilidade e Fadiga.
func (p *Player) GetEffectiveAttackCD() int {
	cd := float64(p.BaseAttackCD)
	if p.AgilityTimer > 0 {
		cd *= 0.60
	}
	if p.FatigueTimer > 0 {
		cd *= 1.50
	}
	if cd < 5 {
		cd = 5
	}
	return int(cd)
}

// ResetAttackCD aplica o cooldown de ataque e inicia a animação de impacto/recuo.
func (p *Player) ResetAttackCD() {
	p.AttackCD = p.GetEffectiveAttackCD()
	p.AttackAnimTimer = 10
}

// TriggerDash ativa a esquiva rápida universal com I-Frames e rastro fantasma.
func (p *Player) TriggerDash() bool {
	if p.DashCD > 0 || p.DashTimer > 0 {
		return false
	}

	dirX, dirY := p.LastMoveDx, p.LastMoveDy
	if dirX == 0 && dirY == 0 {
		dirX = math.Cos(p.FacingAngle)
		dirY = math.Sin(p.FacingAngle)
	} else {
		l := math.Hypot(dirX, dirY)
		if l > 0.001 {
			dirX /= l
			dirY /= l
		}
	}

	dashSpeed := 11.5
	p.DashVx = dirX * dashSpeed
	p.DashVy = dirY * dashSpeed
	p.DashTimer = 12 // ~0.20s de deslocamento rápido
	p.DashCD = p.DashMaxCD
	if p.InvincibleCD < 15 {
		p.InvincibleCD = 15 // ~0.25s de I-frames totais!
	}
	return true
}

// AddRelic adiciona um acúmulo de relíquia passiva à partida atual.
func (p *Player) AddRelic(id RelicID) {
	p.Relics[id]++
	p.TotalRelics++
	if id == RelicIronHeart {
		p.MaxHP += 30
		p.HP += 30
	}
}

// Heal aplica cura respeitando o modificador do Pacto de Sangue (-50% cura).
func (p *Player) Heal(amount float64) float64 {
	if p.Relics[RelicBloodPact] > 0 {
		amount *= 0.5
	}
	prev := p.HP
	p.HP += amount
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
	return p.HP - prev
}

// GetMagnetRadius retorna o raio de atração magnética para drops e almas.
func (p *Player) GetMagnetRadius() float64 {
	return 24.0 + float64(p.Relics[RelicSoulMagnet])*95.0
}

// GetPierceChance retorna a chance (0.0 a 1.0) do projétil perfurar o 1o alvo (Lente de Foco).
func (p *Player) GetPierceChance() float64 {
	chance := float64(p.Relics[RelicFocusLens]) * 0.35
	if chance > 1.0 {
		chance = 1.0
	}
	return chance
}

// GetThornsDamage retorna o dano reflexivo do Manto Espinhoso quando o jogador sofre dano.
func (p *Player) GetThornsDamage() float64 {
	stacks := p.Relics[RelicThornyCloak]
	if stacks <= 0 {
		return 0
	}
	return p.MaxHP * 0.30 * float64(stacks)
}

// GetDamageMultiplier retorna o multiplicador de dano (Cálice Amaldiçoado + Meta upgrades + Pacto de Sangue).
func (p *Player) GetDamageMultiplier() float64 {
	mult := p.DamageMultiplier
	if p.CursedChaliceTimer > 0 {
		mult *= 1.65
	}
	if pact := p.Relics[RelicBloodPact]; pact > 0 {
		mult *= 1.0 + 0.25*float64(pact)
	}
	return mult
}

// TakeDamage aplica dano e ativa invulnerabilidade temporária.
func (p *Player) TakeDamage(amount float64) bool {
	if p.InvincibleCD > 0 {
		return false
	}
	if p.Shield.Active {
		return false
	}

	p.HP -= amount
	if p.HP < 0 {
		p.HP = 0
	}
	p.InvincibleCD = 30
	p.HitFlashCD = 2
	return true
}

// ApplyDrop processa a coleta de itens e aplica buffs/debuffs.
func (p *Player) ApplyDrop(d *Drop) (feedback string, clr color.RGBA) {
	info := d.Info()

	switch d.Type {
	case DropSoulFragment:
		return fmt.Sprintf("+%d Fragmento de Alma!", d.SoulAmount), color.RGBA{R: 190, G: 140, B: 255, A: 255}

	case DropHealthPotion:
		healed := p.Heal(p.MaxHP * 0.35)
		return fmt.Sprintf("+%.0f HP!", healed), color.RGBA{R: 80, G: 255, B: 120, A: 255}

	case DropManaPotion:
		if p.MaxMana > 0 {
			p.Mana += 50
			if p.Mana > p.MaxMana {
				p.Mana = p.MaxMana
			}
			return "+50 MP!", color.RGBA{R: 90, G: 180, B: 255, A: 255}
		}
		if p.ClassType == ClassSheriff {
			p.Ammo = p.MaxAmmo
			p.ReloadTimer = 0
			return "TAMBOR RECARREGADO!", color.RGBA{R: 255, G: 220, B: 80, A: 255}
		}
		healed := p.Heal(20)
		return fmt.Sprintf("+%.0f HP!", healed), color.RGBA{R: 120, G: 220, B: 240, A: 255}

	case DropAgilityStone:
		p.AgilityTimer = 600
		return "AGILIDADE! (+Vel. Ataque)", color.RGBA{R: 255, G: 220, B: 40, A: 255}

	case DropSpeedBoots:
		p.SpeedBoostTimer = 600
		return "IMPETO! (+Vel. Movimento)", color.RGBA{R: 80, G: 220, B: 255, A: 255}

	case DropMasteryTome:
		p.MaxHP += 25
		p.HP += 25
		if p.MaxMana > 0 {
			p.MaxMana += 25
			p.Mana += 25
		}
		return "MAESTRIA! (+Max HP/MP)", color.RGBA{R: 220, G: 140, B: 255, A: 255}

	case DropCursedChalice:
		p.CursedChaliceTimer = 480
		return "CALICE AMALDICOADO! (+65% Dano, -HP/s)", color.RGBA{R: 255, G: 60, B: 60, A: 255}

	case DropPoisonVial:
		p.PoisonTimer = 360
		return "ENVENENADO! (Lentidao e Veneno)", color.RGBA{R: 170, G: 80, B: 240, A: 255}

	case DropFatigueMist:
		p.FatigueTimer = 480
		if p.MaxMana > 0 {
			p.Mana = 0
		}
		return "FADIGA! (Ataques Lentos)", color.RGBA{R: 180, G: 100, B: 220, A: 255}

	case DropChaosSeed:
		return "SEMENTE DO CAOS! (MORCEGOS!)", color.RGBA{R: 255, G: 40, B: 40, A: 255}
	}

	return info.Name, info.Color
}

// GetActiveEffects retorna a lista de buffs/debuffs para o HUD.
func (p *Player) GetActiveEffects() []ActiveEffect {
	var list []ActiveEffect
	if p.SpeedBoostTimer > 0 {
		list = append(list, ActiveEffect{Name: "Ímpeto", IsDebuff: false, Remaining: p.SpeedBoostTimer, Total: 600, Color: color.RGBA{R: 80, G: 220, B: 255, A: 255}})
	}
	if p.AgilityTimer > 0 {
		list = append(list, ActiveEffect{Name: "Agilidade", IsDebuff: false, Remaining: p.AgilityTimer, Total: 600, Color: color.RGBA{R: 255, G: 220, B: 40, A: 255}})
	}
	if p.CursedChaliceTimer > 0 {
		list = append(list, ActiveEffect{Name: "Cálice Maldito", IsDebuff: true, Remaining: p.CursedChaliceTimer, Total: 480, Color: color.RGBA{R: 255, G: 60, B: 60, A: 255}})
	}
	if p.PoisonTimer > 0 {
		list = append(list, ActiveEffect{Name: "Veneno", IsDebuff: true, Remaining: p.PoisonTimer, Total: 360, Color: color.RGBA{R: 60, G: 230, B: 80, A: 255}})
	}
	if p.FatigueTimer > 0 {
		list = append(list, ActiveEffect{Name: "Fadiga", IsDebuff: true, Remaining: p.FatigueTimer, Total: 480, Color: color.RGBA{R: 180, G: 100, B: 220, A: 255}})
	}
	return list
}

func (p *Player) currentSprite() *ebiten.Image {
	switch p.ClassType {
	case ClassWarrior:
		return p.spriteWarrior
	case ClassArcher:
		if p.AttackCD > p.BaseAttackCD-6 {
			return p.spriteArcherShoot
		}
		return p.spriteArcherRest
	case ClassMage:
		return p.spriteMage
	case ClassShotgunner:
		return p.spriteShotgunner
	case ClassSheriff:
		return p.spriteSheriff
	}
	return p.spriteWarrior
}

// Draw renderiza o jogador, pós-imagens do Dash, filtro de veneno, barra de recarga e arma ativa.
func (p *Player) Draw(dst *ebiten.Image) {
	// Renderiza pós-imagens fantasmas do Dash sempre visíveis
	for _, ai := range p.AfterImages {
		if ai.Sprite == nil {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		sw := float64(ai.Sprite.Bounds().Dx())
		sh := float64(ai.Sprite.Bounds().Dy())
		op.GeoM.Translate(-sw/2, -sh/2)
		if ai.FacingLeft {
			op.GeoM.Scale(-2.0, 2.0)
		} else {
			op.GeoM.Scale(2.0, 2.0)
		}
		op.GeoM.Translate(ai.Cx, ai.Cy)
		alpha := float32(ai.Life) / float32(ai.MaxLife) * 0.65
		op.ColorScale.Scale(0.45, 0.85, 1.35, alpha)
		dst.DrawImage(ai.Sprite, op)
	}

	if p.HitFlashCD <= 0 && p.DashTimer <= 0 && p.InvincibleCD > 0 && (p.InvincibleCD/4)%2 == 0 {
		return
	}

	cx, cy := p.Center()

	// Calcula parâmetros de animação procedural (caminhada vs repouso vs recuo de ataque)
	var bobY, tilt, stretchX, stretchY float64
	stretchX = 1.0
	stretchY = 1.0

	if p.DashTimer > 0 {
		// Estica na direção da esquiva
		stretchX = 1.18
		stretchY = 0.85
	} else if p.IsMoving {
		// Ciclo de caminhada rítmico e saltinho
		phase := float64(p.WalkTimer) * 0.24
		bobY = -math.Abs(math.Sin(phase)) * 4.5
		tilt = math.Sin(phase) * 0.11
		squash := math.Sin(phase*2) * 0.08
		stretchX = 1.0 - squash*0.7
		stretchY = 1.0 + squash
	} else {
		// Respiração suave em repouso
		breath := math.Sin(float64(p.FrameCounter)*0.08) * 0.05
		stretchX = 1.0 - breath*0.5
		stretchY = 1.0 + breath
	}

	// Recuo e impacto ao atacar
	var recoilOffsetX, recoilOffsetY float64
	if p.AttackAnimTimer > 0 {
		recoilRatio := float64(p.AttackAnimTimer) / 10.0
		recoilOffsetX = -math.Cos(p.FacingAngle) * 3.5 * recoilRatio
		recoilOffsetY = -math.Sin(p.FacingAngle) * 2.5 * recoilRatio
		if p.FacingLeft {
			tilt += 0.16 * recoilRatio
		} else {
			tilt -= 0.16 * recoilRatio
		}
		stretchX *= 1.0 + 0.10*recoilRatio
		stretchY *= 1.0 - 0.08*recoilRatio
	}

	// Sombra circular sob o herói (encolhe levemente no auge do salto)
	shadowRad := float32(13.0 - (-bobY)*0.5)
	vector.DrawFilledCircle(dst, float32(cx), float32(cy+14), shadowRad, color.RGBA{R: 10, G: 10, B: 18, A: 140}, true)

	spriteImg := p.currentSprite()

	if spriteImg != nil {
		op := &ebiten.DrawImageOptions{}
		baseScale := 2.0
		sw := float64(spriteImg.Bounds().Dx())
		sh := float64(spriteImg.Bounds().Dy())

		// 1. Move origem para o centro do sprite
		op.GeoM.Translate(-sw/2, -sh/2)

		// 2. Aplica flip horizontal e squash & stretch
		sx := baseScale * stretchX
		sy := baseScale * stretchY
		if p.FacingLeft {
			sx = -sx
		}
		op.GeoM.Scale(sx, sy)

		// 3. Aplica rotação de balanço / recuo
		op.GeoM.Rotate(tilt)

		// 4. Translada para a posição final na arena
		op.GeoM.Translate(cx+recoilOffsetX, cy+bobY+recoilOffsetY)

		// Flash branco de impacto ou filtro esverdeado de veneno
		if p.HitFlashCD > 0 {
			op.ColorScale.Scale(255, 255, 255, 1.0)
		} else if p.PoisonTimer > 0 {
			pulse := float32(0.35 + math.Sin(float64(p.FrameCounter)*0.2)*0.25)
			op.ColorScale.Scale(0.6, 1.3+pulse, 0.6, 1.0)
		}

		dst.DrawImage(spriteImg, op)
	}

	// Barra de Recarga Circular / Arco sobre o Xerife durante recarga
	if p.ClassType == ClassSheriff && p.ReloadTimer > 0 {
		progress := 1.0 - float64(p.ReloadTimer)/float64(p.MaxReloadTime)
		barW := float32(34)
		barH := float32(5)
		bx := float32(cx) - barW/2
		by := float32(p.Y - 12)
		vector.DrawFilledRect(dst, bx, by, barW, barH, color.RGBA{30, 25, 20, 220}, false)
		vector.DrawFilledRect(dst, bx, by, barW*float32(progress), barH, color.RGBA{255, 215, 60, 255}, false)
		vector.StrokeRect(dst, bx, by, barW, barH, 1, color.RGBA{200, 180, 140, 255}, false)
	}

	// Postura de Escudo do Guerreiro
	if p.ClassType == ClassWarrior && p.Shield.Active {
		p.drawShieldBarrier(dst, cx, cy)
	}

	// Detalhes procedurais dinâmicos de cada classe
	switch p.ClassType {
	case ClassMage:
		pulse := math.Sin(float64(p.FrameCounter)*0.15) * 2.0
		crystalRadius := float32(4.5 + pulse)
		var crystalX, crystalY float32
		if p.FacingLeft {
			crystalX = float32(p.X + 8)
			crystalY = float32(p.Y + 8)
		} else {
			crystalX = float32(p.X + p.Width - 8)
			crystalY = float32(p.Y + 8)
		}
		vector.DrawFilledCircle(dst, crystalX, crystalY, crystalRadius*1.8, color.RGBA{R: 100, G: 200, B: 255, A: 90}, true)
		vector.DrawFilledCircle(dst, crystalX, crystalY, crystalRadius, color.RGBA{R: 160, G: 240, B: 255, A: 255}, true)
		vector.DrawFilledCircle(dst, crystalX, crystalY, crystalRadius*0.4, color.RGBA{R: 255, G: 255, B: 255, A: 255}, true)
	}
}

func (p *Player) drawShieldBarrier(dst *ebiten.Image, cx, cy float64) {
	steps := 16
	startAngle := p.Shield.Angle - p.Shield.ArcSpan/2
	stepAngle := p.Shield.ArcSpan / float64(steps)
	rad := p.Shield.Radius

	for i := 0; i < steps; i++ {
		a1 := startAngle + float64(i)*stepAngle
		a2 := startAngle + float64(i+1)*stepAngle

		x1 := float32(cx + rad*math.Cos(a1))
		y1 := float32(cy + rad*math.Sin(a1))
		x2 := float32(cx + rad*math.Cos(a2))
		y2 := float32(cy + rad*math.Sin(a2))

		vector.StrokeLine(dst, x1, y1, x2, y2, 4.5, color.RGBA{R: 255, G: 215, B: 70, A: 240}, true)
		vector.StrokeLine(dst, x1, y1, x2, y2, 2.0, color.RGBA{R: 255, G: 255, B: 240, A: 255}, true)

		x1Halo := float32(cx + (rad+4)*math.Cos(a1))
		y1Halo := float32(cy + (rad+4)*math.Sin(a1))
		x2Halo := float32(cx + (rad+4)*math.Cos(a2))
		y2Halo := float32(cy + (rad+4)*math.Sin(a2))
		vector.StrokeLine(dst, x1Halo, y1Halo, x2Halo, y2Halo, 2.0, color.RGBA{R: 120, G: 210, B: 255, A: 120}, true)
	}
}
