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
	ClassWarrior ClassType = "Guerreiro"
	ClassArcher  ClassType = "Arqueiro"
	ClassMage    ClassType = "Mago"
)

// ActiveEffect representa um buff ou debuff ativo no jogador.
type ActiveEffect struct {
	Name      string
	IsDebuff  bool
	Remaining int // frames restantes
	Total     int
	Color     color.RGBA
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
	InvincibleCD int
	FacingAngle  float64
	FacingLeft   bool
	FrameCounter int

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

	// Habilidade Arqueiro: Parâmetros da Chuva de Flechas
	ArrowRainRadius float64
	ArrowRainDamage float64

	// Habilidade Mago: Parâmetros do Pilar de Fogo
	FirePillarRadius   float64
	FirePillarDuration int
	FirePillarDamage   float64
	FirePillarManaCost float64

	// Buffs e Debuffs (Efeitos de Drops)
	SpeedBoostTimer    int
	AgilityTimer       int
	CursedChaliceTimer int
	PoisonTimer        int
	FatigueTimer       int

	// Sprites procedurais cacheados
	spriteWarrior     *ebiten.Image
	spriteArcherRest  *ebiten.Image
	spriteArcherShoot *ebiten.Image
	spriteMage        *ebiten.Image
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
		spriteWarrior:     sprites.GetWarriorSprite(),
		spriteArcherRest:  sprites.GetArcherSprite(false),
		spriteArcherShoot: sprites.GetArcherSprite(true),
		spriteMage:        sprites.GetMageSprite(),
	}

	switch class {
	case ClassWarrior:
		// Meta Upgrades: Vida (+25/nvl), Dano/Alcance (+12% e +6px/nvl), Escudo (-1s recarga, +0.25s duração/nvl)
		bonusHP := float64(save.Warrior.MaxHPLevel) * 25.0
		p.MaxHP = 200 + bonusHP
		p.HP = p.MaxHP
		p.MaxMana = 0
		p.Mana = 0
		p.ManaRegen = 0
		p.BaseSpeed = 3.2
		p.BaseAttackCD = 22

		p.DamageMultiplier = 1.0 + float64(save.Warrior.DamageLevel)*0.12
		p.SlashRange = 72.0 + float64(save.Warrior.DamageLevel)*6.0

		// Postura de Escudo
		shieldCD := 420 - save.Warrior.ShieldLevel*45 // ~7s a ~3.5s
		if shieldCD < 180 {
			shieldCD = 180
		}
		shieldDur := 75 + save.Warrior.ShieldLevel*15 // ~1.25s a ~2.5s
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
		// Meta Upgrades: Velocidade (+0.35/nvl), Recarga/Projétil (+10% recarga, +1.2 vel/nvl), Chuva de Flechas (+15% raio, +20% dano/nvl)
		bonusSpeed := float64(save.Archer.SpeedLevel) * 0.35
		p.MaxHP = 75
		p.HP = p.MaxHP
		p.MaxMana = 0
		p.Mana = 0
		p.ManaRegen = 0
		p.BaseSpeed = 4.0 + bonusSpeed
		p.BaseAttackCD = 13 - int(float64(save.Archer.AttackSpeedLvl)*0.8)
		if p.BaseAttackCD < 7 {
			p.BaseAttackCD = 7
		}
		p.ArrowSpeed = 9.5 + float64(save.Archer.AttackSpeedLvl)*1.2

		// Chuva de Flechas
		p.SpecialMaxCD = 460 // ~7.6 segundos
		p.ArrowRainRadius = 70.0 * (1.0 + float64(save.Archer.ArrowRainLevel)*0.15)
		p.ArrowRainDamage = 13.0 * (1.0 + float64(save.Archer.ArrowRainLevel)*0.20)

	case ClassMage:
		// Meta Upgrades: Mana/Regen (+20 MP, +3.5 regen/nvl), Poder Mágico (+15%/nvl), Pilar de Fogo (+15% raio, +20 frames/nvl)
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

		// Pilar de Fogo
		p.SpecialMaxCD = 360 // ~6 segundos
		p.FirePillarRadius = 55.0 * (1.0 + float64(save.Mage.FirePillarLevel)*0.15)
		p.FirePillarDuration = 90 + save.Mage.FirePillarLevel*20
		p.FirePillarDamage = 24.0 * (1.0 + float64(save.Mage.SpellPowerLevel)*0.15)
		p.FirePillarManaCost = 30.0
	}

	return p
}

// Update processa movimentação, colisão de bordas, efeitos contínuos, escudo e regeneração.
func (p *Player) Update(screenWidth, screenHeight float64) {
	p.FrameCounter++

	// 1. Redução de cooldowns de ataque, especial e invencibilidade
	if p.AttackCD > 0 {
		p.AttackCD--
	}
	if p.SpecialCD > 0 {
		p.SpecialCD--
	}
	if p.InvincibleCD > 0 {
		p.InvincibleCD--
	}

	// 2. Atualiza estado da Postura de Escudo do Guerreiro
	if p.ClassType == ClassWarrior {
		if p.Shield.Timer > 0 {
			p.Shield.Timer--
			p.Shield.Angle = p.FacingAngle
			if p.Shield.Timer <= 0 {
				p.Shield.Active = false
			}
		}
	}

	// 3. Atualiza Duração de Efeitos (Buffs / Debuffs)
	if p.SpeedBoostTimer > 0 {
		p.SpeedBoostTimer--
	}
	if p.AgilityTimer > 0 {
		p.AgilityTimer--
	}
	if p.FatigueTimer > 0 {
		p.FatigueTimer--
	}

	// Cálice Amaldiçoado: drena vida com o tempo (4 HP/segundo)
	if p.CursedChaliceTimer > 0 {
		p.CursedChaliceTimer--
		p.HP -= 4.0 / 60.0
		if p.HP < 1 {
			p.HP = 1
		}
	}

	// Frasco de Veneno: causa dano por veneno contínuo (3 HP/segundo)
	if p.PoisonTimer > 0 {
		p.PoisonTimer--
		p.HP -= 3.0 / 60.0
		if p.HP < 0 {
			p.HP = 0
		}
	}

	// 4. Regeneração de Mana (60 frames por segundo)
	if p.MaxMana > 0 && p.Mana < p.MaxMana {
		p.Mana += p.ManaRegen / 60.0
		if p.Mana > p.MaxMana {
			p.Mana = p.MaxMana
		}
	}

	// 5. Cálculo da Velocidade Efetiva com Buffs/Debuffs
	speed := p.BaseSpeed
	if p.SpeedBoostTimer > 0 {
		speed *= 1.45
	}
	if p.PoisonTimer > 0 {
		speed *= 0.65
	}

	// 6. Movimentação WASD ou Setas
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

	p.X += dx * speed
	p.Y += dy * speed

	// Restringe aos limites da arena
	hudHeight := 52.0
	if p.X < 12 {
		p.X = 12
	}
	if p.X+p.Width > screenWidth-12 {
		p.X = screenWidth - 12 - p.Width
	}
	if p.Y < hudHeight {
		p.Y = hudHeight
	}
	if p.Y+p.Height > screenHeight-12 {
		p.Y = screenHeight - 12 - p.Height
	}

	// 7. Direção do mouse e espelhamento horizontal
	mx, my := ebiten.CursorPosition()
	cx, cy := p.Center()
	p.FacingAngle = math.Atan2(float64(my)-cy, float64(mx)-cx)
	p.FacingLeft = float64(mx) < cx
}

// Center retorna as coordenadas centrais do jogador.
func (p *Player) Center() (float64, float64) {
	return p.X + p.Width/2, p.Y + p.Height/2
}

// Bounds retorna a caixa delimitadora AABB do jogador.
func (p *Player) Bounds() collision.Rect {
	return collision.NewRect(p.X, p.Y, p.Width, p.Height)
}

// CanAttack verifica se o jogador pode golpear ou atirar.
func (p *Player) CanAttack() bool {
	return p.AttackCD <= 0
}

// CanUseSpecial verifica se a habilidade especial está fora de cooldown (e com mana suficiente se Mago).
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

// CheckShieldReflect verifica se um projétil inimigo colidiu com o arco do escudo frontal.
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

	// Está dentro do arco frontal do escudo
	if diff <= p.Shield.ArcSpan/2 {
		// Reflete o projétil!
		proj.IsEnemy = false
		proj.Vx = -proj.Vx * 1.45
		proj.Vy = -proj.Vy * 1.45
		proj.Damage *= p.Shield.ReflectDamage
		proj.Color = color.RGBA{R: 255, G: 235, B: 120, A: 255} // Converte em projétil radiante
		return true
	}

	return false
}

// GetEffectiveAttackCD calcula o cooldown do ataque considerando Agilidade e Fadiga.
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

// ResetAttackCD aplica o tempo de recarga após um ataque.
func (p *Player) ResetAttackCD() {
	p.AttackCD = p.GetEffectiveAttackCD()
}

// GetDamageMultiplier retorna o modificador de dano (Cálice Amaldiçoado + Meta upgrades).
func (p *Player) GetDamageMultiplier() float64 {
	mult := p.DamageMultiplier
	if p.CursedChaliceTimer > 0 {
		mult *= 1.65
	}
	return mult
}

// TakeDamage aplica dano se não estiver invulnerável ou defendendo com o escudo frontalmente.
func (p *Player) TakeDamage(amount float64) bool {
	if p.InvincibleCD > 0 {
		return false
	}
	// Se escudo estiver ativo, absorve dano
	if p.Shield.Active {
		return false
	}

	p.HP -= amount
	if p.HP < 0 {
		p.HP = 0
	}
	p.InvincibleCD = 30
	return true
}

// ApplyDrop processa a coleta de um item e aplica seu buff/debuff.
func (p *Player) ApplyDrop(d *Drop) (feedback string, clr color.RGBA) {
	info := d.Info()

	switch d.Type {
	case DropSoulFragment:
		return fmt.Sprintf("+%d Fragmento de Alma!", d.SoulAmount), color.RGBA{R: 190, G: 140, B: 255, A: 255}

	case DropHealthPotion:
		heal := p.MaxHP * 0.35
		p.HP += heal
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
		return fmt.Sprintf("+%.0f HP!", heal), color.RGBA{R: 80, G: 255, B: 120, A: 255}

	case DropManaPotion:
		if p.MaxMana > 0 {
			p.Mana += 50
			if p.Mana > p.MaxMana {
				p.Mana = p.MaxMana
			}
			return "+50 MP!", color.RGBA{R: 90, G: 180, B: 255, A: 255}
		}
		p.HP += 20
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
		return "+20 HP (Sem Mana)", color.RGBA{R: 120, G: 220, B: 240, A: 255}

	case DropAgilityStone:
		p.AgilityTimer = 600
		return "AGILIDADE! (+Vel. Ataque)", color.RGBA{R: 255, G: 220, B: 40, A: 255}

	case DropSpeedBoots:
		p.SpeedBoostTimer = 600
		return "IMPÉTO! (+Vel. Movimento)", color.RGBA{R: 80, G: 220, B: 255, A: 255}

	case DropMasteryTome:
		p.MaxHP += 25
		p.HP += 25
		if p.MaxMana > 0 {
			p.MaxMana += 25
			p.Mana += 25
		}
		return "MAESTRIA! (+Max HP/MP)", color.RGBA{R: 220, G: 140, B: 255, A: 255}

	// Malignos
	case DropCursedChalice:
		p.CursedChaliceTimer = 480
		return "CALICE AMALDICOADO! (+65% Dano, -HP/s)", color.RGBA{R: 255, G: 60, B: 60, A: 255}

	case DropPoisonVial:
		p.PoisonTimer = 360
		return "ENVENENADO! (Lentidao e Dano)", color.RGBA{R: 60, G: 220, B: 60, A: 255}

	case DropFatigueMist:
		p.FatigueTimer = 480
		if p.MaxMana > 0 {
			p.Mana = 0
		}
		return "FADIGA! (Ataques Lentos e Mana Zerada)", color.RGBA{R: 180, G: 100, B: 220, A: 255}

	case DropChaosSeed:
		return "SEMENTE DO CAOS! (Emboscada!)", color.RGBA{R: 255, G: 40, B: 20, A: 255}
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
		list = append(list, ActiveEffect{Name: "Veneno", IsDebuff: true, Remaining: p.PoisonTimer, Total: 360, Color: color.RGBA{R: 50, G: 220, B: 60, A: 255}})
	}
	if p.FatigueTimer > 0 {
		list = append(list, ActiveEffect{Name: "Fadiga", IsDebuff: true, Remaining: p.FatigueTimer, Total: 480, Color: color.RGBA{R: 180, G: 100, B: 220, A: 255}})
	}
	return list
}

// Draw renderiza o personagem, seus sprites de Pixel Art e efeitos ativos como a Postura de Escudo.
func (p *Player) Draw(dst *ebiten.Image) {
	if p.InvincibleCD > 0 && (p.InvincibleCD/4)%2 == 0 {
		return
	}

	cx, cy := p.Center()

	// Sombra circular suave sob o herói
	vector.DrawFilledCircle(dst, float32(cx), float32(cy+14), 13, color.RGBA{R: 15, G: 15, B: 25, A: 130}, true)

	// Seleção do sprite de acordo com a classe
	var spriteImg *ebiten.Image
	switch p.ClassType {
	case ClassWarrior:
		spriteImg = p.spriteWarrior
	case ClassArcher:
		if p.AttackCD > p.BaseAttackCD-6 {
			spriteImg = p.spriteArcherShoot
		} else {
			spriteImg = p.spriteArcherRest
		}
	case ClassMage:
		spriteImg = p.spriteMage
	}

	if spriteImg != nil {
		op := &ebiten.DrawImageOptions{}
		scale := 2.0
		sw := float64(spriteImg.Bounds().Dx())

		if p.FacingLeft {
			op.GeoM.Scale(-scale, scale)
			op.GeoM.Translate(p.X+sw*scale, p.Y)
		} else {
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(p.X, p.Y)
		}

		dst.DrawImage(spriteImg, op)
	}

	// Desenho da Postura de Escudo do Guerreiro se ativa
	if p.ClassType == ClassWarrior && p.Shield.Active {
		p.drawShieldBarrier(dst, cx, cy)
	}

	// Detalhes procedurais dinâmicos para cada classe
	switch p.ClassType {
	case ClassMage:
		// Animação de pulso no cristal do cajado
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

	case ClassArcher:
		pointerLen := 22.0
		mx := float32(cx + pointerLen*math.Cos(p.FacingAngle))
		my := float32(cy + pointerLen*math.Sin(p.FacingAngle))
		vector.StrokeLine(dst, float32(cx), float32(cy), mx, my, 1.5, color.RGBA{R: 120, G: 220, B: 120, A: 160}, true)
	}
}

// drawShieldBarrier desenha a barreira frontal em arco com runas e reflexão visual.
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

		// Arco externo de barreira dourada de ferro divino
		vector.StrokeLine(dst, x1, y1, x2, y2, 4.5, color.RGBA{R: 255, G: 215, B: 70, A: 240}, true)
		// Núcleo luminoso branco
		vector.StrokeLine(dst, x1, y1, x2, y2, 2.0, color.RGBA{R: 255, G: 255, B: 240, A: 255}, true)

		// Halo de reflexão
		x1Halo := float32(cx + (rad+4)*math.Cos(a1))
		y1Halo := float32(cy + (rad+4)*math.Sin(a1))
		x2Halo := float32(cx + (rad+4)*math.Cos(a2))
		y2Halo := float32(cy + (rad+4)*math.Sin(a2))
		vector.StrokeLine(dst, x1Halo, y1Halo, x2Halo, y2Halo, 2.0, color.RGBA{R: 120, G: 210, B: 255, A: 120}, true)
	}
}
