package scenes

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/entities"
	"joguinho/internal/progression"
	"joguinho/internal/sprites"
	"joguinho/internal/ui"
)

// SanctuaryScene gerencia o santuário de evolução e a árvore de upgrades permanente.
type SanctuaryScene struct {
	navigator     SceneNavigator
	selectedTab   entities.ClassType
	btnTabWarrior *ui.Button
	btnTabArcher  *ui.Button
	btnTabMage    *ui.Button

	btnUpgrade1 *ui.Button
	btnUpgrade2 *ui.Button
	btnUpgrade3 *ui.Button

	btnBack  *ui.Button
	btnPlay  *ui.Button
	soulIcon *ebiten.Image
}

// NewSanctuaryScene inicializa o santuário de upgrades permanente.
func NewSanctuaryScene(nav SceneNavigator) *SanctuaryScene {
	s := &SanctuaryScene{
		navigator:     nav,
		selectedTab:   entities.ClassWarrior,
		btnTabWarrior: ui.NewButton(80, 110, 200, 40, "⚜ GUERREIRO ⚜"),
		btnTabArcher:  ui.NewButton(300, 110, 200, 40, "🏹 ARQUEIRO 🏹"),
		btnTabMage:    ui.NewButton(520, 110, 200, 40, "🔮 MAGO 🔮"),

		btnUpgrade1: ui.NewButton(490, 195, 230, 42, "COMPRAR"),
		btnUpgrade2: ui.NewButton(490, 290, 230, 42, "COMPRAR"),
		btnUpgrade3: ui.NewButton(490, 385, 230, 42, "COMPRAR"),

		btnBack:  ui.NewButton(80, 525, 240, 44, "🏰  VOLTAR AO MENU  🏰"),
		btnPlay:  ui.NewButton(480, 525, 240, 44, "⚔  SELECIONAR E JOGAR  ⚔"),
		soulIcon: sprites.GetSoulFragmentSprite(),
	}
	s.refreshUpgradeButtons()
	return s
}

// Update processa cliques de abas, compras e navegação.
func (s *SanctuaryScene) Update() error {
	s.btnTabWarrior.Update()
	s.btnTabArcher.Update()
	s.btnTabMage.Update()

	s.btnUpgrade1.Update()
	s.btnUpgrade2.Update()
	s.btnUpgrade3.Update()

	s.btnBack.Update()
	s.btnPlay.Update()

	// Alterna abas
	if s.btnTabWarrior.IsClicked() {
		s.selectedTab = entities.ClassWarrior
		s.refreshUpgradeButtons()
	} else if s.btnTabArcher.IsClicked() {
		s.selectedTab = entities.ClassArcher
		s.refreshUpgradeButtons()
	} else if s.btnTabMage.IsClicked() {
		s.selectedTab = entities.ClassMage
		s.refreshUpgradeButtons()
	}

	// Compras de upgrades
	if s.btnUpgrade1.IsClicked() {
		s.buyUpgrade(1)
	} else if s.btnUpgrade2.IsClicked() {
		s.buyUpgrade(2)
	} else if s.btnUpgrade3.IsClicked() {
		s.buyUpgrade(3)
	}

	// Navegação
	if s.btnBack.IsClicked() {
		s.navigator.SwitchToMenu()
	} else if s.btnPlay.IsClicked() {
		s.navigator.SwitchToClassSelect()
	}

	return nil
}

func (s *SanctuaryScene) buyUpgrade(slot int) {
	save := progression.GetSave()

	var currentLvl *int
	switch s.selectedTab {
	case entities.ClassWarrior:
		if slot == 1 {
			currentLvl = &save.Warrior.MaxHPLevel
		} else if slot == 2 {
			currentLvl = &save.Warrior.DamageLevel
		} else {
			currentLvl = &save.Warrior.ShieldLevel
		}
	case entities.ClassArcher:
		if slot == 1 {
			currentLvl = &save.Archer.SpeedLevel
		} else if slot == 2 {
			currentLvl = &save.Archer.AttackSpeedLvl
		} else {
			currentLvl = &save.Archer.ArrowRainLevel
		}
	case entities.ClassMage:
		if slot == 1 {
			currentLvl = &save.Mage.ManaLevel
		} else if slot == 2 {
			currentLvl = &save.Mage.SpellPowerLevel
		} else {
			currentLvl = &save.Mage.FirePillarLevel
		}
	}

	if currentLvl == nil || *currentLvl >= progression.MaxLevel {
		return
	}

	cost := progression.UpgradeCost(*currentLvl)
	if save.SoulFragments >= cost {
		save.SoulFragments -= cost
		*currentLvl++
		_ = progression.Save()
		s.refreshUpgradeButtons()
	}
}

func (s *SanctuaryScene) refreshUpgradeButtons() {
	save := progression.GetSave()

	var lvl1, lvl2, lvl3 int
	switch s.selectedTab {
	case entities.ClassWarrior:
		lvl1 = save.Warrior.MaxHPLevel
		lvl2 = save.Warrior.DamageLevel
		lvl3 = save.Warrior.ShieldLevel
	case entities.ClassArcher:
		lvl1 = save.Archer.SpeedLevel
		lvl2 = save.Archer.AttackSpeedLvl
		lvl3 = save.Archer.ArrowRainLevel
	case entities.ClassMage:
		lvl1 = save.Mage.ManaLevel
		lvl2 = save.Mage.SpellPowerLevel
		lvl3 = save.Mage.FirePillarLevel
	}

	s.formatButton(s.btnUpgrade1, lvl1, save.SoulFragments)
	s.formatButton(s.btnUpgrade2, lvl2, save.SoulFragments)
	s.formatButton(s.btnUpgrade3, lvl3, save.SoulFragments)
}

func (s *SanctuaryScene) formatButton(btn *ui.Button, lvl, availableSouls int) {
	if lvl >= progression.MaxLevel {
		btn.Text = "★ NÍVEL MÁXIMO ★"
		btn.NormalColor = color.RGBA{R: 40, G: 45, B: 50, A: 255}
		return
	}

	cost := progression.UpgradeCost(lvl)
	if availableSouls >= cost {
		btn.Text = fmt.Sprintf("EVOLUIR (%d Essências)", cost)
		btn.NormalColor = color.RGBA{R: 48, G: 34, B: 24, A: 255}
	} else {
		btn.Text = fmt.Sprintf("REQUER %d ESSÊNCIAS", cost)
		btn.NormalColor = color.RGBA{R: 50, G: 30, B: 30, A: 255}
	}
}

// Draw renderiza o santuário com visual medieval.
func (s *SanctuaryScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 20, G: 16, B: 22, A: 255})

	// Título do Santuário
	vector.DrawFilledRect(screen, 180, 20, 440, 68, color.RGBA{R: 40, G: 28, B: 36, A: 250}, false)
	vector.StrokeRect(screen, 180, 20, 440, 68, 2.5, color.RGBA{R: 195, G: 150, B: 80, A: 255}, false)
	ebitenutil.DebugPrintAt(screen, "=== ⚜ SANTUÁRIO DAS ALMAS ⚜ ===", 280, 35)
	ebitenutil.DebugPrintAt(screen, "Evolua os atributos sagrados de vossos guerreiros", 235, 55)

	// Saldo de Fragmentos no Canto Superior Direito
	save := progression.GetSave()
	vector.DrawFilledRect(screen, 630, 25, 145, 36, color.RGBA{R: 30, G: 22, B: 40, A: 240}, false)
	vector.StrokeRect(screen, 630, 25, 145, 36, 1.5, color.RGBA{R: 160, G: 110, B: 255, A: 255}, false)

	if s.soulIcon != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1.8, 1.8)
		op.GeoM.Translate(638, 28)
		screen.DrawImage(s.soulIcon, op)
	}
	soulsStr := fmt.Sprintf("%d Almas", save.SoulFragments)
	ebitenutil.DebugPrintAt(screen, soulsStr, 672, 36)

	// Abas das Classes
	s.btnTabWarrior.Draw(screen)
	s.btnTabArcher.Draw(screen)
	s.btnTabMage.Draw(screen)

	// Painel Central com os 3 upgrades da classe
	s.drawUpgradesPanel(screen)

	// Botões de navegação
	s.btnBack.Draw(screen)
	s.btnPlay.Draw(screen)
}

func (s *SanctuaryScene) drawUpgradesPanel(dst *ebiten.Image) {
	save := progression.GetSave()

	var (
		name1, desc1, stat1 string
		lvl1                int
		name2, desc2, stat2 string
		lvl2                int
		name3, desc3, stat3 string
		lvl3                int
	)

	switch s.selectedTab {
	case entities.ClassWarrior:
		lvl1 = save.Warrior.MaxHPLevel
		name1 = "Vitalidade Titânica"
		desc1 = "Aumenta permanentemente a Vida Máxima da classe."
		stat1 = fmt.Sprintf("+%d HP Máximo (Nvl %d/5)", lvl1*25, lvl1)

		lvl2 = save.Warrior.DamageLevel
		name2 = "Lâmina Dilacerante"
		desc2 = "Aumenta o dano de impacto e o alcance do golpe Melee."
		stat2 = fmt.Sprintf("+%d%% Dano | +%dpx Alcance (Nvl %d/5)", lvl2*12, lvl2*6, lvl2)

		lvl3 = save.Warrior.ShieldLevel
		name3 = "Bastião Divino (Escudo)"
		desc3 = "Reduz o tempo de recarga da Postura de Escudo e amplia reflexão."
		stat3 = fmt.Sprintf("-%.1fs Recarga | +%.0f%% Dano Refletido (Nvl %d/5)", float64(lvl3)*0.8, float64(lvl3)*35, lvl3)

	case entities.ClassArcher:
		lvl1 = save.Archer.SpeedLevel
		name1 = "Passo do Vento"
		desc1 = "Aumenta a velocidade de movimentação base do arqueiro."
		stat1 = fmt.Sprintf("+%.2f Velocidade (Nvl %d/5)", float64(lvl1)*0.35, lvl1)

		lvl2 = save.Archer.AttackSpeedLvl
		name2 = "Disparo Voraz"
		desc2 = "Aumenta a cadência de tiros e a velocidade dos projéteis."
		stat2 = fmt.Sprintf("+%d%% Vel. Disparo | +%.1f Vel. Flechas (Nvl %d/5)", lvl2*10, float64(lvl2)*1.2, lvl2)

		lvl3 = save.Archer.ArrowRainLevel
		name3 = "Chuva Devastadora"
		desc3 = "Aumenta o raio de impacto e o dano contínuo da Chuva de Flechas."
		stat3 = fmt.Sprintf("+%d%% Área | +%d%% Dano da Salva (Nvl %d/5)", lvl3*15, lvl3*20, lvl3)

	case entities.ClassMage:
		lvl1 = save.Mage.ManaLevel
		name1 = "Fonte Mística"
		desc1 = "Aumenta a Mana Máxima e a taxa de regeneração contínua."
		stat1 = fmt.Sprintf("+%d Mana | +%.1f MP/segundo (Nvl %d/5)", lvl1*20, float64(lvl1)*3.5, lvl1)

		lvl2 = save.Mage.SpellPowerLevel
		name2 = "Canalização Arcana"
		desc2 = "Eleva o dano base de todas as magias e orbes de energia."
		stat2 = fmt.Sprintf("+%d%% Dano Mágico (Nvl %d/5)", lvl2*15, lvl2)

		lvl3 = save.Mage.FirePillarLevel
		name3 = "Cataclismo de Chamas"
		desc3 = "Amplia o raio de erupção e a duração ardente do Pilar de Fogo."
		stat3 = fmt.Sprintf("+%d%% Raio | +%d Frames de Fogo (Nvl %d/5)", lvl3*15, lvl3*20, lvl3)
	}

	// Desenha os 3 blocos de upgrades
	s.drawUpgradeRow(dst, 80, 180, 640, 75, name1, desc1, stat1, lvl1)
	s.drawUpgradeRow(dst, 80, 275, 640, 75, name2, desc2, stat2, lvl2)
	s.drawUpgradeRow(dst, 80, 370, 640, 75, name3, desc3, stat3, lvl3)

	s.btnUpgrade1.Draw(dst)
	s.btnUpgrade2.Draw(dst)
	s.btnUpgrade3.Draw(dst)
}

func (s *SanctuaryScene) drawUpgradeRow(dst *ebiten.Image, x, y, w, h float32, title, desc, stat string, lvl int) {
	// Fundo da linha
	vector.DrawFilledRect(dst, x, y, w, h, color.RGBA{R: 35, G: 26, B: 30, A: 230}, false)
	vector.StrokeRect(dst, x, y, w, h, 1.5, color.RGBA{R: 120, G: 90, B: 60, A: 200}, false)

	// Marcadores de Nível (5 pipetas douradas)
	pipX := x + 16
	pipY := y + 14
	for i := 0; i < progression.MaxLevel; i++ {
		pClr := color.RGBA{R: 55, G: 45, B: 40, A: 255}
		if i < lvl {
			pClr = color.RGBA{R: 255, G: 215, B: 60, A: 255}
		}
		vector.DrawFilledRect(dst, pipX+float32(i*18), pipY, 14, 8, pClr, false)
		vector.StrokeRect(dst, pipX+float32(i*18), pipY, 14, 8, 1, color.RGBA{R: 20, G: 15, B: 12, A: 255}, false)
	}

	// Título e Descrição
	ebitenutil.DebugPrintAt(dst, title, int(x)+125, int(y)+10)
	ebitenutil.DebugPrintAt(dst, desc, int(x)+20, int(y)+32)
	ebitenutil.DebugPrintAt(dst, stat, int(x)+20, int(y)+52)
}
