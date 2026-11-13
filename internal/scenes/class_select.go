package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/collision"
	"joguinho/internal/entities"
	"joguinho/internal/sprites"
	"joguinho/internal/ui"
)

type classCardMeta struct {
	class     entities.ClassType
	title     string
	bounds    collision.Rect
	spriteImg *ebiten.Image
	themeClr  color.RGBA
	line1     string
	line2     string
	line3     string
	special   string
}

// ClassSelectScene permite ao jogador escolher entre as 5 classes disponíveis.
type ClassSelectScene struct {
	navigator     SceneNavigator
	selectedClass entities.ClassType
	cards         []classCardMeta
	btnStart      *ui.Button
	btnBack       *ui.Button
}

// NewClassSelectScene inicializa a seleção com as 5 classes organizadas em 2 fileiras.
func NewClassSelectScene(nav SceneNavigator) *ClassSelectScene {
	cards := []classCardMeta{
		// Fileira Superior (3 classes clássicas)
		{
			class:     entities.ClassWarrior,
			title:     "GUERREIRO",
			bounds:    collision.NewRect(45, 72, 220, 195),
			spriteImg: sprites.GetWarriorSprite(),
			themeClr:  color.RGBA{200, 70, 70, 255},
			line1:     "HP: 200 | Vel: Media",
			line2:     "Ataque: Espada Larga (Melee)",
			line3:     "Alto dano e resistencia",
			special:   "Esp: Escudo Refletor",
		},
		{
			class:     entities.ClassArcher,
			title:     "ARQUEIRO",
			bounds:    collision.NewRect(290, 72, 220, 195),
			spriteImg: sprites.GetArcherSprite(false),
			themeClr:  color.RGBA{50, 180, 90, 255},
			line1:     "HP: 75 | Vel: Rapida",
			line2:     "Ataque: Flechas Continuas",
			line3:     "Agil e longo alcance",
			special:   "Esp: Chuva de Flechas",
		},
		{
			class:     entities.ClassMage,
			title:     "MAGO",
			bounds:    collision.NewRect(535, 72, 220, 195),
			spriteImg: sprites.GetMageSprite(),
			themeClr:  color.RGBA{60, 110, 230, 255},
			line1:     "HP: 110 | MP: 100 (Regen)",
			line2:     "Ataque: Orbes Arcanos",
			line3:     "Magias explosivas",
			special:   "Esp: Pilar de Fogo",
		},
		// Fileira Inferior (2 novas classes de fogo)
		{
			class:     entities.ClassShotgunner,
			title:     "ATIRADOR PESADO",
			bounds:    collision.NewRect(145, 282, 240, 195),
			spriteImg: sprites.GetShotgunnerSprite(),
			themeClr:  color.RGBA{210, 90, 40, 255},
			line1:     "HP: 140 | Vel: Baixa",
			line2:     "Ataque: Espingarda Cal.12",
			line3:     "6 balotes com alto recuo",
			special:   "Esp: Disparo Sismico 360",
		},
		{
			class:     entities.ClassSheriff,
			title:     "XERIFE",
			bounds:    collision.NewRect(415, 282, 240, 195),
			spriteImg: sprites.GetSheriffSprite(),
			themeClr:  color.RGBA{235, 185, 50, 255},
			line1:     "HP: 95 | Tambor: 6 Balas",
			line2:     "Ataque: Revolver Preciso",
			line3:     "Recarga [R] e Criticos",
			special:   "Esp: Deadeye (6 Teleg.)",
		},
	}

	return &ClassSelectScene{
		navigator:     nav,
		selectedClass: entities.ClassWarrior,
		cards:         cards,
		btnStart:      ui.NewButton(250, 495, 300, 48, "⚔  INICIAR BATALHA  ⚔"),
		btnBack:       ui.NewButton(35, 542, 160, 40, "🏰  VOLTAR  🏰"),
	}
}

// Update processa clique direto nos cards das classes e nos botões.
func (c *ClassSelectScene) Update() error {
	c.btnStart.Update()
	c.btnBack.Update()

	// Clique em qualquer um dos 5 cards seleciona a classe
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		for _, card := range c.cards {
			if card.bounds.Contains(float64(mx), float64(my)) {
				c.selectedClass = card.class
				break
			}
		}
	}

	if c.btnStart.IsClicked() {
		c.navigator.SwitchToGameplay(c.selectedClass)
	} else if c.btnBack.IsClicked() {
		c.navigator.SwitchToMenu()
	}

	return nil
}

// Draw renderiza a tela com os 5 heróis, prévia em pixel art e atributos.
func (c *ClassSelectScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 20, G: 18, B: 24, A: 255})

	// Título
	ebitenutil.DebugPrintAt(screen, "=== ⚜ SELECAO DE HEROI (CLIQUE PARA ESCOLHER) ⚜ ===", 225, 18)
	ebitenutil.DebugPrintAt(screen, "Escolha seu campeao entre as 5 classes para enfrentar a masmorra:", 185, 38)

	mx, my := ebiten.CursorPosition()
	for _, card := range c.cards {
		isHovered := card.bounds.Contains(float64(mx), float64(my))
		c.drawCard(screen, card, isHovered)
	}

	c.btnStart.Draw(screen)
	c.btnBack.Draw(screen)
}

func (c *ClassSelectScene) drawCard(dst *ebiten.Image, card classCardMeta, isHovered bool) {
	isSelected := c.selectedClass == card.class
	x := float32(card.bounds.X)
	y := float32(card.bounds.Y)
	w := float32(card.bounds.W)
	h := float32(card.bounds.H)

	bg := color.RGBA{R: 32, G: 26, B: 28, A: 245}
	if isSelected {
		bg = color.RGBA{R: 52, G: 40, B: 34, A: 255}
	} else if isHovered {
		bg = color.RGBA{R: 42, G: 34, B: 36, A: 255}
	}
	vector.DrawFilledRect(dst, x, y, w, h, bg, false)

	borderClr := color.RGBA{R: 95, G: 75, B: 60, A: 220}
	stroke := float32(1.5)
	if isSelected {
		borderClr = color.RGBA{R: 255, G: 215, B: 60, A: 255}
		stroke = 3.0
	} else if isHovered {
		borderClr = color.RGBA{R: 190, G: 150, B: 90, A: 255}
	}
	vector.StrokeRect(dst, x, y, w, h, stroke, borderClr, false)

	// Título
	titleX := int(x) + (int(w)-len(card.title)*6)/2
	ebitenutil.DebugPrintAt(dst, card.title, titleX, int(y)+10)

	// Caixa do avatar em Pixel Art
	boxSize := float32(48)
	boxX := x + (w-boxSize)/2
	boxY := y + 30
	vector.DrawFilledRect(dst, boxX, boxY, boxSize, boxSize, color.RGBA{R: 18, G: 16, B: 22, A: 255}, false)
	vector.StrokeRect(dst, boxX, boxY, boxSize, boxSize, 1.5, card.themeClr, false)

	if card.spriteImg != nil {
		op := &ebiten.DrawImageOptions{}
		scale := 2.4
		op.GeoM.Scale(scale, scale)
		sw := float64(card.spriteImg.Bounds().Dx()) * scale
		sh := float64(card.spriteImg.Bounds().Dy()) * scale
		op.GeoM.Translate(float64(boxX)+(float64(boxSize)-sw)/2, float64(boxY)+(float64(boxSize)-sh)/2)
		dst.DrawImage(card.spriteImg, op)
	}

	// Textos de atributos
	ty := int(y) + 88
	ebitenutil.DebugPrintAt(dst, card.line1, int(x)+14, ty)
	ebitenutil.DebugPrintAt(dst, card.line2, int(x)+14, ty+18)
	ebitenutil.DebugPrintAt(dst, card.line3, int(x)+14, ty+36)
	ebitenutil.DebugPrintAt(dst, card.special, int(x)+14, ty+56)

	if isSelected {
		ebitenutil.DebugPrintAt(dst, ">> SELECIONADO <<", int(x)+(int(w)-102)/2, int(y+h)-22)
	}
}
