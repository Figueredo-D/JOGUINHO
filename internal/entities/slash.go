package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
)

// Slash representa o efeito visual e a área de colisão do golpe corpo a corpo do Guerreiro.
type Slash struct {
	OriginX, OriginY float64
	Angle            float64 // radianos para onde o golpe é desferido
	Range            float64
	ArcSpan          float64 // abertura em radianos (ex: PI * 0.75)
	Damage           float64
	Lifetime         int
	MaxLifetime      int
	HitEnemies       map[uint64]bool // evita acertar o mesmo inimigo múltiplos frames do mesmo golpe
	Active           bool
}

// NewSlash cria uma nova rajada de ataque corpo a corpo.
func NewSlash(x, y, targetX, targetY, attackRange, damage float64) *Slash {
	dx := targetX - x
	dy := targetY - y
	angle := math.Atan2(dy, dx)

	return &Slash{
		OriginX:     x,
		OriginY:     y,
		Angle:       angle,
		Range:       attackRange,
		ArcSpan:     math.Pi * 0.75,
		Damage:      damage,
		Lifetime:    12, // dura 12 frames (~0.2s)
		MaxLifetime: 12,
		HitEnemies:  make(map[uint64]bool),
		Active:      true,
	}
}

// Update avança o frame do ataque.
func (s *Slash) Update() {
	if !s.Active {
		return
	}
	s.Lifetime--
	if s.Lifetime <= 0 {
		s.Active = false
	}
}

// CollidesWithEnemy verifica se uma entidade com reticência/retângulo está no alcance e ângulo do arco.
func (s *Slash) CollidesWithEnemy(enemyID uint64, enemyBounds collision.Rect) bool {
	if !s.Active || s.HitEnemies[enemyID] {
		return false
	}

	// Ponto central do inimigo
	cx := enemyBounds.X + enemyBounds.W/2
	cy := enemyBounds.Y + enemyBounds.H/2

	dx := cx - s.OriginX
	dy := cy - s.OriginY
	dist := math.Hypot(dx, dy)

	// Está dentro do alcance do raio
	if dist > s.Range+enemyBounds.W/2 {
		return false
	}

	// Ângulo do inimigo em relação ao ponto de ataque
	angleToEnemy := math.Atan2(dy, dx)
	diff := math.Abs(normalizeAngle(angleToEnemy - s.Angle))

	// Está dentro da abertura angular do corte
	if diff <= s.ArcSpan/2 {
		s.HitEnemies[enemyID] = true
		return true
	}

	return false
}

// Draw renderiza o arco do golpe na tela com efeito dinâmico.
func (s *Slash) Draw(dst *ebiten.Image) {
	if !s.Active {
		return
	}

	progress := 1.0 - float64(s.Lifetime)/float64(s.MaxLifetime)
	currentRadius := float32(s.Range * (0.6 + 0.4*progress))

	// Desenha linhas do arco simulando a lâmina cortando o ar
	steps := 12
	startAngle := s.Angle - s.ArcSpan/2
	stepAngle := s.ArcSpan / float64(steps)

	for i := 0; i < steps; i++ {
		a1 := startAngle + float64(i)*stepAngle
		a2 := startAngle + float64(i+1)*stepAngle

		x1 := float32(s.OriginX + float64(currentRadius)*math.Cos(a1))
		y1 := float32(s.OriginY + float64(currentRadius)*math.Sin(a1))
		x2 := float32(s.OriginX + float64(currentRadius)*math.Cos(a2))
		y2 := float32(s.OriginY + float64(currentRadius)*math.Sin(a2))

		clr := color.RGBA{R: 255, G: 160, B: 40, A: uint8(220 * (1 - progress))}
		vector.StrokeLine(dst, x1, y1, x2, y2, 4.0, clr, true)

		// Rastro interno incandescente
		x1Inner := float32(s.OriginX + float64(currentRadius*0.8)*math.Cos(a1))
		y1Inner := float32(s.OriginY + float64(currentRadius*0.8)*math.Sin(a1))
		x2Inner := float32(s.OriginX + float64(currentRadius*0.8)*math.Cos(a2))
		y2Inner := float32(s.OriginY + float64(currentRadius*0.8)*math.Sin(a2))
		vector.StrokeLine(dst, x1Inner, y1Inner, x2Inner, y2Inner, 2.0, color.RGBA{R: 255, G: 240, B: 180, A: uint8(180 * (1 - progress))}, true)
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
