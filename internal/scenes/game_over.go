package scenes

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/entities"
	"joguinho/internal/ui"
)

// GameOverScene exibe a derrota com estética de cripta medieval e opções de reinício.
type GameOverScene struct {
	navigator SceneNavigator
	score     int
	wave      int
	class     entities.ClassType
	btnRetry  *ui.Button
	btnMenu   *ui.Button
}

// NewGameOverScene cria a tela de fim de jogo.
func NewGameOverScene(nav SceneNavigator, score, wave int, class entities.ClassType) *GameOverScene {
	btnW := 260.0
	btnH := 46.0
	centerX := (800.0 - btnW) / 2

	return &GameOverScene{
		navigator: nav,
		score:     score,
		wave:      wave,
		class:     class,
		btnRetry:  ui.NewButton(centerX, 360, btnW, btnH, "⚔  RENASCER / RETENTAR  ⚔"),
		btnMenu:   ui.NewButton(centerX, 420, btnW, btnH, "🏰  VOLTAR AO MENU  🏰"),
	}
}

// Update processa as ações dos botões.
func (g *GameOverScene) Update() error {
	g.btnRetry.Update()
	g.btnMenu.Update()

	if g.btnRetry.IsClicked() {
		g.navigator.SwitchToGameplay(g.class)
	} else if g.btnMenu.IsClicked() {
		g.navigator.SwitchToMenu()
	}

	return nil
}

// Draw renderiza a tela de derrota com temática de cripta de fortaleza.
func (g *GameOverScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 22, G: 16, B: 18, A: 255})

	// Painel central de tábua de carvalho escurecido
	panelW := float32(480)
	panelH := float32(380)
	panelX := float32((800 - panelW) / 2)
	panelY := float32(100)

	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, color.RGBA{R: 35, G: 22, B: 24, A: 250}, false)
	vector.StrokeRect(screen, panelX, panelY, panelW, panelH, 2.5, color.RGBA{R: 180, G: 50, B: 60, A: 255}, false)
	vector.StrokeRect(screen, panelX+4, panelY+4, panelW-8, panelH-8, 1.0, color.RGBA{R: 90, G: 30, B: 35, A: 200}, false)

	// Título gótico
	ebitenutil.DebugPrintAt(screen, "=== ⛨ SEU HERÓI SUCUMBIU ⛨ ===", int(panelX)+125, int(panelY)+28)
	ebitenutil.DebugPrintAt(screen, "A escuridao da masmorra ceifou vossa jornada...", int(panelX)+85, int(panelY)+52)

	// Linha divisória
	vector.StrokeLine(screen, panelX+30, panelY+80, panelX+panelW-30, panelY+80, 1.5, color.RGBA{R: 120, G: 45, B: 55, A: 200}, false)

	// Estatísticas da batalha
	statClass := fmt.Sprintf("Classe Combatente: %s", g.class)
	statWave := fmt.Sprintf("Onda Sobrevivida: %d", g.wave)
	statScore := fmt.Sprintf("Pontuacao de Honra: %d pts", g.score)

	ebitenutil.DebugPrintAt(screen, statClass, int(panelX)+135, int(panelY)+120)
	ebitenutil.DebugPrintAt(screen, statWave, int(panelX)+135, int(panelY)+150)
	ebitenutil.DebugPrintAt(screen, statScore, int(panelX)+135, int(panelY)+180)

	// Linha divisória antes dos botões
	vector.StrokeLine(screen, panelX+30, panelY+225, panelX+panelW-30, panelY+225, 1.5, color.RGBA{R: 120, G: 45, B: 55, A: 200}, false)

	// Botões
	g.btnRetry.Draw(screen)
	g.btnMenu.Draw(screen)
}
