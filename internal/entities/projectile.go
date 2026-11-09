package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
)

// ProjectileKind define o tipo visual e comportamento do projétil.
type ProjectileKind int

const (
	ProjectileKindArrow ProjectileKind = iota
	ProjectileKindMagic
	ProjectileKindMagicBasic
	ProjectileKindEnemyAcid
)

// Projectile representa um projétil disparado pelo jogador ou por um inimigo atirador.
type Projectile struct {
	X, Y        float64
	Vx, Vy      float64
	Radius      float64
	Damage      float64
	Color       color.RGBA
	Active      bool
	Lifetime    int
	Kind        ProjectileKind
	IsEnemy     bool // Se true, colide com o jogador; se false, colide com os monstros
	FrameCount  int
}

// NewProjectile cria um novo projétil com estilo customizado.
func NewProjectile(x, y, vx, vy, radius, damage float64, clr color.RGBA, lifetime int, kind ProjectileKind, isEnemy bool) *Projectile {
	return &Projectile{
		X:          x,
		Y:          y,
		Vx:         vx,
		Vy:         vy,
		Radius:     radius,
		Damage:     damage,
		Color:      clr,
		Active:     true,
		Lifetime:   lifetime,
		Kind:       kind,
		IsEnemy:    isEnemy,
		FrameCount: 0,
	}
}

// Update move o projétil e decrementa seu tempo de vida.
func (p *Projectile) Update() {
	if !p.Active {
		return
	}
	p.FrameCount++
	p.X += p.Vx
	p.Y += p.Vy
	p.Lifetime--
	if p.Lifetime <= 0 {
		p.Active = false
	}
}

// Bounds retorna a caixa delimitadora do projétil para colisão AABB.
func (p *Projectile) Bounds() collision.Rect {
	return collision.NewRect(p.X-p.Radius, p.Y-p.Radius, p.Radius*2, p.Radius*2)
}

// Draw renderiza o projétil na tela com estética visual correspondente.
func (p *Projectile) Draw(dst *ebiten.Image) {
	if !p.Active {
		return
	}

	switch p.Kind {
	case ProjectileKindArrow:
		// Flecha: ponta afiada com traço de rastro
		vector.StrokeLine(dst, float32(p.X-p.Vx*1.8), float32(p.Y-p.Vy*1.8), float32(p.X), float32(p.Y), 2.5, p.Color, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), float32(p.Radius), color.RGBA{R: 255, G: 255, B: 210, A: 255}, true)

	case ProjectileKindMagic:
		// Esfera arcana de alta energia com partículas cintilantes
		pulse := math.Sin(float64(p.FrameCount)*0.2) * 1.5
		rad := float32(p.Radius) + float32(pulse)
		haloColor := color.RGBA{R: 40, G: 120, B: 255, A: 100}
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad*1.8, haloColor, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad, p.Color, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad*0.45, color.RGBA{R: 255, G: 255, B: 255, A: 250}, true)

	case ProjectileKindMagicBasic:
		// Centelha menor de energia
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), float32(p.Radius), p.Color, true)

	case ProjectileKindEnemyAcid:
		// Cuspe de ácido/fogo do inimigo Atirador
		pulse := math.Sin(float64(p.FrameCount)*0.25) * 1.2
		rad := float32(p.Radius) + float32(pulse)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad*1.5, color.RGBA{R: 40, G: 180, B: 60, A: 120}, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad, color.RGBA{R: 80, G: 250, B: 100, A: 255}, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad*0.4, color.RGBA{R: 240, G: 255, B: 180, A: 255}, true)
	}
}
