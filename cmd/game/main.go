package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"joguinho/internal/game"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Dungeon Survivors - Ebitengine 2D")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil && err != ebiten.Termination {
		log.Fatalf("Falha fatal na execucao do jogo: %v", err)
	}
}
