package world

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// BiomeType define os biomas temáticos que alternam por ciclo de hordas.
type BiomeType int

const (
	BiomeDungeon BiomeType = iota // Masmorra de Pedra Úmida
	BiomeLava                     // Caverna de Lava
	BiomeCrypt                    // Cripta Sombria
)

// BiomeName retorna o nome de exibição do bioma.
func (b BiomeType) String() string {
	switch b {
	case BiomeDungeon:
		return "Masmorra de Pedra Umida"
	case BiomeLava:
		return "Caverna de Lava"
	case BiomeCrypt:
		return "Cripta Sombria"
	default:
		return "Arena Desconhecida"
	}
}

// Arena gerencia o mapa procedural texturizado em tiles e suas variações de bioma.
type Arena struct {
	CurrentBiome BiomeType
	cachedImage  *ebiten.Image
	width        int
	height       int
	seed         int64
}

// NewArena inicializa uma arena texturizada 800x600.
func NewArena(width, height int) *Arena {
	a := &Arena{
		CurrentBiome: BiomeDungeon,
		width:        width,
		height:       height,
		seed:         42069,
	}
	a.buildImage()
	return a
}

// SetBiome troca o bioma e reconstrói a textura de fundo da arena.
func (a *Arena) SetBiome(b BiomeType) {
	if a.CurrentBiome == b && a.cachedImage != nil {
		return
	}
	a.CurrentBiome = b
	a.buildImage()
}

// Draw renderiza a arena pré-renderizada na tela.
func (a *Arena) Draw(screen *ebiten.Image) {
	if a.cachedImage != nil {
		screen.DrawImage(a.cachedImage, nil)
	}
}

func (a *Arena) buildImage() {
	img := ebiten.NewImage(a.width, a.height)
	rng := rand.New(rand.NewSource(a.seed + int64(a.CurrentBiome)*997))

	var baseBg, tileLight, tileDark, crackClr, accentClr, wallBase, wallEdge color.RGBA

	switch a.CurrentBiome {
	case BiomeDungeon:
		baseBg = color.RGBA{26, 30, 38, 255}
		tileLight = color.RGBA{38, 44, 54, 255}
		tileDark = color.RGBA{28, 33, 42, 255}
		crackClr = color.RGBA{18, 22, 28, 255}
		accentClr = color.RGBA{45, 95, 55, 220} // Musgo úmido
		wallBase = color.RGBA{42, 46, 56, 255}
		wallEdge = color.RGBA{75, 82, 98, 255}

	case BiomeLava:
		baseBg = color.RGBA{26, 18, 18, 255}
		tileLight = color.RGBA{42, 28, 26, 255}
		tileDark = color.RGBA{30, 20, 20, 255}
		crackClr = color.RGBA{18, 10, 10, 255}
		accentClr = color.RGBA{235, 95, 20, 210} // Veios de lava brilhante
		wallBase = color.RGBA{55, 25, 22, 255}
		wallEdge = color.RGBA{130, 50, 30, 255}

	case BiomeCrypt:
		baseBg = color.RGBA{20, 18, 28, 255}
		tileLight = color.RGBA{34, 28, 46, 255}
		tileDark = color.RGBA{24, 20, 34, 255}
		crackClr = color.RGBA{14, 10, 20, 255}
		accentClr = color.RGBA{115, 70, 180, 180} // Runas sombrias
		wallBase = color.RGBA{36, 26, 50, 255}
		wallEdge = color.RGBA{95, 70, 130, 255}
	}

	img.Fill(baseBg)

	// 1. Camada de chão em Lajotas / Tiles irregulares (32x32)
	tileSize := 32
	startY := 52
	for y := startY; y < a.height; y += tileSize {
		for x := 0; x < a.width; x += tileSize {
			// Variação pseudo-aleatória de tonalidade
			var clr color.RGBA
			r := rng.Float64()
			if r < 0.45 {
				clr = tileLight
			} else if r < 0.85 {
				clr = tileDark
			} else {
				clr = baseBg
			}

			// Lajota com junta de argamassa
			vector.DrawFilledRect(img, float32(x+1), float32(y+1), float32(tileSize-2), float32(tileSize-2), clr, false)

			// Relevo sutil superior na lajota
			highlight := color.RGBA{clr.R + 6, clr.G + 6, clr.B + 6, 180}
			vector.StrokeLine(img, float32(x+1), float32(y+1), float32(x+tileSize-2), float32(y+1), 1, highlight, false)

			// Fissuras e detalhes aleatórios dentro da lajota
			if rng.Float64() < 0.28 {
				fx := float32(x + 6 + rng.Intn(16))
				fy := float32(y + 6 + rng.Intn(16))
				fx2 := fx + float32(rng.Intn(12)-6)
				fy2 := fy + float32(rng.Intn(12)-6)
				vector.StrokeLine(img, fx, fy, fx2, fy2, 1.2, crackClr, false)
			}

			// Pedrinhas soltas no piso
			if rng.Float64() < 0.22 {
				px := float32(x + 4 + rng.Intn(22))
				py := float32(y + 4 + rng.Intn(22))
				vector.DrawFilledRect(img, px, py, 3, 2, color.RGBA{clr.R + 15, clr.G + 15, clr.B + 15, 200}, false)
			}

			// Elementos temáticos do bioma
			if rng.Float64() < 0.16 {
				ax := float32(x + 5 + rng.Intn(18))
				ay := float32(y + 5 + rng.Intn(18))

				switch a.CurrentBiome {
				case BiomeDungeon:
					// Musgo orgânico nas frestas
					vector.DrawFilledCircle(img, ax, ay, float32(2+rng.Intn(3)), accentClr, true)
					vector.DrawFilledCircle(img, ax+3, ay+1, 2, accentClr, true)
				case BiomeLava:
					// Rachadura com magma brilhando
					vector.StrokeLine(img, ax, ay, ax+float32(rng.Intn(14)-7), ay+float32(rng.Intn(14)-7), 2.0, accentClr, true)
					vector.DrawFilledCircle(img, ax, ay, 1.5, color.RGBA{255, 220, 80, 220}, true)
				case BiomeCrypt:
					// Pequena runa arcana ou lasca de osso
					if rng.Float64() < 0.5 {
						vector.StrokeRect(img, ax, ay, 5, 5, 1, accentClr, false)
					} else {
						vector.StrokeLine(img, ax, ay, ax+6, ay+2, 1.8, color.RGBA{190, 185, 175, 160}, false)
					}
				}
			}
		}
	}

	// 2. Bordas de Paredes de Rocha e Sombras Top-Down
	borderThick := float32(12)
	// Sombras projetadas pelas paredes no chão da arena
	shadowGrad := []color.RGBA{
		{0, 0, 0, 120},
		{0, 0, 0, 75},
		{0, 0, 0, 35},
	}
	for i, sClr := range shadowGrad {
		offset := float32(i * 5)
		// Sombra superior (abaixo do HUD)
		vector.DrawFilledRect(img, borderThick, float32(startY)+borderThick+offset, float32(a.width)-borderThick*2, 5, sClr, false)
		// Sombra esquerda
		vector.DrawFilledRect(img, borderThick+offset, float32(startY)+borderThick, 5, float32(a.height-startY)-borderThick*2, sClr, false)
		// Sombra direita
		vector.DrawFilledRect(img, float32(a.width)-borderThick-offset-5, float32(startY)+borderThick, 5, float32(a.height-startY)-borderThick*2, sClr, false)
		// Sombra inferior
		vector.DrawFilledRect(img, borderThick, float32(a.height)-borderThick-offset-5, float32(a.width)-borderThick*2, 5, sClr, false)
	}

	// Blocos das paredes externas
	// Parede esquerda e direita
	for y := startY; y < a.height; y += 24 {
		vector.DrawFilledRect(img, 0, float32(y), borderThick, 23, wallBase, false)
		vector.StrokeRect(img, 0, float32(y), borderThick, 23, 1, wallEdge, false)

		vector.DrawFilledRect(img, float32(a.width)-borderThick, float32(y), borderThick, 23, wallBase, false)
		vector.StrokeRect(img, float32(a.width)-borderThick, float32(y), borderThick, 23, 1, wallEdge, false)
	}
	// Parede superior (y=52) e inferior
	for x := 0; x < a.width; x += 28 {
		vector.DrawFilledRect(img, float32(x), float32(startY), 27, borderThick, wallBase, false)
		vector.StrokeRect(img, float32(x), float32(startY), 27, borderThick, 1, wallEdge, false)

		vector.DrawFilledRect(img, float32(x), float32(a.height)-borderThick, 27, borderThick, wallBase, false)
		vector.StrokeRect(img, float32(x), float32(a.height)-borderThick, 27, borderThick, 1, wallEdge, false)
	}

	// 3. Pilares Ornamentais nos 4 cantos da arena
	a.drawCornerPillar(img, 14, float32(startY)+14, wallBase, wallEdge, accentClr)
	a.drawCornerPillar(img, float32(a.width)-34, float32(startY)+14, wallBase, wallEdge, accentClr)
	a.drawCornerPillar(img, 14, float32(a.height)-34, wallBase, wallEdge, accentClr)
	a.drawCornerPillar(img, float32(a.width)-34, float32(a.height)-34, wallBase, wallEdge, accentClr)

	a.cachedImage = img
}

func (a *Arena) drawCornerPillar(dst *ebiten.Image, x, y float32, base, edge, accent color.RGBA) {
	// Sombra do pilar
	vector.DrawFilledRect(dst, x+4, y+5, 20, 20, color.RGBA{0, 0, 0, 140}, false)
	// Base do pilar
	vector.DrawFilledRect(dst, x, y, 20, 20, base, false)
	vector.StrokeRect(dst, x, y, 20, 20, 2, edge, false)
	// Detalhe de brasa/musgo/runa no topo do pilar
	vector.DrawFilledCircle(dst, x+10, y+10, 4, accent, true)
}
