package scenes

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/ui"
)

// OptionsScene gerencia as configurações com estética de fortaleza medieval.
type OptionsScene struct {
	navigator     SceneNavigator
	btnFullscreen *ui.Button
	btnVolDown    *ui.Button
	btnVolUp      *ui.Button
	btnBack       *ui.Button
}

// NewOptionsScene instancia a tela de opções medievais.
func NewOptionsScene(nav SceneNavigator) *OptionsScene {
	fsLabel := "⚙  TELA CHEIA: DESATIVADO  ⚙"
	if nav.IsFullscreen() {
		fsLabel = "⚙  TELA CHEIA: ATIVADO  ⚙"
	}

	return &OptionsScene{
		navigator:     nav,
		btnFullscreen: ui.NewButton(270, 220, 260, 46, fsLabel),
		btnVolDown:    ui.NewButton(270, 310, 65, 46, " ⚔ - "),
		btnVolUp:      ui.NewButton(465, 310, 65, 46, " + ⚔ "),
		btnBack:       ui.NewButton(270, 420, 260, 46, "🏰  VOLTAR AO MENU  🏰"),
	}
}

// Update processa cliques nas opções.
func (o *OptionsScene) Update() error {
	o.btnFullscreen.Update()
	o.btnVolDown.Update()
	o.btnVolUp.Update()
	o.btnBack.Update()

	if o.btnFullscreen.IsClicked() {
		o.navigator.ToggleFullscreen()
		if o.navigator.IsFullscreen() {
			o.btnFullscreen.Text = "⚙  TELA CHEIA: ATIVADO  ⚙"
		} else {
			o.btnFullscreen.Text = "⚙  TELA CHEIA: DESATIVADO  ⚙"
		}
	}

	if o.btnVolDown.IsClicked() {
		vol := o.navigator.GetVolume() - 10
		if vol < 0 {
			vol = 0
		}
		o.navigator.SetVolume(vol)
	}

	if o.btnVolUp.IsClicked() {
		vol := o.navigator.GetVolume() + 10
		if vol > 100 {
			vol = 100
		}
		o.navigator.SetVolume(vol)
	}

	if o.btnBack.IsClicked() {
		o.navigator.SwitchToMenu()
	}

	return nil
}

// Draw renderiza as opções do reino.
func (o *OptionsScene) Draw(screen *ebiten.Image) {
	// Fundo de muralha de pedra
	screen.Fill(color.RGBA{R: 24, G: 20, B: 18, A: 255})
	for y := 0; y < 600; y += 36 {
		for x := 0; x < 800; x += 70 {
			vector.DrawFilledRect(screen, float32(x+2), float32(y+2), 66, 32, color.RGBA{R: 38, G: 32, B: 28, A: 255}, false)
			vector.StrokeRect(screen, float32(x+2), float32(y+2), 66, 32, 1, color.RGBA{R: 18, G: 14, B: 12, A: 200}, false)
		}
	}

	// Placa do título em carvalho
	vector.DrawFilledRect(screen, 230, 85, 340, 72, color.RGBA{R: 44, G: 30, B: 20, A: 255}, false)
	vector.StrokeRect(screen, 230, 85, 340, 72, 2.5, color.RGBA{R: 185, G: 145, B: 75, A: 255}, false)
	vector.StrokeRect(screen, 234, 89, 332, 64, 1.0, color.RGBA{R: 90, G: 70, B: 40, A: 200}, false)

	ebitenutil.DebugPrintAt(screen, "=== ⚜ OPCOES DO REINO ⚜ ===", 285, 105)
	ebitenutil.DebugPrintAt(screen, "Ajuste os decretos e parametros visuais", 270, 126)

	// Botão de Fullscreen
	o.btnFullscreen.Draw(screen)

	// Mostrador central de Volume
	vector.DrawFilledRect(screen, 345, 310, 110, 46, color.RGBA{R: 36, G: 25, B: 18, A: 255}, false)
	vector.StrokeRect(screen, 345, 310, 110, 46, 2.0, color.RGBA{R: 160, G: 120, B: 60, A: 255}, false)
	volStr := fmt.Sprintf("Bruma: %d%%", o.navigator.GetVolume())
	ebitenutil.DebugPrintAt(screen, volStr, 365, 327)

	o.btnVolDown.Draw(screen)
	o.btnVolUp.Draw(screen)

	// Botão voltar
	o.btnBack.Draw(screen)
}
