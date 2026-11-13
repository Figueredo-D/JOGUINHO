package entities

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
)

// SlashSpark representa uma faísca ou centelha de impacto/atrito da lâmina cortando o ar.
type SlashSpark struct {
	X, Y     float64
	Vx, Vy   float64
	Lifetime int
	MaxLife  int
	Color    color.RGBA
	Size     float32
}

// Slash representa o efeito visual e a área de colisão do golpe corpo a corpo do Guerreiro.
type Slash struct {
	OriginX, OriginY float64
	Angle            float64 // radianos para onde o golpe é desferido
	Range            float64
	ArcSpan          float64 // abertura em radianos
	Damage           float64
	Lifetime         int
	MaxLifetime      int
	HitEnemies       map[uint64]bool // evita acertar o mesmo inimigo múltiplos frames do mesmo golpe
	Active           bool
	Sparks           []SlashSpark
}

// NewSlash cria uma nova rajada de ataque corpo a corpo com arco dinâmico.
func NewSlash(x, y, targetX, targetY, attackRange, damage float64) *Slash {
	dx := targetX - x
	dy := targetY - y
	angle := math.Atan2(dy, dx)

	return &Slash{
		OriginX:     x,
		OriginY:     y,
		Angle:       angle,
		Range:       attackRange,
		ArcSpan:     math.Pi * 0.85, // Arco amplo de 153 graus
		Damage:      damage,
		Lifetime:    13,             // ~0.22s de animação fluida
		MaxLifetime: 13,
		HitEnemies:  make(map[uint64]bool),
		Active:      true,
		Sparks:      make([]SlashSpark, 0, 16),
	}
}

// Update avança o frame do ataque, calcula estocada e dispersa faíscas.
func (s *Slash) Update() {
	if !s.Active {
		return
	}
	s.Lifetime--
	if s.Lifetime <= 0 {
		s.Active = false
		return
	}

	progress := 1.0 - float64(s.Lifetime)/float64(s.MaxLifetime)

	// Gera faíscas ao longo da ponta da lâmina durante a primeira metade do golpe
	if progress < 0.75 {
		currentAngle := (s.Angle - s.ArcSpan/2) + progress*s.ArcSpan
		edgeRadius := s.Range * 0.95
		bx := s.OriginX + math.Cos(currentAngle)*edgeRadius
		by := s.OriginY + math.Sin(currentAngle)*edgeRadius

		for k := 0; k < 2; k++ {
			sparkAng := currentAngle + (rand.Float64()-0.5)*0.8
			spd := 1.5 + rand.Float64()*2.5
			clr := color.RGBA{R: 255, G: 245, B: 200, A: 255}
			if rand.Float64() < 0.4 {
				clr = color.RGBA{R: 120, G: 220, B: 255, A: 255} // Faísca celeste
			}
			s.Sparks = append(s.Sparks, SlashSpark{
				X:        bx + (rand.Float64()-0.5)*4,
				Y:        by + (rand.Float64()-0.5)*4,
				Vx:       math.Cos(sparkAng) * spd,
				Vy:       math.Sin(sparkAng) * spd,
				Lifetime: 7 + rand.Intn(6),
				MaxLife:  12,
				Color:    clr,
				Size:     float32(1.2 + rand.Float64()*1.2),
			})
		}
	}

	// Atualiza faíscas existentes
	for i := len(s.Sparks) - 1; i >= 0; i-- {
		sp := &s.Sparks[i]
		sp.X += sp.Vx
		sp.Y += sp.Vy
		sp.Vx *= 0.92
		sp.Vy *= 0.92
		sp.Lifetime--
		if sp.Lifetime <= 0 {
			s.Sparks = append(s.Sparks[:i], s.Sparks[i+1:]...)
		}
	}
}

// CollidesWithEnemy verifica se uma entidade está no alcance e ângulo do arco.
func (s *Slash) CollidesWithEnemy(enemyID uint64, enemyBounds collision.Rect) bool {
	if !s.Active || s.HitEnemies[enemyID] {
		return false
	}

	cx := enemyBounds.X + enemyBounds.W/2
	cy := enemyBounds.Y + enemyBounds.H/2

	dx := cx - s.OriginX
	dy := cy - s.OriginY
	dist := math.Hypot(dx, dy)

	if dist > s.Range+enemyBounds.W/2 {
		return false
	}

	angleToEnemy := math.Atan2(dy, dx)
	diff := math.Abs(normalizeAngle(angleToEnemy - s.Angle))

	if diff <= s.ArcSpan/2 {
		s.HitEnemies[enemyID] = true
		return true
	}

	return false
}

// Draw renderiza o arco volumétrico de lâmina semi-circular com gradiente prateado, rastro estocada e faíscas.
func (s *Slash) Draw(dst *ebiten.Image) {
	if !s.Active {
		return
	}

	progress := 1.0 - float64(s.Lifetime)/float64(s.MaxLifetime)
	fade := float32(1.0 - progress)

	// Leve estocada/avanço da base no primeiro terço da animação (impacto físico)
	thrust := math.Sin(math.Min(progress*2.8, 1.0) * math.Pi * 0.5) * 10.0
	baseX := s.OriginX + math.Cos(s.Angle)*thrust
	baseY := s.OriginY + math.Sin(s.Angle)*thrust

	// Ângulo de varredura atual da lâmina
	startAngle := s.Angle - s.ArcSpan/2
	leadingAngle := startAngle + progress*s.ArcSpan
	trailSpan := s.ArcSpan * (0.35 + 0.65*(1.0-progress))
	trailStart := leadingAngle - trailSpan

	steps := 22
	stepAngle := trailSpan / float64(steps)

	for i := 0; i < steps; i++ {
		t1 := float64(i) / float64(steps)
		t2 := float64(i+1) / float64(steps)

		a1 := trailStart + float64(i)*stepAngle
		a2 := trailStart + float64(i+1)*stepAngle

		// Espessura graduada: fino na cauda do rastro, encorpado no ápice (t=0.7) e afilado na ponta (t=1.0)
		tMid := (t1 + t2) * 0.5
		thickFactor := math.Sin(tMid * math.Pi * 0.88)
		outerThick := float32(2.0 + 8.5*thickFactor)
		innerThick := float32(1.5 + 4.0*thickFactor)

		r1 := float32(s.Range * (0.80 + 0.20*math.Sin(t1*math.Pi*0.9)))
		r2 := float32(s.Range * (0.80 + 0.20*math.Sin(t2*math.Pi*0.9)))

		x1 := float32(baseX + float64(r1)*math.Cos(a1))
		y1 := float32(baseY + float64(r1)*math.Sin(a1))
		x2 := float32(baseX + float64(r2)*math.Cos(a2))
		y2 := float32(baseY + float64(r2)*math.Sin(a2))

		// Camada 1: Rastro etéreo azul-celeste externo (aura de velocidade)
		alphaAura := uint8(float32(110) * fade * float32(tMid))
		vector.StrokeLine(dst, x1, y1, x2, y2, outerThick+3.0, color.RGBA{R: 70, G: 165, B: 255, A: alphaAura}, true)

		// Camada 2: Corpo da lâmina de aço polido (Silver Blade)
		alphaSteel := uint8(float32(230) * fade * float32(0.3+0.7*tMid))
		vector.StrokeLine(dst, x1, y1, x2, y2, outerThick, color.RGBA{R: 195, G: 215, B: 245, A: alphaSteel}, true)

		// Camada 3: Fio cortante branco reluzente no centro
		alphaCore := uint8(float32(255) * fade * float32(tMid))
		vector.StrokeLine(dst, x1, y1, x2, y2, innerThick, color.RGBA{R: 255, G: 255, B: 255, A: alphaCore}, true)
	}

	// Brilho intenso de corte na ponta da lâmina (Razor Glint)
	if progress < 0.90 {
		tipR := s.Range * 0.98
		tipX := float32(baseX + tipR*math.Cos(leadingAngle))
		tipY := float32(baseY + tipR*math.Sin(leadingAngle))
		glintLen := float32(6.5 * (1.0 - progress))
		vector.StrokeLine(dst, tipX-glintLen, tipY, tipX+glintLen, tipY, 2.0, color.RGBA{255, 255, 255, 255}, true)
		vector.StrokeLine(dst, tipX, tipY-glintLen, tipX, tipY+glintLen, 2.0, color.RGBA{255, 255, 255, 255}, true)
		vector.DrawFilledCircle(dst, tipX, tipY, 3.0, color.RGBA{140, 230, 255, 240}, true)
	}

	// Renderiza as micro-faíscas de atrito
	for _, sp := range s.Sparks {
		sparkFade := float32(sp.Lifetime) / float32(sp.MaxLife)
		c := sp.Color
		c.A = uint8(float32(c.A) * sparkFade)
		vector.DrawFilledCircle(dst, float32(sp.X), float32(sp.Y), sp.Size*sparkFade, c, true)
	}
}

func normalizeAngle(a float64) float64 {
	for a > math.Pi {
		a -= 2 * math.Pi
	}
	for a < -math.Pi {
		a += 2 * math.Pi
	}
	return a
}
