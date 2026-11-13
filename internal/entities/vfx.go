package entities

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Particle representa uma partícula visual leve na arena.
type Particle struct {
	X, Y     float64
	Vx, Vy   float64
	Size     float32
	Color    color.RGBA
	Life     int
	MaxLife  int
	Friction float64
	Gravity  float64
	IsBubble bool
}

// NewParticle cria uma partícula simples.
func NewParticle(x, y, vx, vy float64, size float32, clr color.RGBA, life int) *Particle {
	return &Particle{
		X:        x,
		Y:        y,
		Vx:       vx,
		Vy:       vy,
		Size:     size,
		Color:    clr,
		Life:     life,
		MaxLife:  life,
		Friction: 0.94,
	}
}

// Update atualiza posição e decaimento.
func (p *Particle) Update() bool {
	p.Vx *= p.Friction
	p.Vy *= p.Friction
	p.Vy += p.Gravity
	p.X += p.Vx
	p.Y += p.Vy
	p.Life--
	return p.Life > 0
}

// Draw renderiza a partícula com transparência proporcional à vida restante.
func (p *Particle) Draw(dst *ebiten.Image) {
	if p.Life <= 0 {
		return
	}
	alphaRatio := float64(p.Life) / float64(p.MaxLife)
	clr := p.Color
	clr.A = uint8(float64(p.Color.A) * alphaRatio)

	if p.IsBubble {
		vector.StrokeCircle(dst, float32(p.X), float32(p.Y), p.Size, 1.2, clr, true)
	} else {
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), p.Size, clr, true)
	}
}

// ShockwaveEffect representa uma onda de choque expansiva circular (Disparo Sísmico, Pisão do Boss, Rachadura).
type ShockwaveEffect struct {
	X, Y       float64
	CurrentRad float64
	MaxRad     float64
	Color      color.RGBA
	Life       int
	MaxLife    int
	IsCrack    bool // Se true, desenha rachaduras radiais avermelhadas no chão
}

// NewShockwave cria uma onda expansiva.
func NewShockwave(x, y, maxRad float64, clr color.RGBA, life int, isCrack bool) *ShockwaveEffect {
	return &ShockwaveEffect{
		X:          x,
		Y:          y,
		CurrentRad: 10,
		MaxRad:     maxRad,
		Color:      clr,
		Life:       life,
		MaxLife:    life,
		IsCrack:    isCrack,
	}
}

// Update expande a onda e decrementa a duração.
func (s *ShockwaveEffect) Update() bool {
	s.Life--
	progress := 1.0 - float64(s.Life)/float64(s.MaxLife)
	s.CurrentRad = 10 + (s.MaxRad-10)*progress
	return s.Life > 0
}

// Draw renderiza a onda de choque ou rachaduras.
func (s *ShockwaveEffect) Draw(dst *ebiten.Image) {
	if s.Life <= 0 {
		return
	}
	alphaRatio := float64(s.Life) / float64(s.MaxLife)
	clr := s.Color
	clr.A = uint8(float64(s.Color.A) * alphaRatio)

	if s.IsCrack {
		// Círculo de rachadura e linhas radiais de impacto
		vector.StrokeCircle(dst, float32(s.X), float32(s.Y), float32(s.MaxRad*0.85), 2.5, clr, true)
		for i := 0; i < 8; i++ {
			ang := float64(i) * (math.Pi / 4.0)
			x2 := float32(s.X + s.CurrentRad*math.Cos(ang))
			y2 := float32(s.Y + s.CurrentRad*math.Sin(ang))
			vector.StrokeLine(dst, float32(s.X), float32(s.Y), x2, y2, 2.0, clr, true)
		}
	} else {
		vector.StrokeCircle(dst, float32(s.X), float32(s.Y), float32(s.CurrentRad), 3.5, clr, true)
		innerClr := clr
		innerClr.A = uint8(float64(clr.A) * 0.4)
		vector.DrawFilledCircle(dst, float32(s.X), float32(s.Y), float32(s.CurrentRad*0.92), innerClr, true)
	}
}

// SpawnBurst cria uma explosão de partículas em torno de (x, y).
func SpawnBurst(x, y float64, count int, speed float64, size float32, clr color.RGBA, life int) []*Particle {
	out := make([]*Particle, 0, count)
	for i := 0; i < count; i++ {
		ang := rand.Float64() * math.Pi * 2
		spd := (0.3 + rand.Float64()*0.7) * speed
		p := NewParticle(x, y, math.Cos(ang)*spd, math.Sin(ang)*spd, size, clr, life+rand.Intn(10))
		out = append(out, p)
	}
	return out
}
