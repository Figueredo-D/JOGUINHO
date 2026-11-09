package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
)

// Button representa um botão interativo na tela com estética medieval.
type Button struct {
	Bounds       collision.Rect
	Text         string
	NormalColor  color.RGBA
	HoverColor   color.RGBA
	PressedColor color.RGBA
	BorderColor  color.RGBA
	TextColor    color.RGBA

	isHovered bool
	isPressed bool
	clicked   bool
}

// NewButton cria um novo botão estilizado no padrão medieval (madeira, ferro forjado e ouro antigo).
func NewButton(x, y, w, h float64, text string) *Button {
	return &Button{
		Bounds:       collision.NewRect(x, y, w, h),
		Text:         text,
		NormalColor:  color.RGBA{R: 48, G: 34, B: 24, A: 255},  // Madeira de carvalho nobre escura
		HoverColor:   color.RGBA{R: 72, G: 48, B: 32, A: 255},  // Madeira aquecida pela tocha
		PressedColor: color.RGBA{R: 35, G: 24, B: 16, A: 255},  // Pressionado
		BorderColor:  color.RGBA{R: 160, G: 125, B: 65, A: 255}, // Bronze/Ferro forjado
		TextColor:    color.RGBA{R: 245, G: 230, B: 200, A: 255}, // Pergaminho claro
	}
}

// Update atualiza o estado do botão (hover, pressed, click).
func (b *Button) Update() {
	mx, my := ebiten.CursorPosition()
	b.isHovered = b.Bounds.Contains(float64(mx), float64(my))
	b.clicked = false

	if b.isHovered {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.isPressed = true
		} else {
			if b.isPressed && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
				b.clicked = true
			}
			b.isPressed = false
		}
	} else {
		b.isPressed = false
	}
}

// IsClicked retorna true se o botão foi clicado no frame atual.
func (b *Button) IsClicked() bool {
	return b.clicked
}

// Draw renderiza o botão na tela com acabamento de tabuleta medieval de carvalho com cravos de ferro.
func (b *Button) Draw(dst *ebiten.Image) {
	bx := float32(b.Bounds.X)
	by := float32(b.Bounds.Y)
	bw := float32(b.Bounds.W)
	bh := float32(b.Bounds.H)

	offsetY := float32(0)
	bg := b.NormalColor
	border := b.BorderColor
	textColor := b.TextColor

	if b.isPressed {
		bg = b.PressedColor
		offsetY = 1.5
	} else if b.isHovered {
		bg = b.HoverColor
		border = color.RGBA{R: 245, G: 195, B: 60, A: 255} // Dourado brilhante de destaque
		textColor = color.RGBA{R: 255, G: 245, B: 215, A: 255}
	}

	drawY := by + offsetY

	// 1. Sombra da tábua de madeira
	vector.DrawFilledRect(dst, bx+2, drawY+3, bw, bh, color.RGBA{R: 15, G: 12, B: 10, A: 160}, false)

	// 2. Preenchimento da tábua de carvalho
	vector.DrawFilledRect(dst, bx, drawY, bw, bh, bg, false)

	// Linhas sutis simulando ranhuras de tábua de madeira
	vector.StrokeLine(dst, bx+4, drawY+bh*0.33, bx+bw-4, drawY+bh*0.33, 1, color.RGBA{R: 35, G: 25, B: 18, A: 140}, false)
	vector.StrokeLine(dst, bx+4, drawY+bh*0.66, bx+bw-4, drawY+bh*0.66, 1, color.RGBA{R: 35, G: 25, B: 18, A: 140}, false)

	// 3. Moldura externa reforçada de ferro/bronze
	vector.StrokeRect(dst, bx, drawY, bw, bh, 2, border, false)

	// Moldura interna ornamental sutil
	vector.StrokeRect(dst, bx+3, drawY+3, bw-6, bh-6, 1, color.RGBA{R: border.R / 2, G: border.G / 2, B: border.B / 2, A: 180}, false)

	// 4. Quatro rebites / cravos de metal nos cantos
	rivetColor := color.RGBA{R: 210, G: 175, B: 110, A: 255}
	if b.isHovered {
		rivetColor = color.RGBA{R: 255, G: 220, B: 80, A: 255}
	}
	drawRivet(dst, bx+6, drawY+6, rivetColor)
	drawRivet(dst, bx+bw-6, drawY+6, rivetColor)
	drawRivet(dst, bx+6, drawY+bh-6, rivetColor)
	drawRivet(dst, bx+bw-6, drawY+bh-6, rivetColor)

	// 5. Texto centralizado com visual nítido
	approxCharWidth := 6.0
	textWidth := float64(len(b.Text)) * approxCharWidth
	tx := int(b.Bounds.X + (b.Bounds.W-textWidth)/2)
	ty := int(drawY + (bh-12)/2)

	// Sombra do texto para legibilidade
	ebitenutil.DebugPrintAt(dst, b.Text, tx+1, ty+1)
	_ = textColor
	ebitenutil.DebugPrintAt(dst, b.Text, tx, ty)
}

func drawRivet(dst *ebiten.Image, x, y float32, clr color.RGBA) {
	// Cravo metálico com relevo
	vector.DrawFilledCircle(dst, x, y, 2.2, color.RGBA{R: 20, G: 15, B: 10, A: 220}, true)
	vector.DrawFilledCircle(dst, x-0.5, y-0.5, 1.8, clr, true)
	vector.DrawFilledCircle(dst, x-0.8, y-0.8, 0.8, color.RGBA{R: 255, G: 255, B: 240, A: 230}, true)
}
