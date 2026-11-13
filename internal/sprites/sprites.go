package sprites

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Palette mapeia um caractere a uma cor RGBA.
type Palette map[rune]color.RGBA

// CreateImageFromPattern converte um padrão de strings ASCII em um *ebiten.Image nítido em pixel art.
func CreateImageFromPattern(pattern []string, palette Palette) *ebiten.Image {
	height := len(pattern)
	if height == 0 {
		return ebiten.NewImage(1, 1)
	}
	width := len(pattern[0])

	img := ebiten.NewImage(width, height)
	for y, row := range pattern {
		for x, ch := range row {
			if clr, ok := palette[ch]; ok && clr.A > 0 {
				img.Set(x, y, clr)
			}
		}
	}
	return img
}

// Cores comuns para pixel art
var (
	Trans = color.RGBA{0, 0, 0, 0}
	Black = color.RGBA{15, 15, 20, 255}
	White = color.RGBA{250, 250, 255, 255}

	// Metais e Armaduras
	MetalDark  = color.RGBA{70, 75, 90, 255}
	MetalMid   = color.RGBA{140, 150, 170, 255}
	MetalLight = color.RGBA{200, 215, 235, 255}
	Gold       = color.RGBA{245, 195, 45, 255}
	GoldDark   = color.RGBA{170, 125, 20, 255}

	// Pele e Couro
	SkinLight = color.RGBA{245, 200, 160, 255}
	SkinDark  = color.RGBA{195, 145, 105, 255}
	Leather   = color.RGBA{115, 65, 35, 255}
	LeatherDk = color.RGBA{70, 35, 15, 255}

	// Cores temáticas
	RedLight    = color.RGBA{240, 60, 60, 255}
	RedDark     = color.RGBA{160, 25, 25, 255}
	GreenLight  = color.RGBA{60, 200, 90, 255}
	GreenMid    = color.RGBA{35, 140, 60, 255}
	GreenDark   = color.RGBA{20, 75, 30, 255}
	BlueLight   = color.RGBA{90, 180, 255, 255}
	BlueMid     = color.RGBA{40, 110, 230, 255}
	BlueDark    = color.RGBA{20, 50, 140, 255}
	PurpleLight = color.RGBA{210, 130, 255, 255}
	PurpleMid   = color.RGBA{145, 65, 215, 255}
	PurpleDark  = color.RGBA{80, 25, 130, 255}
)

// Geradores de Sprites dos Jogadores

// GetWarriorSprite gera o sprite do Guerreiro (elmo com fenda, peitoral metálico e espada larga).
func GetWarriorSprite() *ebiten.Image {
	pattern := []string{
		"....MMMMMMMM....",
		"...MMLLLLLLMM...",
		"...MLKKKKKKLM...", // Fenda escura de visão do elmo
		"...MLKKGGKKLM...", // Fenda com brilho dourado
		"...MMDDDDDDMM...",
		"....MMMMMMMM....",
		"..G.MMLLLLMM.W..", // W = lâmina de espada larga
		"..S.MLMMMMMM.W..",
		"..G.MMMMMMMM.W..",
		"..G.MLDDDDLM.W..",
		"....MLLLLLLM.W..",
		"....LLDCCDLL.W..", // C = cinto
		"....DDDDDLDD.G..", // G = guarda da espada
		"....DD....DD.H..", // H = empunhadura
		"...MMD....DMM...",
		"...MMM....MMM...",
	}

	palette := Palette{
		'.': Trans,
		'M': MetalMid,
		'L': MetalLight,
		'D': MetalDark,
		'K': Black,
		'G': Gold,
		'S': SkinLight,
		'C': Leather,
		'W': color.RGBA{230, 240, 255, 255}, // Lâmina prateada reluzente
		'H': LeatherDk,
	}

	return CreateImageFromPattern(pattern, palette)
}

// GetArcherSprite gera o sprite do Arqueiro (capuz verde, arco de madeira curvado com corda).
func GetArcherSprite(firing bool) *ebiten.Image {
	// Se firing == true, arco tensionado puxando corda
	var pattern []string
	if firing {
		pattern = []string{
			"....GGGGGG......",
			"...GLLLLLLG.....",
			"...GDDDDDDG.....",
			"...GLSKKSLL.B...",
			"...GLSSSSLL..B..",
			"....GGGGGG...C..",
			"...GGLLLLGG.B.C.",
			"..S.GDDDDGG.B..C", // B = madeira do arco curvado, C = corda tensionada
			"..S.GLLLLGG..B.C",
			"..L.GDDDDGG...C.",
			"....GDDDDGG..B..",
			"....LLCCLL......",
			"....DDLLDD......",
			"....DD..DD......",
			"...LLD..DLL.....",
			"...LLD..DLL.....",
		}
	} else {
		pattern = []string{
			"....GGGGGG......",
			"...GLLLLLLG.....",
			"...GDDDDDDG.....",
			"...GLSKKSLL.B...",
			"...GLSSSSLL.BC..",
			"....GGGGGG..BC..",
			"...GGLLLLGG.BC..",
			"..S.GDDDDGG.BC..",
			"..S.GLLLLGG.BC..",
			"..L.GDDDDGG.BC..",
			"....GDDDDGG.B...",
			"....LLCCLL......",
			"....DDLLDD......",
			"....DD..DD......",
			"...LLD..DLL.....",
			"...LLD..DLL.....",
		}
	}

	palette := Palette{
		'.': Trans,
		'G': GreenMid,
		'L': GreenLight,
		'D': GreenDark,
		'S': SkinLight,
		'K': Black,
		'C': White, // Corda do arco
		'B': Leather, // Madeira do arco
	}

	return CreateImageFromPattern(pattern, palette)
}

// GetMageSprite gera o sprite do Mago (chapéu cônico pontudo, manto e cajado com cristal).
func GetMageSprite() *ebiten.Image {
	pattern := []string{
		".......PP.......", // Ponta do chapéu
		"......PPPP......",
		".....PPLLPP.OO..", // O = Cristal mágico do topo do cajado
		"....PPLLLLPP.O..",
		"...PPPPPPPPPP...",
		"...PDDDDDDDDP...",
		"...PDSSKKSSDP.B.", // B = haste de madeira do cajado
		"...PDSSSSSSDP.B.",
		"...PPPPPPPPPP.B.",
		"..PPDDLLLLDDP.B.",
		"..PPDDLLLLDDP.B.",
		"...PPDDDDDDDP.B.",
		"...PPDDLLDDPP.B.",
		"...PPDDLLDDPP...",
		"....DDDDDDDD....",
		"....LLDDDDLL....",
	}

	palette := Palette{
		'.': Trans,
		'P': PurpleDark,
		'L': PurpleLight,
		'D': PurpleMid,
		'S': SkinLight,
		'K': Black,
		'O': BlueLight, // Cristal brilhante
		'B': Leather,   // Cajado
	}

	return CreateImageFromPattern(pattern, palette)
}

// Sprites dos Inimigos

// GetZombieSprite gera o sprite do Inimigo Comum (pele verde apodrecida, olhos vazios, trapos).
func GetZombieSprite() *ebiten.Image {
	pattern := []string{
		"....GGGGGG......",
		"...GLYYYYLG.....",
		"...GYYKKYYG.....", // Olhos escurecidos amarelados
		"...GLYDDYLG.....",
		"...GGGDGGGG.....",
		"....GGGGGG......",
		"..SS.GLLLLG.SS..", // Braços estendidos pra frente
		"..SS.GLDDLG.SS..",
		"..SS.GDDDDG.SS..",
		"....GLLDDLG.....",
		"....GDDDDDG.....",
		"....LLCCLL......",
		"....DDRRDD......",
		"....DD..DD......",
		"...DLD..DLD.....",
		"...DDD..DDD.....",
	}

	palette := Palette{
		'.': Trans,
		'G': color.RGBA{50, 110, 60, 255},
		'L': color.RGBA{80, 155, 90, 255},
		'D': color.RGBA{30, 75, 40, 255},
		'Y': color.RGBA{220, 220, 60, 255},
		'K': Black,
		'S': color.RGBA{95, 175, 105, 255},
		'C': LeatherDk,
		'R': color.RGBA{85, 45, 30, 255}, // Trapos rasgados
	}

	return CreateImageFromPattern(pattern, palette)
}

// GetHellhoundSprite gera o sprite do Corredor (cão infernal veloz, olhos em brasa, presas).
func GetHellhoundSprite() *ebiten.Image {
	pattern := []string{
		"................",
		".....RRR........",
		"....RLLLR.......",
		"...RLLFFF.......", // F = Fogo / Olho incandescente
		"..RLLLLFFF......",
		"..RLLLLKWWF.....", // W = presas brancas
		"..RDDDDDDDD.....",
		"...RDDDDRDD.....",
		"...RDD..RDD.....",
		"..RRDD..RRDD....",
		"..RRDD..RRDD....",
		"..R..D..R..D....",
		"................",
		"................",
		"................",
		"................",
	}

	palette := Palette{
		'.': Trans,
		'R': RedDark,
		'L': color.RGBA{200, 40, 40, 255},
		'D': color.RGBA{100, 15, 15, 255},
		'F': color.RGBA{255, 210, 30, 255}, // Olhos de fogo
		'K': Black,
		'W': White,
	}

	return CreateImageFromPattern(pattern, palette)
}

// GetGolemSprite gera o sprite do Tanque (golem colossal de rocha com runas místicas, 24x24).
func GetGolemSprite() *ebiten.Image {
	pattern := []string{
		"......DDDDDDDDDDDD......",
		".....DDMMMMMMMMMMDD.....",
		"....DDMLLLLLLLLMDDD....",
		"....DDMKKOOOOOOKMDDD....", // O = Olhos / Núcleo brilhante
		"....DDMMLLLLLLMMDDDD....",
		"...DDDDDDDDDDDDDDDDDD...",
		"..DDMMMLLLLLLLLMMMDDDD..",
		"..DDMLLLLRRRRLLLLMDDDD..", // R = Runa mágica no peitoral
		".DDMMLLLRROORRLLMMMDDDD.",
		".DDMMLLLRROORRLLMMMDDDD.",
		".DDMMLLLLRRRRLLLLMMDDDD.",
		".DDMMMLLLLLLLLMMMMDDDDD.",
		"..DDMMMMMMMMMMMMDDDDDD..",
		"..DDDDMMMMMMMMDDDDDDDD..",
		"...DDDDDDDDDDDDDDDDDD...",
		"....DDDDDD....DDDDDD....",
		"....DDMMDD....DDMMDD....",
		"....DDLLDD....DDLLDD....",
		"....DDLLDD....DDLLDD....",
		"....DDMMDD....DDMMDD....",
		"...DDDDDDDD..DDDDDDDD...",
		"...DDMMMMDD..DDMMMMDD...",
		"...DDDDDDDD..DDDDDDDD...",
		"........................",
	}

	palette := Palette{
		'.': Trans,
		'D': color.RGBA{50, 50, 65, 255},
		'M': color.RGBA{85, 90, 110, 255},
		'L': color.RGBA{130, 135, 160, 255},
		'K': Black,
		'O': color.RGBA{140, 240, 255, 255}, // Brilho rúnico ciano
		'R': color.RGBA{70, 180, 230, 255},
	}

	return CreateImageFromPattern(pattern, palette)
}

// GetShooterSprite gera o sprite do Inimigo Atirador/Bombardeiro (necromante / criatura cuspidora).
func GetShooterSprite() *ebiten.Image {
	pattern := []string{
		"....KKKKKK......",
		"...KDDDDDDK.....",
		"...KDVVVVDK.....", // V = olhos tóxicos venenosos
		"...KDKKKKDK.....",
		"...KDDDDDDK.....",
		"....KKKKKK......",
		"...KKDDDDKK.....",
		"..KVDDDDDDKV....", // V = orbes nas mãos
		"..KVDLLLLDKV....",
		"...KDDDDDDDK....",
		"...KDDDDDDDK....",
		"...KDD...DDK....",
		"...KDD...DDK....",
		"...KDD...DDK....",
		"..KKDD...DDKK...",
		"..KKKK...KKKK...",
	}

	palette := Palette{
		'.': Trans,
		'K': color.RGBA{25, 20, 35, 255},
		'D': color.RGBA{85, 40, 100, 255},
		'L': color.RGBA{135, 65, 155, 255},
		'V': color.RGBA{80, 245, 110, 255}, // Verde ácido brilhante
	}

	return CreateImageFromPattern(pattern, palette)
}

// Sprites de Itens / Drops (12x12)

// GetHealthPotionSprite gera o frasco de cura em pixel art detalhado (com cruz/coração em relevo e brilhos de vidro).
func GetHealthPotionSprite() *ebiten.Image {
	pattern := []string{
		".....LLKKLL.....", // Rolha de cortiça com anel de couro
		".....LLKKLL.....",
		"......WWWW......", // Gargalo de vidro transparente
		"....WWWWWWWW....",
		"...WWLLLLLLWW...", // Ombro do frasco com reflexos
		"..WWLRRRRRRLWW..",
		"..WLRRRCCRRRLW..", // CC = Cruz sagrada em relevo
		".WWLRRCCCCRRLWW.",
		".WLLRRCCCCRRLLW.",
		".WLLRRRCCRRRLLW.",
		".WWLRRRRRRRRLWW.",
		"..WWLRRRRRRLWW..",
		"..WWLLRRRRLLWW..",
		"...WWWWWWWWWW...", // Fundo do frasco reforçado
		"................",
		"................",
	}
	palette := Palette{
		'.': Trans,
		'K': LeatherDk,
		'L': Leather,                        // Cortiça da rolha
		'W': color.RGBA{225, 240, 255, 230}, // Vidro cristalino com reflexo especular
		'R': color.RGBA{220, 25, 35, 255},   // Elixir carmesim profundo
		'C': White,                          // Cruz/emblema sagrado em relevo
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetManaPotionSprite gera o frasco de mana detalhado (com runa mágica cintilante e vidro refratário).
func GetManaPotionSprite() *ebiten.Image {
	pattern := []string{
		".....LLKKLL.....",
		".....LLKKLL.....",
		"......WWWW......",
		"....WWWWWWWW....",
		"...WWLLLLLLWW...",
		"..WWLBBBBBBLLW..",
		"..WLBBBCCBBBBLW.", // CC = Runa mágica de mana
		".WWLBBCCCCBBLLW.",
		".WLLBBCCCCBBLLW.",
		".WLLBBBCCBBBLLW.",
		".WWLBBBBBBBBLLW.",
		"..WWLBBBBBBLLW..",
		"..WWLLBBBBLLWW..",
		"...WWWWWWWWWW...",
		"................",
		"................",
	}
	palette := Palette{
		'.': Trans,
		'K': MetalDark,
		'L': MetalLight,
		'W': color.RGBA{210, 235, 255, 230},
		'B': color.RGBA{25, 100, 235, 255},  // Elixir de mana safira
		'C': color.RGBA{130, 240, 255, 255}, // Emblema místico radiante
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetAgilityStoneSprite gera a pedra de agilidade.
func GetAgilityStoneSprite() *ebiten.Image {
	pattern := []string{
		"....YYYY....",
		"...YLLLLY...",
		"..YLLLLLLYY.",
		".YLLLLLLLLLY",
		".YLLLLLLLLLY",
		".YLLGGGGLLLY",
		".YLLGGGGLLLY",
		"..YLLGGLLYY.",
		"...YLLLLYY..",
		"....YYYY....",
		"............",
		"............",
	}
	palette := Palette{
		'.': Trans,
		'Y': color.RGBA{255, 215, 0, 255},
		'L': color.RGBA{255, 245, 140, 255},
		'G': color.RGBA{200, 150, 0, 255},
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetBootsSprite gera as botas de ímpeto com asinhas.
func GetBootsSprite() *ebiten.Image {
	pattern := []string{
		"W..........W",
		"WW........WW",
		".WW.BBBB.WW.",
		"..W.BLLBB.W.",
		"...BBLLLLB..",
		"...BBLLLLB..",
		"...BBLLLLBBB",
		"...BBLLLLLLB",
		"...BBBBBBBBB",
		"............",
		"............",
		"............",
	}
	palette := Palette{
		'.': Trans,
		'W': color.RGBA{240, 245, 255, 255}, // Asas
		'B': color.RGBA{50, 120, 200, 255},  // Botas azuis de velocidade
		'L': color.RGBA{100, 180, 255, 255},
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetTomeSprite gera o tomo de maestria.
func GetTomeSprite() *ebiten.Image {
	pattern := []string{
		"...PPPPPP...",
		"..PPDDDDPP..",
		".PPDGGGGDPP.", // Símbolo dourado de poder
		".PPDGLLGDPP.",
		".PPDGLLGDPP.",
		".PPDGGGGDPP.",
		".PPDDDDDDPP.",
		".WWWWWWWWWW.", // Páginas envelhecidas
		"..WWWWWWWW..",
		"............",
		"............",
		"............",
	}
	palette := Palette{
		'.': Trans,
		'P': PurpleDark,
		'D': PurpleMid,
		'G': Gold,
		'L': color.RGBA{255, 240, 160, 255},
		'W': color.RGBA{220, 210, 180, 255},
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetChaliceSprite gera o cálice amaldiçoado.
func GetChaliceSprite() *ebiten.Image {
	pattern := []string{
		".GGGGGGGGGG.",
		".GKRKKKKRKGG",
		"..GRRRRRRGG.", // Líquido carmesim maldito
		"...GRRRRGG..",
		"....GGGGG...",
		".....GGG....",
		".....GGG....",
		"....GGGGG...",
		"...GGGGGGG..",
		"............",
		"............",
		"............",
	}
	palette := Palette{
		'.': Trans,
		'G': color.RGBA{80, 75, 70, 255}, // Ouro negro envelhecido
		'K': Black,
		'R': color.RGBA{200, 20, 40, 255},
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetPoisonSprite gera o frasco de veneno.
func GetPoisonSprite() *ebiten.Image {
	pattern := []string{
		"....CCCC....",
		".....WW.....",
		"....WWWW....",
		"...WLLLLW...",
		"..WLLGGLLW..",
		"..WLGKKGLW..", // Caveira verde de veneno
		".WLLGKKGLLW.",
		".WLLGGGGLLW.",
		".WLLGGGGLLW.",
		"..WGGGGGGW..",
		"...WWWWWW...",
		"............",
	}
	palette := Palette{
		'.': Trans,
		'C': Leather,
		'W': color.RGBA{200, 240, 200, 220},
		'L': color.RGBA{140, 255, 140, 255},
		'G': color.RGBA{40, 180, 50, 255},
		'K': Black,
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetFatigueMistSprite gera a névoa de fadiga.
func GetFatigueMistSprite() *ebiten.Image {
	pattern := []string{
		"....DDDD....",
		"...DLLLLD...",
		"..DLLKKLLD..",
		".DLLK..KLLD.",
		".DLLK..KLLD.",
		"..DLLKKLLD..",
		"...DLLLLD...",
		"....DDDD....",
		"............",
		"............",
		"............",
		"............",
	}
	palette := Palette{
		'.': Trans,
		'D': color.RGBA{90, 70, 110, 255},
		'L': color.RGBA{160, 130, 190, 255},
		'K': color.RGBA{45, 35, 55, 255},
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetChaosSeedSprite gera a semente do caos.
func GetChaosSeedSprite() *ebiten.Image {
	pattern := []string{
		".....RR.....",
		"....RLLR....",
		"...RLFFLR...",
		"..RLFFFFLR..",
		"..RLFFFFLR..",
		"...RLFFLR...",
		"....RLLR....",
		".....RR.....",
		"............",
		"............",
		"............",
		"............",
	}
	palette := Palette{
		'.': Trans,
		'R': Black,
		'L': color.RGBA{220, 20, 60, 255},
		'F': color.RGBA{255, 120, 20, 255}, // Centro flamejante
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetSoulFragmentSprite gera o cristal de Essência / Fragmento de Alma (moeda meta).
func GetSoulFragmentSprite() *ebiten.Image {
	pattern := []string{
		".....SS.....",
		"....SAAS....",
		"...SAAAAS...",
		"..SAAWWAAS..",
		".SAAWWWWAAAS",
		"SAAWWWWWWAAAS",
		"SAAWWWWWWAAAS",
		".SAAWWWWAAAS",
		"..SAAWWAAS..",
		"...SAAAAS...",
		"....SAAS....",
		".....SS.....",
	}
	palette := Palette{
		'.': Trans,
		'S': color.RGBA{80, 40, 160, 255},   // Borda violeta etérea
		'A': color.RGBA{130, 90, 255, 255},  // Azul arcano brilhante
		'W': color.RGBA{235, 230, 255, 255}, // Núcleo de luz pura
	}
	return CreateImageFromPattern(pattern, palette)
}
