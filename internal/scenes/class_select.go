package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/entities"
	"joguinho/internal/sprites"
	"joguinho/internal/ui"
)

// ClassSelectScene permite ao jogador escolher entre Guerreiro, Arqueiro e Mago.
type ClassSelectScene struct {
	navigator     SceneNavigator
	selectedClass entities.ClassType
	btnWarrior    *ui.Button
	btnArcher     *ui.Button
	btnMage       *ui.Button
	btnStart      *ui.Button
	btnBack       *ui.Button
	frameCounter  int

	spriteWarrior *ebiten.Image
	spriteArcher  *ebiten.Image
	spriteMage    *ebiten.Image
}

// NewClassSelectScene inicializa a tela de seleção com os 3 cards.
func NewClassSelectScene(nav SceneNavigator) *ClassSelectScene {
	return &ClassSelectScene{
		navigator:     nav,
		selectedClass: entities.ClassWarrior,
		btnWarrior:    ui.NewButton(60, 390, 200, 42, "SELECIONAR"),
		btnArcher:     ui.NewButton(300, 390, 200, 42, "SELECIONAR"),
		btnMage:       ui.NewButton(540, 390, 200, 42, "SELECIONAR"),
		btnStart:      ui.NewButton(260, 480, 280, 50, ">> INICIAR PARTIDA <<"),
		btnBack:       ui.NewButton(40, 530, 140, 36, "VOLTAR"),
		spriteWarrior: sprites.GetWarriorSprite(),
		spriteArcher:  sprites.GetArcherSprite(false),
		spriteMage:    sprites.GetMageSprite(),
	}
}

// Update processa seleção e navegação.
func (c *ClassSelectScene) Update() error {
	c.frameCounter++

	c.btnWarrior.Update()
	c.btnArcher.Update()
	c.btnMage.Update()
	c.btnStart.Update()
	c.btnBack.Update()

	if c.btnWarrior.IsClicked() {
		c.selectedClass = entities.ClassWarrior
	} else if c.btnArcher.IsClicked() {
		c.selectedClass = entities.ClassArcher
	} else if c.btnMage.IsClicked() {
		c.selectedClass = entities.ClassMage
	}

	if c.btnStart.IsClicked() {
		c.navigator.SwitchToGameplay(c.selectedClass)
	}

	if c.btnBack.IsClicked() {
		c.navigator.SwitchToMenu()
	}

	return nil
}

// Draw renderiza os cards com visualização em pixel art procedural, atributos e descrições.
func (c *ClassSelectScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 18, G: 22, B: 32, A: 255})

	// Título
	ebitenutil.DebugPrintAt(screen, "=== SELECAO DE CLASSE ===", 320, 25)
	ebitenutil.DebugPrintAt(screen, "Escolha seu heroi para adentrar na arena:", 280, 48)

	// Desenha os 3 cards de classes
	c.drawClassCard(screen, 50, 80, 220, 365, entities.ClassWarrior, "GUERREIRO",
		c.spriteWarrior,
		color.RGBA{R: 200, G: 70, B: 70, A: 255},
		"HP: 200 (Alto)",
		"Mana: Nenhuma",
		"Velocidade: Media",
		"Ataque: Melee Espada Larga",
		"Golpe em arco frontal",
		"com alto dano em area.",
	)

	c.drawClassCard(screen, 290, 80, 220, 365, entities.ClassArcher, "ARQUEIRO",
		c.spriteArcher,
		color.RGBA{R: 50, G: 180, B: 90, A: 255},
		"HP: 75 (Baixo)",
		"Mana: Nenhuma",
		"Velocidade: Rapida",
		"Ataque: Flechas Continuas",
		"Disparos velozes e",
		"constantes no cursor.",
	)

	c.drawClassCard(screen, 530, 80, 220, 365, entities.ClassMage, "MAGO",
		c.spriteMage,
		color.RGBA{R: 60, G: 110, B: 230, A: 255},
		"HP: 110 (Medio)",
		"Mana: 100 (Regenera)",
		"Velocidade: Normal",
		"Ataque: Magia Arcana",
		"Orbes arcanos poderosos",
		"com gasto de mana.",
	)

	// Renderiza botões dos cards
	c.btnWarrior.Draw(screen)
	c.btnArcher.Draw(screen)
	c.btnMage.Draw(screen)

	// Botões de ação principal
	c.btnStart.Draw(screen)
	c.btnBack.Draw(screen)
}

func (c *ClassSelectScene) drawClassCard(
	dst *ebiten.Image,
	x, y, w, h float32,
	class entities.ClassType,
	title string,
	spriteImg *ebiten.Image,
	themeColor color.RGBA,
	hp, mana, speed, atkType, desc1, desc2 string,
) {
	isSelected := c.selectedClass == class

	// Fundo do card
	bgColor := color.RGBA{R: 28, G: 34, B: 48, A: 240}
	if isSelected {
		bgColor = color.RGBA{R: 38, G: 48, B: 72, A: 255}
	}
	vector.DrawFilledRect(dst, x, y, w, h, bgColor, false)

	// Borda destacada se selecionado
	borderColor := color.RGBA{R: 80, G: 100, B: 140, A: 200}
	borderStroke := float32(1.5)
	if isSelected {
		borderColor = color.RGBA{R: 255, G: 215, B: 80, A: 255}
		borderStroke = 3.0
	}
	vector.StrokeRect(dst, x, y, w, h, borderStroke, borderColor, false)

	// Título da classe
	ebitenutil.DebugPrintAt(dst, " [ "+title+" ] ", int(x)+55, int(y)+16)

	// Caixa do Avatar em Pixel Art
	avatarBoxX := x + (w-56)/2
	avatarBoxY := y + 46
	vector.DrawFilledRect(dst, avatarBoxX, avatarBoxY, 56, 56, color.RGBA{R: 20, G: 24, B: 36, A: 255}, false)
	vector.StrokeRect(dst, avatarBoxX, avatarBoxY, 56, 56, 1.5, themeColor, false)

	if spriteImg != nil {
		op := &ebiten.DrawImageOptions{}
		scale := 2.6
		op.GeoM.Scale(scale, scale)
		sw := float64(spriteImg.Bounds().Dx()) * scale
		sh := float64(spriteImg.Bounds().Dy()) * scale
		op.GeoM.Translate(float64(avatarBoxX)+(56-sw)/2, float64(avatarBoxY)+(56-sh)/2)
		dst.DrawImage(spriteImg, op)
	}

	// Atributos
	ty := int(y) + 116
	ebitenutil.DebugPrintAt(dst, hp, int(x)+20, ty)
	ebitenutil.DebugPrintAt(dst, mana, int(x)+20, ty+18)
	ebitenutil.DebugPrintAt(dst, speed, int(x)+20, ty+36)
	ebitenutil.DebugPrintAt(dst, atkType, int(x)+20, ty+54)

	// Linha divisória
	vector.StrokeLine(dst, x+15, float32(ty+78), x+w-15, float32(ty+78), 1, color.RGBA{R: 70, G: 85, B: 120, A: 180}, false)

	// Descrição
	ebitenutil.DebugPrintAt(dst, desc1, int(x)+20, ty+90)
	ebitenutil.DebugPrintAt(dst, desc2, int(x)+20, ty+108)

	if isSelected {
		ebitenutil.DebugPrintAt(dst, ">> CLASSE ATIVA <<", int(x)+45, ty+140)
	}
}
