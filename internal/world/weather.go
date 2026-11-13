package world

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// WeatherType define os modificadores climáticos aleatórios da arena.
type WeatherType int

const (
	WeatherNone WeatherType = iota
	WeatherFog              // Névoa Espessa (reduz visão periférica ao redor do jogador)
	WeatherAcidRain         // Chuva Ácida (gotas + poças ácidas que causam dano contínuo a todos)
	WeatherBloodMoon        // Lua Carmesim (monstros +25% vel., porém 2x drop de almas)
)

// Name retorna o título de exibição do evento climático.
func (w WeatherType) Name() string {
	switch w {
	case WeatherFog:
		return "NEVOA ESPESSA (Visao Reduzida)"
	case WeatherAcidRain:
		return "CHUVA ACIDA (Pocas Toxicas)"
	case WeatherBloodMoon:
		return "LUA CARMESIM (Inimigos Velozes, 2x Almas!)"
	default:
		return ""
	}
}

// AcidPuddle representa uma poça de ácido no chão durante a Chuva Ácida.
type AcidPuddle struct {
	X, Y     float64
	Radius   float64
	Life     int
	MaxLife  int
	Frame    int
}

type rainDrop struct {
	x, y  float64
	speed float64
}

// WeatherSystem gerencia a duração, partículas e efeitos de tela do clima ativo.
type WeatherSystem struct {
	Active       WeatherType
	Timer        int
	MaxDuration  int
	Puddles      []*AcidPuddle
	RainDrops    []rainDrop
	fogOverlay   *ebiten.Image
	frameCounter int
}

// NewWeatherSystem inicializa o gerenciador climático e pré-renderiza a máscara de névoa.
func NewWeatherSystem() *WeatherSystem {
	ws := &WeatherSystem{
		Active:    WeatherNone,
		Puddles:   make([]*AcidPuddle, 0, 10),
		RainDrops: make([]rainDrop, 0, 60),
	}
	ws.buildFogOverlay()
	return ws
}

// buildFogOverlay pré-renderiza uma máscara radial 1600x1200 com abertura central suave para visão limitada.
func (ws *WeatherSystem) buildFogOverlay() {
	w, h := 1600, 1200
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	cx, cy := float64(w)/2, float64(h)/2

	clearRad := 135.0
	fadeEnd := 275.0
	maxAlpha := 228.0

	// Preenche toda a máscara com névoa escura usando acesso direto a Pix para velocidade instantânea
	for y := 0; y < h; y++ {
		dy := float64(y) - cy
		dy2 := dy * dy
		for x := 0; x < w; x++ {
			alpha := uint8(maxAlpha)
			if math.Abs(dy) <= fadeEnd && math.Abs(float64(x)-cx) <= fadeEnd {
				dx := float64(x) - cx
				dist := math.Sqrt(dx*dx + dy2)
				if dist <= clearRad {
					alpha = 0
				} else if dist < fadeEnd {
					t := (dist - clearRad) / (fadeEnd - clearRad)
					alpha = uint8(t * maxAlpha)
				}
			}
			idx := (y*w + x) * 4
			rgba.Pix[idx] = 10
			rgba.Pix[idx+1] = 14
			rgba.Pix[idx+2] = 18
			rgba.Pix[idx+3] = alpha
		}
	}

	ws.fogOverlay = ebiten.NewImageFromImage(rgba)
}

// TriggerRandom inicia um evento climático aleatório por 35 segundos (2100 frames).
func (ws *WeatherSystem) TriggerRandom() WeatherType {
	options := []WeatherType{WeatherFog, WeatherAcidRain, WeatherBloodMoon}
	chosen := options[rand.Intn(len(options))]
	ws.Active = chosen
	ws.MaxDuration = 2100 // 35 segundos
	ws.Timer = ws.MaxDuration
	ws.Puddles = ws.Puddles[:0]
	ws.RainDrops = ws.RainDrops[:0]
	return chosen
}

// Update atualiza o evento climático, gotas de chuva e poças de ácido.
func (ws *WeatherSystem) Update() {
	if ws.Active == WeatherNone {
		return
	}
	ws.frameCounter++
	ws.Timer--
	if ws.Timer <= 0 {
		ws.Active = WeatherNone
		ws.Puddles = ws.Puddles[:0]
		ws.RainDrops = ws.RainDrops[:0]
		return
	}

	if ws.Active == WeatherAcidRain {
		// Gera gotas de chuva caindo
		if len(ws.RainDrops) < 55 {
			ws.RainDrops = append(ws.RainDrops, rainDrop{
				x:     rand.Float64() * 840,
				y:     52,
				speed: 7.0 + rand.Float64()*4.0,
			})
		}
		for i := len(ws.RainDrops) - 1; i >= 0; i-- {
			ws.RainDrops[i].x -= 1.5
			ws.RainDrops[i].y += ws.RainDrops[i].speed
			if ws.RainDrops[i].y > 600 {
				ws.RainDrops[i].y = 52
				ws.RainDrops[i].x = rand.Float64() * 840
			}
		}

		// Spawna poças ácidas periodicamente no chão
		if ws.frameCounter%120 == 0 && len(ws.Puddles) < 7 {
			ws.Puddles = append(ws.Puddles, &AcidPuddle{
				X:       50 + rand.Float64()*700,
				Y:       90 + rand.Float64()*460,
				Radius:  28 + rand.Float64()*14,
				Life:    420, // 7 segundos
				MaxLife: 420,
			})
		}

		for i := len(ws.Puddles) - 1; i >= 0; i-- {
			p := ws.Puddles[i]
			p.Frame++
			p.Life--
			if p.Life <= 0 {
				ws.Puddles = append(ws.Puddles[:i], ws.Puddles[i+1:]...)
			}
		}
	}
}

// DrawGroundEffects renderiza as poças de ácido no chão (abaixo dos personagens).
func (ws *WeatherSystem) DrawGroundEffects(dst *ebiten.Image) {
	if ws.Active != WeatherAcidRain {
		return
	}
	for _, p := range ws.Puddles {
		alpha := uint8(140)
		if p.Life < 60 {
			alpha = uint8(float64(alpha) * (float64(p.Life) / 60.0))
		}
		pulse := float32(math.Sin(float64(p.Frame)*0.12) * 2.5)
		vector.DrawFilledCircle(dst, float32(p.X), float32(p.Y), float32(p.Radius)+pulse, color.RGBA{R: 50, G: 210, B: 60, A: alpha}, true)
		vector.StrokeCircle(dst, float32(p.X), float32(p.Y), float32(p.Radius)+pulse, 1.5, color.RGBA{R: 130, G: 255, B: 110, A: alpha}, true)
		// Bolha interna
		bx := float32(p.X + math.Cos(float64(p.Frame)*0.08)*p.Radius*0.4)
		by := float32(p.Y + math.Sin(float64(p.Frame)*0.11)*p.Radius*0.4)
		vector.DrawFilledCircle(dst, bx, by, 3, color.RGBA{R: 200, G: 255, B: 170, A: alpha}, true)
	}
}

// DrawOverlayEffects renderiza chuva ácida no ar, iluminação Lua Carmesim ou Névoa Espessa centrada em (px, py).
func (ws *WeatherSystem) DrawOverlayEffects(dst *ebiten.Image, px, py float64) {
	switch ws.Active {
	case WeatherAcidRain:
		for _, rd := range ws.RainDrops {
			vector.StrokeLine(dst, float32(rd.x), float32(rd.y), float32(rd.x-3), float32(rd.y+10), 1.5, color.RGBA{R: 90, G: 245, B: 100, A: 160}, false)
		}

	case WeatherBloodMoon:
		// Tonalidade avermelhada sombria sobre toda a arena
		vector.DrawFilledRect(dst, 0, 52, 800, 548, color.RGBA{R: 130, G: 10, B: 20, A: 55}, false)

	case WeatherFog:
		if ws.fogOverlay != nil {
			op := &ebiten.DrawImageOptions{}
			// Centraliza o furo transparente da vinheta exatamente na posição (px, py) do jogador
			op.GeoM.Translate(px-800, py-600)
			dst.DrawImage(ws.fogOverlay, op)
		}
	}
}
