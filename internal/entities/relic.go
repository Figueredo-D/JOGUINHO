package entities

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"joguinho/internal/sprites"
)

// RelicID identifica cada relíquia passiva acumulável.
type RelicID string

const (
	RelicFocusLens   RelicID = "focus_lens"
	RelicThornyCloak RelicID = "thorny_cloak"
	RelicStaticSpark RelicID = "static_spark"
	RelicBloodPact   RelicID = "blood_pact"
	RelicSoulMagnet  RelicID = "soul_magnet"
	RelicIronHeart   RelicID = "iron_heart"
)

// RelicDef contém os metadados de exibição de uma relíquia passiva.
type RelicDef struct {
	ID          RelicID
	Name        string
	Description string
	EffectDesc  string
	Color       color.RGBA
}

// GetSprite retorna o sprite em pixel art correspondente à relíquia.
func (r RelicDef) GetSprite() *ebiten.Image {
	return GetRelicSprite(r.ID)
}

// GetRelicSprite retorna o sprite de pixel art para um dado ID de relíquia.
func GetRelicSprite(id RelicID) *ebiten.Image {
	switch id {
	case RelicFocusLens:
		return sprites.GetFocusLensRelicSprite()
	case RelicThornyCloak:
		return sprites.GetThornyCloakRelicSprite()
	case RelicStaticSpark:
		return sprites.GetStaticSparkRelicSprite()
	case RelicBloodPact:
		return sprites.GetBloodPactRelicSprite()
	case RelicSoulMagnet:
		return sprites.GetSoulMagnetRelicSprite()
	case RelicIronHeart:
		return sprites.GetIronHeartRelicSprite()
	default:
		return nil
	}
}

// AllRelics define o catálogo completo de relíquias da partida.
var AllRelics = []RelicDef{
	{
		ID:          RelicFocusLens,
		Name:        "Lente de Foco",
		Description: "Cristal optico que concentra a energia dos projeteis.",
		EffectDesc:  "+35% chance dos projeteis perfurarem o 1o alvo.",
		Color:       color.RGBA{R: 100, G: 210, B: 255, A: 255},
	},
	{
		ID:          RelicThornyCloak,
		Name:        "Manto Espinhoso",
		Description: "Tecido encantado com espinhos de aco retorcido.",
		EffectDesc:  "Reflete 30% da sua Vida Max como dano ao ser atingido.",
		Color:       color.RGBA{R: 90, G: 200, B: 100, A: 255},
	},
	{
		ID:          RelicStaticSpark,
		Name:        "Faisca Eletrostatica",
		Description: "Orbe carregado que libera raios em cadeia.",
		EffectDesc:  "A cada 5 hits, solta um raio em ate 3 inimigos proximos.",
		Color:       color.RGBA{R: 255, G: 230, B: 60, A: 255},
	},
	{
		ID:          RelicBloodPact,
		Name:        "Pacto de Sangue",
		Description: "Poder proibido selado em tinta carmesim.",
		EffectDesc:  "+25% de Dano Global, mas curas recebidas caem pela metade.",
		Color:       color.RGBA{R: 225, G: 40, B: 60, A: 255},
	},
	{
		ID:          RelicSoulMagnet,
		Name:        "Ima de Almas",
		Description: "Reliquia magnetica que atrai itens e essencias.",
		EffectDesc:  "+90px de raio de atracao para drops e fragmentos.",
		Color:       color.RGBA{R: 180, G: 110, B: 255, A: 255},
	},
	{
		ID:          RelicIronHeart,
		Name:        "Coracao de Ferro",
		Description: "Nucleo forjado que fortalece o vigor vital.",
		EffectDesc:  "+30 HP Maximo imediato e regeneracao lenta continua.",
		Color:       color.RGBA{R: 240, G: 150, B: 60, A: 255},
	},
}

// PickRandomRelics sorteia 3 relíquias distintas para a janela de escolha.
func PickRandomRelics(count int) []RelicDef {
	pool := make([]RelicDef, len(AllRelics))
	copy(pool, AllRelics)
	rand.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})
	if count > len(pool) {
		count = len(pool)
	}
	return pool[:count]
}
