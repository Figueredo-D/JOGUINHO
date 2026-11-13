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

// TorchSpark representa uma faísca ou brasa flutuando no menu inicial.
type TorchSpark struct {
	x, y   float32
	vx, vy float32
	life   int
	size   float32
	clr    color.RGBA
}

// MenuScene gerencia a tela inicial da Versão 1.0 com estética medieval e brasas flutuantes.
type MenuScene struct {
	navigator       SceneNavigator
	btnPlay         *ui.Button
	btnSanctuary    *ui.Button
	btnAchievements *ui.Button
	btnHowToPlay    *ui.Button
	btnOptions      *ui.Button
	btnQuit         *ui.Button
	frameCounter    int
	sparks          []*TorchSpark
	embers          []*TorchSpark
}

// NewMenuScene instancia a tela de menu oficial de "As Catacumbas de Dorfieegu".
func NewMenuScene(nav SceneNavigator) *MenuScene {
	embers := make([]*TorchSpark, 0, 45)
	for i := 0; i < 45; i++ {
		embers = append(embers, &TorchSpark{
			x:    float32(rand.Float64() * 800),
			y:    float32(rand.Float64() * 600),
			vx:   float32((rand.Float64() - 0.5) * 0.4),
			vy:   float32(-0.5 - rand.Float64()*0.9),
			size: float32(1.2 + rand.Float64()*1.8),
			clr:  color.RGBA{R: 255, G: uint8(110 + rand.Intn(110)), B: 30, A: uint8(120 + rand.Intn(110))},
		})
	}

	return &MenuScene{
		navigator:       nav,
		btnPlay:         ui.NewButton(245, 255, 310, 46, "⚔  INICIAR EXPEDIÇÃO  ⚔"),
		btnSanctuary:    ui.NewButton(185, 314, 205, 42, "⚜  SANTUÁRIO  ⚜"),
		btnAchievements: ui.NewButton(410, 314, 205, 42, "🏆  CONQUISTAS  🏆"),
		btnHowToPlay:    ui.NewButton(185, 366, 205, 42, "📜  COMO JOGAR  📜"),
		btnOptions:      ui.NewButton(410, 366, 205, 42, "⚙  OPÇÕES / ÁUDIO  ⚙"),
		btnQuit:         ui.NewButton(270, 424, 260, 42, "⛨  ABANDONAR CASTELO  ⛨"),
		sparks:          make([]*TorchSpark, 0, 40),
		embers:          embers,
	}
}

// Update processa interações e animação das tochas, faíscas e cinzas ascendentes.
func (m *MenuScene) Update() error {
	m.frameCounter++

	m.btnPlay.Update()
	m.btnSanctuary.Update()
	m.btnAchievements.Update()
	m.btnHowToPlay.Update()
	m.btnOptions.Update()
	m.btnQuit.Update()

	if m.btnPlay.IsClicked() {
		m.navigator.SwitchToClassSelect()
	} else if m.btnSanctuary.IsClicked() {
		m.navigator.SwitchToSanctuary()
	} else if m.btnAchievements.IsClicked() {
		m.navigator.SwitchToAchievements()
	} else if m.btnHowToPlay.IsClicked() {
		m.navigator.SwitchToHowToPlay()
	} else if m.btnOptions.IsClicked() {
		m.navigator.SwitchToOptions()
	} else if m.btnQuit.IsClicked() {
		m.navigator.QuitGame()
	}

	// Faíscas das tochas laterais
	if m.frameCounter%3 == 0 {
		m.spawnTorchSparks(120, 240)
		m.spawnTorchSparks(680, 240)
	}
	for i := len(m.sparks) - 1; i >= 0; i-- {
		s := m.sparks[i]
		s.x += s.vx
		s.y += s.vy
		s.life--
		if s.life <= 0 {
			m.sparks = append(m.sparks[:i], m.sparks[i+1:]...)
		}
	}

	// Brasas e cinzas atmosféricas subindo suavemente do fundo da tela
	for _, e := range m.embers {
		e.x += e.vx + float32(math.Sin(float64(m.frameCounter)*0.03+float64(e.y)*0.02)*0.25)
		e.y += e.vy
		if e.y < -10 {
			e.y = 605
			e.x = float32(rand.Float64() * 800)
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
		size: 1.5,
		clr: color.RGBA{
			R: 255,
			G: uint8(140 + rand.Intn(100)),
			B: 30,
			A: 240,
		},
	})
}

// Draw renderiza o menu medieval completo da Versão 1.0.
func (m *MenuScene) Draw(screen *ebiten.Image) {
	// 1. Parede de pedras de fortaleza medieval
	m.drawCastleStoneWall(screen)

	// 2. Brasas e cinzas subindo pelo ambiente
	for _, em := range m.embers {
		vector.DrawFilledCircle(screen, em.x, em.y, em.size, em.clr, true)
	}

	// 3. Iluminação e tochas laterais
	m.drawTorchGlow(screen, 120, 235)
	m.drawTorchGlow(screen, 680, 235)
	m.drawTorch(screen, 120, 240)
	m.drawTorch(screen, 680, 240)

	for _, sp := range m.sparks {
		vector.DrawFilledCircle(screen, sp.x, sp.y, sp.size, sp.clr, true)
	}

	// 4. Brasão heráldico medieval
	m.drawHeraldicCrest(screen, 400, 72)

	// 5. Logotipo Oficial Estilizado: AS CATACUMBAS DE DORFIEEGU (Versão 1.0)
	m.drawTitleBanner(screen, 400, 188)

	// 6. Botões
	m.btnPlay.Draw(screen)
	m.btnSanctuary.Draw(screen)
	m.btnAchievements.Draw(screen)
	m.btnHowToPlay.Draw(screen)
	m.btnOptions.Draw(screen)
	m.btnQuit.Draw(screen)

	// 7. Rodapé
	m.drawFooterScroll(screen)
}

func (m *MenuScene) drawCastleStoneWall(dst *ebiten.Image) {
	dst.Fill(color.RGBA{R: 24, G: 20, B: 18, A: 255})

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
			vector.StrokeLine(dst, x+2, y+2, x+float32(brickW-2), y+2, 1, color.RGBA{R: shade + 18, G: shade + 14, B: shade + 10, A: 200}, false)
			vector.StrokeLine(dst, x+2, y+float32(brickH-2), x+float32(brickW-2), y+float32(brickH-2), 1, color.RGBA{R: 16, G: 12, B: 10, A: 220}, false)
		}
	}

	vector.DrawFilledRect(dst, 0, 0, 80, 600, color.RGBA{R: 20, G: 17, B: 15, A: 220}, false)
	vector.StrokeLine(dst, 80, 0, 80, 600, 3, color.RGBA{R: 70, G: 55, B: 40, A: 255}, false)
	vector.DrawFilledRect(dst, 720, 0, 80, 600, color.RGBA{R: 20, G: 17, B: 15, A: 220}, false)
	vector.StrokeLine(dst, 720, 0, 720, 600, 3, color.RGBA{R: 70, G: 55, B: 40, A: 255}, false)
}

func (m *MenuScene) drawTorchGlow(dst *ebiten.Image, tx, ty float32) {
	pulse := float32(math.Sin(float64(m.frameCounter)*0.12) * 5)
	vector.DrawFilledCircle(dst, tx, ty-10, 85+pulse, color.RGBA{R: 220, G: 110, B: 20, A: 25}, true)
	vector.DrawFilledCircle(dst, tx, ty-10, 55+pulse*0.8, color.RGBA{R: 255, G: 160, B: 40, A: 35}, true)
	vector.DrawFilledCircle(dst, tx, ty-10, 30+pulse*0.5, color.RGBA{R: 255, G: 210, B: 70, A: 50}, true)
}

func (m *MenuScene) drawTorch(dst *ebiten.Image, tx, ty float32) {
	vector.DrawFilledRect(dst, tx-5, ty, 10, 28, color.RGBA{R: 45, G: 42, B: 40, A: 255}, false)
	vector.StrokeRect(dst, tx-5, ty, 10, 28, 1.5, color.RGBA{R: 90, G: 80, B: 70, A: 255}, false)
	vector.DrawFilledRect(dst, tx-9, ty-4, 18, 8, color.RGBA{R: 60, G: 55, B: 50, A: 255}, false)
	vector.StrokeRect(dst, tx-9, ty-4, 18, 8, 1, color.RGBA{R: 130, G: 115, B: 90, A: 255}, false)
	vector.StrokeLine(dst, tx, ty+18, tx+14, ty+32, 2.5, color.RGBA{R: 50, G: 45, B: 40, A: 255}, true)

	flicker1 := float32(math.Sin(float64(m.frameCounter)*0.25) * 3)
	flicker2 := float32(math.Cos(float64(m.frameCounter)*0.35) * 2)

	vector.DrawFilledCircle(dst, tx+flicker1*0.4, ty-14, 14+flicker2, color.RGBA{R: 230, G: 40, B: 20, A: 220}, true)
	vector.DrawFilledCircle(dst, tx+flicker2*0.3, ty-16, 10+flicker1*0.5, color.RGBA{R: 255, G: 140, B: 20, A: 240}, true)
	vector.DrawFilledCircle(dst, tx, ty-18, 6, color.RGBA{R: 255, G: 245, B: 140, A: 255}, true)
}

func (m *MenuScene) drawHeraldicCrest(dst *ebiten.Image, cx, cy float32) {
	vector.StrokeLine(dst, cx-48, cy-42, cx+48, cy+42, 3.5, color.RGBA{R: 210, G: 220, B: 235, A: 255}, true)
	vector.StrokeLine(dst, cx-36, cy-23, cx-21, cy-38, 4.0, color.RGBA{R: 215, G: 170, B: 50, A: 255}, true)
	vector.DrawFilledCircle(dst, cx-46, cy-40, 4.0, color.RGBA{R: 190, G: 140, B: 30, A: 255}, true)

	vector.StrokeLine(dst, cx+48, cy-42, cx-48, cy+42, 3.5, color.RGBA{R: 210, G: 220, B: 235, A: 255}, true)
	vector.StrokeLine(dst, cx+36, cy-23, cx+21, cy-38, 4.0, color.RGBA{R: 215, G: 170, B: 50, A: 255}, true)
	vector.DrawFilledCircle(dst, cx+46, cy-40, 4.0, color.RGBA{R: 190, G: 140, B: 30, A: 255}, true)

	shieldW := float32(48)
	shieldH := float32(58)
	sx := cx - shieldW/2
	sy := cy - shieldH/2

	vector.DrawFilledRect(dst, sx, sy, shieldW/2, shieldH*0.7, color.RGBA{R: 160, G: 30, B: 35, A: 255}, false)
	vector.DrawFilledRect(dst, sx+shieldW/2, sy, shieldW/2, shieldH*0.7, color.RGBA{R: 25, G: 65, B: 140, A: 255}, false)
	vector.DrawFilledRect(dst, sx+6, sy+shieldH*0.7, shieldW-12, 11, color.RGBA{R: 160, G: 30, B: 35, A: 255}, false)
	vector.DrawFilledCircle(dst, cx, sy+shieldH*0.8, 11, color.RGBA{R: 90, G: 45, B: 85, A: 255}, true)
	vector.StrokeRect(dst, sx, sy, shieldW, shieldH*0.75, 2.5, color.RGBA{R: 245, G: 195, B: 55, A: 255}, false)

	vector.StrokeLine(dst, cx, sy+8, cx, sy+shieldH*0.7, 3, color.RGBA{R: 255, G: 225, B: 90, A: 255}, false)
	vector.StrokeLine(dst, sx+9, sy+20, sx+shieldW-9, sy+20, 3, color.RGBA{R: 255, G: 225, B: 90, A: 255}, false)
}

func (m *MenuScene) drawTitleBanner(dst *ebiten.Image, cx, cy float32) {
	bannerW := float32(540)
	bannerH := float32(84)
	bx := cx - bannerW/2
	by := cy - bannerH/2

	vector.DrawFilledRect(dst, bx+4, by+5, bannerW, bannerH, color.RGBA{R: 10, G: 8, B: 6, A: 205}, false)
	vector.DrawFilledRect(dst, bx, by, bannerW, bannerH, color.RGBA{R: 45, G: 30, B: 20, A: 255}, false)

	vector.StrokeLine(dst, bx+8, by+28, bx+bannerW-8, by+28, 1.5, color.RGBA{R: 30, G: 20, B: 14, A: 160}, false)
	vector.StrokeLine(dst, bx+8, by+56, bx+bannerW-8, by+56, 1.5, color.RGBA{R: 30, G: 20, B: 14, A: 160}, false)

	vector.StrokeRect(dst, bx, by, bannerW, bannerH, 2.5, color.RGBA{R: 220, G: 175, B: 70, A: 255}, false)
	vector.StrokeRect(dst, bx+4, by+4, bannerW-8, bannerH-8, 1.0, color.RGBA{R: 90, G: 70, B: 45, A: 200}, false)

	drawRivet(dst, bx+10, by+10, color.RGBA{R: 240, G: 200, B: 100, A: 255})
	drawRivet(dst, bx+bannerW-10, by+10, color.RGBA{R: 240, G: 200, B: 100, A: 255})
	drawRivet(dst, bx+10, by+bannerH-10, color.RGBA{R: 240, G: 200, B: 100, A: 255})
	drawRivet(dst, bx+bannerW-10, by+bannerH-10, color.RGBA{R: 240, G: 200, B: 100, A: 255})

	t1 := "=== ⚜ AS CATACUMBAS DE DORFIEEGU ⚜ ==="
	t2 := "VERSAO 1.0 OFICIAL  |  ROGUE-LITE MEDIEVAL"
	t3 := "'Purifique os 15 andares e derrote o Soberano das Sombras!'"

	ebitenutil.DebugPrintAt(dst, t1, int(cx)-130, int(by)+15)
	ebitenutil.DebugPrintAt(dst, t2, int(cx)-135, int(by)+37)
	ebitenutil.DebugPrintAt(dst, t3, int(cx)-182, int(by)+58)
}

func (m *MenuScene) drawFooterScroll(dst *ebiten.Image) {
	scrollW := float32(580)
	scrollH := float32(42)
	sx := float32(400) - scrollW/2
	sy := float32(536)

	vector.DrawFilledRect(dst, sx, sy, scrollW, scrollH, color.RGBA{R: 35, G: 28, B: 20, A: 240}, false)
	vector.StrokeRect(dst, sx, sy, scrollW, scrollH, 1.5, color.RGBA{R: 140, G: 110, B: 65, A: 220}, false)

	ebitenutil.DebugPrintAt(dst, "[ESC / P] Pausa  |  [Espaco / Shift] Dash  |  [E / Botao Dir.] Habilidade", int(sx)+38, int(sy)+8)
	ebitenutil.DebugPrintAt(dst, "WASD mover | Mouse mirar & Esq. atacar | [R] Recarregar (Xerife)", int(sx)+45, int(sy)+24)
}
