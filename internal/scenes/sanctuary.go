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

// SanctuaryScene gerencia o santuário de evolução para todas as 5 classes.
type SanctuaryScene struct {
	navigator        SceneNavigator
	selectedTab      entities.ClassType
	btnTabWarrior    *ui.Button
	btnTabArcher     *ui.Button
	btnTabMage       *ui.Button
	btnTabShotgunner *ui.Button
	btnTabSheriff    *ui.Button

	btnUpgrade1 *ui.Button
	btnUpgrade2 *ui.Button
	btnUpgrade3 *ui.Button

	btnBack  *ui.Button
	btnPlay  *ui.Button
	soulIcon *ebiten.Image
}

// NewSanctuaryScene inicializa o santuário de upgrades permanentes.
func NewSanctuaryScene(nav SceneNavigator) *SanctuaryScene {
	s := &SanctuaryScene{
		navigator:        nav,
		selectedTab:      entities.ClassWarrior,
		btnTabWarrior:    ui.NewButton(20, 108, 144, 38, "GUERREIRO"),
		btnTabArcher:     ui.NewButton(172, 108, 144, 38, "ARQUEIRO"),
		btnTabMage:       ui.NewButton(324, 108, 144, 38, "MAGO"),
		btnTabShotgunner: ui.NewButton(476, 108, 154, 38, "ATIRADOR PESADO"),
		btnTabSheriff:    ui.NewButton(638, 108, 142, 38, "XERIFE"),

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
	s.btnTabShotgunner.Update()
	s.btnTabSheriff.Update()

	s.btnUpgrade1.Update()
	s.btnUpgrade2.Update()
	s.btnUpgrade3.Update()

	s.btnBack.Update()
	s.btnPlay.Update()

	if s.btnTabWarrior.IsClicked() {
		s.selectedTab = entities.ClassWarrior
		s.refreshUpgradeButtons()
	} else if s.btnTabArcher.IsClicked() {
		s.selectedTab = entities.ClassArcher
		s.refreshUpgradeButtons()
	} else if s.btnTabMage.IsClicked() {
		s.selectedTab = entities.ClassMage
		s.refreshUpgradeButtons()
	} else if s.btnTabShotgunner.IsClicked() {
		s.selectedTab = entities.ClassShotgunner
		s.refreshUpgradeButtons()
	} else if s.btnTabSheriff.IsClicked() {
		s.selectedTab = entities.ClassSheriff
		s.refreshUpgradeButtons()
	}

	if s.btnUpgrade1.IsClicked() {
		s.buyUpgrade(1)
	} else if s.btnUpgrade2.IsClicked() {
		s.buyUpgrade(2)
	} else if s.btnUpgrade3.IsClicked() {
		s.buyUpgrade(3)
	}

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
	case entities.ClassShotgunner:
		if slot == 1 {
			currentLvl = &save.Shotgunner.SpreadPelletsLevel
		} else if slot == 2 {
			currentLvl = &save.Shotgunner.ShrapnelDmgLevel
		} else {
			currentLvl = &save.Shotgunner.ShockwaveLevel
		}
	case entities.ClassSheriff:
		if slot == 1 {
			currentLvl = &save.Sheriff.ReloadSpeedLevel
		} else if slot == 2 {
			currentLvl = &save.Sheriff.CritDamageLevel
		} else {
			currentLvl = &save.Sheriff.ExtraAmmoLevel
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
	case entities.ClassShotgunner:
		lvl1 = save.Shotgunner.SpreadPelletsLevel
		lvl2 = save.Shotgunner.ShrapnelDmgLevel
		lvl3 = save.Shotgunner.ShockwaveLevel
	case entities.ClassSheriff:
		lvl1 = save.Sheriff.ReloadSpeedLevel
		lvl2 = save.Sheriff.CritDamageLevel
		lvl3 = save.Sheriff.ExtraAmmoLevel
	}

	s.formatButton(s.btnUpgrade1, lvl1, save.SoulFragments)
	s.formatButton(s.btnUpgrade2, lvl2, save.SoulFragments)
	s.formatButton(s.btnUpgrade3, lvl3, save.SoulFragments)

	// Destaca a aba selecionada
	s.styleTab(s.btnTabWarrior, s.selectedTab == entities.ClassWarrior)
	s.styleTab(s.btnTabArcher, s.selectedTab == entities.ClassArcher)
	s.styleTab(s.btnTabMage, s.selectedTab == entities.ClassMage)
	s.styleTab(s.btnTabShotgunner, s.selectedTab == entities.ClassShotgunner)
	s.styleTab(s.btnTabSheriff, s.selectedTab == entities.ClassSheriff)
}

func (s *SanctuaryScene) styleTab(btn *ui.Button, active bool) {
	if active {
		btn.BorderColor = color.RGBA{R: 255, G: 215, B: 60, A: 255}
		btn.NormalColor = color.RGBA{R: 78, G: 52, B: 34, A: 255}
	} else {
		btn.BorderColor = color.RGBA{R: 130, G: 100, B: 55, A: 255}
		btn.NormalColor = color.RGBA{R: 48, G: 34, B: 24, A: 255}
	}
}

func (s *SanctuaryScene) formatButton(btn *ui.Button, lvl, availableSouls int) {
	if lvl >= progression.MaxLevel {
		btn.Text = "★ NIVEL MAXIMO ★"
		btn.NormalColor = color.RGBA{R: 40, G: 45, B: 50, A: 255}
		return
	}

	cost := progression.UpgradeCost(lvl)
	if availableSouls >= cost {
		btn.Text = fmt.Sprintf("EVOLUIR (%d Essencias)", cost)
		btn.NormalColor = color.RGBA{R: 48, G: 34, B: 24, A: 255}
	} else {
		btn.Text = fmt.Sprintf("REQUER %d ESSENCIAS", cost)
		btn.NormalColor = color.RGBA{R: 50, G: 30, B: 30, A: 255}
	}
}

// Draw renderiza o santuário com visual medieval.
func (s *SanctuaryScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 20, G: 16, B: 22, A: 255})

	// Título do Santuário
	vector.DrawFilledRect(screen, 180, 20, 440, 68, color.RGBA{R: 40, G: 28, B: 36, A: 250}, false)
	vector.StrokeRect(screen, 180, 20, 440, 68, 2.5, color.RGBA{R: 195, G: 150, B: 80, A: 255}, false)
	ebitenutil.DebugPrintAt(screen, "=== ⚜ SANTUARIO DAS ALMAS ⚜ ===", 280, 35)
	ebitenutil.DebugPrintAt(screen, "Evolua os atributos sagrados de vossos campeoes", 235, 55)

	// Saldo de Fragmentos
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

	// Abas das 5 Classes
	s.btnTabWarrior.Draw(screen)
	s.btnTabArcher.Draw(screen)
	s.btnTabMage.Draw(screen)
	s.btnTabShotgunner.Draw(screen)
	s.btnTabSheriff.Draw(screen)

	// Painel Central com os 3 upgrades da classe selecionada
	s.drawUpgradesPanel(screen)

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
		name1 = "Vitalidade Titanica"
		desc1 = "Aumenta permanentemente a Vida Maxima da classe."
		stat1 = fmt.Sprintf("+%d HP Maximo (Nvl %d/5)", lvl1*25, lvl1)

		lvl2 = save.Warrior.DamageLevel
		name2 = "Lamina Dilacerante"
		desc2 = "Aumenta o dano de impacto e o alcance do golpe Melee."
		stat2 = fmt.Sprintf("+%d%% Dano | +%dpx Alcance (Nvl %d/5)", lvl2*12, lvl2*6, lvl2)

		lvl3 = save.Warrior.ShieldLevel
		name3 = "Bastiao Divino (Escudo)"
		desc3 = "Reduz recarga da Postura de Escudo e amplia reflexao."
		stat3 = fmt.Sprintf("-%.1fs Recarga | +%.0f%% Dano Refletido (Nvl %d/5)", float64(lvl3)*0.75, float64(lvl3)*35, lvl3)

	case entities.ClassArcher:
		lvl1 = save.Archer.SpeedLevel
		name1 = "Passo do Vento"
		desc1 = "Aumenta a velocidade de movimentacao base do arqueiro."
		stat1 = fmt.Sprintf("+%.2f Velocidade (Nvl %d/5)", float64(lvl1)*0.35, lvl1)

		lvl2 = save.Archer.AttackSpeedLvl
		name2 = "Disparo Voraz"
		desc2 = "Aumenta a cadencia de tiros e a velocidade das flechas."
		stat2 = fmt.Sprintf("+%d%% Cadencia | +%.1f Vel. Flechas (Nvl %d/5)", lvl2*10, float64(lvl2)*1.2, lvl2)

		lvl3 = save.Archer.ArrowRainLevel
		name3 = "Chuva Devastadora"
		desc3 = "Aumenta o raio de impacto e dano da Chuva de Flechas."
		stat3 = fmt.Sprintf("+%d%% Area | +%d%% Dano da Salva (Nvl %d/5)", lvl3*15, lvl3*20, lvl3)

	case entities.ClassMage:
		lvl1 = save.Mage.ManaLevel
		name1 = "Fonte Mistica"
		desc1 = "Aumenta a Mana Maxima e a taxa de regeneracao continua."
		stat1 = fmt.Sprintf("+%d Mana | +%.1f MP/s (Nvl %d/5)", lvl1*20, float64(lvl1)*3.5, lvl1)

		lvl2 = save.Mage.SpellPowerLevel
		name2 = "Canalizacao Arcana"
		desc2 = "Eleva o dano base de todas as magias e orbes."
		stat2 = fmt.Sprintf("+%d%% Dano Magico (Nvl %d/5)", lvl2*15, lvl2)

		lvl3 = save.Mage.FirePillarLevel
		name3 = "Cataclismo de Chamas"
		desc3 = "Amplia o raio e duracao ardente do Pilar de Fogo."
		stat3 = fmt.Sprintf("+%d%% Raio | +%d Frames Fogo (Nvl %d/5)", lvl3*15, lvl3*20, lvl3)

	case entities.ClassShotgunner:
		lvl1 = save.Shotgunner.SpreadPelletsLevel
		name1 = "Cano Reforcado (Dispersao)"
		desc1 = "Adiciona balotes ao cone de tiro e melhora a precisao."
		stat1 = fmt.Sprintf("+%d Balotes extras por tiro (Nvl %d/5)", lvl1, lvl1)

		lvl2 = save.Shotgunner.ShrapnelDmgLevel
		name2 = "Polvora Pesada (Estilhacos)"
		desc2 = "Aumenta o dano dos estilhacos e a forca de recuo."
		stat2 = fmt.Sprintf("+%d%% Dano | +%.1f Empurrao (Nvl %d/5)", lvl2*15, float64(lvl2)*0.8, lvl2)

		lvl3 = save.Shotgunner.ShockwaveLevel
		name3 = "Detonacao Sismica"
		desc3 = "Amplia o raio e dano do Choque Sismico 360 graus."
		stat3 = fmt.Sprintf("+%dpx Raio | +%d%% Dano Choque (Nvl %d/5)", lvl3*20, lvl3*20, lvl3)

	case entities.ClassSheriff:
		lvl1 = save.Sheriff.ReloadSpeedLevel
		name1 = "Maos Velozes (Recarga)"
		desc1 = "Reduz drasticamente o tempo de recarga do tambor."
		stat1 = fmt.Sprintf("-%.2fs Tempo de Recarga (Nvl %d/5)", float64(lvl1)*0.15, lvl1)

		lvl2 = save.Sheriff.CritDamageLevel
		name2 = "Olho de Aguia (Critico)"
		desc2 = "Aumenta a chance de tiro critico e o dano multiplicado."
		stat2 = fmt.Sprintf("+%d%% Chance | +%d%% Dano Critico (Nvl %d/5)", lvl2*8, lvl2*22, lvl2)

		lvl3 = save.Sheriff.ExtraAmmoLevel
		name3 = "Tambor Estendido"
		desc3 = "Adiciona mais capacidade de balas no tambor do revolver."
		stat3 = fmt.Sprintf("+%d Balas no Tambor (Nvl %d/5)", lvl3, lvl3)
	}

	s.drawUpgradeRow(dst, 80, 180, 640, 75, name1, desc1, stat1, lvl1)
	s.drawUpgradeRow(dst, 80, 275, 640, 75, name2, desc2, stat2, lvl2)
	s.drawUpgradeRow(dst, 80, 370, 640, 75, name3, desc3, stat3, lvl3)

	s.btnUpgrade1.Draw(dst)
	s.btnUpgrade2.Draw(dst)
	s.btnUpgrade3.Draw(dst)
}

func (s *SanctuaryScene) drawUpgradeRow(dst *ebiten.Image, x, y, w, h float32, title, desc, stat string, lvl int) {
	vector.DrawFilledRect(dst, x, y, w, h, color.RGBA{R: 35, G: 26, B: 30, A: 230}, false)
	vector.StrokeRect(dst, x, y, w, h, 1.5, color.RGBA{R: 120, G: 90, B: 60, A: 200}, false)

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

	ebitenutil.DebugPrintAt(dst, title, int(x)+125, int(y)+10)
	ebitenutil.DebugPrintAt(dst, desc, int(x)+20, int(y)+32)
	ebitenutil.DebugPrintAt(dst, stat, int(x)+20, int(y)+52)
}
