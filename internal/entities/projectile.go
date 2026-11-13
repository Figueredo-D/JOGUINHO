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
	ProjectileKindPellet    // Balote de espingarda calibre 12
	ProjectileKindRevolver  // Tiro rápido de revólver do Xerife
	ProjectileKindDeadeye   // Tiro traçante dourado teleguiado (Deadeye)
	ProjectileKindBossOrb   // Orbe sombrio teleguiado lento do Chefe
)

// Projectile representa um projétil disparado pelo jogador ou por um inimigo/chefe.
type Projectile struct {
	X, Y        float64
	Vx, Vy      float64
	Radius      float64
	Damage      float64
	Knockback   float64
	IsCrit      bool
	IsHoming    bool
	HomingTurn  float64
	Color       color.RGBA
	Active      bool
	Lifetime    int
	Kind        ProjectileKind
	IsEnemy     bool
	FrameCount  int
	TargetID    uint64          // Opcional para tiros teleguiados do Xerife
	PierceLeft  int             // Quantos alvos extras pode perfurar (Lente de Foco)
	HitIDs      map[uint64]bool // Evita acertar o mesmo inimigo múltiplas vezes ao perfurar
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

// SteerTowards ajusta gradualmente a trajetória do projétil em direção a (tx, ty) mantendo a velocidade escalar.
func (p *Projectile) SteerTowards(tx, ty, turnFactor float64) {
	if !p.Active {
		return
	}
	spd := math.Hypot(p.Vx, p.Vy)
	if spd < 0.001 {
		return
	}
	dx := tx - p.X
	dy := ty - p.Y
	dist := math.Hypot(dx, dy)
	if dist < 1.0 {
		return
	}

	desiredVx := (dx / dist) * spd
	desiredVy := (dy / dist) * spd

	p.Vx = p.Vx*(1-turnFactor) + desiredVx*turnFactor
	p.Vy = p.Vy*(1-turnFactor) + desiredVy*turnFactor

	// Re-normaliza velocidade
	newSpd := math.Hypot(p.Vx, p.Vy)
	if newSpd > 0.001 {
		p.Vx = (p.Vx / newSpd) * spd
		p.Vy = (p.Vy / newSpd) * spd
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
		vector.StrokeLine(dst, float32(p.X-p.Vx*1.8), float32(p.Y-p.Vy*1.8), float32(p.X), float32(p.Y), 2.5, p.Color, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), float32(p.Radius), color.RGBA{R: 255, G: 255, B: 210, A: 255}, true)

	case ProjectileKindMagic:
		pulse := math.Sin(float64(p.FrameCount)*0.2) * 1.5
		rad := float32(p.Radius) + float32(pulse)
		haloColor := color.RGBA{R: 40, G: 120, B: 255, A: 100}
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad*1.8, haloColor, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad, p.Color, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad*0.45, color.RGBA{R: 255, G: 255, B: 255, A: 250}, true)

	case ProjectileKindMagicBasic:
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), float32(p.Radius), p.Color, true)

	case ProjectileKindEnemyAcid:
		pulse := math.Sin(float64(p.FrameCount)*0.25) * 1.2
		rad := float32(p.Radius) + float32(pulse)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad*1.5, color.RGBA{R: 40, G: 180, B: 60, A: 120}, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad, color.RGBA{R: 80, G: 250, B: 100, A: 255}, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad*0.4, color.RGBA{R: 240, G: 255, B: 180, A: 255}, true)

	case ProjectileKindPellet:
		// Balote incandescente curto e grosso
		vector.StrokeLine(dst, float32(p.X-p.Vx*1.1), float32(p.Y-p.Vy*1.1), float32(p.X), float32(p.Y), 3.0, color.RGBA{R: 255, G: 160, B: 40, A: 220}, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), float32(p.Radius), color.RGBA{R: 255, G: 240, B: 180, A: 255}, true)

	case ProjectileKindRevolver:
		// Tiro prateado de alta velocidade
		vector.StrokeLine(dst, float32(p.X-p.Vx*1.5), float32(p.Y-p.Vy*1.5), float32(p.X), float32(p.Y), 2.2, color.RGBA{R: 210, G: 225, B: 255, A: 230}, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), float32(p.Radius), color.RGBA{R: 255, G: 255, B: 255, A: 255}, true)

	case ProjectileKindDeadeye:
		// Traçante dourado luminoso ultra-rápido do Deadeye!
		vector.StrokeLine(dst, float32(p.X-p.Vx*2.8), float32(p.Y-p.Vy*2.8), float32(p.X), float32(p.Y), 3.2, color.RGBA{R: 255, G: 215, B: 40, A: 240}, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), float32(p.Radius*1.4), color.RGBA{R: 255, G: 240, B: 120, A: 140}, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), float32(p.Radius), color.RGBA{R: 255, G: 255, B: 230, A: 255}, true)

	case ProjectileKindBossOrb:
		// Orbe sombrio roxo/carmesim teleguiado lento invocado pelo Chefe
		pulse := math.Sin(float64(p.FrameCount)*0.18) * 2.0
		rad := float32(p.Radius) + float32(pulse)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad*1.7, color.RGBA{R: 180, G: 20, B: 60, A: 100}, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad, color.RGBA{R: 230, G: 40, B: 90, A: 255}, true)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), rad*0.45, color.RGBA{R: 255, G: 200, B: 220, A: 255}, true)
	}
}
