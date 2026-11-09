package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
	"joguinho/internal/sprites"
)

// DropType define os tipos de itens dropáveis.
type DropType int

const (
	// Moeda Meta Roguelite
	DropSoulFragment DropType = iota

	// Benignos
	DropHealthPotion
	DropManaPotion
	DropAgilityStone
	DropSpeedBoots
	DropMasteryTome

	// Malignos
	DropCursedChalice
	DropPoisonVial
	DropFatigueMist
	DropChaosSeed
)

// DropInfo contém metadados de exibição e tipo do drop.
type DropInfo struct {
	Name      string
	IsDebuff  bool
	Color     color.RGBA
	SpriteImg *ebiten.Image
}

var dropInfoMap map[DropType]DropInfo

func initDropInfo() {
	if dropInfoMap != nil {
		return
	}
	dropInfoMap = map[DropType]DropInfo{
		DropSoulFragment: {
			Name:      "Fragmento de Alma",
			IsDebuff:  false,
			Color:     color.RGBA{R: 160, G: 110, B: 255, A: 255},
			SpriteImg: sprites.GetSoulFragmentSprite(),
		},
		DropHealthPotion: {
			Name:      "Frasco de Cura (+HP)",
			IsDebuff:  false,
			Color:     color.RGBA{R: 240, G: 60, B: 60, A: 255},
			SpriteImg: sprites.GetHealthPotionSprite(),
		},
		DropManaPotion: {
			Name:      "Elixir de Mana (+MP)",
			IsDebuff:  false,
			Color:     color.RGBA{R: 60, G: 140, B: 255, A: 255},
			SpriteImg: sprites.GetManaPotionSprite(),
		},
		DropAgilityStone: {
			Name:      "Pedra da Agilidade (+Vel. Ataque)",
			IsDebuff:  false,
			Color:     color.RGBA{R: 255, G: 220, B: 50, A: 255},
			SpriteImg: sprites.GetAgilityStoneSprite(),
		},
		DropSpeedBoots: {
			Name:      "Botas de Ímpeto (+Vel. Movimento)",
			IsDebuff:  false,
			Color:     color.RGBA{R: 80, G: 200, B: 255, A: 255},
			SpriteImg: sprites.GetBootsSprite(),
		},
		DropMasteryTome: {
			Name:      "Tomo de Maestria (+Max HP/MP)",
			IsDebuff:  false,
			Color:     color.RGBA{R: 210, G: 130, B: 255, A: 255},
			SpriteImg: sprites.GetTomeSprite(),
		},

		// Malignos
		DropCursedChalice: {
			Name:      "Cálice Amaldiçoado (+Dano, -Vida/s)",
			IsDebuff:  true,
			Color:     color.RGBA{R: 220, G: 30, B: 50, A: 255},
			SpriteImg: sprites.GetChaliceSprite(),
		},
		DropPoisonVial: {
			Name:      "Frasco de Veneno (Lentidão e Veneno)",
			IsDebuff:  true,
			Color:     color.RGBA{R: 50, G: 210, B: 60, A: 255},
			SpriteImg: sprites.GetPoisonSprite(),
		},
		DropFatigueMist: {
			Name:      "Névoa de Fadiga (Ataques Lentos)",
			IsDebuff:  true,
			Color:     color.RGBA{R: 160, G: 90, B: 200, A: 255},
			SpriteImg: sprites.GetFatigueMistSprite(),
		},
		DropChaosSeed: {
			Name:      "Semente do Caos (Invocação Súbita!)",
			IsDebuff:  true,
			Color:     color.RGBA{R: 255, G: 70, B: 20, A: 255},
			SpriteImg: sprites.GetChaosSeedSprite(),
		},
	}
}

// Drop representa um item deixado no chão por um inimigo.
type Drop struct {
	Type       DropType
	X, Y       float64
	Lifetime   int
	MaxLife    int
	Active     bool
	Frame      int
	SoulAmount int // Quantidade de almas se for DropSoulFragment
}

// NewDrop cria uma nova instância de Drop no chão.
func NewDrop(dType DropType, x, y float64) *Drop {
	initDropInfo()
	return &Drop{
		Type:       dType,
		X:          x,
		Y:          y,
		Lifetime:   720, // 12 segundos a 60fps
		MaxLife:    720,
		Active:     true,
		SoulAmount: 1,
	}
}

// NewSoulDrop cria um drop específico de fragmento de alma com quantia definida.
func NewSoulDrop(x, y float64, amount int) *Drop {
	initDropInfo()
	return &Drop{
		Type:       DropSoulFragment,
		X:          x,
		Y:          y,
		Lifetime:   900, // 15 segundos
		MaxLife:    900,
		Active:     true,
		SoulAmount: amount,
	}
}

// Update avança o contador de vida do drop no chão.
func (d *Drop) Update() {
	if !d.Active {
		return
	}
	d.Frame++
	d.Lifetime--
	if d.Lifetime <= 0 {
		d.Active = false
	}
}

// Bounds retorna o AABB de coleta do item.
func (d *Drop) Bounds() collision.Rect {
	return collision.NewRect(d.X-10, d.Y-10, 20, 20)
}

// Info retorna metadados informativos do drop.
func (d *Drop) Info() DropInfo {
	initDropInfo()
	return dropInfoMap[d.Type]
}

// Draw renderiza o item no chão com animação de flutuação e cintilância.
func (d *Drop) Draw(dst *ebiten.Image) {
	if !d.Active {
		return
	}

	// Efeito de piscar quando está perto de expirar (últimos 3 segundos)
	if d.Lifetime < 180 && (d.Lifetime/6)%2 == 0 {
		return
	}

	info := d.Info()

	// Flutuação vertical senoidal suave
	floatOffset := math.Sin(float64(d.Frame)*0.08) * 4.0
	drawY := d.Y + floatOffset

	// Aura cintilante no chão abaixo do item
	auraRadius := float32(8.0 + math.Sin(float64(d.Frame)*0.12)*2.0)
	auraColor := info.Color
	auraColor.A = 80
	vector.DrawFilledCircle(dst, float32(d.X), float32(d.Y+6), auraRadius, auraColor, true)

	// Desenho do sprite do item (12x12 escalonado x1.5 para ficar bem nítido)
	if info.SpriteImg != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1.4, 1.4)
		// Centraliza o sprite
		w := float64(info.SpriteImg.Bounds().Dx()) * 1.4
		h := float64(info.SpriteImg.Bounds().Dy()) * 1.4
		op.GeoM.Translate(d.X-w/2, drawY-h/2)
		dst.DrawImage(info.SpriteImg, op)
	}
}
