package sprites

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// GetStaticSparkRelicSprite gera o ícone em pixel art da Faísca Eletrostática (relâmpago sobre escudo místico).
func GetStaticSparkRelicSprite() *ebiten.Image {
	pattern := []string{
		"....GGGGGGGG....",
		"..GGSSBBBBSSGG..",
		".GSBBBBBBBBBBBSG",
		".GBBBBBWWYBBBBBG",
		"GBBBBBWWYYBBBBBG",
		"GBBBBWWYYBBBBBBG",
		"GBBBWWYYBBBBBBBG",
		"GBBWWYYYYYYYBBBG",
		"GBBBBYYYWWYYBBBG",
		"GBBBBBBYWWYBBBBG",
		"GBBBBBWWYYBBBBBG",
		"GBBBBWWYYBBBBBBG",
		".GBBWWYYBBBBBBG.",
		".GSBBYYBBBBBBSG.",
		"..GGSBBBBBBSSG..",
		"....GGGGGGGG....",
	}
	palette := Palette{
		'.': Trans,
		'G': Gold,                          // Borda ornamental dourada
		'S': color.RGBA{40, 60, 100, 255},  // Chanfro do escudo
		'B': color.RGBA{18, 24, 45, 255},   // Fundo azul abissal
		'W': White,                         // Núcleo elétrico branco
		'Y': color.RGBA{255, 225, 45, 255}, // Raio amarelo incandescente
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetSoulMagnetRelicSprite gera o ícone em pixel art do Ímã de Almas (ferradura com pólos e ondas magnéticas).
func GetSoulMagnetRelicSprite() *ebiten.Image {
	pattern := []string{
		".....MMMMMM.....",
		"...MMLLLLLLMM...",
		"..MLLMMMMMMLLM..",
		".MLMM......MMLLM",
		".MLM........MLLM",
		".MLM........MLLM",
		".RRR........BBBB",
		".RWR........BWB.",
		".RRR........BBBB",
		"................",
		"...PP......PP...",
		"..P..P....P..P..",
		".P....P..P....P.",
		"................",
		"...CC......CC...",
		"................",
	}
	palette := Palette{
		'.': Trans,
		'M': MetalDark,                      // Corpo de metal do ímã
		'L': MetalLight,                     // Brilho na curva
		'R': color.RGBA{225, 40, 50, 255},   // Polo magnético vermelho (Norte)
		'B': color.RGBA{45, 110, 235, 255},  // Polo magnético azul (Sul)
		'W': White,                          // Faixa reflexiva nos polos
		'P': color.RGBA{210, 120, 255, 220}, // Linhas de campo de almas (violeta)
		'C': color.RGBA{120, 240, 255, 200}, // Ondas etéreas ciano
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetFocusLensRelicSprite gera o ícone da Lente de Foco (cristal óptico com aro dourado e refração).
func GetFocusLensRelicSprite() *ebiten.Image {
	pattern := []string{
		"....GGGGGG......",
		"..GGDDDDDDGG....",
		".GDDLLLLLLDDG...",
		".GDLWWWWCCLDG...",
		"GDLWWCCCCCLDDG..",
		"GDLCCCCCCCLLDG..",
		"GDLCCCCCCCLLDG..",
		"GDLCCCCCCCLLDG..",
		"GDLLCCCCCLLDDG..",
		".GDLLLLLLLDDG...",
		"..GGDDDDDDGG....",
		"....GGGGGG.GG...",
		"..........GDDG..",
		"...........GDDG.",
		"............GG..",
		"................",
	}
	palette := Palette{
		'.': Trans,
		'G': Gold,                          // Aro e cabo de latão forjado
		'D': GoldDark,                      // Sombra do metal
		'L': color.RGBA{140, 220, 255, 255},// Borda reflexiva do vidro
		'W': White,                         // Ponto de luz concentrado da lente
		'C': color.RGBA{50, 160, 230, 240}, // Corpo do cristal óptico
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetThornyCloakRelicSprite gera o ícone do Manto Espinhoso (capa com espinhos pontiagudos de aço).
func GetThornyCloakRelicSprite() *ebiten.Image {
	pattern := []string{
		"......SS..SS....",
		".....SLLSSLLS...",
		"....SLLMMMMLLS..",
		"..SSLLGGGGGGLLSS",
		".SLMGGGGGGGGMLS.",
		".SMGGGEEEEEEGMMS",
		"SSMGGEEEEEEEGMSS",
		"S.MGGEEEEEEEGM.S",
		"S.MGGEEEREEEGM.S", // R = broche / gema carmesim
		"..MGGEEEREEEGM..",
		"..MGGEEEEEEEGM..",
		"..MGGEEEEEEEGM..",
		"..MGGGEEEEEGM...",
		"...MMMGEEEGM....",
		".....MMMMMM.....",
		"................",
	}
	palette := Palette{
		'.': Trans,
		'S': MetalLight,                     // Pontas pontiagudas de aço
		'L': MetalMid,                       // Espinhos metálicos
		'M': MetalDark,                      // Base dos cravos
		'G': color.RGBA{25, 90, 45, 255},    // Tecido verde musgo escuro
		'E': color.RGBA{45, 145, 75, 255},   // Tecido da capa verde floresta
		'R': color.RGBA{220, 35, 50, 255},   // Broche de rubi da capa
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetBloodPactRelicSprite gera o ícone do Pacto de Sangue (pergaminho carmesim com selo de sangue e adaga).
func GetBloodPactRelicSprite() *ebiten.Image {
	pattern := []string{
		"...PPPPPPPPPP...",
		"..PDDDDDDDDDDP..",
		".PDDDDDDDDDDDDP.",
		".PDMMDDDDDDDDDP.",
		".PDDMLLDDDDDDDP.", // M/L = Lâmina da adaga cerimonial
		".PDDDDMLLDRRDDP.", // R = Selo de sangue
		".PDDDDDDMDRRDDP.",
		".PDDRRDDDGDDMDP.", // G = Empunhadura de ouro da adaga
		".PDDRRRRDRRMDDP.",
		".PDDRRDDRRRDDDP.",
		".PDDDRRRRRDDDDP.",
		".PDDDDDDDDDDDDP.",
		"..PDDDDDDDDDDP..",
		"...PPPPPPPPPP...",
		"................",
		"................",
	}
	palette := Palette{
		'.': Trans,
		'P': color.RGBA{135, 95, 55, 255},   // Borda do pergaminho antigo
		'D': color.RGBA{215, 185, 135, 255}, // Pergaminho envelhecido
		'M': MetalDark,
		'L': MetalLight,                     // Lâmina de aço
		'G': Gold,                          // Guarda da adaga
		'R': color.RGBA{195, 20, 35, 255},   // Selo e sangue fresco
	}
	return CreateImageFromPattern(pattern, palette)
}

// GetIronHeartRelicSprite gera o ícone do Coração de Ferro (coração mecânico com núcleo engrenado incandescente).
func GetIronHeartRelicSprite() *ebiten.Image {
	pattern := []string{
		"....MM....MM....",
		"..MMLLMMMMLLMM..",
		".MLLDDMMMMDDLLM.",
		".MLDDOOLLFFDDLM.",
		"MLDDOOWWFFFFDDLM", // W = núcleo fundido branco incandescente
		"MLDDOOWWFFFFDDLM", // O/F = fogo e chamas do núcleo
		"MLDDDDOOFFDDDDLM",
		"MLDDDDGGGGDDDDLM", // G = engrenagens de latão
		".MLDDDGGGGDDDLM.",
		".MLDDDDGGDDDDLM.",
		"..MLDDDDDDDDLM..",
		"...MLDDDDDDLM...",
		"....MLDDDDLM....",
		".....MLDDLM.....",
		"......MLMM......",
		".......MM.......",
	}
	palette := Palette{
		'.': Trans,
		'M': MetalDark,
		'L': MetalLight,                     // Placas externas de metal reforçado
		'D': color.RGBA{85, 75, 95, 255},    // Estrutura de aço escuro
		'G': Gold,                          // Engrenagens internas de bronze
		'O': color.RGBA{255, 120, 20, 255},  // Brilho de fogo no núcleo
		'F': color.RGBA{255, 60, 30, 255},   // Fornalha vital
		'W': White,                          // Ponto mais quente
	}
	return CreateImageFromPattern(pattern, palette)
}
