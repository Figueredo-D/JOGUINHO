package entities

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
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

// Player representa o herói controlado pelo jogador.
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

// NewPlayer inicializa o jogador com atributos e sprites procedurais.
func NewPlayer(class ClassType, startX, startY float64) *Player {
	p := &Player{
		ClassType:         class,
		X:                 startX,
		Y:                 startY,
		Width:             32,
		Height:            32,
		spriteWarrior:     sprites.GetWarriorSprite(),
		spriteArcherRest:  sprites.GetArcherSprite(false),
		spriteArcherShoot: sprites.GetArcherSprite(true),
		spriteMage:        sprites.GetMageSprite(),
	}

	switch class {
	case ClassWarrior:
		p.MaxHP = 200
		p.HP = 200
		p.MaxMana = 0
		p.Mana = 0
		p.ManaRegen = 0
		p.BaseSpeed = 3.2
		p.BaseAttackCD = 22

	case ClassArcher:
		p.MaxHP = 75
		p.HP = 75
		p.MaxMana = 0
		p.Mana = 0
		p.ManaRegen = 0
		p.BaseSpeed = 4.0
		p.BaseAttackCD = 13

	case ClassMage:
		p.MaxHP = 110
		p.HP = 110
		p.MaxMana = 100
		p.Mana = 100
		p.ManaRegen = 16
		p.BaseSpeed = 3.0
		p.BaseAttackCD = 20
	}

	return p
}

// Update processa movimentação, colisão de bordas, efeitos contínuos e regeneração.
func (p *Player) Update(screenWidth, screenHeight float64) {
	p.FrameCounter++

	// 1. Redução de cooldowns de ataque e invencibilidade
	if p.AttackCD > 0 {
		p.AttackCD--
	}
	if p.InvincibleCD > 0 {
		p.InvincibleCD--
	}

	// 2. Atualiza Duração de Efeitos (Buffs / Debuffs)
	if p.SpeedBoostTimer > 0 {
		p.SpeedBoostTimer--
	}
	if p.AgilityTimer > 0 {
		p.AgilityTimer--
	}
	if p.FatigueTimer > 0 {
		p.FatigueTimer--
	}

	// Cálice Amaldiçoado: drena vida com o tempo (4 HP/segundo = ~0.066 HP/frame)
	if p.CursedChaliceTimer > 0 {
		p.CursedChaliceTimer--
		p.HP -= 4.0 / 60.0
		if p.HP < 1 {
			p.HP = 1 // Não mata sozinho pelo cálice para não frustrar o jogador
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

	// 3. Regeneração de Mana (60 frames por segundo)
	if p.MaxMana > 0 && p.Mana < p.MaxMana {
		p.Mana += p.ManaRegen / 60.0
		if p.Mana > p.MaxMana {
			p.Mana = p.MaxMana
		}
	}

	// 4. Cálculo da Velocidade Efetiva com Buffs/Debuffs
	speed := p.BaseSpeed
	if p.SpeedBoostTimer > 0 {
		speed *= 1.45 // +45% de velocidade
	}
	if p.PoisonTimer > 0 {
		speed *= 0.65 // -35% de lentidão
	}

	// 5. Movimentação WASD ou Setas
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

	// 6. Direção do mouse e espelhamento horizontal
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

// GetEffectiveAttackCD calcula o cooldown do ataque considerando Agilidade e Fadiga.
func (p *Player) GetEffectiveAttackCD() int {
	cd := float64(p.BaseAttackCD)
	if p.AgilityTimer > 0 {
		cd *= 0.60 // -40% de recarga (ataque muito mais rápido)
	}
	if p.FatigueTimer > 0 {
		cd *= 1.50 // +50% de recarga
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

// GetDamageMultiplier retorna o modificador de dano (Cálice Amaldiçoado).
func (p *Player) GetDamageMultiplier() float64 {
	mult := 1.0
	if p.CursedChaliceTimer > 0 {
		mult *= 1.65 // +65% de dano brutal!
	}
	return mult
}

// TakeDamage aplica dano ao jogador e ativa invulnerabilidade temporária.
func (p *Player) TakeDamage(amount float64) bool {
	if p.InvincibleCD > 0 {
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
		// Se classe sem mana, converte em pequena cura
		p.HP += 20
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
		return "+20 HP (Sem Mana)", color.RGBA{R: 120, G: 220, B: 240, A: 255}

	case DropAgilityStone:
		p.AgilityTimer = 600 // 10 segundos
		return "AGILIDADE! (+Vel. Ataque)", color.RGBA{R: 255, G: 220, B: 40, A: 255}

	case DropSpeedBoots:
		p.SpeedBoostTimer = 600 // 10 segundos
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
		p.CursedChaliceTimer = 480 // 8 segundos
		return "CALICE AMALDICOADO! (+65% Dano, -HP/s)", color.RGBA{R: 255, G: 60, B: 60, A: 255}

	case DropPoisonVial:
		p.PoisonTimer = 360 // 6 segundos
		return "ENVENENADO! (Lentidao e Dano)", color.RGBA{R: 60, G: 220, B: 60, A: 255}

	case DropFatigueMist:
		p.FatigueTimer = 480 // 8 segundos
		if p.MaxMana > 0 {
			p.Mana = 0
		}
		return "FADIGA! (Ataques Lentos e Mana Zerada)", color.RGBA{R: 180, G: 100, B: 220, A: 255}

	case DropChaosSeed:
		// Retorno especial tratado na GameplayScene para spawn de inimigos extras
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

// Draw renderiza o personagem com seu sprite de Pixel Art procedural e detalhes visuais da classe.
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
		// Se estiver no cooldown de ataque, usa sprite tensionando o arco
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
		// Escala 2x para pixel art 16x16 virar 32x32 nítido
		scale := 2.0
		sw := float64(spriteImg.Bounds().Dx())

		// Se o personagem estiver olhando para a esquerda, espelha horizontalmente
		if p.FacingLeft {
			op.GeoM.Scale(-scale, scale)
			op.GeoM.Translate(p.X+sw*scale, p.Y)
		} else {
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(p.X, p.Y)
		}

		dst.DrawImage(spriteImg, op)
	}

	// Detalhes procedurais dinâmicos para cada classe:
	switch p.ClassType {
	case ClassMage:
		// Animação de pulso e brilho no cristal do topo do cajado!
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
		// Halo de luz arcana
		vector.DrawFilledCircle(dst, crystalX, crystalY, crystalRadius*1.8, color.RGBA{R: 100, G: 200, B: 255, A: 90}, true)
		vector.DrawFilledCircle(dst, crystalX, crystalY, crystalRadius, color.RGBA{R: 160, G: 240, B: 255, A: 255}, true)
		vector.DrawFilledCircle(dst, crystalX, crystalY, crystalRadius*0.4, color.RGBA{R: 255, G: 255, B: 255, A: 255}, true)

	case ClassArcher:
		// Indicador de mira sutil para arco
		pointerLen := 22.0
		mx := float32(cx + pointerLen*math.Cos(p.FacingAngle))
		my := float32(cy + pointerLen*math.Sin(p.FacingAngle))
		vector.StrokeLine(dst, float32(cx), float32(cy), mx, my, 1.5, color.RGBA{R: 120, G: 220, B: 120, A: 160}, true)
	}
}
