package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"joguinho/internal/entities"
	"joguinho/internal/scenes"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

// Game implementa a interface ebiten.Game e gerencia a máquina de estados de cenas.
type Game struct {
	currentScene scenes.Scene
	isFullscreen bool
	volume       int
	shouldQuit   bool
}

// NewGame instancia o ciclo principal com a tela inicial de Menu.
func NewGame() *Game {
	g := &Game{
		volume: 80,
	}
	g.currentScene = scenes.NewMenuScene(g)
	return g
}

// Update executa o ciclo de lógica da cena ativa.
func (g *Game) Update() error {
	if g.shouldQuit {
		return ebiten.Termination
	}
	if g.currentScene != nil {
		return g.currentScene.Update()
	}
	return nil
}

// Draw renderiza a cena atual.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentScene != nil {
		g.currentScene.Draw(screen)
	}
}

// Layout define a resolução lógica interna do jogo (800x600) independente do tamanho da janela.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// --- Implementação de scenes.SceneNavigator ---

func (g *Game) SwitchToMenu() {
	g.currentScene = scenes.NewMenuScene(g)
}

func (g *Game) SwitchToClassSelect() {
	g.currentScene = scenes.NewClassSelectScene(g)
}

func (g *Game) SwitchToSanctuary() {
	g.currentScene = scenes.NewSanctuaryScene(g)
}

func (g *Game) SwitchToGameplay(class entities.ClassType) {
	g.currentScene = scenes.NewGameplayScene(g, class)
}

func (g *Game) SwitchToGameOver(score int, wave int, class entities.ClassType, soulsGained int) {
	g.currentScene = scenes.NewGameOverScene(g, score, wave, class, soulsGained)
}

func (g *Game) SwitchToOptions() {
	g.currentScene = scenes.NewOptionsScene(g)
}

func (g *Game) ToggleFullscreen() {
	g.isFullscreen = !g.isFullscreen
	ebiten.SetFullscreen(g.isFullscreen)
}

func (g *Game) IsFullscreen() bool {
	return g.isFullscreen
}

func (g *Game) SetVolume(vol int) {
	g.volume = vol
}

func (g *Game) GetVolume() int {
	return g.volume
}

func (g *Game) QuitGame() {
	g.shouldQuit = true
}
