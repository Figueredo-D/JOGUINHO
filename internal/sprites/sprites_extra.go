package sprites

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// GetShotgunnerSprite gera o sprite do Atirador Pesado (jaleco escuro, bandoleira de cartuchos vermelhos e espingarda dupla).
func GetShotgunnerSprite() *ebiten.Image {
	pattern := []string{
		"....CCCCCCCC....",
		"...CCDDDDDDCC...",
		"...CDSSSSSSDC...",
		"...CDSKSSKSDC...",
		"...CDSSSSSSDC...",
		"....CCCCCCCC....",
		"..CCCRDDDDDCCC..", // R = Cartuchos vermelhos na bandoleira cruzando o peito
		"..CDDRDDDDDDDC..",
		"..CDDCRDDDDWMM..", // W = coronha, M = cano duplo de metal da espingarda
		"..CDDDCRDDWWMM..",
		"..CDDDDCRDDWMM..",
		"...DDDDDDDDDC...",
		"...DDDD..DDDD...",
		"...BBBB..BBBB...",
		"..BBBBB..BBBBB..",
		"................",
	}
	palette := Palette{
		'.': Trans,
		'C': color.RGBA{45, 50, 58, 255},   // Boina / jaqueta cinza escuro
		'D': color.RGBA{28, 30, 36, 255},   // Jaleco pesado preto-carvão
		'S': SkinLight,
		'K': Black,
		'R': color.RGBA{220, 40, 40, 255},  // Cartuchos calibre 12 vermelhos
		'W': Leather,                       // Coronha de madeira
		'M': MetalLight,                    // Cano duplo da espingarda
		'B': LeatherDk,                     // Botas pesadas
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetSheriffSprite gera o sprite do Xerife (chapéu de cowboy aba larga, lenço vermelho, poncho e revólver prateado).
func GetSheriffSprite() *ebiten.Image {
	pattern := []string{
		".....BBBBBB.....",
		"....BBBBBBBB....",
		".BBBBBBBBBBBBBB.", // Aba larga do chapéu de cowboy
		"...BBSSSSSSBB...",
		"...BDSKSSKSDB...",
		"...BDSSSSSSDB...",
		"....RRRRRRRR....", // Lenço vermelho no pescoço
		"...PPPRRRRPPP...", // P = poncho marrom claro/areia
		"..PPPPPPPPPPP...",
		"..PPPPPGPPPP.MM.", // G = estrela de xerife dourada, M = revólver prateado
		"..PPPPPPPPPP.MW.", // W = coronha do revólver
		"...PPPPPPPP.....",
		"....LL..LL......",
		"....LL..LL......",
		"...BBB..BBB.....",
		"................",
	}
	palette := Palette{
		'.': Trans,
		'B': color.RGBA{95, 55, 28, 255},   // Chapéu e botas de couro
		'D': SkinDark,
		'S': SkinLight,
		'K': Black,
		'R': color.RGBA{210, 35, 35, 255},  // Lenço vermelho
		'P': color.RGBA{175, 130, 75, 255}, // Poncho do velho oeste
		'G': Gold,                          // Distintivo de xerife
		'M': MetalLight,                    // Revólver prateado
		'W': Leather,
		'L': color.RGBA{55, 65, 90, 255},   // Calça jeans/índigo escura
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetBatSprite gera o sprite do Morcego Carmesim da emboscada (asas animadas).
func GetBatSprite(wingUp bool) *ebiten.Image {
	var pattern []string
	if wingUp {
		pattern = []string{
			".R............R.",
			".RR..........RR.",
			".RRL........LRR.",
			"..RRL..KK..LRR..",
			"..RRLLKFFKLLRR..", // F = olhos vermelhos/amarelos
			"...RRLLLLLLRR...",
			"....RLLWWLLR....", // W = presas
			".....RRRRRR.....",
			"......RRRR......",
			".......RR.......",
			"................",
			"................",
		}
	} else {
		pattern = []string{
			"................",
			"................",
			"......KKKK......",
			".....KKFFKK.....",
			"....KKLLLLKK....",
			"...RLLLLLLLLR...",
			"..RRL.LWWL.LRR..",
			".RRR..RRRR..RRR.",
			".RR....RR....RR.",
			".R............R.",
			"................",
			"................",
		}
	}

	palette := Palette{
		'.': Trans,
		'R': color.RGBA{140, 15, 30, 255},
		'L': color.RGBA{210, 35, 50, 255},
		'K': Black,
		'F': color.RGBA{255, 220, 60, 255},
		'W': White,
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetBossSprite gera o sprite colossal 32x32 do Guardião / Senhor da Masmorra.
func GetBossSprite() *ebiten.Image {
	pattern := []string{
		"....GG....................GG....",
		"...GGGG..................GGGG...",
		"...GGDDGG..............GGDDGG...",
		"....GGDDDGG..........GGDDDGG....",
		".....GGDDDDDDDDDDDDDDDDDDGG.....",
		"......DDDDDDDDDDDDDDDDDDDD......",
		".....DDDMMMMMMMMMMMMMMMMDDD.....",
		"....DDMMLLLLLLLLLLLLLLLLMMDD....",
		"....DDMLLKKKKKFFFFKKKKKLLMDD....", // Olhos flamejantes
		"....DDMLLKKKKFFFFFFKKKKLLMDD....",
		"....DDMMLLLLLLLLLLLLLLLLMMDD....",
		".....DDDDDDDDDDDDDDDDDDDDDD.....",
		"...DDDDMMDDDDDDDDDDDDDDMMDDDD...",
		"..DDDMMMMDDDDDDDDDDDDDDMMMMDDD..",
		".DDMMMMMMDDDDGGOOGGDDDDMMMMMMDD.", // Núcleo rúnico dourado/ciano no peitoral
		".DDMMLLLLDDDGOOOOOOGDDDLLLLMMDD.",
		".DDMMLLLLDDDGOOOOOOGDDDLLLLMMDD.",
		".DDMMMMMMDDDDGGOOGGDDDDMMMMMMDD.",
		"..DDDMMMMDDDDDDDDDDDDDDMMMMDDD..",
		"...DDDDDDDDDDDDDDDDDDDDDDDDDD...",
		"....DDDDDDDDDDDDDDDDDDDDDDDD....",
		"....DDDMMMMDDDDDDDDDDMMMMDDD....",
		"....DDDMMMMDDDDDDDDDDMMMMDDD....",
		"....DDDDDDDDDDDDDDDDDDDDDDDD....",
		".....DDDDDDDDDD..DDDDDDDDDD.....",
		".....DDDDDDDDD....DDDDDDDDD.....",
		".....DDMMMMDDD....DDDMMMMDD.....",
		".....DDMMMMDDD....DDDMMMMDD.....",
		"....DDDDDDDDDD....DDDDDDDDDD....",
		"....DDDDDDDDDD....DDDDDDDDDD....",
		"...DDDDDDDDDDD....DDDDDDDDDDD...",
		"................................",
	}
	palette := Palette{
		'.': Trans,
		'G': Gold,
		'D': color.RGBA{35, 28, 45, 255},   // Armadura abissal escura
		'M': color.RGBA{75, 65, 95, 255},   // Placas de aço sombrio
		'L': color.RGBA{125, 115, 155, 255},
		'K': Black,
		'F': color.RGBA{255, 60, 40, 255},  // Brilho demoníaco nos olhos
		'O': color.RGBA{80, 230, 255, 255}, // Núcleo rúnico
	}
	return CreateImageFromPattern(pattern, palette)
}
