package scenes

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/progression"
	"joguinho/internal/ui"
)

// AchievementsScene exibe o quadro de honra e progresso das conquistas persistidas.
type AchievementsScene struct {
	navigator SceneNavigator
	btnBack   *ui.Button
}

// NewAchievementsScene cria a tela de conquistas.
func NewAchievementsScene(nav SceneNavigator) *AchievementsScene {
	return &AchievementsScene{
		navigator: nav,
		btnBack:   ui.NewButton(280, 528, 240, 44, "🏰  VOLTAR AO MENU  🏰"),
	}
}

// Update processa clique no botão de retorno.
func (a *AchievementsScene) Update() error {
	a.btnBack.Update()
	if a.btnBack.IsClicked() {
		a.navigator.SwitchToMenu()
	}
	return nil
}

// Draw renderiza os cartões de conquista com status e barras de progresso.
func (a *AchievementsScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 18, G: 16, B: 22, A: 255})

	// Painel do título
	vector.DrawFilledRect(screen, 170, 20, 460, 62, color.RGBA{R: 44, G: 30, B: 22, A: 250}, false)
	vector.StrokeRect(screen, 170, 20, 460, 62, 2.5, color.RGBA{R: 215, G: 175, B: 70, A: 255}, false)
	ebitenutil.DebugPrintAt(screen, "=== 🏆 SALAO DAS CONQUISTAS 🏆 ===", 275, 34)

	save := progression.GetSave()
	unlockedCount := 0
	for _, def := range progression.AllAchievements {
		if save.UnlockedAchievements[def.ID] {
			unlockedCount++
		}
	}
	summary := fmt.Sprintf("Desbloqueadas: %d / %d  |  Total Abates: %d  |  Melhor Onda: %d",
		unlockedCount, len(progression.AllAchievements), save.TotalKills, save.MaxWaveReached)
	ebitenutil.DebugPrintAt(screen, summary, 200, 56)

	// Renderiza os 6 cards em 2 colunas x 3 linhas
	cardW := float32(360)
	cardH := float32(120)
	startX := float32(30)
	startY := float32(102)
	gapX := float32(20)
	gapY := float32(16)

	for i, def := range progression.AllAchievements {
		col := i % 2
		row := i / 2
		cx := startX + float32(col)*(cardW+gapX)
		cy := startY + float32(row)*(cardH+gapY)

		unlocked := save.UnlockedAchievements[def.ID]
		a.drawAchievementCard(screen, cx, cy, cardW, cardH, def, unlocked, save)
	}

	a.btnBack.Draw(screen)
}

func (a *AchievementsScene) drawAchievementCard(dst *ebiten.Image, x, y, w, h float32, def progression.AchievementDef, unlocked bool, save *progression.SaveData) {
	bg := color.RGBA{R: 30, G: 26, B: 28, A: 240}
	border := color.RGBA{R: 85, G: 75, B: 70, A: 220}
	if unlocked {
		bg = color.RGBA{R: 52, G: 40, B: 24, A: 255}
		border = color.RGBA{R: 245, G: 200, B: 60, A: 255}
	}

	vector.DrawFilledRect(dst, x, y, w, h, bg, false)
	vector.StrokeRect(dst, x, y, w, h, 2, border, false)

	// Ícone lateral (Troféu Dourado se concluído, Cadeado Cinza se bloqueado)
	iconX := x + 28
	iconY := y + 40
	if unlocked {
		vector.DrawFilledCircle(dst, iconX, iconY, 16, color.RGBA{R: 245, G: 195, B: 45, A: 255}, true)
		vector.StrokeCircle(dst, iconX, iconY, 16, 2, color.RGBA{R: 255, G: 245, B: 190, A: 255}, true)
		vector.DrawFilledCircle(dst, iconX, iconY, 7, color.RGBA{R: 255, G: 255, B: 220, A: 255}, true)
	} else {
		vector.DrawFilledCircle(dst, iconX, iconY, 16, color.RGBA{R: 48, G: 45, B: 50, A: 255}, true)
		vector.StrokeCircle(dst, iconX, iconY, 16, 1.5, color.RGBA{R: 100, G: 95, B: 105, A: 255}, true)
	}

	// Status e Título
	statusStr := "[BLOQUEADO]"
	if unlocked {
		statusStr = "[CONCLUIDO ★]"
	}
	ebitenutil.DebugPrintAt(dst, def.Title+"  "+statusStr, int(x)+58, int(y)+14)

	// Descrição quebrada em até 2 linhas se necessário
	desc := def.Description
	if len(desc) > 46 {
		cut := 46
		for cut > 20 && desc[cut] != ' ' {
			cut--
		}
		ebitenutil.DebugPrintAt(dst, desc[:cut], int(x)+58, int(y)+36)
		ebitenutil.DebugPrintAt(dst, desc[cut+1:], int(x)+58, int(y)+52)
	} else {
		ebitenutil.DebugPrintAt(dst, desc, int(x)+58, int(y)+40)
	}

	// Progresso numérico / Barra
	cur := 0
	switch def.ID {
	case "survivor":
		cur = save.MaxWaveReached
	case "iron_wall":
		cur = save.ShieldReflectCount
	case "blood_moon_hunter":
		cur = save.BloodMoonKills
	default:
		if unlocked {
			cur = def.Goal
		}
	}
	if unlocked || cur > def.Goal {
		cur = def.Goal
	}

	barX := x + 58
	barY := y + h - 28
	barW := w - 76
	barH := float32(12)
	vector.DrawFilledRect(dst, barX, barY, barW, barH, color.RGBA{R: 20, G: 18, B: 20, A: 220}, false)

	ratio := float32(cur) / float32(def.Goal)
	fillClr := color.RGBA{R: 110, G: 130, B: 170, A: 255}
	if unlocked {
		fillClr = color.RGBA{R: 80, G: 210, B: 90, A: 255}
	}
	vector.DrawFilledRect(dst, barX, barY, barW*ratio, barH, fillClr, false)
	vector.StrokeRect(dst, barX, barY, barW, barH, 1, border, false)

	progTxt := fmt.Sprintf("%d / %d", cur, def.Goal)
	ebitenutil.DebugPrintAt(dst, progTxt, int(barX+barW/2)-18, int(barY))
}
