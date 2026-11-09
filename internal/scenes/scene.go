package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"

	"joguinho/internal/entities"
)

// Scene define a interface comum para qualquer tela do jogo.
type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

// SceneNavigator provê métodos para transição entre cenas e ações globais de jogo.
type SceneNavigator interface {
	SwitchToMenu()
	SwitchToClassSelect()
	SwitchToGameplay(class entities.ClassType)
	SwitchToGameOver(score int, wave int, class entities.ClassType)
	SwitchToOptions()
	ToggleFullscreen()
	IsFullscreen() bool
	SetVolume(vol int)
	GetVolume() int
	QuitGame()
}
