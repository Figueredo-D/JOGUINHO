package scenes

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"joguinho/internal/entities"
	"joguinho/internal/progression"
	"joguinho/internal/sprites"
	"joguinho/internal/ui"
)

// FloatingText representa um número de dano ou texto flutuante temporário.
type FloatingText struct {
	Text     string
	X, Y     float64
	Color    color.RGBA
	Lifetime int
}

// GameplayScene gerencia toda a partida ativa do jogo, habilidades especiais e o menu de pausa (ESC).
type GameplayScene struct {
	navigator     SceneNavigator
	player        *entities.Player
	enemies       []*entities.Enemy
	projectiles   []*entities.Projectile
	slashes       []*entities.Slash
	drops         []*entities.Drop
	arrowRains    []*entities.ArrowRainEffect
	firePillars   []*entities.FirePillarEffect
	floatingTexts []*FloatingText

	hpBar      *ui.ProgressBar
	manaBar    *ui.ProgressBar
	specialBar *ui.ProgressBar
	soulIcon   *ebiten.Image

	nextEnemyID    uint64
	spawnTimer     int
	spawnDelay     int
	wave           int
	killsInWave    int
	killsToWave    int
	score          int
	totalKills     int
	soulsCollected int
	survivalFrames int

	// Menu de Pausa (ESC)
	isPaused    bool
	btnResume   *ui.Button
	btnOpts     *ui.Button
	btnMainMenu *ui.Button
}

// NewGameplayScene inicia uma nova partida para a classe informada aplicando os meta upgrades.
func NewGameplayScene(nav SceneNavigator, class entities.ClassType) *GameplayScene {
	p := entities.NewPlayer(class, 400, 320)

	hpBar := ui.NewProgressBar(20, 10, 190, 16, "HP", p.MaxHP, color.RGBA{R: 220, G: 50, B: 50, A: 255})
	var manaBar *ui.ProgressBar
	if p.MaxMana > 0 {
		manaBar = ui.NewProgressBar(20, 29, 150, 13, "MP", p.MaxMana, color.RGBA{R: 50, G: 120, B: 240, A: 255})
	}

	specialLabel := "Especial [E/Dir]"
	switch class {
	case entities.ClassWarrior:
		specialLabel = "Escudo [E/Dir]"
	case entities.ClassArcher:
		specialLabel = "Chuva [E/Dir]"
	case entities.ClassMage:
		specialLabel = "Pilar [E/Dir]"
	}

	specialBar := ui.NewProgressBar(220, 10, 150, 16, specialLabel, float64(p.SpecialMaxCD), color.RGBA{R: 245, G: 190, B: 50, A: 255})
	specialBar.ShowFraction = false

	btnW := 280.0
	btnH := 44.0
	centerX := (800.0 - btnW) / 2

	fsLabel := "⚙  TELA CHEIA: DESATIVADO  ⚙"
	if nav.IsFullscreen() {
		fsLabel = "⚙  TELA CHEIA: ATIVADO  ⚙"
	}

	return &GameplayScene{
		navigator:      nav,
		player:         p,
		hpBar:          hpBar,
		manaBar:        manaBar,
		specialBar:     specialBar,
		soulIcon:       sprites.GetSoulFragmentSprite(),
		spawnTimer:     0,
		spawnDelay:     72,
		wave:           1,
		killsInWave:    0,
		killsToWave:    12,
		score:          0,
		soulsCollected: 0,
		survivalFrames: 0,
		isPaused:       false,
		btnResume:      ui.NewButton(centerX, 250, btnW, btnH, "⚔  RETOMAR BATALHA  ⚔"),
		btnOpts:        ui.NewButton(centerX, 305, btnW, btnH, fsLabel),
		btnMainMenu:    ui.NewButton(centerX, 360, btnW, btnH, "🏰  VOLTAR AO MENU  🏰"),
	}
}

// Update executa o ciclo de lógica ou o menu de pausa ao pressionar ESC.
func (g *GameplayScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.isPaused = !g.isPaused
	}

	if g.isPaused {
		g.btnResume.Update()
		g.btnOpts.Update()
		g.btnMainMenu.Update()

		if g.btnResume.IsClicked() {
			g.isPaused = false
		} else if g.btnOpts.IsClicked() {
			g.navigator.ToggleFullscreen()
			if g.navigator.IsFullscreen() {
				g.btnOpts.Text = "⚙  TELA CHEIA: ATIVADO  ⚙"
			} else {
				g.btnOpts.Text = "⚙  TELA CHEIA: DESATIVADO  ⚙"
			}
		} else if g.btnMainMenu.IsClicked() {
			g.navigator.SwitchToMenu()
		}

		return nil
	}

	g.survivalFrames++

	// 1. Atualiza Player
	g.player.Update(800, 600)

	// Atualiza barras do HUD
	g.hpBar.SetValues(g.player.HP, g.player.MaxHP)
	if g.manaBar != nil {
		g.manaBar.SetValues(g.player.Mana, g.player.MaxMana)
	}

	// Atualiza barra da habilidade especial
	specCurrent := float64(g.player.SpecialMaxCD - g.player.SpecialCD)
	if specCurrent < 0 {
		specCurrent = 0
	}
	g.specialBar.SetValues(specCurrent, float64(g.player.SpecialMaxCD))

	// Verifica se jogador morreu
	if g.player.HP <= 0 {
		g.handleGameOver()
		return nil
	}

	// 2. Ataque Básico (Clique Esquerdo ou Tecla Espaço)
	attackTriggered := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || ebiten.IsKeyPressed(ebiten.KeySpace)
	if attackTriggered && g.player.CanAttack() {
		g.handlePlayerAttack()
	}

	// 3. Habilidade Especial Ativa (Tecla E ou Botão Direito do Mouse)
	specialTriggered := inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	if specialTriggered && g.player.CanUseSpecial() {
		g.handleSpecialSkill()
	}

	// 4. Spawner de Inimigos dinâmico
	g.spawnTimer++
	if g.spawnTimer >= g.spawnDelay {
		g.spawnTimer = 0
		g.spawnEnemy()
	}

	// 5. Atualiza Slashes (Melee do Guerreiro)
	for i := len(g.slashes) - 1; i >= 0; i-- {
		s := g.slashes[i]
		s.Update()
		if !s.Active {
			g.slashes = append(g.slashes[:i], g.slashes[i+1:]...)
		}
	}

	// 6. Atualiza Habilidades Especiais em Área (Chuva de Flechas e Pilar de Fogo)
	for i := len(g.arrowRains) - 1; i >= 0; i-- {
		ar := g.arrowRains[i]
		ar.Update()
		if !ar.Active {
			g.arrowRains = append(g.arrowRains[:i], g.arrowRains[i+1:]...)
		}
	}

	for i := len(g.firePillars) - 1; i >= 0; i-- {
		fp := g.firePillars[i]
		fp.Update()
		if !fp.Active {
			g.firePillars = append(g.firePillars[:i], g.firePillars[i+1:]...)
		}
	}

	// 7. Atualiza Projéteis e Reflexão por Escudo
	pBounds := g.player.Bounds()
	for i := len(g.projectiles) - 1; i >= 0; i-- {
		proj := g.projectiles[i]
		proj.Update()

		// Remove projéteis fora da tela ou expirados
		if !proj.Active || proj.X < 0 || proj.X > 800 || proj.Y < 45 || proj.Y > 600 {
			g.projectiles = append(g.projectiles[:i], g.projectiles[i+1:]...)
			continue
		}

		// Se for projétil inimigo: checa se é refletido pela Postura de Escudo do Guerreiro!
		if proj.IsEnemy && g.player.ClassType == entities.ClassWarrior && g.player.Shield.Active {
			if g.player.CheckShieldReflect(proj) {
				g.addFloatingText("REFLETIDO!", proj.X, proj.Y-12, color.RGBA{R: 255, G: 235, B: 100, A: 255})
				continue
			}
		}

		// Projétil inimigo colidindo com o jogador
		if proj.IsEnemy && proj.Bounds().Intersects(pBounds) {
			proj.Active = false
			if g.player.TakeDamage(proj.Damage) {
				g.addFloatingText(fmt.Sprintf("-%.0f", proj.Damage), g.player.X, g.player.Y-10, color.RGBA{R: 255, G: 80, B: 80, A: 255})
			}
			g.projectiles = append(g.projectiles[:i], g.projectiles[i+1:]...)
		}
	}

	// 8. Atualiza Inimigos & Colisões com ataques e habilidades
	px, py := g.player.Center()
	for i := len(g.enemies) - 1; i >= 0; i-- {
		e := g.enemies[i]
		newEnemyProj := e.Update(px, py)
		if newEnemyProj != nil {
			g.projectiles = append(g.projectiles, newEnemyProj)
		}

		eBounds := e.Bounds()

		// Colisão com ataque Melee (Slash do Guerreiro)
		for _, slash := range g.slashes {
			if slash.CollidesWithEnemy(e.ID, eBounds) {
				dmg := slash.Damage * g.player.GetDamageMultiplier()
				killed := e.TakeDamage(dmg)
				g.addFloatingText(fmt.Sprintf("-%.0f", dmg), e.X, e.Y-6, color.RGBA{R: 255, G: 210, B: 50, A: 255})
				if killed {
					g.onEnemyDefeated(e)
				}
			}
		}

		// Colisão com Projéteis do Jogador (ou Projéteis Refletidos pelo Escudo!)
		for _, proj := range g.projectiles {
			if proj.Active && !proj.IsEnemy && proj.Bounds().Intersects(eBounds) {
				proj.Active = false
				dmg := proj.Damage * g.player.GetDamageMultiplier()
				killed := e.TakeDamage(dmg)
				g.addFloatingText(fmt.Sprintf("-%.0f", dmg), e.X, e.Y-6, color.RGBA{R: 255, G: 255, B: 110, A: 255})
				if killed {
					g.onEnemyDefeated(e)
				}
				break
			}
		}

		// Dano contínuo da Chuva de Flechas (Arqueiro)
		for _, ar := range g.arrowRains {
			if ar.CollidesWith(eBounds) && g.survivalFrames%10 == 0 {
				dmg := ar.DamagePerTick * g.player.GetDamageMultiplier()
				killed := e.TakeDamage(dmg)
				g.addFloatingText(fmt.Sprintf("-%.0f", dmg), e.X, e.Y-4, color.RGBA{R: 120, G: 240, B: 140, A: 255})
				// Aplica desaceleração temporária
				e.Speed *= 0.65
				if killed {
					g.onEnemyDefeated(e)
				}
			}
		}

		// Dano contínuo e intenso do Pilar de Fogo (Mago)
		for _, fp := range g.firePillars {
			if fp.CollidesWith(eBounds) && g.survivalFrames%8 == 0 {
				dmg := fp.DamagePerTick * g.player.GetDamageMultiplier()
				killed := e.TakeDamage(dmg)
				g.addFloatingText(fmt.Sprintf("-%.0f", dmg), e.X, e.Y-4, color.RGBA{R: 255, G: 140, B: 30, A: 255})
				if killed {
					g.onEnemyDefeated(e)
				}
			}
		}

		// Colisão Inimigo vs Jogador (Dano por contato)
		if e.Alive && eBounds.Intersects(pBounds) {
			if e.CanDamagePlayer() {
				if g.player.TakeDamage(e.Damage) {
					e.ResetAttackCD()
					g.addFloatingText(fmt.Sprintf("-%.0f", e.Damage), g.player.X, g.player.Y-10, color.RGBA{R: 255, G: 60, B: 60, A: 255})
				}
			}
		}

		// Remove inimigos mortos
		if !e.Alive {
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
		}
	}

	// 9. Atualiza e Coleta Drops no chão
	for i := len(g.drops) - 1; i >= 0; i-- {
		d := g.drops[i]
		d.Update()
		if !d.Active {
			g.drops = append(g.drops[:i], g.drops[i+1:]...)
			continue
		}

		if d.Bounds().Intersects(pBounds) {
			d.Active = false
			feedback, clr := g.player.ApplyDrop(d)
			g.addFloatingText(feedback, g.player.X-20, g.player.Y-20, clr)

			if d.Type == entities.DropSoulFragment {
				g.soulsCollected += d.SoulAmount
			} else if d.Type == entities.DropChaosSeed {
				g.spawnChaosAmbush()
			}

			g.drops = append(g.drops[:i], g.drops[i+1:]...)
		}
	}

	// 10. Atualiza Textos Flutuantes
	for i := len(g.floatingTexts) - 1; i >= 0; i-- {
		ft := g.floatingTexts[i]
		ft.Y -= 0.65
		ft.Lifetime--
		if ft.Lifetime <= 0 {
			g.floatingTexts = append(g.floatingTexts[:i], g.floatingTexts[i+1:]...)
		}
	}

	return nil
}

// handlePlayerAttack executa o tipo de golpe padrão conforme a classe.
func (g *GameplayScene) handlePlayerAttack() {
	px, py := g.player.Center()
	mx, my := ebiten.CursorPosition()
	fmx, fmy := float64(mx), float64(my)

	switch g.player.ClassType {
	case entities.ClassWarrior:
		slash := entities.NewSlash(px, py, fmx, fmy, g.player.SlashRange, 68)
		g.slashes = append(g.slashes, slash)
		g.player.ResetAttackCD()

	case entities.ClassArcher:
		dx := fmx - px
		dy := fmy - py
		dist := math.Hypot(dx, dy)
		if dist > 0.001 {
			vx := (dx / dist) * g.player.ArrowSpeed
			vy := (dy / dist) * g.player.ArrowSpeed
			proj := entities.NewProjectile(px, py, vx, vy, 3.5, 28, color.RGBA{R: 120, G: 255, B: 120, A: 255}, 85, entities.ProjectileKindArrow, false)
			g.projectiles = append(g.projectiles, proj)
			g.player.ResetAttackCD()
		}

	case entities.ClassMage:
		dx := fmx - px
		dy := fmy - py
		dist := math.Hypot(dx, dy)
		if dist > 0.001 {
			manaCost := 20.0
			speed := 6.8
			vx := (dx / dist) * speed
			vy := (dy / dist) * speed

			if g.player.Mana >= manaCost {
				g.player.Mana -= manaCost
				dmg := 55.0 * g.player.SpellMultiplier
				proj := entities.NewProjectile(px, py, vx, vy, 8.5, dmg, color.RGBA{R: 70, G: 160, B: 255, A: 255}, 90, entities.ProjectileKindMagic, false)
				g.projectiles = append(g.projectiles, proj)
			} else {
				dmg := 15.0 * g.player.SpellMultiplier
				proj := entities.NewProjectile(px, py, vx*0.9, vy*0.9, 4.0, dmg, color.RGBA{R: 160, G: 190, B: 255, A: 200}, 70, entities.ProjectileKindMagicBasic, false)
				g.projectiles = append(g.projectiles, proj)
			}
			g.player.ResetAttackCD()
		}
	}
}

// handleSpecialSkill executa a habilidade especial ativa com base na classe do herói.
func (g *GameplayScene) handleSpecialSkill() {
	mx, my := ebiten.CursorPosition()
	fmx, fmy := float64(mx), float64(my)

	switch g.player.ClassType {
	case entities.ClassWarrior:
		// Guerreiro: Postura de Escudo (Reflexão e Barreira)
		g.player.ActivateShield()
		g.addFloatingText("POSTURA DE ESCUDO!", g.player.X-20, g.player.Y-25, color.RGBA{R: 255, G: 220, B: 80, A: 255})

	case entities.ClassArcher:
		// Arqueiro: Chuva de Flechas em Área
		rain := entities.NewArrowRain(fmx, fmy, g.player.ArrowRainRadius, g.player.ArrowRainDamage, 130)
		g.arrowRains = append(g.arrowRains, rain)
		g.player.SpecialCD = g.player.SpecialMaxCD
		g.addFloatingText("CHUVA DE FLECHAS!", fmx-40, fmy-20, color.RGBA{R: 120, G: 240, B: 140, A: 255})

	case entities.ClassMage:
		// Mago: Pilar de Fogo (Erupção Mágica)
		if g.player.Mana >= g.player.FirePillarManaCost {
			g.player.Mana -= g.player.FirePillarManaCost
			pillar := entities.NewFirePillar(fmx, fmy, g.player.FirePillarRadius, g.player.FirePillarDamage, g.player.FirePillarDuration)
			g.firePillars = append(g.firePillars, pillar)
			g.player.SpecialCD = g.player.SpecialMaxCD
			g.addFloatingText("PILAR DE FOGO!", fmx-35, fmy-25, color.RGBA{R: 255, G: 140, B: 40, A: 255})
		}
	}
}

// spawnEnemy cria um inimigo nas bordas da tela baseado na composição da onda atual.
func (g *GameplayScene) spawnEnemy() {
	g.nextEnemyID++
	edge := rand.Intn(4)
	var sx, sy float64

	switch edge {
	case 0:
		sx = rand.Float64() * 800
		sy = 55
	case 1:
		sx = 785
		sy = 60 + rand.Float64()*525
	case 2:
		sx = rand.Float64() * 800
		sy = 585
	case 3:
		sx = 15
		sy = 60 + rand.Float64()*525
	}

	waveMult := 1.0 + float64(g.wave-1)*0.18
	eType := entities.EnemyTypeCommon
	r := rand.Float64()

	switch {
	case g.wave >= 4:
		if r < 0.22 {
			eType = entities.EnemyTypeTank
		} else if r < 0.55 {
			eType = entities.EnemyTypeRunner
		} else if r < 0.80 {
			eType = entities.EnemyTypeShooter
		} else {
			eType = entities.EnemyTypeCommon
		}

	case g.wave == 3:
		if r < 0.15 {
			eType = entities.EnemyTypeTank
		} else if r < 0.45 {
			eType = entities.EnemyTypeShooter
		} else if r < 0.75 {
			eType = entities.EnemyTypeRunner
		} else {
			eType = entities.EnemyTypeCommon
		}

	case g.wave == 2:
		if r < 0.40 {
			eType = entities.EnemyTypeRunner
		} else if r < 0.65 {
			eType = entities.EnemyTypeShooter
		} else {
			eType = entities.EnemyTypeCommon
		}

	default:
		if r < 0.20 {
			eType = entities.EnemyTypeRunner
		} else {
			eType = entities.EnemyTypeCommon
		}
	}

	enemy := entities.NewEnemy(g.nextEnemyID, eType, sx, sy, waveMult)
	g.enemies = append(g.enemies, enemy)
}

// spawnChaosAmbush invoca 4 corredores rápidos ao redor do jogador pela Semente do Caos.
func (g *GameplayScene) spawnChaosAmbush() {
	px, py := g.player.Center()
	angles := []float64{0, math.Pi / 2, math.Pi, 3 * math.Pi / 2}
	dist := 140.0

	for _, angle := range angles {
		g.nextEnemyID++
		sx := px + dist*math.Cos(angle)
		sy := py + dist*math.Sin(angle)
		if sx < 20 {
			sx = 20
		}
		if sx > 780 {
			sx = 780
		}
		if sy < 60 {
			sy = 60
		}
		if sy > 580 {
			sy = 580
		}

		waveMult := 1.0 + float64(g.wave-1)*0.15
		ambusher := entities.NewEnemy(g.nextEnemyID, entities.EnemyTypeRunner, sx, sy, waveMult)
		g.enemies = append(g.enemies, ambusher)
	}
}

// onEnemyDefeated contabiliza pontuação, drops e avanço de onda.
func (g *GameplayScene) onEnemyDefeated(e *entities.Enemy) {
	g.score += e.ScoreValue * g.wave
	g.totalKills++
	g.killsInWave++

	// 1. Chance de Drop de Fragmento de Alma Meta (~45% para comuns/corredores/atiradores, 100% para tanques)
	soulDropChance := 0.45
	soulAmount := 1
	if e.Type == entities.EnemyTypeTank {
		soulDropChance = 1.0
		soulAmount = 3 + rand.Intn(3) // 3 a 5 fragmentos!
	} else if e.Type == entities.EnemyTypeShooter || e.Type == entities.EnemyTypeRunner {
		soulAmount = 1 + rand.Intn(2)
	}

	if rand.Float64() < soulDropChance {
		soulDrop := entities.NewSoulDrop(e.X+e.Width/2, e.Y+e.Height/2, soulAmount)
		g.drops = append(g.drops, soulDrop)
	}

	// 2. Chance de Drop de Itens Regulares (~30%)
	if rand.Float64() < 0.30 {
		g.rollRegularDrop(e.X+e.Width/2, e.Y+e.Height/2)
	}

	// Transição de onda
	if g.killsInWave >= g.killsToWave {
		g.wave++
		g.killsInWave = 0
		g.killsToWave += 6

		if g.spawnDelay > 20 {
			g.spawnDelay -= 7
		}

		g.player.HP += g.player.MaxHP * 0.25
		if g.player.HP > g.player.MaxHP {
			g.player.HP = g.player.MaxHP
		}
		if g.player.MaxMana > 0 {
			g.player.Mana = g.player.MaxMana
		}

		g.addFloatingText(fmt.Sprintf("=== ONDA %d! ===", g.wave), 340, 180, color.RGBA{R: 255, G: 230, B: 60, A: 255})
	}
}

// rollRegularDrop sorteia um item benigno ou maligno tradicional.
func (g *GameplayScene) rollRegularDrop(x, y float64) {
	isBuff := rand.Float64() < 0.62
	var dType entities.DropType

	if isBuff {
		buffs := []entities.DropType{
			entities.DropHealthPotion,
			entities.DropManaPotion,
			entities.DropAgilityStone,
			entities.DropSpeedBoots,
			entities.DropMasteryTome,
		}
		dType = buffs[rand.Intn(len(buffs))]
	} else {
		debuffs := []entities.DropType{
			entities.DropCursedChalice,
			entities.DropPoisonVial,
			entities.DropFatigueMist,
			entities.DropChaosSeed,
		}
		dType = debuffs[rand.Intn(len(debuffs))]
	}

	drop := entities.NewDrop(dType, x, y)
	g.drops = append(g.drops, drop)
}

func (g *GameplayScene) handleGameOver() {
	// Bônus de conversão da pontuação em fragmentos adicionais
	bonusSouls := g.score / 500
	totalSouls := g.soulsCollected + bonusSouls

	// Persiste no save
	progression.AddSoulFragments(totalSouls)
	save := progression.GetSave()
	save.TotalRuns++
	save.TotalKills += g.totalKills
	_ = progression.Save()

	g.navigator.SwitchToGameOver(g.score, g.wave, g.player.ClassType, totalSouls)
}

// addFloatingText cria uma notificação visual flutuante.
func (g *GameplayScene) addFloatingText(text string, x, y float64, clr color.RGBA) {
	g.floatingTexts = append(g.floatingTexts, &FloatingText{
		Text:     text,
		X:        x,
		Y:        y,
		Color:    clr,
		Lifetime: 50,
	})
}

// Draw renderiza o jogo, chão, áreas de feitiços, entidades, HUD e o menu de pausa.
func (g *GameplayScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 20, G: 23, B: 32, A: 255})

	// Piso da masmorra em blocos de pedra
	gridClr := color.RGBA{R: 30, G: 35, B: 48, A: 255}
	for x := 0; x < 800; x += 40 {
		vector.StrokeLine(screen, float32(x), 52, float32(x), 600, 1, gridClr, false)
	}
	for y := 52; y < 600; y += 40 {
		vector.StrokeLine(screen, 0, float32(y), 800, float32(y), 1, gridClr, false)
	}

	vector.StrokeRect(screen, 6, 52, 788, 542, 3, color.RGBA{R: 70, G: 80, B: 110, A: 255}, false)

	// 1. Áreas de Habilidades Especiais (Chuva de Flechas e Pilar de Fogo) no chão
	for _, ar := range g.arrowRains {
		ar.Draw(screen)
	}
	for _, fp := range g.firePillars {
		fp.Draw(screen)
	}

	// 2. Drops no Chão
	for _, d := range g.drops {
		d.Draw(screen)
	}

	// 3. Inimigos
	for _, e := range g.enemies {
		e.Draw(screen)
	}

	// 4. Projéteis
	for _, p := range g.projectiles {
		p.Draw(screen)
	}

	// 5. Slashes (Melee)
	for _, s := range g.slashes {
		s.Draw(screen)
	}

	// 6. Jogador (inclui a barreira do escudo se ativa)
	g.player.Draw(screen)

	// 7. Textos Flutuantes
	for _, ft := range g.floatingTexts {
		ebitenutil.DebugPrintAt(screen, ft.Text, int(ft.X), int(ft.Y))
	}

	// 8. HUD Superior Completo
	g.drawHUD(screen)

	// 9. Menu de Pausa Medieval sobreposto
	if g.isPaused {
		g.drawPauseMenu(screen)
	}
}

func (g *GameplayScene) drawHUD(screen *ebiten.Image) {
	// Painel superior
	vector.DrawFilledRect(screen, 0, 0, 800, 52, color.RGBA{R: 14, G: 16, B: 24, A: 252}, false)
	vector.StrokeLine(screen, 0, 52, 800, 52, 2, color.RGBA{R: 65, G: 80, B: 115, A: 255}, false)

	// Barras de HP e Mana
	g.hpBar.Draw(screen)
	if g.manaBar != nil {
		g.manaBar.Draw(screen)
	}

	// Barra da Habilidade Especial Ativa
	g.specialBar.Draw(screen)
	if g.player.SpecialCD <= 0 {
		ebitenutil.DebugPrintAt(screen, "PRONTO!", 315, 12)
	} else {
		secs := float64(g.player.SpecialCD) / 60.0
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.1fs", secs), 325, 12)
	}

	// Informações de Onda
	waveInfo := fmt.Sprintf("ONDA %d (Faltam %d)", g.wave, g.killsToWave-g.killsInWave)
	ebitenutil.DebugPrintAt(screen, waveInfo, 385, 12)

	// Contador de Fragmentos de Alma coletados na run
	if g.soulIcon != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1.5, 1.5)
		op.GeoM.Translate(530, 8)
		screen.DrawImage(g.soulIcon, op)
	}
	soulsStr := fmt.Sprintf("+%d Almas", g.soulsCollected)
	ebitenutil.DebugPrintAt(screen, soulsStr, 555, 12)

	// Tempo e Score
	totalSeconds := g.survivalFrames / 60
	mins := totalSeconds / 60
	secs := totalSeconds % 60
	timeStr := fmt.Sprintf("%02d:%02d", mins, secs)
	scoreStr := fmt.Sprintf("Score: %d", g.score)

	ebitenutil.DebugPrintAt(screen, timeStr, 640, 12)
	ebitenutil.DebugPrintAt(screen, scoreStr, 700, 12)

	// Linha inferior do HUD para status e dicas
	effects := g.player.GetActiveEffects()
	if len(effects) > 0 {
		effX := 220
		for _, eff := range effects {
			secsLeft := (eff.Remaining + 59) / 60
			label := fmt.Sprintf("%s (%ds)", eff.Name, secsLeft)

			badgeW := float32(len(label)*6 + 12)
			vector.DrawFilledRect(screen, float32(effX), 32, badgeW, 16, color.RGBA{R: 28, G: 32, B: 45, A: 220}, false)
			vector.StrokeRect(screen, float32(effX), 32, badgeW, 16, 1, eff.Color, false)

			ebitenutil.DebugPrintAt(screen, label, effX+6, 34)
			effX += int(badgeW) + 8
		}
	} else {
		ebitenutil.DebugPrintAt(screen, "[E ou Botão Direito]: Usar Habilidade Especial  |  [ESC]: Pausar", 220, 34)
	}
}

// drawPauseMenu renderiza a janela modal de pausa medieval.
func (g *GameplayScene) drawPauseMenu(dst *ebiten.Image) {
	vector.DrawFilledRect(dst, 0, 0, 800, 600, color.RGBA{R: 10, G: 8, B: 6, A: 215}, false)

	panelW := float32(480)
	panelH := float32(350)
	px := float32((800 - panelW) / 2)
	py := float32((600 - panelH) / 2)

	vector.DrawFilledRect(dst, px+5, py+6, panelW, panelH, color.RGBA{R: 5, G: 4, B: 3, A: 230}, false)
	vector.DrawFilledRect(dst, px, py, panelW, panelH, color.RGBA{R: 44, G: 30, B: 20, A: 255}, false)

	vector.StrokeLine(dst, px+8, py+65, px+panelW-8, py+65, 1.5, color.RGBA{R: 28, G: 18, B: 12, A: 160}, false)
	vector.StrokeLine(dst, px+8, py+235, px+panelW-8, py+235, 1.5, color.RGBA{R: 28, G: 18, B: 12, A: 160}, false)

	vector.StrokeRect(dst, px, py, panelW, panelH, 3.0, color.RGBA{R: 190, G: 150, B: 75, A: 255}, false)
	vector.StrokeRect(dst, px+5, py+5, panelW-10, panelH-10, 1.2, color.RGBA{R: 90, G: 70, B: 40, A: 200}, false)

	drawRivet(dst, px+12, py+12, color.RGBA{R: 245, G: 205, B: 100, A: 255})
	drawRivet(dst, px+panelW-12, py+12, color.RGBA{R: 245, G: 205, B: 100, A: 255})
	drawRivet(dst, px+12, py+panelH-12, color.RGBA{R: 245, G: 205, B: 100, A: 255})
	drawRivet(dst, px+panelW-12, py+panelH-12, color.RGBA{R: 245, G: 205, B: 100, A: 255})

	t1 := "=== ⚜ BATALHA PAUSADA ⚜ ==="
	t2 := "O reino aguarda vosso retorno ao combate..."
	ebitenutil.DebugPrintAt(dst, t1, int(px)+135, int(py)+22)
	ebitenutil.DebugPrintAt(dst, t2, int(px)+100, int(py)+42)

	totalSeconds := g.survivalFrames / 60
	mins := totalSeconds / 60
	secs := totalSeconds % 60
	statsStr := fmt.Sprintf("Classe: %s  |  Onda: %d  |  Tempo: %02d:%02d  |  Almas: %d",
		g.player.ClassType, g.wave, mins, secs, g.soulsCollected)
	ebitenutil.DebugPrintAt(dst, statsStr, int(px)+35, int(py)+82)

	g.btnResume.Draw(dst)
	g.btnOpts.Draw(dst)
	g.btnMainMenu.Draw(dst)

	ebitenutil.DebugPrintAt(dst, "Pressione [ESC] novamente para despausar e retornar", int(px)+75, int(py+panelH)-26)
}
