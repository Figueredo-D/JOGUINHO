package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"joguinho/internal/audio"
	"joguinho/internal/entities"
	"joguinho/internal/progression"
	"joguinho/internal/scenes"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

// Game implementa a interface ebiten.Game e gerencia cenas, áudio BGM e configurações persistidas.
type Game struct {
	currentScene scenes.Scene
	isFullscreen bool
	masterVolume int
	sfxVolume    int
	musicVolume  int
	shouldQuit   bool
}

// NewGame carrega as configurações salvas, inicializa o sintetizador/BGM e abre o Menu Principal.
func NewGame() *Game {
	save := progression.GetSave()
	cfg := save.Settings

	audio.Init()
	audio.SetMasterVolume(cfg.MasterVolume)
	audio.SetSFXVolume(cfg.SFXVolume)
	audio.SetMusicVolume(cfg.MusicVolume)

	if cfg.Fullscreen {
		ebiten.SetFullscreen(true)
	}

	g := &Game{
		isFullscreen: cfg.Fullscreen,
		masterVolume: cfg.MasterVolume,
		sfxVolume:    cfg.SFXVolume,
		musicVolume:  cfg.MusicVolume,
	}
	g.SwitchToMenu()
	return g
}

// Update executa o ciclo de lógica da cena ativa.
func (g *Game) Update() error {
	audio.Tick()
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

// Layout define a resolução lógica interna do jogo (800x600).
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// --- Implementação de scenes.SceneNavigator ---

func (g *Game) SwitchToMenu() {
	audio.PlayBGM(audio.TrackMenu)
	g.currentScene = scenes.NewMenuScene(g)
}

func (g *Game) SwitchToClassSelect() {
	g.currentScene = scenes.NewClassSelectScene(g)
}

func (g *Game) SwitchToSanctuary() {
	g.currentScene = scenes.NewSanctuaryScene(g)
}

func (g *Game) SwitchToAchievements() {
	g.currentScene = scenes.NewAchievementsScene(g)
}

func (g *Game) SwitchToHowToPlay() {
	g.currentScene = scenes.NewHowToPlayScene(g)
}

func (g *Game) SwitchToGameplay(class entities.ClassType) {
	audio.PlayBGM(audio.TrackExploration)
	g.currentScene = scenes.NewGameplayScene(g, class)
}

func (g *Game) SwitchToGameOver(score int, wave int, class entities.ClassType, soulsGained int) {
	audio.PlayBGM(audio.TrackMenu)
	g.currentScene = scenes.NewGameOverScene(g, score, wave, class, soulsGained)
}

func (g *Game) SwitchToVictory(summary scenes.RunSummary) {
	audio.PlayBGM(audio.TrackMenu)
	g.currentScene = scenes.NewVictoryScene(g, summary)
}

func (g *Game) SwitchToOptions() {
	g.currentScene = scenes.NewOptionsScene(g)
}

func (g *Game) ToggleFullscreen() {
	g.isFullscreen = !g.isFullscreen
	ebiten.SetFullscreen(g.isFullscreen)

	save := progression.GetSave()
	save.Settings.Fullscreen = g.isFullscreen
	_ = progression.Save()
}

func (g *Game) IsFullscreen() bool {
	return g.isFullscreen
}

func (g *Game) SetMasterVolume(vol int) {
	g.masterVolume = vol
	audio.SetMasterVolume(vol)

	save := progression.GetSave()
	save.Settings.MasterVolume = vol
	_ = progression.Save()
}

func (g *Game) GetMasterVolume() int {
	return g.masterVolume
}

func (g *Game) SetSFXVolume(vol int) {
	g.sfxVolume = vol
	audio.SetSFXVolume(vol)

	save := progression.GetSave()
	save.Settings.SFXVolume = vol
	_ = progression.Save()
}

func (g *Game) GetSFXVolume() int {
	return g.sfxVolume
}

func (g *Game) SetMusicVolume(vol int) {
	g.musicVolume = vol
	audio.SetMusicVolume(vol)

	save := progression.GetSave()
	save.Settings.MusicVolume = vol
	_ = progression.Save()
}

func (g *Game) GetMusicVolume() int {
	return g.musicVolume
}

func (g *Game) QuitGame() {
	g.shouldQuit = true
}
