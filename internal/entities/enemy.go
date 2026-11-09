package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
	"joguinho/internal/sprites"
)

// EnemyType define os 4 arquétipos exigidos.
type EnemyType int

const (
	EnemyTypeCommon EnemyType = iota // 1. Comum (Zumbi/Goblin)
	EnemyTypeRunner                  // 2. Corredor (Cão Infernal veloz)
	EnemyTypeTank                    // 3. Tanque (Golem colossal de pedra)
	EnemyTypeShooter                 // 4. Atirador (Bombardeiro que mantém distância e dispara projéteis)
)

// Enemy representa um inimigo na arena.
type Enemy struct {
	ID            uint64
	Type          EnemyType
	X, Y          float64
	Width         float64
	Height        float64
	Speed         float64
	HP            float64
	MaxHP         float64
	Damage        float64
	ScoreValue    int
	AttackCD      int // Cooldown de dano por contato
	ShootCD       int // Cooldown de disparo para o atirador
	HitFlashCD    int
	Alive         bool
	FacingLeft    bool
	FrameCount    int

	spriteImg     *ebiten.Image
}

// NewEnemy inicializa o inimigo com base no arquétipo e na escala de onda.
func NewEnemy(id uint64, enemyType EnemyType, startX, startY float64, waveMultiplier float64) *Enemy {
	e := &Enemy{
		ID:         id,
		Type:       enemyType,
		X:          startX,
		Y:          startY,
		Alive:      true,
		FrameCount: 0,
	}

	switch enemyType {
	case EnemyTypeCommon:
		// 1. Comum (Zumbi/Goblin): Vida e velocidade médias
		e.Width = 32
		e.Height = 32
		e.MaxHP = 38 * waveMultiplier
		e.HP = e.MaxHP
		e.Speed = 1.9 + (waveMultiplier-1)*0.15
		e.Damage = 12
		e.ScoreValue = 100
		e.spriteImg = sprites.GetZombieSprite()

	case EnemyTypeRunner:
		// 2. Corredor/Rápido: Vida muito baixa, velocidade extrema (2.5x a 3x o normal)
		e.Width = 30
		e.Height = 24
		e.MaxHP = 16 * waveMultiplier
		e.HP = e.MaxHP
		e.Speed = 5.2 + (waveMultiplier-1)*0.2 // Extremamente veloz!
		e.Damage = 8
		e.ScoreValue = 160
		e.spriteImg = sprites.GetHellhoundSprite()

	case EnemyTypeTank:
		// 3. Tanque/Bruto: Visual avantajado (48x48), lento, ~10x vida e alto dano
		e.Width = 48
		e.Height = 48
		e.MaxHP = 360 * waveMultiplier
		e.HP = e.MaxHP
		e.Speed = 1.0 + (waveMultiplier-1)*0.08 // Metade da velocidade base
		e.Damage = 35
		e.ScoreValue = 450
		e.spriteImg = sprites.GetGolemSprite()

	case EnemyTypeShooter:
		// 4. Atirador/Bombardeiro: Mantém distância média e cospe projéteis
		e.Width = 32
		e.Height = 32
		e.MaxHP = 45 * waveMultiplier
		e.HP = e.MaxHP
		e.Speed = 1.6
		e.Damage = 14
		e.ScoreValue = 220
		e.ShootCD = 100 // Atira a cada ~1.6 segundos
		e.spriteImg = sprites.GetShooterSprite()
	}

	return e
}

// Update executa a inteligência artificial do inimigo.
// Retorna um novo projétil se o inimigo atirador disparar neste frame.
func (e *Enemy) Update(targetX, targetY float64) *Projectile {
	if !e.Alive {
		return nil
	}

	e.FrameCount++
	if e.AttackCD > 0 {
		e.AttackCD--
	}
	if e.HitFlashCD > 0 {
		e.HitFlashCD--
	}

	cx := e.X + e.Width/2
	cy := e.Y + e.Height/2
	dx := targetX - cx
	dy := targetY - cy
	dist := math.Hypot(dx, dy)

	e.FacingLeft = dx < 0

	var newProj *Projectile

	switch e.Type {
	case EnemyTypeCommon, EnemyTypeRunner, EnemyTypeTank:
		// Perseguição direta ao jogador
		if dist > 2.0 {
			e.X += (dx / dist) * e.Speed
			e.Y += (dy / dist) * e.Speed
		}

	case EnemyTypeShooter:
		// Atirador: lógica de manutenção de distância segura
		safeMinDist := 160.0
		safeMaxDist := 280.0

		if dist < safeMinDist {
			// Muito perto: foge na direção oposta
			e.X -= (dx / dist) * e.Speed * 1.2
			e.Y -= (dy / dist) * e.Speed * 1.2
		} else if dist > safeMaxDist {
			// Muito longe: aproxima-se do jogador
			e.X += (dx / dist) * e.Speed
			e.Y += (dy / dist) * e.Speed
		} else {
			// Na distância ideal: circula ligeiramente
			e.X += (-dy / dist) * e.Speed * 0.7
			e.Y += (dx / dist) * e.Speed * 0.7
		}

		// Disparo de projétil ácido se estiver na distância de combate
		if e.ShootCD > 0 {
			e.ShootCD--
		} else if dist < 360.0 {
			e.ShootCD = 110 // recarga
			// Dispara projétil em direção ao jogador
			speed := 4.2
			vx := (dx / dist) * speed
			vy := (dy / dist) * speed
			newProj = NewProjectile(cx, cy, vx, vy, 5.0, 15, color.RGBA{R: 70, G: 240, B: 90, A: 255}, 140, ProjectileKindEnemyAcid, true)
		}
	}

	return newProj
}

// Bounds retorna a caixa delimitadora do monstro.
func (e *Enemy) Bounds() collision.Rect {
	return collision.NewRect(e.X, e.Y, e.Width, e.Height)
}

// TakeDamage aplica dano e retorna se o monstro morreu.
func (e *Enemy) TakeDamage(amount float64) bool {
	e.HP -= amount
	e.HitFlashCD = 8
	if e.HP <= 0 {
		e.HP = 0
		e.Alive = false
		return true
	}
	return false
}

// CanDamagePlayer verifica se pode causar dano por contato.
func (e *Enemy) CanDamagePlayer() bool {
	return e.AttackCD <= 0
}

// ResetAttackCD reseta o cooldown de dano de contato.
func (e *Enemy) ResetAttackCD() {
	e.AttackCD = 40
}

// Draw renderiza o inimigo com seu sprite de pixel art procedural e efeitos de dano.
func (e *Enemy) Draw(dst *ebiten.Image) {
	if !e.Alive {
		return
	}

	cx := float32(e.X + e.Width/2)
	cy := float32(e.Y + e.Height/2)

	// Sombra circular suave
	shadowRadius := float32(e.Width * 0.4)
	vector.DrawFilledCircle(dst, cx, cy+float32(e.Height*0.42), shadowRadius, color.RGBA{R: 15, G: 15, B: 25, A: 120}, true)

	// Desenho do sprite procedural
	if e.spriteImg != nil {
		op := &ebiten.DrawImageOptions{}

		scale := float64(e.Width) / float64(e.spriteImg.Bounds().Dx())
		sw := float64(e.spriteImg.Bounds().Dx())

		if e.FacingLeft {
			op.GeoM.Scale(-scale, scale)
			op.GeoM.Translate(e.X+sw*scale, e.Y)
		} else {
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(e.X, e.Y)
		}

		// Se estiver no flash de hit, aplica tonalidade branca/amarela de impacto
		if e.HitFlashCD > 0 {
			op.ColorScale.Scale(3.0, 3.0, 3.0, 1.0)
		}

		dst.DrawImage(e.spriteImg, op)
	}

	// Barra de vida miniatura proporcional
	if e.HP < e.MaxHP {
		barW := float32(e.Width)
		barH := float32(4)
		barX := float32(e.X)
		barY := float32(e.Y - 7)

		vector.DrawFilledRect(dst, barX, barY, barW, barH, color.RGBA{R: 25, G: 25, B: 30, A: 210}, false)
		ratio := float32(e.HP / e.MaxHP)
		if ratio < 0 {
			ratio = 0
		}

		// Cor da barra: Tanque tem tom roxo/dourado, normais têm tom vermelho/laranja
		barColor := color.RGBA{R: 240, G: 50, B: 50, A: 255}
		if e.Type == EnemyTypeTank {
			barColor = color.RGBA{R: 180, G: 80, B: 255, A: 255}
		}

		vector.DrawFilledRect(dst, barX, barY, barW*ratio, barH, barColor, false)
		vector.StrokeRect(dst, barX, barY, barW, barH, 1, color.RGBA{R: 190, G: 190, B: 200, A: 200}, false)
	}
}
