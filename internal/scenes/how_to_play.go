package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/ui"
)

// HowToPlayScene exibe o guia visual de comandos, classes e sistemas do jogo.
type HowToPlayScene struct {
	navigator SceneNavigator
	btnBack   *ui.Button
}

// NewHowToPlayScene inicializa a tela de guia de comandos.
func NewHowToPlayScene(nav SceneNavigator) *HowToPlayScene {
	return &HowToPlayScene{
		navigator: nav,
		btnBack:   ui.NewButton(280, 530, 240, 44, "🏰  VOLTAR AO MENU  🏰"),
	}
}

// Update processa clique no botão voltar.
func (h *HowToPlayScene) Update() error {
	h.btnBack.Update()
	if h.btnBack.IsClicked() {
		h.navigator.SwitchToMenu()
	}
	return nil
}

// Draw renderiza as teclas estilizadas e descrições de cada ação.
func (h *HowToPlayScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 20, G: 18, B: 22, A: 255})

	// Título
	vector.DrawFilledRect(screen, 180, 18, 440, 56, color.RGBA{R: 44, G: 30, B: 22, A: 250}, false)
	vector.StrokeRect(screen, 180, 18, 440, 56, 2.5, color.RGBA{R: 215, G: 175, B: 70, A: 255}, false)
	ebitenutil.DebugPrintAt(screen, "=== 📜 GUIA DE COMANDOS E SOBREVIVENCIA 📜 ===", 255, 32)
	ebitenutil.DebugPrintAt(screen, "Domine suas habilidades para purificar As Catacumbas de Dorfieegu", 205, 52)

	// Bloco Esquerdo: Controles
	vector.DrawFilledRect(screen, 35, 92, 360, 415, color.RGBA{R: 32, G: 25, B: 22, A: 245}, false)
	vector.StrokeRect(screen, 35, 92, 360, 415, 2, color.RGBA{R: 150, G: 115, B: 65, A: 230}, false)
	ebitenutil.DebugPrintAt(screen, "[ CONTROLES DE COMBATE ]", 145, 108)

	h.drawKeyRow(screen, 52, 138, "W A S D / Setas", "Movimentar o heroi pela arena")
	h.drawKeyRow(screen, 52, 198, "Mouse + Esq.", "Mirar e disparar ataque basico")
	h.drawKeyRow(screen, 52, 258, "Espaco / Shift", "Esquiva (Dash) com invulnerabilidade")
	h.drawKeyRow(screen, 52, 318, "E / Botao Dir.", "Ativar Habilidade Especial da classe")
	h.drawKeyRow(screen, 52, 378, "Tecla R", "Recarregar tambor de balas (Xerife)")
	h.drawKeyRow(screen, 52, 438, "ESC ou P", "Pausar partida e ver ficha da run")

	// Bloco Direito: Jornada e Dicas
	vector.DrawFilledRect(screen, 415, 92, 350, 415, color.RGBA{R: 32, G: 25, B: 22, A: 245}, false)
	vector.StrokeRect(screen, 415, 92, 350, 415, 2, color.RGBA{R: 150, G: 115, B: 65, A: 230}, false)
	ebitenutil.DebugPrintAt(screen, "[ OBJETIVO DA EXPEDICAO ]", 505, 108)

	tips := []string{
		"1. CAMPANHA DE 15 ONDAS:",
		"   Enfrente o Carrasco na Onda 5,",
		"   o Golem de Almas na Onda 10 e",
		"   O Soberano de Dorfieegu na Onda 15!",
		"",
		"2. RELIQUIAS PASSIVAS:",
		"   Ao vencer Chefes ou marcos de onda,",
		"   escolha 1 de 3 artefatos poderosos",
		"   que acumulam durante a partida.",
		"",
		"3. FRAGMENTOS DE ALMA & SANTUARIO:",
		"   Colete almas caidas e invista no",
		"   Santuario para evoluir suas classes",
		"   permanentemente entre as partidas.",
		"",
		"4. CLIMA E DROPS:",
		"   Fique atento a pochetes de veneno,",
		"   Sementes do Caos e eventos climaticos",
		"   como Chuva Acida e Lua Carmesim!",
	}

	for i, line := range tips {
		ebitenutil.DebugPrintAt(screen, line, 432, 138+i*18)
	}

	h.btnBack.Draw(screen)
}

func (h *HowToPlayScene) drawKeyRow(dst *ebiten.Image, x, y float32, keyLabel, desc string) {
	// Caixinha visual imitando tecla medieval
	boxW := float32(120)
	boxH := float32(34)
	vector.DrawFilledRect(dst, x, y, boxW, boxH, color.RGBA{56, 42, 32, 255}, false)
	vector.StrokeRect(dst, x, y, boxW, boxH, 2, color.RGBA{220, 180, 80, 255}, false)
	tx := int(x) + (int(boxW)-len(keyLabel)*6)/2
	ebitenutil.DebugPrintAt(dst, keyLabel, tx, int(y)+11)

	// Descrição ao lado
	if len(desc) > 32 {
		cut := 32
		for cut > 12 && desc[cut] != ' ' {
			cut--
		}
		ebitenutil.DebugPrintAt(dst, desc[:cut], int(x+boxW)+12, int(y)+4)
		ebitenutil.DebugPrintAt(dst, desc[cut+1:], int(x+boxW)+12, int(y)+20)
	} else {
		ebitenutil.DebugPrintAt(dst, desc, int(x+boxW)+12, int(y)+11)
	}
}
