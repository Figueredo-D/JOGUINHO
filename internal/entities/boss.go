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

// BossState define os padrões multifásicos do Chefe.
type BossState int

const (
	BossStateChase BossState = iota
	BossStateChargePrep
	BossStateCharging
	BossStateStompPrep
	BossStateSummonOrbs
	BossStateLaserPrep
	BossStateLaserFire
)

// LaserBeam representa um feixe de laser etéreo disparado pelo Soberano de Dorfieegu.
type LaserBeam struct {
	X1, Y1 float64
	X2, Y2 float64
	Angle  float64
	Width  float64
}

// BossActionOutput retorna projéteis, ondas de choque, lasers ou partículas geradas pelo Chefe num frame.
type BossActionOutput struct {
	Projectiles []*Projectile
	Shockwaves  []*ShockwaveEffect
	Particles   []*Particle
	StompDamage float64
	StompRadius float64
	LaserDamage float64 // Dano ativo dos feixes neste frame
}

// Boss representa o Chefe de Onda (Ondas 5, 10 e Chefe Final na Onda 15).
type Boss struct {
	Name       string
	IsFinal    bool
	X, Y       float64
	Width      float64
	Height     float64
	HP         float64
	MaxHP      float64
	Speed      float64
	Damage     float64
	State      BossState
	StateTimer int
	PatternIdx int
	HitFlashCD int
	AttackCD   int
	Alive      bool
	FacingLeft bool
	FrameCount int
	WaveLevel  int

	// Vetores de ataque e fases
	ChargeVx, ChargeVy float64
	TargetX, TargetY   float64
	ChargesLeft        int
	Lasers             []LaserBeam

	spriteImg *ebiten.Image
}

// NewBoss instancia o Chefe rebalanceado para as ondas 5, 10, 15 ou ondas superiores (2.5x vida e 2x dano).
func NewBoss(wave int) *Boss {
	var name string
	var maxHP, spd, dmg, size float64
	var isFinal bool

	switch {
	case wave == 5:
		name = "O CARRASCO DAS CATACUMBAS"
		maxHP = 3750 // 1500 * 2.5
		spd = 1.85
		dmg = 56     // 28 * 2
		size = 66
	case wave == 10:
		name = "O GOLEM DE ALMAS"
		maxHP = 6750 // 2700 * 2.5
		spd = 1.60
		dmg = 76     // 38 * 2
		size = 76
	case wave == 15:
		name = "O SOBERANO DE DORFIEEGU"
		maxHP = 10500 // 4200 * 2.5
		spd = 2.00
		dmg = 84     // 42 * 2
		size = 86
		isFinal = true
	default: // Ondas 20 e subsequentes
		name = fmt.Sprintf("ARQUI-DEMÔNIO ANCESTRAL (ONDA %d)", wave)
		extra := float64(wave - 15)
		maxHP = 10500 + extra*1800
		spd = 2.15
		dmg = 84 + extra*6
		size = 90
		isFinal = true
	}

	return &Boss{
		Name:       name,
		IsFinal:    isFinal,
		X:          400 - size/2,
		Y:          110,
		Width:      size,
		Height:     size,
		HP:         maxHP,
		MaxHP:      maxHP,
		Speed:      spd,
		Damage:     dmg,
		State:      BossStateChase,
		StateTimer: 35, // Intervalo inicial reduzido para ação imediata
		Alive:      true,
		WaveLevel:  wave,
		spriteImg:  sprites.GetBossSprite(),
	}
}

// Phase retorna a fase atual do Chefe (1, 2 ou 3) baseada na vida restante.
func (b *Boss) Phase() int {
	ratio := b.HP / b.MaxHP
	if ratio <= 0.33 {
		return 3
	}
	if ratio <= 0.66 {
		return 2
	}
	return 1
}

// Center retorna o centro do Boss.
func (b *Boss) Center() (float64, float64) {
	return b.X + b.Width/2, b.Y + b.Height/2
}

// Bounds retorna a caixa delimitadora AABB.
func (b *Boss) Bounds() collision.Rect {
	return collision.NewRect(b.X, b.Y, b.Width, b.Height)
}

// TakeDamage aplica dano e ativa flash branco de impacto.
func (b *Boss) TakeDamage(amount float64) bool {
	b.HP -= amount
	b.HitFlashCD = 6
	if b.HP <= 0 {
		b.HP = 0
		b.Alive = false
		return true
	}
	return false
}

// CanDamagePlayer verifica cooldown de contato.
func (b *Boss) CanDamagePlayer() bool {
	return b.AttackCD <= 0
}

// ResetAttackCD reseta cooldown de contato.
func (b *Boss) ResetAttackCD() {
	b.AttackCD = 40
}

// Update executa a máquina de estados multifásica do Chefe.
func (b *Boss) Update(playerX, playerY float64) BossActionOutput {
	var out BossActionOutput
	if !b.Alive {
		return out
	}

	b.FrameCount++
	if b.HitFlashCD > 0 {
		b.HitFlashCD--
	}
	if b.AttackCD > 0 {
		b.AttackCD--
	}

	cx, cy := b.Center()
	dx := playerX - cx
	dy := playerY - cy
	dist := math.Hypot(dx, dy)
	b.FacingLeft = dx < 0

	phase := b.Phase()
	hpRatio := b.HP / b.MaxHP
	if hpRatio < 0 {
		hpRatio = 0
	}
	// Aceleração dinâmica conforme a vida do Chefe decai (até +75% de velocidade de perseguição no ápice do desespero)
	enrageSpeedBonus := (1.0 - hpRatio) * 0.75
	effSpeed := b.Speed * (1.0 + float64(phase-1)*0.18 + enrageSpeedBonus)
	b.StateTimer--

	switch b.State {
	case BossStateChase:
		if dist > 4.0 {
			b.X += (dx / dist) * effSpeed
			b.Y += (dy / dist) * effSpeed
		}

		// Padrões agressivos com menos de 50% de HP: rajadas em leque e espirais contínuas!
		if hpRatio < 0.50 {
			// Disparo em leque (5 projéteis) direcionados ao jogador a cada 46 frames
			if b.FrameCount%46 == 0 {
				baseAng := math.Atan2(dy, dx)
				fanSpd := 3.6 + (1.0-hpRatio)*1.4
				fanDmg := b.Damage * 0.45
				for k := -2; k <= 2; k++ {
					ang := baseAng + float64(k)*0.22
					vx := math.Cos(ang) * fanSpd
					vy := math.Sin(ang) * fanSpd
					fanOrb := NewProjectile(cx, cy, vx, vy, 6.0, fanDmg, color.RGBA{255, 65, 45, 255}, 160, ProjectileKindBossOrb, true)
					out.Projectiles = append(out.Projectiles, fanOrb)
				}
				out.Particles = append(out.Particles, SpawnBurst(cx, cy, 10, 2.8, 2.5, color.RGBA{255, 120, 30, 255}, 18)...)
			}

			// Na Fase 3 (HP < 33%), dispara espiral rotativa quádrupla de almas sombrias
			if hpRatio < 0.33 && b.FrameCount%28 == 0 {
				spiralBase := float64(b.FrameCount) * 0.14
				spiralDmg := b.Damage * 0.40
				for k := 0; k < 4; k++ {
					ang := spiralBase + float64(k)*(math.Pi/2)
					vx := math.Cos(ang) * 3.4
					vy := math.Sin(ang) * 3.4
					spiralOrb := NewProjectile(cx, cy, vx, vy, 5.5, spiralDmg, color.RGBA{180, 50, 240, 255}, 150, ProjectileKindBossOrb, true)
					out.Projectiles = append(out.Projectiles, spiralOrb)
				}
			}
		}

		if b.StateTimer <= 0 {
			// O Soberano de Dorfieegu (Onda 15+) usa todos os 4 padrões (incluindo Feixes de Laser!)
			mod := 3
			if b.IsFinal {
				mod = 4
			}
			switch b.PatternIdx % mod {
			case 0:
				b.State = BossStateChargePrep
				b.StateTimer = 26 // Telegrafia mais ágil
				b.ChargesLeft = 1
				if b.IsFinal && phase >= 2 {
					b.ChargesLeft = 2 // Investidas encadeadas na Fase 2/3!
				}
				b.TargetX = playerX
				b.TargetY = playerY
			case 1:
				b.State = BossStateSummonOrbs
				b.StateTimer = 20
			case 2:
				b.State = BossStateStompPrep
				b.StateTimer = 34
			case 3:
				b.State = BossStateLaserPrep
				b.StateTimer = 32
				b.prepareLasers(cx, cy, playerX, playerY, phase)
			}
			b.PatternIdx++
		}

	case BossStateChargePrep:
		if b.StateTimer == 10 {
			b.TargetX = playerX
			b.TargetY = playerY
		}
		if b.StateTimer <= 0 {
			cdx := b.TargetX - cx
			cdy := b.TargetY - cy
			cdist := math.Hypot(cdx, cdy)
			if cdist < 1 {
				cdist = 1
			}
			chargeSpeed := 9.6 + float64(phase-1)*1.4
			b.ChargeVx = (cdx / cdist) * chargeSpeed
			b.ChargeVy = (cdy / cdist) * chargeSpeed
			b.State = BossStateCharging
			b.StateTimer = 26
		}

	case BossStateCharging:
		b.X += b.ChargeVx
		b.Y += b.ChargeVy

		if b.FrameCount%2 == 0 {
			out.Particles = append(out.Particles, SpawnBurst(cx, cy+20, 3, 2.2, 3.5, color.RGBA{255, 110, 40, 220}, 22)...)
		}

		if b.StateTimer <= 0 {
			b.ChargesLeft--
			if b.ChargesLeft > 0 {
				b.State = BossStateChargePrep
				b.StateTimer = 16
				b.TargetX = playerX
				b.TargetY = playerY
			} else {
				b.State = BossStateChase
				b.StateTimer = 32 // Intervalo de descanso reduzido
			}
		}

	case BossStateSummonOrbs:
		if b.StateTimer == 10 {
			count := 6
			if b.WaveLevel == 10 {
				count = 8
			} else if b.IsFinal {
				count = 8 + phase*2 // 10, 12 ou 14 orbes no Soberano!
			}
			for i := 0; i < count; i++ {
				ang := float64(i)*(math.Pi*2/float64(count)) + float64(b.FrameCount)*0.05
				spd := 2.6 + float64(phase-1)*0.35
				vx := math.Cos(ang) * spd
				vy := math.Sin(ang) * spd
				orb := NewProjectile(cx, cy, vx, vy, 7.5, 36, color.RGBA{230, 40, 90, 255}, 240, ProjectileKindBossOrb, true)
				orb.IsHoming = true
				orb.HomingTurn = 0.034
				out.Projectiles = append(out.Projectiles, orb)
			}
			out.Particles = append(out.Particles, SpawnBurst(cx, cy, 14, 3.5, 3.0, color.RGBA{220, 60, 120, 255}, 25)...)
		}
		if b.StateTimer <= 0 {
			b.State = BossStateChase
			b.StateTimer = 28 // Intervalo de descanso reduzido
		}

	case BossStateStompPrep:
		if b.StateTimer <= 0 {
			stompRad := 135.0
			if b.IsFinal {
				stompRad = 165.0
			}
			out.Shockwaves = append(out.Shockwaves, NewShockwave(cx, cy+10, stompRad, color.RGBA{255, 70, 40, 240}, 28, false))
			out.Particles = append(out.Particles, SpawnBurst(cx, cy+10, 24, 4.8, 4.0, color.RGBA{255, 160, 60, 255}, 30)...)
			out.StompDamage = 72 // 2x dano do pisão
			out.StompRadius = stompRad

			b.State = BossStateChase
			b.StateTimer = 32 // Intervalo de descanso reduzido
		}

	case BossStateLaserPrep:
		// Mantém origem dos feixes no Boss enquanto telegrafa
		for i := range b.Lasers {
			b.Lasers[i].X1 = cx
			b.Lasers[i].Y1 = cy
			b.Lasers[i].X2 = cx + math.Cos(b.Lasers[i].Angle)*900
			b.Lasers[i].Y2 = cy + math.Sin(b.Lasers[i].Angle)*900
		}
		if b.StateTimer <= 0 {
			b.State = BossStateLaserFire
			b.StateTimer = 24 // Disparo ativo por 24 frames
		}

	case BossStateLaserFire:
		// Causa dano contínuo se jogador cruzar qualquer laser ativo
		for _, l := range b.Lasers {
			if pointToSegmentDist(playerX, playerY, l.X1, l.Y1, l.X2, l.Y2) <= l.Width*0.6+12 {
				out.LaserDamage = 56 // 2x dano de contato do feixe
				break
			}
		}
		if b.StateTimer <= 0 {
			b.Lasers = nil
			b.State = BossStateChase
			b.StateTimer = 28 // Intervalo de descanso reduzido
		}
	}

	// Limites da arena
	if b.X < 14 {
		b.X = 14
	}
	if b.X+b.Width > 786 {
		b.X = 786 - b.Width
	}
	if b.Y < 56 {
		b.Y = 56
	}
	if b.Y+b.Height > 586 {
		b.Y = 586 - b.Height
	}

	return out
}

func (b *Boss) prepareLasers(cx, cy, px, py float64, phase int) {
	baseAngle := math.Atan2(py-cy, px-cx)
	b.Lasers = nil
	angles := []float64{baseAngle}
	if phase == 2 {
		angles = []float64{baseAngle - 0.30, baseAngle + 0.30}
	} else if phase >= 3 {
		angles = []float64{baseAngle - 0.42, baseAngle, baseAngle + 0.42}
	}

	for _, ang := range angles {
		b.Lasers = append(b.Lasers, LaserBeam{
			X1:    cx,
			Y1:    cy,
			X2:    cx + math.Cos(ang)*900,
			Y2:    cy + math.Sin(ang)*900,
			Angle: ang,
			Width: 16,
		})
	}
}

func pointToSegmentDist(px, py, x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	l2 := dx*dx + dy*dy
	if l2 == 0 {
		return math.Hypot(px-x1, py-y1)
	}
	t := ((px-x1)*dx + (py-y1)*dy) / l2
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	projX := x1 + t*dx
	projY := y1 + t*dy
	return math.Hypot(px-projX, py-projY)
}

// Draw renderiza o Chefe, feixes de laser etéreos, telegrafias e flash branco puro.
func (b *Boss) Draw(dst *ebiten.Image) {
	if !b.Alive {
		return
	}

	cx, cy := b.Center()

	// 1. Telegrafias e Feixes de Laser
	switch b.State {
	case BossStateChargePrep:
		vector.StrokeLine(dst, float32(cx), float32(cy), float32(b.TargetX), float32(b.TargetY), 2.5, color.RGBA{255, 50, 50, 160}, true)
		vector.StrokeCircle(dst, float32(b.TargetX), float32(b.TargetY), 10, 2, color.RGBA{255, 50, 50, 200}, true)

	case BossStateStompPrep:
		progress := 1.0 - float64(b.StateTimer)/46.0
		stompRad := float32(130.0)
		if b.IsFinal {
			stompRad = 155.0
		}
		vector.DrawFilledCircle(dst, float32(cx), float32(cy+10), stompRad*float32(progress), color.RGBA{255, 50, 30, 55}, true)
		vector.StrokeCircle(dst, float32(cx), float32(cy+10), stompRad, 2.5, color.RGBA{255, 80, 40, 220}, true)

	case BossStateLaserPrep:
		for _, l := range b.Lasers {
			vector.StrokeLine(dst, float32(l.X1), float32(l.Y1), float32(l.X2), float32(l.Y2), 2.0, color.RGBA{200, 90, 255, 160}, true)
		}

	case BossStateLaserFire:
		for _, l := range b.Lasers {
			vector.StrokeLine(dst, float32(l.X1), float32(l.Y1), float32(l.X2), float32(l.Y2), float32(l.Width), color.RGBA{180, 60, 255, 210}, true)
			vector.StrokeLine(dst, float32(l.X1), float32(l.Y1), float32(l.X2), float32(l.Y2), float32(l.Width*0.45), color.RGBA{255, 240, 255, 255}, true)
		}
	}

	// 2. Animação conforme o estado
	var bobY, shakeX, tilt, stretchX, stretchY float64
	stretchX = 1.0
	stretchY = 1.0
	t := float64(b.FrameCount)

	switch b.State {
	case BossStateChase:
		bobY = -math.Abs(math.Sin(t*0.11)) * 5.0
		tilt = math.Sin(t*0.11) * 0.08
	case BossStateChargePrep, BossStateLaserPrep:
		shakeX = math.Sin(t*1.4) * 3.0
		stretchX = 1.08
		stretchY = 0.92
	case BossStateCharging:
		if b.FacingLeft {
			tilt = -0.22
		} else {
			tilt = 0.22
		}
		stretchX = 1.12
		stretchY = 0.92
	case BossStateStompPrep:
		progress := 1.0 - float64(b.StateTimer)/46.0
		bobY = -progress * 22.0
	case BossStateSummonOrbs, BossStateLaserFire:
		pulseScale := math.Sin(t*0.3) * 0.08
		stretchX = 1.0 + pulseScale
		stretchY = 1.0 + pulseScale
	}

	shadowRad := float32(b.Width*0.46) - float32(-bobY)*0.4
	vector.DrawFilledCircle(dst, float32(cx), float32(cy+b.Height*0.42), shadowRad, color.RGBA{10, 8, 16, 175}, true)

	// Aura ancestral diferenciada para o Soberano de Dorfieegu
	auraClr := color.RGBA{180, 40, 70, 120}
	if b.IsFinal {
		auraClr = color.RGBA{215, 165, 50, 170}
	}
	pulse := float32(math.Sin(t*0.1) * 4)
	vector.StrokeCircle(dst, float32(cx), float32(cy+bobY), float32(b.Width*0.56)+pulse, 2.5, auraClr, true)

	// 3. Sprite com Flash Branco Puro nos 2 primeiros frames do impacto
	if b.spriteImg != nil {
		op := &ebiten.DrawImageOptions{}
		sw := float64(b.spriteImg.Bounds().Dx())
		sh := float64(b.spriteImg.Bounds().Dy())
		scaleX := (b.Width / sw) * stretchX
		scaleY := (b.Height / sh) * stretchY

		op.GeoM.Translate(-sw/2, -sh/2)
		if b.FacingLeft {
			op.GeoM.Scale(-scaleX, scaleY)
		} else {
			op.GeoM.Scale(scaleX, scaleY)
		}
		op.GeoM.Rotate(tilt)
		op.GeoM.Translate(cx+shakeX, cy+bobY)

		if b.HitFlashCD >= 4 {
			op.ColorScale.Scale(255, 255, 255, 1.0) // Flash branco puro!
		} else if b.HitFlashCD > 0 {
			op.ColorScale.Scale(2.5, 2.5, 2.5, 1.0)
		} else if b.IsFinal && b.Phase() == 3 {
			op.ColorScale.Scale(1.35, 0.85, 1.25, 1.0) // Brilho sombrio na Fase 3 do Soberano
		}

		dst.DrawImage(b.spriteImg, op)
	}
}
