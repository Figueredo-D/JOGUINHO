package scenes

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/ui"
)

// TorchSpark representa uma faísca ou brasa flutuando de uma tocha medieval.
type TorchSpark struct {
	x, y   float32
	vx, vy float32
	life   int
	clr    color.RGBA
}

// MenuScene gerencia a tela inicial com ambientação e design inspirados na Idade Média.
type MenuScene struct {
	navigator    SceneNavigator
	btnPlay      *ui.Button
	btnOptions   *ui.Button
	btnQuit      *ui.Button
	frameCounter int
	sparks       []*TorchSpark
}

// NewMenuScene instancia a tela de menu no estilo castelo medieval.
func NewMenuScene(nav SceneNavigator) *MenuScene {
	btnW := 260.0
	btnH := 48.0
	centerX := (800.0 - btnW) / 2

	return &MenuScene{
		navigator:  nav,
		btnPlay:    ui.NewButton(centerX, 290, btnW, btnH, "⚔  INICIAR JORNADA  ⚔"),
		btnOptions: ui.NewButton(centerX, 355, btnW, btnH, "⚙  OPÇÕES DO REINO  ⚙"),
		btnQuit:    ui.NewButton(centerX, 420, btnW, btnH, "⛨  ABANDONAR CASTELO  ⛨"),
		sparks:     make([]*TorchSpark, 0, 40),
	}
}

// Update processa interações e animação das labaredas e faíscas medievais.
func (m *MenuScene) Update() error {
	m.frameCounter++

	m.btnPlay.Update()
	m.btnOptions.Update()
	m.btnQuit.Update()

	if m.btnPlay.IsClicked() {
		m.navigator.SwitchToClassSelect()
	} else if m.btnOptions.IsClicked() {
		m.navigator.SwitchToOptions()
	} else if m.btnQuit.IsClicked() {
		m.navigator.QuitGame()
	}

	// Gera faíscas periódicas subindo das tochas
	if m.frameCounter%3 == 0 {
		m.spawnTorchSparks(120, 240)
		m.spawnTorchSparks(680, 240)
	}

	// Atualiza faíscas
	for i := len(m.sparks) - 1; i >= 0; i-- {
		s := m.sparks[i]
		s.x += s.vx
		s.y += s.vy
		s.life--
		if s.life <= 0 {
			m.sparks = append(m.sparks[:i], m.sparks[i+1:]...)
		}
	}

	return nil
}

func (m *MenuScene) spawnTorchSparks(tx, ty float32) {
	if len(m.sparks) >= 50 {
		return
	}
	m.sparks = append(m.sparks, &TorchSpark{
		x:    tx + float32(rand.Float64()*12-6),
		y:    ty - 15,
		vx:   float32((rand.Float64() - 0.5) * 0.8),
		vy:   float32(-1.2 - rand.Float64()*1.5),
		life: 25 + rand.Intn(25),
		clr: color.RGBA{
			R: 255,
			G: uint8(140 + rand.Intn(100)),
			B: 30,
			A: 240,
		},
	})
}

// Draw renderiza o menu medieval com muralhas de pedra, tochas, brasão e ornamentos.
func (m *MenuScene) Draw(screen *ebiten.Image) {
	// 1. Parede de pedras de fortaleza medieval
	m.drawCastleStoneWall(screen)

	// 2. Iluminação e clarão das tochas nas paredes
	m.drawTorchGlow(screen, 120, 235)
	m.drawTorchGlow(screen, 680, 235)

	// 3. Tochas de ferro forjado e fogo crepitante
	m.drawTorch(screen, 120, 240)
	m.drawTorch(screen, 680, 240)

	// 4. Faíscas/brasas das tochas
	for _, sp := range m.sparks {
		vector.DrawFilledCircle(screen, sp.x, sp.y, 1.5, sp.clr, true)
	}

	// 5. Brasão heráldico medieval com espadas cruzadas
	m.drawHeraldicCrest(screen, 400, 78)

	// 6. Placa de Título em Carvalho e Ferro Forjado
	m.drawTitleBanner(screen, 400, 195)

	// 7. Botões no estilo medieval
	m.btnPlay.Draw(screen)
	m.btnOptions.Draw(screen)
	m.btnQuit.Draw(screen)

	// 8. Rodapé de manuscrito e guia de teclas
	m.drawFooterScroll(screen)
}

// drawCastleStoneWall desenha uma parede de pedras talhadas de fortaleza.
func (m *MenuScene) drawCastleStoneWall(dst *ebiten.Image) {
	dst.Fill(color.RGBA{R: 24, G: 20, B: 18, A: 255}) // Argamassa escura

	brickH := 36
	brickW := 70

	for row := 0; row*brickH < 600; row++ {
		y := float32(row * brickH)
		shift := float32(0)
		if row%2 == 1 {
			shift = float32(brickW / 2)
		}

		for col := -1; col*brickW < 800+brickW; col++ {
			x := float32(col*brickW) + shift

			// Variação de tonalidade de pedra medieval
			var shade uint8
			seed := (row*17 + col*31) % 5
			switch seed {
			case 0:
				shade = 36
			case 1:
				shade = 42
			case 2:
				shade = 48
			case 3:
				shade = 54
			default:
				shade = 40
			}

			stoneClr := color.RGBA{R: shade, G: uint8(float64(shade) * 0.88), B: uint8(float64(shade) * 0.78), A: 255}
			vector.DrawFilledRect(dst, x+2, y+2, float32(brickW-4), float32(brickH-4), stoneClr, false)

			// Chanfro de iluminação superior da pedra
			vector.StrokeLine(dst, x+2, y+2, x+float32(brickW-2), y+2, 1, color.RGBA{R: shade + 18, G: shade + 14, B: shade + 10, A: 200}, false)
			// Sombra inferior da pedra
			vector.StrokeLine(dst, x+2, y+float32(brickH-2), x+float32(brickW-2), y+float32(brickH-2), 1, color.RGBA{R: 16, G: 12, B: 10, A: 220}, false)
		}
	}

	// Pilares laterais de sustentação da fortaleza
	vector.DrawFilledRect(dst, 0, 0, 80, 600, color.RGBA{R: 20, G: 17, B: 15, A: 220}, false)
	vector.StrokeLine(dst, 80, 0, 80, 600, 3, color.RGBA{R: 70, G: 55, B: 40, A: 255}, false)

	vector.DrawFilledRect(dst, 720, 0, 80, 600, color.RGBA{R: 20, G: 17, B: 15, A: 220}, false)
	vector.StrokeLine(dst, 720, 0, 720, 600, 3, color.RGBA{R: 70, G: 55, B: 40, A: 255}, false)
}

// drawTorchGlow desenha a iluminação quente da tocha se espalhando nas pedras.
func (m *MenuScene) drawTorchGlow(dst *ebiten.Image, tx, ty float32) {
	pulse := float32(math.Sin(float64(m.frameCounter)*0.12) * 5)
	vector.DrawFilledCircle(dst, tx, ty-10, 85+pulse, color.RGBA{R: 220, G: 110, B: 20, A: 25}, true)
	vector.DrawFilledCircle(dst, tx, ty-10, 55+pulse*0.8, color.RGBA{R: 255, G: 160, B: 40, A: 35}, true)
	vector.DrawFilledCircle(dst, tx, ty-10, 30+pulse*0.5, color.RGBA{R: 255, G: 210, B: 70, A: 50}, true)
}

// drawTorch renderiza uma tocha medieval com suporte de ferro forjado e fogo animado.
func (m *MenuScene) drawTorch(dst *ebiten.Image, tx, ty float32) {
	// Suporte de ferro na parede
	vector.DrawFilledRect(dst, tx-5, ty, 10, 28, color.RGBA{R: 45, G: 42, B: 40, A: 255}, false)
	vector.StrokeRect(dst, tx-5, ty, 10, 28, 1.5, color.RGBA{R: 90, G: 80, B: 70, A: 255}, false)

	// Anel de ferro da tocha
	vector.DrawFilledRect(dst, tx-9, ty-4, 18, 8, color.RGBA{R: 60, G: 55, B: 50, A: 255}, false)
	vector.StrokeRect(dst, tx-9, ty-4, 18, 8, 1, color.RGBA{R: 130, G: 115, B: 90, A: 255}, false)

	// Braçadeira angular
	vector.StrokeLine(dst, tx, ty+18, tx+14, ty+32, 2.5, color.RGBA{R: 50, G: 45, B: 40, A: 255}, true)

	// Labareda animada da tocha
	flicker1 := float32(math.Sin(float64(m.frameCounter)*0.25) * 3)
	flicker2 := float32(math.Cos(float64(m.frameCounter)*0.35) * 2)

	// Fogo externo (Vermelho rubi)
	vector.DrawFilledCircle(dst, tx+flicker1*0.4, ty-14, 14+flicker2, color.RGBA{R: 230, G: 40, B: 20, A: 220}, true)
	// Fogo intermediário (Laranja ardente)
	vector.DrawFilledCircle(dst, tx+flicker2*0.3, ty-16, 10+flicker1*0.5, color.RGBA{R: 255, G: 140, B: 20, A: 240}, true)
	// Núcleo incandescente (Amarelo / Branco)
	vector.DrawFilledCircle(dst, tx, ty-18, 6, color.RGBA{R: 255, G: 245, B: 140, A: 255}, true)
}

// drawHeraldicCrest desenha o brasão medieval com espadas cruzadas.
func (m *MenuScene) drawHeraldicCrest(dst *ebiten.Image, cx, cy float32) {
	// Espada 1 (Diagonal esquerda -> direita)
	vector.StrokeLine(dst, cx-50, cy-45, cx+50, cy+45, 3.5, color.RGBA{R: 210, G: 220, B: 235, A: 255}, true)
	// Guarda da espada 1
	vector.StrokeLine(dst, cx-38, cy-25, cx-22, cy-40, 4.0, color.RGBA{R: 215, G: 170, B: 50, A: 255}, true)
	// Pomo
	vector.DrawFilledCircle(dst, cx-48, cy-43, 4.0, color.RGBA{R: 190, G: 140, B: 30, A: 255}, true)

	// Espada 2 (Diagonal direita -> esquerda)
	vector.StrokeLine(dst, cx+50, cy-45, cx-50, cy+45, 3.5, color.RGBA{R: 210, G: 220, B: 235, A: 255}, true)
	// Guarda da espada 2
	vector.StrokeLine(dst, cx+38, cy-25, cx+22, cy-40, 4.0, color.RGBA{R: 215, G: 170, B: 50, A: 255}, true)
	// Pomo
	vector.DrawFilledCircle(dst, cx+48, cy-43, 4.0, color.RGBA{R: 190, G: 140, B: 30, A: 255}, true)

	// Escudo Heráldico Central
	shieldW := float32(52)
	shieldH := float32(64)
	sx := cx - shieldW/2
	sy := cy - shieldH/2

	// Corpo do escudo (bipartido vermelho carmesim e azul nobre)
	vector.DrawFilledRect(dst, sx, sy, shieldW/2, shieldH*0.7, color.RGBA{R: 160, G: 30, B: 35, A: 255}, false)
	vector.DrawFilledRect(dst, sx+shieldW/2, sy, shieldW/2, shieldH*0.7, color.RGBA{R: 25, G: 65, B: 140, A: 255}, false)

	// Bico inferior do escudo
	vector.DrawFilledRect(dst, sx+6, sy+shieldH*0.7, shieldW-12, 12, color.RGBA{R: 160, G: 30, B: 35, A: 255}, false)
	vector.DrawFilledCircle(dst, cx, sy+shieldH*0.8, 12, color.RGBA{R: 90, G: 45, B: 85, A: 255}, true)

	// Borda dourada reforçada do escudo
	vector.StrokeRect(dst, sx, sy, shieldW, shieldH*0.75, 2.5, color.RGBA{R: 245, G: 195, B: 55, A: 255}, false)

	// Cruz dourada heráldica no centro
	vector.StrokeLine(dst, cx, sy+8, cx, sy+shieldH*0.7, 3, color.RGBA{R: 255, G: 225, B: 90, A: 255}, false)
	vector.StrokeLine(dst, sx+10, sy+22, sx+shieldW-10, sy+22, 3, color.RGBA{R: 255, G: 225, B: 90, A: 255}, false)
}

// drawTitleBanner desenha a grande placa de tábua de carvalho do castelo com o título.
func (m *MenuScene) drawTitleBanner(dst *ebiten.Image, cx, cy float32) {
	bannerW := float32(520)
	bannerH := float32(84)
	bx := cx - bannerW/2
	by := cy - bannerH/2

	// Sombra da placa
	vector.DrawFilledRect(dst, bx+4, by+5, bannerW, bannerH, color.RGBA{R: 10, G: 8, B: 6, A: 200}, false)

	// Placa de carvalho nobre
	vector.DrawFilledRect(dst, bx, by, bannerW, bannerH, color.RGBA{R: 45, G: 30, B: 20, A: 255}, false)

	// Ranhuras horizontais da madeira
	vector.StrokeLine(dst, bx+8, by+28, bx+bannerW-8, by+28, 1.5, color.RGBA{R: 30, G: 20, B: 14, A: 160}, false)
	vector.StrokeLine(dst, bx+8, by+56, bx+bannerW-8, by+56, 1.5, color.RGBA{R: 30, G: 20, B: 14, A: 160}, false)

	// Moldura dupla de ferro forjado e ouro antigo
	vector.StrokeRect(dst, bx, by, bannerW, bannerH, 2.5, color.RGBA{R: 180, G: 140, B: 70, A: 255}, false)
	vector.StrokeRect(dst, bx+4, by+4, bannerW-8, bannerH-8, 1.0, color.RGBA{R: 90, G: 70, B: 45, A: 200}, false)

	// Rebites de ferro nos 4 cantos
	drawRivet(dst, bx+10, by+10, color.RGBA{R: 240, G: 200, B: 100, A: 255})
	drawRivet(dst, bx+bannerW-10, by+10, color.RGBA{R: 240, G: 200, B: 100, A: 255})
	drawRivet(dst, bx+10, by+bannerH-10, color.RGBA{R: 240, G: 200, B: 100, A: 255})
	drawRivet(dst, bx+bannerW-10, by+bannerH-10, color.RGBA{R: 240, G: 200, B: 100, A: 255})

	// Textos com visual gótico e medieval
	t1 := "=== ⚜ DUNGEON SURVIVORS ⚜ ==="
	t2 := "CRONICAS DA FORTALEZA DE ACO"
	t3 := "'Pela honra e gloria, que a lamina corte as sombras!'"

	ebitenutil.DebugPrintAt(dst, t1, int(cx)-110, int(by)+16)
	ebitenutil.DebugPrintAt(dst, t2, int(cx)-115, int(by)+38)
	ebitenutil.DebugPrintAt(dst, t3, int(cx)-185, int(by)+58)
}

// drawFooterScroll desenha o rodapé com visual de pergaminho antigo.
func (m *MenuScene) drawFooterScroll(dst *ebiten.Image) {
	scrollW := float32(560)
	scrollH := float32(40)
	sx := float32(400) - scrollW/2
	sy := float32(538)

	// Fundo de pergaminho envelhecido
	vector.DrawFilledRect(dst, sx, sy, scrollW, scrollH, color.RGBA{R: 35, G: 28, B: 20, A: 240}, false)
	vector.StrokeRect(dst, sx, sy, scrollW, scrollH, 1.5, color.RGBA{R: 140, G: 110, B: 65, A: 220}, false)

	ebitenutil.DebugPrintAt(dst, "[ESC] Durante a batalha para pausar e abrir o menu do reino", int(sx)+65, int(sy)+8)
	ebitenutil.DebugPrintAt(dst, "Controles: WASD mover | Mouse mirar | Clique Esquerdo ou Espaco atacar", int(sx)+30, int(sy)+22)
}
