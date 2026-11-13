package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
	"joguinho/internal/sprites"
)

// EnemyType define os arquétipos de inimigos comuns e especiais.
type EnemyType int

const (
	EnemyTypeCommon  EnemyType = iota // 1. Comum (Zumbi/Goblin)
	EnemyTypeRunner                   // 2. Corredor (Cão Infernal veloz)
	EnemyTypeTank                     // 3. Tanque (Golem colossal de pedra)
	EnemyTypeShooter                  // 4. Atirador (Bombardeiro que mantém distância)
	EnemyTypeBat                      // 5. Morcego Carmesim (Emboscada da Semente do Caos)
)

// Enemy representa um inimigo na arena.
type Enemy struct {
	ID         uint64
	Type       EnemyType
	X, Y       float64
	Width      float64
	Height     float64
	BaseSpeed  float64
	Speed      float64
	HP         float64
	MaxHP      float64
	Damage     float64
	ScoreValue int
	AttackCD   int
	ShootCD    int
	HitFlashCD int
	StunTimer  int // Atordoamento (ex: Disparo Sísmico)
	SlowTimer  int // Desaceleração (ex: Chuva de Flechas)
	Alive      bool
	FacingLeft bool
	FrameCount int

	// Knockback físico
	KnockVx, KnockVy float64

	// Comportamento exclusivo de Morcego Carmesim
	IsAttached    bool
	AttachOffsetX float64
	AttachOffsetY float64

	spriteImg     *ebiten.Image
	spriteBatUp   *ebiten.Image
	spriteBatDown *ebiten.Image
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
		e.Width = 32
		e.Height = 32
		e.MaxHP = 38 * waveMultiplier
		e.HP = e.MaxHP
		e.BaseSpeed = 1.9 + (waveMultiplier-1)*0.15
		e.Damage = 12
		e.ScoreValue = 100
		e.spriteImg = sprites.GetZombieSprite()

	case EnemyTypeRunner:
		e.Width = 30
		e.Height = 24
		e.MaxHP = 16 * waveMultiplier
		e.HP = e.MaxHP
		e.BaseSpeed = 5.2 + (waveMultiplier-1)*0.2
		e.Damage = 8
		e.ScoreValue = 160
		e.spriteImg = sprites.GetHellhoundSprite()

	case EnemyTypeTank:
		e.Width = 48
		e.Height = 48
		e.MaxHP = 360 * waveMultiplier
		e.HP = e.MaxHP
		e.BaseSpeed = 1.0 + (waveMultiplier-1)*0.08
		e.Damage = 35
		e.ScoreValue = 450
		e.spriteImg = sprites.GetGolemSprite()

	case EnemyTypeShooter:
		e.Width = 32
		e.Height = 32
		e.MaxHP = 45 * waveMultiplier
		e.HP = e.MaxHP
		e.BaseSpeed = 1.6
		e.Damage = 14
		e.ScoreValue = 220
		e.ShootCD = 100
		e.spriteImg = sprites.GetShooterSprite()

	case EnemyTypeBat:
		// Morcego Carmesim da Semente do Caos: rápido, voo curvo e gruda no alvo
		e.Width = 24
		e.Height = 18
		e.MaxHP = 22 * waveMultiplier
		e.HP = e.MaxHP
		e.BaseSpeed = 3.8
		e.Damage = 6
		e.ScoreValue = 130
		e.spriteBatUp = sprites.GetBatSprite(true)
		e.spriteBatDown = sprites.GetBatSprite(false)
		e.AttachOffsetX = (float64(id%5) - 2) * 8
		e.AttachOffsetY = (float64(id%3) - 1) * 8
	}

	e.Speed = e.BaseSpeed
	return e
}

// ApplyKnockback aplica impulso físico de empurrão no inimigo (resiste um pouco no Tanque).
func (e *Enemy) ApplyKnockback(vx, vy float64) {
	resist := 1.0
	if e.Type == EnemyTypeTank {
		resist = 0.25 // Tanque resiste a 75% do recuo
	}
	e.KnockVx += vx * resist
	e.KnockVy += vy * resist
	// Se for um morcego grudado, o knockback solta ele!
	e.IsAttached = false
}

// Update executa a IA, atordoamento, lentidão e física de knockback.
func (e *Enemy) Update(targetX, targetY, speedMult float64) *Projectile {
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

	// Física de Knockback residual
	if math.Abs(e.KnockVx) > 0.1 || math.Abs(e.KnockVy) > 0.1 {
		e.X += e.KnockVx
		e.Y += e.KnockVy
		e.KnockVx *= 0.80
		e.KnockVy *= 0.80
	}

	// Mantém dentro da arena (mesmo com knockback forte)
	if e.X < 12 {
		e.X = 12
	}
	if e.X+e.Width > 788 {
		e.X = 788 - e.Width
	}
	if e.Y < 54 {
		e.Y = 54
	}
	if e.Y+e.Height > 588 {
		e.Y = 588 - e.Height
	}

	// Atordoado: não se move nem atira
	if e.StunTimer > 0 {
		e.StunTimer--
		return nil
	}

	// Calcula velocidade efetiva considerando lentidão e clima (Lua Carmesim)
	effSpeed := e.BaseSpeed * speedMult
	if e.SlowTimer > 0 {
		e.SlowTimer--
		effSpeed *= 0.55
	}
	e.Speed = effSpeed

	cx := e.X + e.Width/2
	cy := e.Y + e.Height/2
	dx := targetX - cx
	dy := targetY - cy
	dist := math.Hypot(dx, dy)

	e.FacingLeft = dx < 0

	var newProj *Projectile

	switch e.Type {
	case EnemyTypeCommon, EnemyTypeRunner, EnemyTypeTank:
		if dist > 2.0 {
			e.X += (dx / dist) * effSpeed
			e.Y += (dy / dist) * effSpeed
		}

	case EnemyTypeBat:
		if e.IsAttached {
			// Grudado no personagem, orbitando em cima dele
			orbitAngle := float64(e.FrameCount) * 0.18
			e.X = targetX - e.Width/2 + math.Cos(orbitAngle)*14 + e.AttachOffsetX*0.3
			e.Y = targetY - e.Height/2 + math.Sin(orbitAngle)*14 + e.AttachOffsetY*0.3
		} else if dist > 18.0 {
			// Voo em padrão helicoidal / curvado
			dirX := dx / dist
			dirY := dy / dist
			perpX := -dirY
			perpY := dirX
			wave := math.Sin(float64(e.FrameCount)*0.18+float64(e.ID)) * 2.2
			e.X += dirX*effSpeed + perpX*wave
			e.Y += dirY*effSpeed + perpY*wave
		} else {
			// Chegou perto: gruda no jogador!
			e.IsAttached = true
		}

	case EnemyTypeShooter:
		safeMinDist := 160.0
		safeMaxDist := 280.0

		if dist < safeMinDist {
			e.X -= (dx / dist) * effSpeed * 1.2
			e.Y -= (dy / dist) * effSpeed * 1.2
		} else if dist > safeMaxDist {
			e.X += (dx / dist) * effSpeed
			e.Y += (dy / dist) * effSpeed
		} else {
			e.X += (-dy / dist) * effSpeed * 0.7
			e.Y += (dx / dist) * effSpeed * 0.7
		}

		if e.ShootCD > 0 {
			e.ShootCD--
		} else if dist < 360.0 {
			e.ShootCD = 110
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
	return e.AttackCD <= 0 && e.StunTimer <= 0
}

// ResetAttackCD reseta o cooldown de dano de contato.
func (e *Enemy) ResetAttackCD() {
	if e.Type == EnemyTypeBat {
		e.AttackCD = 32 // Morcegos mordem rápido quando grudados!
		return
	}
	e.AttackCD = 40
}

// Draw renderiza o inimigo com seu sprite de pixel art, barras e indicadores de Stun/Slow.
func (e *Enemy) Draw(dst *ebiten.Image) {
	if !e.Alive {
		return
	}

	cx := float32(e.X + e.Width/2)
	cy := float32(e.Y + e.Height/2)

	// Calcula parâmetros de animação específicos por arquétipo
	var bobY, tilt, stretchX, stretchY float64
	stretchX = 1.0
	stretchY = 1.0

	if e.StunTimer <= 0 {
		t := float64(e.FrameCount) + float64(e.ID)*1.7
		switch e.Type {
		case EnemyTypeCommon:
			// Andar cambaleante de zumbi
			phase := t * 0.16
			bobY = -math.Abs(math.Sin(phase)) * 3.5
			tilt = math.Sin(phase) * 0.16
			squash := math.Sin(phase*2) * 0.06
			stretchX = 1.0 - squash*0.6
			stretchY = 1.0 + squash

		case EnemyTypeRunner:
			// Galope frenético de Cão Infernal (estica no pulo!)
			phase := t * 0.36
			bobY = -math.Abs(math.Sin(phase)) * 6.0
			tilt = math.Sin(phase) * 0.14
			stretch := math.Cos(phase) * 0.14
			stretchX = 1.0 + stretch
			stretchY = 1.0 - stretch*0.8

		case EnemyTypeTank:
			// Passos colossais lentos e pesados do Golem
			phase := t * 0.09
			bobY = -math.Abs(math.Sin(phase)) * 4.0
			tilt = math.Sin(phase) * 0.08
			squash := math.Sin(phase*2) * 0.07
			stretchX = 1.0 - squash*0.8
			stretchY = 1.0 + squash

		case EnemyTypeShooter:
			// Levitação fantasmagórica do Atirador
			bobY = -4.0 + math.Sin(t*0.10)*4.0
			tilt = math.Sin(t*0.07) * 0.08
			// Pulsa antes de atirar
			if e.ShootCD < 25 {
				pulse := math.Sin(t*0.5) * 0.10
				stretchX = 1.0 + pulse
				stretchY = 1.0 + pulse
			}

		case EnemyTypeBat:
			// Voo agitado do Morcego Carmesim
			bobY = -3.0 + math.Sin(t*0.38)*4.5
			tilt = math.Cos(t*0.22) * 0.20
		}
	}

	// Sombra circular suave (encolhe proporcionalmente ao salto/levitação)
	shadowRadius := float32(e.Width*0.4) - float32(-bobY)*0.4
	if shadowRadius < 4 {
		shadowRadius = 4
	}
	vector.DrawFilledCircle(dst, cx, cy+float32(e.Height*0.42), shadowRadius, color.RGBA{R: 15, G: 15, B: 25, A: 120}, true)

	// Escolhe sprite (para o Morcego, alterna asas animadas rapidamente)
	imgToDraw := e.spriteImg
	if e.Type == EnemyTypeBat {
		if (e.FrameCount/4)%2 == 0 {
			imgToDraw = e.spriteBatUp
		} else {
			imgToDraw = e.spriteBatDown
		}
	}

	if imgToDraw != nil {
		op := &ebiten.DrawImageOptions{}

		sw := float64(imgToDraw.Bounds().Dx())
		sh := float64(imgToDraw.Bounds().Dy())
		scaleX := (float64(e.Width) / sw) * stretchX
		scaleY := (float64(e.Height) / sh) * stretchY

		// 1. Centraliza origem
		op.GeoM.Translate(-sw/2, -sh/2)

		// 2. Escala + Flip
		if e.FacingLeft {
			op.GeoM.Scale(-scaleX, scaleY)
		} else {
			op.GeoM.Scale(scaleX, scaleY)
		}

		// 3. Rotação de balanço/caminhada
		op.GeoM.Rotate(tilt)

		// 4. Translada para a arena com deslocamento vertical da animação
		op.GeoM.Translate(float64(cx), float64(cy)+bobY)

		if e.HitFlashCD >= 6 {
			op.ColorScale.Scale(255, 255, 255, 1.0) // Flash branco puro de dano
		} else if e.HitFlashCD > 0 {
			op.ColorScale.Scale(2.5, 2.5, 2.5, 1.0)
		} else if e.SlowTimer > 0 {
			op.ColorScale.Scale(0.7, 0.9, 1.5, 1.0)
		}

		dst.DrawImage(imgToDraw, op)
	}

	// Indicador visual de Atordoamento (estrelas girando acima da cabeça)
	if e.StunTimer > 0 {
		for k := 0; k < 3; k++ {
			ang := float64(e.FrameCount)*0.2 + float64(k)*(math.Pi*2/3)
			sx := cx + float32(math.Cos(ang)*10)
			sy := float32(e.Y-6) + float32(math.Sin(ang)*3)
			vector.DrawFilledCircle(dst, sx, sy, 2.2, color.RGBA{R: 255, G: 230, B: 60, A: 255}, true)
		}
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

		barColor := color.RGBA{R: 240, G: 50, B: 50, A: 255}
		if e.Type == EnemyTypeTank {
			barColor = color.RGBA{R: 180, G: 80, B: 255, A: 255}
		}

		vector.DrawFilledRect(dst, barX, barY, barW*ratio, barH, barColor, false)
		vector.StrokeRect(dst, barX, barY, barW, barH, 1, color.RGBA{R: 190, G: 190, B: 200, A: 200}, false)
	}
}
