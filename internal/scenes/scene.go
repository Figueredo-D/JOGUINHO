package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/entities"
)

// Scene define a interface comum para qualquer tela do jogo.
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

// SceneNavigator provê métodos para transição entre cenas e ações globais de jogo/áudio.
type SceneNavigator interface {
	SwitchToMenu()
	SwitchToClassSelect()
	SwitchToSanctuary()
	SwitchToAchievements()
	SwitchToHowToPlay()
	SwitchToGameplay(class entities.ClassType)
	SwitchToGameOver(score int, wave int, class entities.ClassType, soulsGained int)
	SwitchToVictory(summary RunSummary)
	SwitchToOptions()

	ToggleFullscreen()
	IsFullscreen() bool

	SetMasterVolume(vol int)
	GetMasterVolume() int
	SetSFXVolume(vol int)
	GetSFXVolume() int
	SetMusicVolume(vol int)
	GetMusicVolume() int

	QuitGame()
}

// drawRivet desenha um cravo/rebite metálico medieval com relevo.
func drawRivet(dst *ebiten.Image, x, y float32, clr color.RGBA) {
	vector.DrawFilledCircle(dst, x, y, 2.5, color.RGBA{R: 20, G: 15, B: 10, A: 220}, true)
	vector.DrawFilledCircle(dst, x-0.5, y-0.5, 2.0, clr, true)
	vector.DrawFilledCircle(dst, x-0.8, y-0.8, 0.9, color.RGBA{R: 255, G: 255, B: 240, A: 230}, true)
}
