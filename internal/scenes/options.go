package scenes

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/audio"
	"joguinho/internal/collision"
	"joguinho/internal/ui"
)

// OptionsScene gerencia áudio (Master/SFX/BGM) e vídeo (Tela Cheia) com persistência.
type OptionsScene struct {
	navigator     SceneNavigator
	btnFullscreen *ui.Button

	btnMasterDown *ui.Button
	btnMasterUp   *ui.Button
	sliderMaster  collision.Rect

	btnSFXDown *ui.Button
	btnSFXUp   *ui.Button
	sliderSFX  collision.Rect

	btnMusicDown *ui.Button
	btnMusicUp   *ui.Button
	sliderMusic  collision.Rect

	btnBack *ui.Button
}

// NewOptionsScene instancia a tela de configurações completas.
func NewOptionsScene(nav SceneNavigator) *OptionsScene {
	fsLabel := "⚙  TELA CHEIA: DESATIVADO  ⚙"
	if nav.IsFullscreen() {
		fsLabel = "⚙  TELA CHEIA: ATIVADO  ⚙"
	}

	return &OptionsScene{
		navigator:     nav,
		btnFullscreen: ui.NewButton(240, 145, 320, 44, fsLabel),

		btnMasterDown: ui.NewButton(200, 225, 52, 38, " - "),
		sliderMaster:  collision.NewRect(262, 225, 276, 38),
		btnMasterUp:   ui.NewButton(548, 225, 52, 38, " + "),

		btnSFXDown: ui.NewButton(200, 295, 52, 38, " - "),
		sliderSFX:  collision.NewRect(262, 295, 276, 38),
		btnSFXUp:   ui.NewButton(548, 295, 52, 38, " + "),

		btnMusicDown: ui.NewButton(200, 365, 52, 38, " - "),
		sliderMusic:  collision.NewRect(262, 365, 276, 38),
		btnMusicUp:   ui.NewButton(548, 365, 52, 38, " + "),

		btnBack: ui.NewButton(270, 470, 260, 46, "🏰  VOLTAR AO MENU  🏰"),
	}
}

// Update processa ajustes nos botões e barras interativas de volume.
func (o *OptionsScene) Update() error {
	o.btnFullscreen.Update()
	o.btnMasterDown.Update()
	o.btnMasterUp.Update()
	o.btnSFXDown.Update()
	o.btnSFXUp.Update()
	o.btnMusicDown.Update()
	o.btnMusicUp.Update()
	o.btnBack.Update()

	if o.btnFullscreen.IsClicked() {
		o.navigator.ToggleFullscreen()
		if o.navigator.IsFullscreen() {
			o.btnFullscreen.Text = "⚙  TELA CHEIA: ATIVADO  ⚙"
		} else {
			o.btnFullscreen.Text = "⚙  TELA CHEIA: DESATIVADO  ⚙"
		}
	}

	// Ajuste por botões
	if o.btnMasterDown.IsClicked() {
		o.navigator.SetMasterVolume(clampVol(o.navigator.GetMasterVolume() - 10))
	} else if o.btnMasterUp.IsClicked() {
		o.navigator.SetMasterVolume(clampVol(o.navigator.GetMasterVolume() + 10))
	}

	if o.btnSFXDown.IsClicked() {
		o.navigator.SetSFXVolume(clampVol(o.navigator.GetSFXVolume() - 10))
		audio.Play(audio.SFXPickup)
	} else if o.btnSFXUp.IsClicked() {
		o.navigator.SetSFXVolume(clampVol(o.navigator.GetSFXVolume() + 10))
		audio.Play(audio.SFXPickup)
	}

	if o.btnMusicDown.IsClicked() {
		o.navigator.SetMusicVolume(clampVol(o.navigator.GetMusicVolume() - 10))
	} else if o.btnMusicUp.IsClicked() {
		o.navigator.SetMusicVolume(clampVol(o.navigator.GetMusicVolume() + 10))
	}

	// Ajuste direto arrastando ou clicando nos sliders com mouse
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		fmx, fmy := float64(mx), float64(my)
		if o.sliderMaster.Contains(fmx, fmy) {
			pct := int(((fmx - o.sliderMaster.X) / o.sliderMaster.W) * 100)
			o.navigator.SetMasterVolume(clampVol(pct))
		} else if o.sliderSFX.Contains(fmx, fmy) {
			pct := int(((fmx - o.sliderSFX.X) / o.sliderSFX.W) * 100)
			o.navigator.SetSFXVolume(clampVol(pct))
		} else if o.sliderMusic.Contains(fmx, fmy) {
			pct := int(((fmx - o.sliderMusic.X) / o.sliderMusic.W) * 100)
			o.navigator.SetMusicVolume(clampVol(pct))
		}
	}

	if o.btnBack.IsClicked() {
		o.navigator.SwitchToMenu()
	}

	return nil
}

func clampVol(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

// Draw renderiza a janela de configurações de áudio e tela cheia.
func (o *OptionsScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 24, G: 20, B: 18, A: 255})
	for y := 0; y < 600; y += 36 {
		for x := 0; x < 800; x += 70 {
			vector.DrawFilledRect(screen, float32(x+2), float32(y+2), 66, 32, color.RGBA{R: 38, G: 32, B: 28, A: 255}, false)
			vector.StrokeRect(screen, float32(x+2), float32(y+2), 66, 32, 1, color.RGBA{R: 18, G: 14, B: 12, A: 200}, false)
		}
	}

	// Placa do título
	vector.DrawFilledRect(screen, 210, 45, 380, 70, color.RGBA{R: 44, G: 30, B: 20, A: 255}, false)
	vector.StrokeRect(screen, 210, 45, 380, 70, 2.5, color.RGBA{R: 185, G: 145, B: 75, A: 255}, false)
	ebitenutil.DebugPrintAt(screen, "=== ⚜ OPCOES E AUDIO DO REINO ⚜ ===", 270, 64)
	ebitenutil.DebugPrintAt(screen, "Clique ou arraste nas barras para ajustar os volumes", 230, 86)

	o.btnFullscreen.Draw(screen)

	o.drawSlider(screen, o.sliderMaster, "VOLUME GERAL (MASTER)", o.navigator.GetMasterVolume(), color.RGBA{235, 185, 60, 255})
	o.btnMasterDown.Draw(screen)
	o.btnMasterUp.Draw(screen)

	o.drawSlider(screen, o.sliderSFX, "EFEITOS SONOROS (SFX)", o.navigator.GetSFXVolume(), color.RGBA{90, 200, 120, 255})
	o.btnSFXDown.Draw(screen)
	o.btnSFXUp.Draw(screen)

	o.drawSlider(screen, o.sliderMusic, "TRILHA SONORA (BGM)", o.navigator.GetMusicVolume(), color.RGBA{90, 165, 245, 255})
	o.btnMusicDown.Draw(screen)
	o.btnMusicUp.Draw(screen)

	o.btnBack.Draw(screen)
}

func (o *OptionsScene) drawSlider(dst *ebiten.Image, rect collision.Rect, label string, value int, fillClr color.RGBA) {
	x, y, w, h := float32(rect.X), float32(rect.Y), float32(rect.W), float32(rect.H)

	// Rótulo acima da barra
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%s: %d%%", label, value), int(x)+45, int(y)-16)

	// Fundo da barra
	vector.DrawFilledRect(dst, x, y, w, h, color.RGBA{R: 32, G: 24, B: 20, A: 255}, false)
	// Preenchimento proporcional
	fillW := w * float32(value) / 100.0
	vector.DrawFilledRect(dst, x+2, y+2, fillW-4, h-4, fillClr, false)
	// Marcador deslizante (ponteiro)
	if fillW > 6 && fillW < w-4 {
		vector.DrawFilledRect(dst, x+fillW-3, y-2, 6, h+4, color.RGBA{255, 245, 215, 255}, false)
	}
	// Borda
	vector.StrokeRect(dst, x, y, w, h, 2, color.RGBA{R: 165, G: 130, B: 70, A: 255}, false)
}
