package entities

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
)

// ShieldState gerencia a Postura de Escudo do Guerreiro.
type ShieldState struct {
	Active        bool
	Timer         int     // Frames restantes ativo
	MaxDuration   int     // Duração total (ex: 75 frames ~1.25s)
	Cooldown      int     // Cooldown atual
	MaxCooldown   int     // Cooldown total (ex: 420 frames ~7s)
	Angle         float64 // Direção do escudo
	ArcSpan       float64 // Abertura em radianos (ex: PI * 0.7)
	Radius        float64 // Raio de cobertura frontal
	ReflectDamage float64 // Multiplicador de dano dos tiros refletidos
}

// ArrowRainEffect gerencia a Chuva de Flechas do Arqueiro.
type ArrowRainEffect struct {
	TargetX, TargetY float64
	Radius           float64
	Duration         int
	MaxDuration      int
	DamagePerTick    float64
	Active           bool
	Arrows           []rainArrow
}

type rainArrow struct {
	x, y   float64
	target float64
	speed  float64
}

// NewArrowRain cria uma área de chuva de flechas do Arqueiro.
func NewArrowRain(tx, ty, radius, damage float64, duration int) *ArrowRainEffect {
	return &ArrowRainEffect{
		TargetX:       tx,
		TargetY:       ty,
		Radius:        radius,
		Duration:      duration,
		MaxDuration:   duration,
		DamagePerTick: damage,
		Active:        true,
		Arrows:        make([]rainArrow, 0, 20),
	}
}

// Update avança a animação da chuva de flechas e gera novas flechas caindo.
func (ar *ArrowRainEffect) Update() {
	if !ar.Active {
		return
	}
	ar.Duration--
	if ar.Duration <= 0 {
		ar.Active = false
		return
	}

	// Adiciona novas flechas caindo do alto
	if rand.Float64() < 0.65 {
		angle := rand.Float64() * math.Pi * 2
		dist := rand.Float64() * ar.Radius
		spawnX := ar.TargetX + dist*math.Cos(angle)
		spawnY := ar.TargetY + dist*math.Sin(angle)
		ar.Arrows = append(ar.Arrows, rainArrow{
			x:      spawnX,
			y:      spawnY - 90 - rand.Float64()*40,
			target: spawnY,
			speed:  12.0 + rand.Float64()*4.0,
		})
	}

	// Atualiza flechas
	for i := len(ar.Arrows) - 1; i >= 0; i-- {
		ar.Arrows[i].y += ar.Arrows[i].speed
		if ar.Arrows[i].y >= ar.Arrows[i].target {
			ar.Arrows = append(ar.Arrows[:i], ar.Arrows[i+1:]...)
		}
	}
}

// CollidesWith verifica se um inimigo está dentro da área da chuva de flechas.
func (ar *ArrowRainEffect) CollidesWith(eBounds collision.Rect) bool {
	if !ar.Active {
		return false
	}
	cx := eBounds.X + eBounds.W/2
	cy := eBounds.Y + eBounds.H/2
	dist := math.Hypot(cx-ar.TargetX, cy-ar.TargetY)
	return dist <= ar.Radius+eBounds.W/2
}

// Draw renderiza a área no chão e as flechas caindo proceduralmente.
func (ar *ArrowRainEffect) Draw(dst *ebiten.Image) {
	if !ar.Active {
		return
	}

	// Círculo de área no chão com efeito pulsante
	alpha := uint8(60)
	if ar.Duration < 30 {
		alpha = uint8(float64(alpha) * (float64(ar.Duration) / 30.0))
	}
	vector.DrawFilledCircle(dst, float32(ar.TargetX), float32(ar.TargetY), float32(ar.Radius), color.RGBA{R: 50, G: 160, B: 80, A: alpha}, true)
	vector.StrokeCircle(dst, float32(ar.TargetX), float32(ar.TargetY), float32(ar.Radius), 2, color.RGBA{R: 120, G: 240, B: 140, A: alpha * 2}, true)

	// Desenha as flechas caindo em diagonal acentuada
	for _, a := range ar.Arrows {
		vector.StrokeLine(dst, float32(a.x), float32(a.y), float32(a.x), float32(a.y+14), 2.0, color.RGBA{R: 240, G: 245, B: 200, A: 240}, true)
		vector.DrawFilledCircle(dst, float32(a.x), float32(a.y+14), 1.5, color.RGBA{R: 255, G: 255, B: 255, A: 255}, true)
	}
}

// FirePillarEffect gerencia o Pilar de Fogo do Mago.
type FirePillarEffect struct {
	TargetX, TargetY float64
	Radius           float64
	Duration         int
	MaxDuration      int
	DamagePerTick    float64
	Active           bool
	FlameParticles   []flameParticle
}

type flameParticle struct {
	x, y   float64
	vx, vy float64
	life   int
	clr    color.RGBA
}

// NewFirePillar cria uma erupção de pilar de fogo do Mago.
func NewFirePillar(tx, ty, radius, damage float64, duration int) *FirePillarEffect {
	return &FirePillarEffect{
		TargetX:        tx,
		TargetY:        ty,
		Radius:         radius,
		Duration:       duration,
		MaxDuration:    duration,
		DamagePerTick:  damage,
		Active:         true,
		FlameParticles: make([]flameParticle, 0, 40),
	}
}

// Update avança o ciclo do pilar de fogo e gera labaredas verticais ardentes.
func (fp *FirePillarEffect) Update() {
	if !fp.Active {
		return
	}
	fp.Duration--
	if fp.Duration <= 0 {
		fp.Active = false
		return
	}

	// Gera labaredas e fagulhas subindo
	for i := 0; i < 3; i++ {
		offset := (rand.Float64() - 0.5) * fp.Radius * 1.5
		fp.FlameParticles = append(fp.FlameParticles, flameParticle{
			x:    fp.TargetX + offset,
			y:    fp.TargetY + (rand.Float64()-0.5)*15,
			vx:   (rand.Float64() - 0.5) * 1.5,
			vy:   -4.5 - rand.Float64()*4.0,
			life: 18 + rand.Intn(16),
			clr: color.RGBA{
				R: 255,
				G: uint8(100 + rand.Intn(140)),
				B: 20,
				A: 255,
			},
		})
	}

	for i := len(fp.FlameParticles) - 1; i >= 0; i-- {
		fp.FlameParticles[i].x += fp.FlameParticles[i].vx
		fp.FlameParticles[i].y += fp.FlameParticles[i].vy
		fp.FlameParticles[i].life--
		if fp.FlameParticles[i].life <= 0 {
			fp.FlameParticles = append(fp.FlameParticles[:i], fp.FlameParticles[i+1:]...)
		}
	}
}

// CollidesWith verifica se o inimigo está dentro da coluna de fogo.
func (fp *FirePillarEffect) CollidesWith(eBounds collision.Rect) bool {
	if !fp.Active {
		return false
	}
	cx := eBounds.X + eBounds.W/2
	cy := eBounds.Y + eBounds.H/2
	dist := math.Hypot(cx-fp.TargetX, cy-fp.TargetY)
	return dist <= fp.Radius+eBounds.W/2
}

// Draw renderiza o pilar vertical de fogo e as runas incandescentes.
func (fp *FirePillarEffect) Draw(dst *ebiten.Image) {
	if !fp.Active {
		return
	}

	// 1. Runa mágica no chão
	vector.DrawFilledCircle(dst, float32(fp.TargetX), float32(fp.TargetY), float32(fp.Radius), color.RGBA{R: 220, G: 40, B: 20, A: 90}, true)
	vector.StrokeCircle(dst, float32(fp.TargetX), float32(fp.TargetY), float32(fp.Radius), 2.5, color.RGBA{R: 255, G: 160, B: 30, A: 200}, true)
	vector.DrawFilledCircle(dst, float32(fp.TargetX), float32(fp.TargetY), float32(fp.Radius*0.5), color.RGBA{R: 255, G: 220, B: 50, A: 140}, true)

	// 2. Coluna de chamas vertical
	pillarHeight := float32(110.0)
	pillarWidth := float32(fp.Radius * 1.6)
	px := float32(fp.TargetX) - pillarWidth/2
	py := float32(fp.TargetY) - pillarHeight

	// Gradiente/fogo vertical simulado
	vector.DrawFilledRect(dst, px, py, pillarWidth, pillarHeight, color.RGBA{R: 255, G: 60, B: 20, A: 95}, false)
	vector.DrawFilledRect(dst, px+pillarWidth*0.2, py, pillarWidth*0.6, pillarHeight, color.RGBA{R: 255, G: 180, B: 30, A: 130}, false)
	vector.DrawFilledRect(dst, px+pillarWidth*0.35, py, pillarWidth*0.3, pillarHeight, color.RGBA{R: 255, G: 255, B: 150, A: 180}, false)

	// 3. Fagulhas e línguas de fogo subindo
	for _, p := range fp.FlameParticles {
		rad := float32(p.life) * 0.28
		if rad < 1.0 {
			rad = 1.0
		}
		vector.DrawFilledCircle(dst, float32(p.x), float32(p.y), rad, p.clr, true)
	}
}
