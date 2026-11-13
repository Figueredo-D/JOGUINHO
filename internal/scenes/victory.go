package scenes

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/entities"
	"joguinho/internal/sprites"
	"joguinho/internal/ui"
)

// RunSummary contém o relatório completo de desempenho da partida concluída com vitória.
type RunSummary struct {
	Class         entities.ClassType
	Score         int
	Kills         int
	RelicsCount   int
	DamageDealt   float64
	SurvivalSecs  int
	BaseSouls     int
	BonusSouls    int
	TotalSouls    int
}

type triumphSpark struct {
	x, y float32
	vy   float32
	size float32
	clr  color.RGBA
}

// VictoryScene exibe a tela triunfal ao derrotar o Soberano de Dorfieegu na Onda 15.
type VictoryScene struct {
	navigator    SceneNavigator
	summary      RunSummary
	btnSanctuary *ui.Button
	btnPlayAgain *ui.Button
	btnMenu      *ui.Button
	soulIcon     *ebiten.Image
	sparks       []triumphSpark
	frame        int
}

// NewVictoryScene instancia a tela de vitória triunfal.
func NewVictoryScene(nav SceneNavigator, summary RunSummary) *VictoryScene {
	btnW := 260.0
	btnH := 42.0
	centerX := (800.0 - btnW) / 2

	sparks := make([]triumphSpark, 50)
	for i := range sparks {
		sparks[i] = triumphSpark{
			x:    float32(rand.Float64() * 800),
			y:    float32(rand.Float64() * 600),
			vy:   float32(-0.8 - rand.Float64()*1.5),
			size: float32(1.5 + rand.Float64()*2.5),
			clr:  color.RGBA{R: 255, G: uint8(190 + rand.Intn(65)), B: 60, A: 220},
		}
	}

	return &VictoryScene{
		navigator:    nav,
		summary:      summary,
		btnSanctuary: ui.NewButton(centerX, 395, btnW, btnH, "⚜  SANTUÁRIO / EVOLUÇÃO  ⚜"),
		btnPlayAgain: ui.NewButton(centerX, 446, btnW, btnH, "⚔  NOVA EXPEDIÇÃO  ⚔"),
		btnMenu:      ui.NewButton(centerX, 497, btnW, btnH, "🏰  VOLTAR AO MENU  🏰"),
		soulIcon:     sprites.GetSoulFragmentSprite(),
		sparks:       sparks,
	}
}

// Update anima as partículas douradas e processa os botões.
func (v *VictoryScene) Update() error {
	v.frame++

	for i := range v.sparks {
		v.sparks[i].y += v.sparks[i].vy
		if v.sparks[i].y < -10 {
			v.sparks[i].y = 610
			v.sparks[i].x = float32(rand.Float64() * 800)
		}
	}

	v.btnSanctuary.Update()
	v.btnPlayAgain.Update()
	v.btnMenu.Update()

	if v.btnSanctuary.IsClicked() {
		v.navigator.SwitchToSanctuary()
	} else if v.btnPlayAgain.IsClicked() {
		v.navigator.SwitchToClassSelect()
	} else if v.btnMenu.IsClicked() {
		v.navigator.SwitchToMenu()
	}

	return nil
}

// Draw renderiza a celebração, o relatório da partida e os botões.
func (v *VictoryScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 18, G: 22, B: 30, A: 255})

	// Partículas douradas de celebração subindo
	for _, sp := range v.sparks {
		vector.DrawFilledCircle(screen, sp.x, sp.y, sp.size, sp.clr, true)
	}

	// Painel Triunfal Central
	panelW := float32(540)
	panelH := float32(340)
	px := float32((800 - panelW) / 2)
	py := float32(36)

	vector.DrawFilledRect(screen, px+5, py+6, panelW, panelH, color.RGBA{10, 8, 12, 210}, false)
	vector.DrawFilledRect(screen, px, py, panelW, panelH, color.RGBA{R: 38, G: 30, B: 22, A: 245}, false)
	vector.StrokeRect(screen, px, py, panelW, panelH, 3.0, color.RGBA{R: 255, G: 215, B: 60, A: 255}, false)
	vector.StrokeRect(screen, px+5, py+5, panelW-10, panelH-10, 1.2, color.RGBA{R: 160, G: 120, B: 45, A: 220}, false)

	drawRivet(screen, px+12, py+12, color.RGBA{255, 220, 90, 255})
	drawRivet(screen, px+panelW-12, py+12, color.RGBA{255, 220, 90, 255})
	drawRivet(screen, px+12, py+panelH-12, color.RGBA{255, 220, 90, 255})
	drawRivet(screen, px+panelW-12, py+panelH-12, color.RGBA{255, 220, 90, 255})

	ebitenutil.DebugPrintAt(screen, "=== ⚜ VITORIA: CATACUMBAS PURIFICADAS ⚜ ===", int(px)+130, int(py)+22)
	ebitenutil.DebugPrintAt(screen, "O Soberano de Dorfieegu tombou! A luz retorna ao reino.", int(px)+100, int(py)+44)

	vector.StrokeLine(screen, px+30, py+68, px+panelW-30, py+68, 1.5, color.RGBA{200, 160, 60, 200}, false)

	// Relatório Pós-Partida Detalhado
	mins := v.summary.SurvivalSecs / 60
	secs := v.summary.SurvivalSecs % 60
	ry := int(py) + 86
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Classe Campea:        %s", v.summary.Class), int(px)+95, ry)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Tempo Sobrevivido:    %02d:%02d", mins, secs), int(px)+95, ry+24)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Inimigos Eliminados:  %d monstros", v.summary.Kills), int(px)+95, ry+48)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Reliquias Coletadas:  %d artefatos", v.summary.RelicsCount), int(px)+95, ry+72)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Dano Total Causado:   %.0f", v.summary.DamageDealt), int(px)+95, ry+96)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pontuacao de Honra:   %d pts", v.summary.Score), int(px)+95, ry+120)

	// Destaque para Almas Resgatadas + Bônus 1.5x
	boxY := py + 235
	vector.DrawFilledRect(screen, px+45, boxY, panelW-90, 52, color.RGBA{45, 25, 58, 240}, false)
	vector.StrokeRect(screen, px+45, boxY, panelW-90, 52, 1.8, color.RGBA{180, 120, 255, 255}, false)

	if v.soulIcon != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2.0, 2.0)
		op.GeoM.Translate(float64(px)+60, float64(boxY)+14)
		screen.DrawImage(v.soulIcon, op)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Almas da Run: %d  +  Bonus Vitoria (1.5x): +%d", v.summary.BaseSouls, v.summary.BonusSouls), int(px)+100, int(boxY)+10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TOTAL RESGATADO PARA O SANTUARIO: +%d ALMAS!", v.summary.TotalSouls), int(px)+100, int(boxY)+30)

	v.btnSanctuary.Draw(screen)
	v.btnPlayAgain.Draw(screen)
	v.btnMenu.Draw(screen)
}
