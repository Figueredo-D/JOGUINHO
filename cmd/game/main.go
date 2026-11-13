package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"joguinho/internal/game"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("As Catacumbas de Dorfieegu")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowIcon([]image.Image{createAppIcon(32)})

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil && err != ebiten.Termination {
		log.Fatalf("Falha fatal na execucao do jogo: %v", err)
	}
}

// createAppIcon gera um ícone 32x32 procedural para a janela do executável.
func createAppIcon(size int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	gold := color.RGBA{245, 195, 55, 255}
	crimson := color.RGBA{165, 30, 40, 255}
	dark := color.RGBA{25, 20, 24, 255}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if x >= 3 && x < size-3 && y >= 3 && y < size-3 {
				if x == 3 || x == size-4 || y == 3 || y == size-4 {
					img.Set(x, y, gold)
				} else if x == size/2 || y == size/2 {
					img.Set(x, y, gold)
				} else {
					img.Set(x, y, crimson)
				}
			} else {
				img.Set(x, y, dark)
			}
		}
	}
	return img
}
