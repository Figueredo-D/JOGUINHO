package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
)

// ProgressBar exibe uma barra de progresso (vida, mana, cooldown).
type ProgressBar struct {
	Bounds       collision.Rect
	Label        string
	Current      float64
	Max          float64
	FillColor    color.RGBA
	BgColor      color.RGBA
	BorderColor  color.RGBA
	ShowFraction bool
}

// NewProgressBar cria uma nova barra de progresso customizada.
func NewProgressBar(x, y, w, h float64, label string, max float64, fillColor color.RGBA) *ProgressBar {
	return &ProgressBar{
		Bounds:       collision.NewRect(x, y, w, h),
		Label:        label,
		Current:      max,
		Max:          max,
		FillColor:    fillColor,
		BgColor:      color.RGBA{R: 30, G: 30, B: 35, A: 220},
		BorderColor:  color.RGBA{R: 180, G: 190, B: 210, A: 255},
		ShowFraction: true,
	}
}

// SetValues atualiza o valor atual e o máximo da barra.
func (p *ProgressBar) SetValues(current, max float64) {
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}
	p.Current = current
	p.Max = max
}

// Draw renderiza a barra na tela.
func (p *ProgressBar) Draw(dst *ebiten.Image) {
	// Fundo da barra
	vector.DrawFilledRect(dst, float32(p.Bounds.X), float32(p.Bounds.Y), float32(p.Bounds.W), float32(p.Bounds.H), p.BgColor, false)

	// Preenchimento
	if p.Max > 0 && p.Current > 0 {
		ratio := p.Current / p.Max
		if ratio > 1 {
			ratio = 1
		}
		fillW := float32(p.Bounds.W) * float32(ratio)
		vector.DrawFilledRect(dst, float32(p.Bounds.X), float32(p.Bounds.Y), fillW, float32(p.Bounds.H), p.FillColor, false)
	}

	// Borda
	vector.StrokeRect(dst, float32(p.Bounds.X), float32(p.Bounds.Y), float32(p.Bounds.W), float32(p.Bounds.H), 1.5, p.BorderColor, false)

	// Texto/Rótulo
	var text string
	if p.ShowFraction {
		text = fmt.Sprintf("%s: %.0f/%.0f", p.Label, p.Current, p.Max)
	} else {
		text = p.Label
	}

	tx := int(p.Bounds.X + 8)
	ty := int(p.Bounds.Y + (p.Bounds.H-12)/2)
	ebitenutil.DebugPrintAt(dst, text, tx, ty)
}
