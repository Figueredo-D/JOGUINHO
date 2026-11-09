package scenes

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/entities"
	"joguinho/internal/sprites"
	"joguinho/internal/ui"
)

// GameOverScene exibe a derrota com estatísticas, fragmentos resgatados e navegação ao Santuário.
type GameOverScene struct {
	navigator    SceneNavigator
	score        int
	wave         int
	class        entities.ClassType
	soulsGained  int
	btnRetry     *ui.Button
	btnSanctuary *ui.Button
	btnMenu      *ui.Button
	soulIcon     *ebiten.Image
}

// NewGameOverScene cria a tela de fim de jogo com resgate de almas meta.
func NewGameOverScene(nav SceneNavigator, score, wave int, class entities.ClassType, soulsGained int) *GameOverScene {
	btnW := 260.0
	btnH := 44.0
	centerX := (800.0 - btnW) / 2

	return &GameOverScene{
		navigator:    nav,
		score:        score,
		wave:         wave,
		class:        class,
		soulsGained:  soulsGained,
		btnRetry:     ui.NewButton(centerX, 335, btnW, btnH, "⚔  RENASCER / RETENTAR  ⚔"),
		btnSanctuary: ui.NewButton(centerX, 390, btnW, btnH, "⚜  SANTUÁRIO / EVOLUÇÃO  ⚜"),
		btnMenu:      ui.NewButton(centerX, 445, btnW, btnH, "🏰  VOLTAR AO MENU  🏰"),
		soulIcon:     sprites.GetSoulFragmentSprite(),
	}
}

// Update processa as ações dos botões.
func (g *GameOverScene) Update() error {
	g.btnRetry.Update()
	g.btnSanctuary.Update()
	g.btnMenu.Update()

	if g.btnRetry.IsClicked() {
		g.navigator.SwitchToGameplay(g.class)
	} else if g.btnSanctuary.IsClicked() {
		g.navigator.SwitchToSanctuary()
	} else if g.btnMenu.IsClicked() {
		g.navigator.SwitchToMenu()
	}

	return nil
}

// Draw renderiza a tela de derrota com temática de cripta de fortaleza.
func (g *GameOverScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 22, G: 16, B: 18, A: 255})

	// Painel central de tábua de carvalho escurecido
	panelW := float32(500)
	panelH := float32(415)
	panelX := float32((800 - panelW) / 2)
	panelY := float32(85)

	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, color.RGBA{R: 35, G: 22, B: 24, A: 250}, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 2.5, color.RGBA{R: 180, G: 50, B: 60, A: 255}, false)
	vector.StrokeRect(screen, panelX+4, panelY+4, panelW-8, panelH-8, 1.0, color.RGBA{R: 90, G: 30, B: 35, A: 200}, false)

	// Título gótico
	ebitenutil.DebugPrintAt(screen, "=== ⛨ SEU HERÓI SUCUMBIU ⛨ ===", int(panelX)+135, int(panelY)+25)
	ebitenutil.DebugPrintAt(screen, "A escuridao da masmorra ceifou vossa jornada...", int(panelX)+95, int(panelY)+48)

	// Linha divisória
	vector.StrokeLine(screen, panelX+30, panelY+72, panelX+panelW-30, panelY+72, 1.5, color.RGBA{R: 120, G: 45, B: 55, A: 200}, false)

	// Estatísticas da batalha
	statClass := fmt.Sprintf("Classe Combatente: %s", g.class)
	statWave := fmt.Sprintf("Onda Sobrevivida: %d", g.wave)
	statScore := fmt.Sprintf("Pontuacao de Honra: %d pts", g.score)
	statSouls := fmt.Sprintf("Essencias de Alma Resgatadas: +%d", g.soulsGained)

	ebitenutil.DebugPrintAt(screen, statClass, int(panelX)+110, int(panelY)+95)
	ebitenutil.DebugPrintAt(screen, statWave, int(panelX)+110, int(panelY)+120)
	ebitenutil.DebugPrintAt(screen, statScore, int(panelX)+110, int(panelY)+145)

	// Destaque para Essências de Alma resgatadas
	vector.DrawFilledRect(screen, panelX+40, panelY+170, panelW-80, 36, color.RGBA{R: 45, G: 25, B: 55, A: 240}, false)
	vector.StrokeRect(screen, panelX+40, panelY+170, panelW-80, 36, 1.5, color.RGBA{R: 170, G: 110, B: 255, A: 255}, false)

	if g.soulIcon != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1.8, 1.8)
		op.GeoM.Translate(float64(panelX)+50, float64(panelY)+173)
		screen.DrawImage(g.soulIcon, op)
	}
	ebitenutil.DebugPrintAt(screen, statSouls, int(panelX)+85, int(panelY)+182)

	// Linha divisória antes dos botões
	vector.StrokeLine(screen, panelX+30, panelY+225, panelX+panelW-30, panelY+225, 1.5, color.RGBA{R: 120, G: 45, B: 55, A: 200}, false)

	// Botões
	g.btnRetry.Draw(screen)
	g.btnSanctuary.Draw(screen)
	g.btnMenu.Draw(screen)
}
