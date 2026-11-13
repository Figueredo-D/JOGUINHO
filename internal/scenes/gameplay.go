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

	"joguinho/internal/audio"
	"joguinho/internal/collision"
	"joguinho/internal/entities"
	"joguinho/internal/progression"
	"joguinho/internal/sprites"
	"joguinho/internal/ui"
	"joguinho/internal/world"
)

const TotalCampaignWaves = 15

// FloatingText representa um número de dano ou texto flutuante temporário.
type FloatingText struct {
	Text     string
	X, Y     float64
	Color    color.RGBA
	Lifetime int
}

// LightningArc representa um raio elétrico saltando entre dois pontos (Faísca Eletrostática).
type LightningArc struct {
	X1, Y1 float64
	X2, Y2 float64
	Life   int
}

// GameplayScene gerencia a expedição de 15 ondas, game feel (shake/hitstop), áudio, clima e pausa.
type GameplayScene struct {
	navigator     SceneNavigator
	arena         *world.Arena
	weather       *world.WeatherSystem
	worldCanvas   *ebiten.Image
	player        *entities.Player
	enemies       []*entities.Enemy
	boss          *entities.Boss
	projectiles   []*entities.Projectile
	slashes       []*entities.Slash
	drops         []*entities.Drop
	arrowRains    []*entities.ArrowRainEffect
	firePillars   []*entities.FirePillarEffect
	particles     []*entities.Particle
	shockwaves    []*entities.ShockwaveEffect
	lightningArcs []*LightningArc
	floatingTexts []*FloatingText

	hpBar      *ui.ProgressBar
	manaBar    *ui.ProgressBar
	specialBar *ui.ProgressBar
	dashBar    *ui.ProgressBar
	bossBar    *ui.ProgressBar
	soulIcon   *ebiten.Image

	nextEnemyID      uint64
	spawnTimer       int
	spawnDelay       int
	spawnBatchCnt    int
	wave             int
	killsInWave      int
	killsToWave      int
	score            int
	totalKills       int
	totalDamageDealt float64
	soulsCollected   int
	survivalFrames   int

	// Game Feel / Juice: Screen Shake e Hitstop
	shakeIntensity float64
	camOffsetX     float64
	camOffsetY     float64
	hitstopFrames  int

	// Rastreamento de conquistas na run
	bossDamagedPlayer  bool
	recentFirearmKills []int
	achieveBannerTxt   string
	achieveBannerTimer int

	// Efeitos de tela e banners de aviso
	screenFlashClr     color.RGBA
	screenFlashTimer   int
	bossAlertText      string
	bossAlertTimer     int
	weatherBannerText  string
	weatherBannerTimer int

	// Seleção de Relíquias (Escolha 1 de 3)
	isChoosingRelic bool
	relicChoices    []entities.RelicDef
	relicCardBounds [3]collision.Rect

	// Menu de Pausa (ESC ou P) e Submenu de Opções inline
	isPaused      bool
	inPauseOpts   bool
	btnResume     *ui.Button
	btnPauseOpts  *ui.Button
	btnRestart    *ui.Button
	btnMainMenu   *ui.Button

	// Controles de Opções dentro da Pausa
	btnPauseFS         *ui.Button
	btnPauseMasterDown *ui.Button
	btnPauseMasterUp   *ui.Button
	btnPauseSFXDown    *ui.Button
	btnPauseSFXUp      *ui.Button
	btnPauseBGMDown    *ui.Button
	btnPauseBGMUp      *ui.Button
	btnPauseBack       *ui.Button
}

// NewGameplayScene inicia a campanha de 15 ondas com a classe escolhida.
func NewGameplayScene(nav SceneNavigator, class entities.ClassType) *GameplayScene {
	p := entities.NewPlayer(class, 400, 320)

	hpBar := ui.NewProgressBar(12, 8, 175, 16, "HP", p.MaxHP, color.RGBA{R: 220, G: 50, B: 50, A: 255})
	var manaBar *ui.ProgressBar
	if p.MaxMana > 0 {
		manaBar = ui.NewProgressBar(12, 28, 145, 14, "MP", p.MaxMana, color.RGBA{R: 50, G: 120, B: 240, A: 255})
	}

	specialLabel := "Esp [E/Dir]"
	switch class {
	case entities.ClassWarrior:
		specialLabel = "Escudo [E/Dir]"
	case entities.ClassArcher:
		specialLabel = "Chuva [E/Dir]"
	case entities.ClassMage:
		specialLabel = "Pilar [E/Dir]"
	case entities.ClassShotgunner:
		specialLabel = "Sismico [E/Dir]"
	case entities.ClassSheriff:
		specialLabel = "Deadeye [E/Dir]"
	}

	specialBar := ui.NewProgressBar(195, 8, 145, 16, specialLabel, float64(p.SpecialMaxCD), color.RGBA{R: 245, G: 190, B: 50, A: 255})
	specialBar.ShowFraction = false

	dashBar := ui.NewProgressBar(195, 28, 145, 14, "Dash [Esp/Sh]", float64(p.DashMaxCD), color.RGBA{R: 80, G: 210, B: 240, A: 255})
	dashBar.ShowFraction = false

	bossBar := ui.NewProgressBar(175, 58, 450, 18, "CHEFE", 1000, color.RGBA{R: 200, G: 30, B: 60, A: 255})

	fsLabel := "⚙ TELA CHEIA: DESATIVADO"
	if nav.IsFullscreen() {
		fsLabel = "⚙ TELA CHEIA: ATIVADO"
	}

	cardBounds := [3]collision.Rect{
		collision.NewRect(55, 175, 210, 260),
		collision.NewRect(295, 175, 210, 260),
		collision.NewRect(535, 175, 210, 260),
	}

	// Botões à esquerda no modal de pausa (para caber a Ficha da Run ao lado direito!)
	pBtnX := 75.0
	pBtnW := 260.0
	pBtnH := 40.0

	ebiten.SetCursorMode(ebiten.CursorModeVisible)
	return &GameplayScene{
		navigator:       nav,
		arena:           world.NewArena(800, 600),
		weather:         world.NewWeatherSystem(),
		worldCanvas:     ebiten.NewImage(800, 600),
		player:          p,
		hpBar:           hpBar,
		manaBar:         manaBar,
		specialBar:      specialBar,
		dashBar:         dashBar,
		bossBar:         bossBar,
		soulIcon:        sprites.GetSoulFragmentSprite(),
		spawnTimer:      0,
		spawnDelay:      95,
		wave:            1,
		killsInWave:     0,
		killsToWave:     20,
		relicCardBounds: cardBounds,

		btnResume:    ui.NewButton(pBtnX, 195, pBtnW, pBtnH, "⚔  CONTINUAR  ⚔"),
		btnPauseOpts: ui.NewButton(pBtnX, 245, pBtnW, pBtnH, "⚙  OPÇÕES / ÁUDIO  ⚙"),
		btnRestart:   ui.NewButton(pBtnX, 295, pBtnW, pBtnH, "🔄  REINICIAR PARTIDA  🔄"),
		btnMainMenu:  ui.NewButton(pBtnX, 345, pBtnW, pBtnH, "🏰  VOLTAR AO MENU  🏰"),

		btnPauseFS:         ui.NewButton(250, 165, 300, 38, fsLabel),
		btnPauseMasterDown: ui.NewButton(250, 220, 50, 34, " - "),
		btnPauseMasterUp:   ui.NewButton(500, 220, 50, 34, " + "),
		btnPauseSFXDown:    ui.NewButton(250, 270, 50, 34, " - "),
		btnPauseSFXUp:      ui.NewButton(500, 270, 50, 34, " + "),
		btnPauseBGMDown:    ui.NewButton(250, 320, 50, 34, " - "),
		btnPauseBGMUp:      ui.NewButton(500, 320, 50, 34, " + "),
		btnPauseBack:       ui.NewButton(280, 380, 240, 40, "VOLTAR À PAUSA"),
	}
}

func (g *GameplayScene) addScreenShake(intensity float64) {
	if intensity > g.shakeIntensity {
		g.shakeIntensity = intensity
	}
}

func (g *GameplayScene) triggerHitstop(frames int) {
	if frames > g.hitstopFrames {
		g.hitstopFrames = frames
	}
}

// Update executa toda a simulação, Hitstop, Screen Shake, Relíquias e Pausa.
func (g *GameplayScene) Update() error {
	// 1. Seleção de Relíquia (congela simulação e mantém UI responsiva)
	if g.isChoosingRelic {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			for idx, bounds := range g.relicCardBounds {
				if idx < len(g.relicChoices) && bounds.Contains(float64(mx), float64(my)) {
					chosen := g.relicChoices[idx]
					g.player.AddRelic(chosen.ID)
					audio.Play(audio.SFXRelic)
					g.isChoosingRelic = false
					g.addFloatingText("RELIQUIA: "+chosen.Name, g.player.X-30, g.player.Y-25, chosen.Color)

					if g.player.TotalRelics >= 4 {
						g.tryUnlockAchievement("relic_hoarder")
					}
					break
				}
			}
		}
		return nil
	}

	// 2. Pausa com ESC ou P
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if g.isPaused && g.inPauseOpts {
			g.inPauseOpts = false
		} else {
			g.isPaused = !g.isPaused
			g.inPauseOpts = false
		}
	}

	if g.isPaused {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
		if g.inPauseOpts {
			g.updatePauseOptions()
			return nil
		}

		g.btnResume.Update()
		g.btnPauseOpts.Update()
		g.btnRestart.Update()
		g.btnMainMenu.Update()

		if g.btnResume.IsClicked() {
			g.isPaused = false
		} else if g.btnPauseOpts.IsClicked() {
			g.inPauseOpts = true
		} else if g.btnRestart.IsClicked() {
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
			g.navigator.SwitchToGameplay(g.player.ClassType)
		} else if g.btnMainMenu.IsClicked() {
			g.navigator.SwitchToMenu()
		}
		return nil
	}

	// 3. Decaimento do Screen Shake (continua suave mesmo durante Hitstop!)
	if g.shakeIntensity > 0.2 {
		g.camOffsetX = (rand.Float64()*2 - 1) * g.shakeIntensity
		g.camOffsetY = (rand.Float64()*2 - 1) * g.shakeIntensity
		g.shakeIntensity *= 0.84
	} else {
		g.shakeIntensity = 0
		g.camOffsetX = 0
		g.camOffsetY = 0
	}

	// 4. Hitstop / Micro-congelamento (Frame Freeze em golpes pesados!)
	if g.hitstopFrames > 0 {
		g.hitstopFrames--
		return nil
	}

	g.survivalFrames++
	if g.screenFlashTimer > 0 {
		g.screenFlashTimer--
	}
	if g.bossAlertTimer > 0 {
		g.bossAlertTimer--
	}
	if g.weatherBannerTimer > 0 {
		g.weatherBannerTimer--
	}
	if g.achieveBannerTimer > 0 {
		g.achieveBannerTimer--
	}

	// Atualiza eventos climáticos
	g.weather.Update()
	if g.weather.Active == world.WeatherNone && g.survivalFrames%2600 == 0 {
		g.triggerWeatherEvent()
	}

	// 5. Atualiza Player
	g.player.Update(800, 600)
	px, py := g.player.Center()

	// Dash Universal (Espaço ou Shift)
	dashTriggered := inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyShiftRight)
	if dashTriggered && g.player.TriggerDash() {
		audio.Play(audio.SFXDash)
		g.particles = append(g.particles, entities.SpawnBurst(px, py+8, 8, 2.8, 2.5, color.RGBA{120, 220, 255, 220}, 16)...)
	}

	// Poeira nos pés
	if g.player.JustStepped {
		for k := 0; k < 2; k++ {
			dust := entities.NewParticle(
				px+(rand.Float64()-0.5)*10,
				py+12,
				(rand.Float64()-0.5)*0.8,
				-0.4-rand.Float64()*0.5,
				float32(2.0+rand.Float64()*1.5),
				color.RGBA{R: 150, G: 145, B: 135, A: 140},
				18,
			)
			g.particles = append(g.particles, dust)
		}
	}

	// Bolhas de veneno & Fumaça
	if g.player.PoisonTimer > 0 {
		if g.survivalFrames%5 == 0 {
			pClr := color.RGBA{R: 165, G: 70, B: 230, A: 220}
			if rand.Float64() < 0.5 {
				pClr = color.RGBA{R: 60, G: 230, B: 90, A: 220}
			}
			bubble := entities.NewParticle(
				px+(rand.Float64()-0.5)*24,
				py+(rand.Float64()-0.5)*20,
				(rand.Float64()-0.5)*0.4,
				-1.0-rand.Float64()*0.8,
				float32(2.0+rand.Float64()*2.0),
				pClr,
				28,
			)
			bubble.IsBubble = true
			g.particles = append(g.particles, bubble)
		}
		if g.survivalFrames%35 == 0 {
			g.addFloatingText("-2 (Veneno)", g.player.X, g.player.Y-14, color.RGBA{R: 185, G: 85, B: 255, A: 255})
		}
	}

	if g.player.MuzzleSmokeTimer > 0 && g.survivalFrames%2 == 0 {
		smokeX := px + math.Cos(g.player.FacingAngle)*20
		smokeY := py + math.Sin(g.player.FacingAngle)*20
		smoke := entities.NewParticle(
			smokeX,
			smokeY,
			math.Cos(g.player.FacingAngle)*1.2+(rand.Float64()-0.5)*0.6,
			math.Sin(g.player.FacingAngle)*1.2-0.5,
			float32(3.0+rand.Float64()*2.5),
			color.RGBA{R: 170, G: 170, B: 180, A: 170},
			22,
		)
		g.particles = append(g.particles, smoke)
	}

	// Dano das poças ácidas
	if g.weather.Active == world.WeatherAcidRain && g.survivalFrames%28 == 0 && g.player.InvincibleCD <= 0 {
		for _, puddle := range g.weather.Puddles {
			if math.Hypot(px-puddle.X, py-puddle.Y) <= puddle.Radius {
				g.player.HP -= 4
				g.player.HitFlashCD = 2
				g.addScreenShake(3.0)
				if g.player.HP < 0 {
					g.player.HP = 0
				}
				g.addFloatingText("-4 (Acido)", g.player.X, g.player.Y-12, color.RGBA{90, 240, 100, 255})
				break
			}
		}
	}

	if g.player.ClassType == entities.ClassSheriff && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.player.StartReload()
	}

	if g.player.ClassType == entities.ClassSheriff && g.player.DeadeyeQueue > 0 {
		if g.player.DeadeyeTimer > 0 {
			g.player.DeadeyeTimer--
		} else {
			g.fireDeadeyeShot()
			g.player.DeadeyeQueue--
			g.player.DeadeyeTimer = 5
		}
	}

	// Atualiza barras do HUD
	g.hpBar.SetValues(g.player.HP, g.player.MaxHP)
	if g.manaBar != nil {
		g.manaBar.SetValues(g.player.Mana, g.player.MaxMana)
	}
	specCurrent := float64(g.player.SpecialMaxCD - g.player.SpecialCD)
	if specCurrent < 0 {
		specCurrent = 0
	}
	g.specialBar.SetValues(specCurrent, float64(g.player.SpecialMaxCD))

	dashCurrent := float64(g.player.DashMaxCD - g.player.DashCD)
	if dashCurrent < 0 {
		dashCurrent = 0
	}
	g.dashBar.SetValues(dashCurrent, float64(g.player.DashMaxCD))

	if g.player.HP <= 0 {
		g.handleGameOver()
		return nil
	}

	// Ataque Básico (Clique Esquerdo)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.player.CanAttack() {
		g.handlePlayerAttack()
	}

	// Habilidade Especial Ativa (E ou Botão Direito)
	specialTriggered := inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	if specialTriggered && g.player.CanUseSpecial() {
		g.handleSpecialSkill()
	}

	// Spawner de Inimigos
	if g.boss == nil || !g.boss.Alive {
		g.spawnTimer++
		if g.spawnTimer >= g.spawnDelay {
			g.spawnTimer = 0
			g.spawnEnemy()
			g.spawnBatchCnt++
			if g.spawnBatchCnt%8 == 0 {
				g.spawnTimer = -65
			}
		}
	}

	// Atualiza Efeitos Visuais
	for i := len(g.slashes) - 1; i >= 0; i-- {
		s := g.slashes[i]
		s.Update()
		if !s.Active {
			g.slashes = append(g.slashes[:i], g.slashes[i+1:]...)
		}
	}

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

	for i := len(g.shockwaves) - 1; i >= 0; i-- {
		if !g.shockwaves[i].Update() {
			g.shockwaves = append(g.shockwaves[:i], g.shockwaves[i+1:]...)
		}
	}

	for i := len(g.particles) - 1; i >= 0; i-- {
		if !g.particles[i].Update() {
			g.particles = append(g.particles[:i], g.particles[i+1:]...)
		}
	}

	for i := len(g.lightningArcs) - 1; i >= 0; i-- {
		g.lightningArcs[i].Life--
		if g.lightningArcs[i].Life <= 0 {
			g.lightningArcs = append(g.lightningArcs[:i], g.lightningArcs[i+1:]...)
		}
	}

	// Atualiza Projéteis
	pBounds := g.player.Bounds()
	for i := len(g.projectiles) - 1; i >= 0; i-- {
		if i >= len(g.projectiles) {
			continue
		}
		proj := g.projectiles[i]

		if proj.IsHoming {
			if proj.IsEnemy {
				proj.SteerTowards(px, py, proj.HomingTurn)
			} else {
				tx, ty, found := g.findNearestTarget(proj.X, proj.Y)
				if found {
					proj.SteerTowards(tx, ty, 0.28)
				}
			}
		}

		proj.Update()

		if !proj.Active || proj.X < 0 || proj.X > 800 || proj.Y < 45 || proj.Y > 600 {
			g.projectiles = append(g.projectiles[:i], g.projectiles[i+1:]...)
			continue
		}

		// Reflexão pelo Escudo do Guerreiro (+ Hitstop e Shake!)
		if proj.IsEnemy && g.player.ClassType == entities.ClassWarrior && g.player.Shield.Active {
			if g.player.CheckShieldReflect(proj) {
				audio.Play(audio.SFXShieldReflect)
				g.triggerHitstop(3)
				g.addScreenShake(5.0)
				g.addFloatingText("REFLETIDO!", proj.X, proj.Y-12, color.RGBA{R: 255, G: 235, B: 100, A: 255})
				g.particles = append(g.particles, entities.SpawnBurst(proj.X, proj.Y, 6, 3.0, 2.5, color.RGBA{255, 220, 90, 255}, 18)...)

				save := progression.GetSave()
				save.ShieldReflectCount++
				if save.ShieldReflectCount >= 30 {
					g.tryUnlockAchievement("iron_wall")
				}
				_ = progression.Save()
				continue
			}
		}

		if proj.IsEnemy && proj.Bounds().Intersects(pBounds) {
			proj.Active = false
			if g.player.TakeDamage(proj.Damage) {
				if g.boss != nil && g.boss.Alive {
					g.bossDamagedPlayer = true
				}
				audio.Play(audio.SFXHit)
				g.addScreenShake(7.0)
				g.addFloatingText(fmt.Sprintf("-%.0f", proj.Damage), g.player.X, g.player.Y-10, color.RGBA{R: 255, G: 80, B: 80, A: 255})
				g.particles = append(g.particles, entities.SpawnBurst(px, py, 6, 3.0, 2.5, color.RGBA{255, 60, 60, 255}, 18)...)
			}
			g.projectiles = append(g.projectiles[:i], g.projectiles[i+1:]...)
		}
	}

	// Atualiza Chefe
	if g.boss != nil && g.boss.Alive {
		g.updateBoss(px, py, pBounds)
		// Se o Soberano na Onda 15 morreu, já transicionou para VictoryScene!
		if g.boss == nil && g.wave == TotalCampaignWaves {
			return nil
		}
	}

	// Atualiza Inimigos Comuns
	enemySpeedMult := 1.0
	if g.weather.Active == world.WeatherBloodMoon {
		enemySpeedMult = 1.25
	}

	for i := len(g.enemies) - 1; i >= 0; i-- {
		if i >= len(g.enemies) {
			continue
		}
		e := g.enemies[i]
		newEnemyProj := e.Update(px, py, enemySpeedMult)
		if newEnemyProj != nil {
			g.projectiles = append(g.projectiles, newEnemyProj)
		}

		eBounds := e.Bounds()

		if g.weather.Active == world.WeatherAcidRain && g.survivalFrames%28 == 0 {
			ex, ey := e.X+e.Width/2, e.Y+e.Height/2
			for _, puddle := range g.weather.Puddles {
				if math.Hypot(ex-puddle.X, ey-puddle.Y) <= puddle.Radius {
					g.damageEnemy(e, 7, false, 0, 0)
					break
				}
			}
		}
		if len(g.enemies) == 0 {
			break
		}

		if e.Type == entities.EnemyTypeBat && e.IsAttached && g.survivalFrames%12 == 0 {
			g.particles = append(g.particles, entities.SpawnBurst(px, py, 3, 2.5, 2.0, color.RGBA{220, 20, 40, 255}, 15)...)
		}

		// Colisão com Slash do Guerreiro (aplica Knockback responsivo na direção do golpe!)
		for _, slash := range g.slashes {
			if e.Alive && slash.CollidesWithEnemy(e.ID, eBounds) {
				dmg := slash.Damage * g.player.GetDamageMultiplier()
				kvx := math.Cos(slash.Angle) * 7.5
				kvy := math.Sin(slash.Angle) * 7.5
				g.damageEnemy(e, dmg, false, kvx, kvy)
				g.onPlayerAttackHit(e.X+e.Width/2, e.Y+e.Height/2)
			}
		}
		if len(g.enemies) == 0 {
			break
		}

		// Colisão com Projéteis do Jogador
		for _, proj := range g.projectiles {
			if !e.Alive || !proj.Active || proj.IsEnemy {
				continue
			}
			if proj.HitIDs != nil && proj.HitIDs[e.ID] {
				continue
			}
			if proj.Bounds().Intersects(eBounds) {
				if proj.PierceLeft > 0 {
					proj.PierceLeft--
					if proj.HitIDs == nil {
						proj.HitIDs = make(map[uint64]bool)
					}
					proj.HitIDs[e.ID] = true
				} else {
					proj.Active = false
				}

				dmg := proj.Damage * g.player.GetDamageMultiplier()
				kvx, kvy := 0.0, 0.0
				if proj.Knockback > 0 {
					spd := math.Hypot(proj.Vx, proj.Vy)
					if spd > 0.01 {
						kvx = (proj.Vx / spd) * proj.Knockback
						kvy = (proj.Vy / spd) * proj.Knockback
					}
				}
				// Hitstop curto em tiros de Calibre 12 à queima-roupa!
				if proj.Kind == entities.ProjectileKindPellet && math.Hypot(e.X-px, e.Y-py) < 115 {
					g.triggerHitstop(2)
					g.addScreenShake(4.0)
				}

				g.damageEnemy(e, dmg, proj.IsCrit, kvx, kvy)
				g.onPlayerAttackHit(e.X+e.Width/2, e.Y+e.Height/2)

				if !proj.Active {
					break
				}
			}
		}
		if len(g.enemies) == 0 {
			break
		}

		for _, ar := range g.arrowRains {
			if e.Alive && ar.CollidesWith(eBounds) && g.survivalFrames%10 == 0 {
				dmg := ar.DamagePerTick * g.player.GetDamageMultiplier()
				e.SlowTimer = 120
				g.damageEnemy(e, dmg, false, 0, 0)
			}
		}
		if len(g.enemies) == 0 {
			break
		}

		for _, fp := range g.firePillars {
			if e.Alive && fp.CollidesWith(eBounds) && g.survivalFrames%8 == 0 {
				dmg := fp.DamagePerTick * g.player.GetDamageMultiplier()
				g.damageEnemy(e, dmg, false, 0, 0)
			}
		}
		if len(g.enemies) == 0 {
			break
		}

		if e.Alive && eBounds.Intersects(pBounds) {
			if e.CanDamagePlayer() {
				if g.player.TakeDamage(e.Damage) {
					e.ResetAttackCD()
					if g.boss != nil && g.boss.Alive {
						g.bossDamagedPlayer = true
					}
					audio.Play(audio.SFXHit)
					g.addScreenShake(7.5)
					g.addFloatingText(fmt.Sprintf("-%.0f", e.Damage), g.player.X, g.player.Y-10, color.RGBA{R: 255, G: 60, B: 60, A: 255})

					if thorns := g.player.GetThornsDamage(); thorns > 0 {
						g.damageEnemy(e, thorns, false, 0, 0)
						g.addFloatingText(fmt.Sprintf("ESPINHOS -%.0f!", thorns), e.X, e.Y-16, color.RGBA{90, 210, 100, 255})
					}
				}
			}
		}
		if len(g.enemies) == 0 {
			break
		}

		if !e.Alive && i < len(g.enemies) && g.enemies[i] == e {
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
		}
	}

	// Atualiza e coleta Drops
	magnetRad := g.player.GetMagnetRadius()
	for i := len(g.drops) - 1; i >= 0; i-- {
		if i >= len(g.drops) {
			continue
		}
		d := g.drops[i]
		d.Update()
		if !d.Active {
			g.drops = append(g.drops[:i], g.drops[i+1:]...)
			continue
		}

		distToPlayer := math.Hypot(px-d.X, py-d.Y)
		if distToPlayer > 1.0 && distToPlayer <= magnetRad {
			pullSpeed := 4.8
			d.X += ((px - d.X) / distToPlayer) * pullSpeed
			d.Y += ((py - d.Y) / distToPlayer) * pullSpeed
		}

		if d.Bounds().Intersects(pBounds) {
			d.Active = false
			audio.Play(audio.SFXPickup)
			feedback, clr := g.player.ApplyDrop(d)
			g.addFloatingText(feedback, g.player.X-20, g.player.Y-20, clr)
			g.particles = append(g.particles, entities.SpawnBurst(d.X, d.Y, 10, 3.0, 2.5, clr, 22)...)

			if d.Type == entities.DropSoulFragment {
				g.soulsCollected += d.SoulAmount
			} else if d.Type == entities.DropChaosSeed {
				g.triggerChaosAmbush(px, py)
			}

			g.drops = append(g.drops[:i], g.drops[i+1:]...)
		}
	}

	// Textos Flutuantes
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

func (g *GameplayScene) updatePauseOptions() {
	g.btnPauseFS.Update()
	g.btnPauseMasterDown.Update()
	g.btnPauseMasterUp.Update()
	g.btnPauseSFXDown.Update()
	g.btnPauseSFXUp.Update()
	g.btnPauseBGMDown.Update()
	g.btnPauseBGMUp.Update()
	g.btnPauseBack.Update()

	if g.btnPauseFS.IsClicked() {
		g.navigator.ToggleFullscreen()
		if g.navigator.IsFullscreen() {
			g.btnPauseFS.Text = "⚙ TELA CHEIA: ATIVADO"
		} else {
			g.btnPauseFS.Text = "⚙ TELA CHEIA: DESATIVADO"
		}
	}

	if g.btnPauseMasterDown.IsClicked() {
		g.navigator.SetMasterVolume(clampVol(g.navigator.GetMasterVolume() - 10))
	} else if g.btnPauseMasterUp.IsClicked() {
		g.navigator.SetMasterVolume(clampVol(g.navigator.GetMasterVolume() + 10))
	}

	if g.btnPauseSFXDown.IsClicked() {
		g.navigator.SetSFXVolume(clampVol(g.navigator.GetSFXVolume() - 10))
		audio.Play(audio.SFXPickup)
	} else if g.btnPauseSFXUp.IsClicked() {
		g.navigator.SetSFXVolume(clampVol(g.navigator.GetSFXVolume() + 10))
		audio.Play(audio.SFXPickup)
	}

	if g.btnPauseBGMDown.IsClicked() {
		g.navigator.SetMusicVolume(clampVol(g.navigator.GetMusicVolume() - 10))
	} else if g.btnPauseBGMUp.IsClicked() {
		g.navigator.SetMusicVolume(clampVol(g.navigator.GetMusicVolume() + 10))
	}

	if g.btnPauseBack.IsClicked() {
		g.inPauseOpts = false
	}
}

func (g *GameplayScene) onPlayerAttackHit(hitX, hitY float64) {
	stacks := g.player.Relics[entities.RelicStaticSpark]
	if stacks <= 0 {
		return
	}
	g.player.HitCounter++
	if g.player.HitCounter%5 != 0 {
		return
	}

	maxJumps := 2 + stacks
	sparkDmg := 28.0 * float64(stacks) * g.player.GetDamageMultiplier()
	currX, currY := hitX, hitY
	hitSet := make(map[uint64]bool)
	jumped := 0

	for j := 0; j < maxJumps; j++ {
		var bestEnemy *entities.Enemy
		bestDist := 175.0
		for _, e := range g.enemies {
			if !e.Alive || hitSet[e.ID] {
				continue
			}
			ex, ey := e.X+e.Width/2, e.Y+e.Height/2
			d := math.Hypot(ex-currX, ey-currY)
			if d > 2.0 && d < bestDist {
				bestDist = d
				bestEnemy = e
			}
		}
		if bestEnemy == nil {
			break
		}
		hitSet[bestEnemy.ID] = true
		ex, ey := bestEnemy.X+bestEnemy.Width/2, bestEnemy.Y+bestEnemy.Height/2
		g.lightningArcs = append(g.lightningArcs, &LightningArc{X1: currX, Y1: currY, X2: ex, Y2: ey, Life: 12})
		g.damageEnemy(bestEnemy, sparkDmg, false, 0, 0)
		currX, currY = ex, ey
		jumped++
	}

	if jumped > 0 {
		audio.Play(audio.SFXLightning)
	}
}

func (g *GameplayScene) damageEnemy(e *entities.Enemy, dmg float64, isCrit bool, kvx, kvy float64) {
	if kvx != 0 || kvy != 0 {
		e.ApplyKnockback(kvx, kvy)
	}
	g.totalDamageDealt += dmg
	killed := e.TakeDamage(dmg)
	audio.Play(audio.SFXHit)

	label := fmt.Sprintf("-%.0f", dmg)
	clr := color.RGBA{R: 255, G: 235, B: 110, A: 255}
	if isCrit {
		label = fmt.Sprintf("CRIT -%.0f!", dmg)
		clr = color.RGBA{R: 255, G: 190, B: 40, A: 255}
	}
	g.addFloatingText(label, e.X, e.Y-6, clr)
	g.particles = append(g.particles, entities.SpawnBurst(e.X+e.Width/2, e.Y+e.Height/2, 4, 2.6, 2.0, clr, 14)...)

	if killed {
		g.onEnemyDefeated(e)
	}
}

func (g *GameplayScene) updateBoss(px, py float64, pBounds collision.Rect) {
	b := g.boss
	out := b.Update(px, py)

	if len(out.Projectiles) > 0 {
		g.projectiles = append(g.projectiles, out.Projectiles...)
	}
	if len(out.Shockwaves) > 0 {
		g.shockwaves = append(g.shockwaves, out.Shockwaves...)
	}
	if len(out.Particles) > 0 {
		g.particles = append(g.particles, out.Particles...)
	}

	bcx, bcy := b.Center()
	if out.StompDamage > 0 {
		audio.Play(audio.SFXExplosion)
		g.addScreenShake(10.0)
		dist := math.Hypot(px-bcx, py-bcy)
		if dist <= out.StompRadius {
			if g.player.TakeDamage(out.StompDamage) {
				g.bossDamagedPlayer = true
				g.addFloatingText(fmt.Sprintf("-%.0f (PISAO!)", out.StompDamage), g.player.X, g.player.Y-15, color.RGBA{255, 60, 40, 255})
			}
		}
	}

	if out.LaserDamage > 0 {
		if g.player.TakeDamage(out.LaserDamage) {
			g.bossDamagedPlayer = true
			audio.Play(audio.SFXHit)
			g.addScreenShake(8.0)
			g.addFloatingText(fmt.Sprintf("-%.0f (LASER!)", out.LaserDamage), g.player.X, g.player.Y-15, color.RGBA{210, 80, 255, 255})
		}
	}

	bBounds := b.Bounds()
	g.bossBar.Label = b.Name
	g.bossBar.SetValues(b.HP, b.MaxHP)

	for _, slash := range g.slashes {
		if b.Alive && slash.CollidesWithEnemy(999999, bBounds) {
			dmg := slash.Damage * g.player.GetDamageMultiplier()
			g.damageBoss(dmg, false)
		}
	}

	for _, proj := range g.projectiles {
		if !b.Alive {
			break
		}
		if proj.Active && !proj.IsEnemy && proj.Bounds().Intersects(bBounds) {
			proj.Active = false
			dmg := proj.Damage * g.player.GetDamageMultiplier()
			g.damageBoss(dmg, proj.IsCrit)
			break
		}
	}

	for _, ar := range g.arrowRains {
		if b.Alive && ar.CollidesWith(bBounds) && g.survivalFrames%10 == 0 {
			dmg := ar.DamagePerTick * g.player.GetDamageMultiplier()
			g.damageBoss(dmg, false)
		}
	}

	for _, fp := range g.firePillars {
		if b.Alive && fp.CollidesWith(bBounds) && g.survivalFrames%8 == 0 {
			dmg := fp.DamagePerTick * g.player.GetDamageMultiplier()
			g.damageBoss(dmg, false)
		}
	}

	if b.Alive && bBounds.Intersects(pBounds) {
		if b.CanDamagePlayer() {
			if g.player.TakeDamage(b.Damage) {
				b.ResetAttackCD()
				g.bossDamagedPlayer = true
				audio.Play(audio.SFXHit)
				g.addScreenShake(8.5)
				g.addFloatingText(fmt.Sprintf("-%.0f", b.Damage), g.player.X, g.player.Y-12, color.RGBA{255, 50, 50, 255})

				if thorns := g.player.GetThornsDamage(); thorns > 0 {
					g.damageBoss(thorns, false)
				}
			}
		}
	}
}

func (g *GameplayScene) damageBoss(dmg float64, isCrit bool) {
	if g.boss == nil || !g.boss.Alive {
		return
	}
	g.totalDamageDealt += dmg
	killed := g.boss.TakeDamage(dmg)
	audio.Play(audio.SFXHit)

	bcx, bcy := g.boss.Center()
	label := fmt.Sprintf("-%.0f", dmg)
	clr := color.RGBA{255, 220, 100, 255}
	if isCrit {
		label = fmt.Sprintf("CRIT -%.0f!", dmg)
		clr = color.RGBA{255, 160, 40, 255}
	}
	g.addFloatingText(label, bcx-15, g.boss.Y-10, clr)
	g.particles = append(g.particles, entities.SpawnBurst(bcx, bcy, 4, 3.0, 2.5, clr, 16)...)

	if killed {
		g.onBossDefeated()
	}
}

func (g *GameplayScene) onBossDefeated() {
	bcx, bcy := g.boss.Center()
	audio.Play(audio.SFXExplosion)
	g.triggerHitstop(5)
	g.addScreenShake(14.0)

	if !g.bossDamagedPlayer {
		g.tryUnlockAchievement("giant_slayer")
	}

	g.particles = append(g.particles, entities.SpawnBurst(bcx, bcy, 45, 6.5, 4.5, color.RGBA{255, 200, 60, 255}, 45)...)
	g.particles = append(g.particles, entities.SpawnBurst(bcx, bcy, 30, 5.0, 3.5, color.RGBA{180, 90, 255, 255}, 40)...)
	g.shockwaves = append(g.shockwaves, entities.NewShockwave(bcx, bcy, 180, color.RGBA{255, 215, 70, 255}, 35, false))

	g.screenFlashClr = color.RGBA{255, 230, 150, 140}
	g.screenFlashTimer = 25

	tier := g.boss.WaveLevel / 5
	g.score += 2000 * tier
	g.totalKills++
	isFinalBoss := g.boss.IsFinal || g.wave >= TotalCampaignWaves
	g.boss = nil

	// Se derrotou o Soberano de Dorfieegu na Onda 15 -> VITÓRIA!
	if isFinalBoss {
		g.handleVictory()
		return
	}

	// Caso contrário, dropa recompensas, volta música de exploração e abre escolha de relíquia!
	audio.PlayBGM(audio.TrackExploration)

	for i := 0; i < 4; i++ {
		ox := (rand.Float64() - 0.5) * 90
		oy := (rand.Float64() - 0.5) * 90
		g.drops = append(g.drops, entities.NewSoulDrop(bcx+ox, bcy+oy, 3))
	}
	g.drops = append(g.drops, entities.NewDrop(entities.DropHealthPotion, bcx-30, bcy+20))
	g.drops = append(g.drops, entities.NewDrop(entities.DropMasteryTome, bcx+30, bcy+20))

	g.advanceToWave(g.wave + 1)
	g.openRelicSelection()
}

func (g *GameplayScene) handleVictory() {
	// Todos os projéteis inimigos ativos viram fragmentos de alma!
	convertedSouls := 0
	for _, proj := range g.projectiles {
		if proj.IsEnemy && proj.Active {
			convertedSouls++
			g.particles = append(g.particles, entities.SpawnBurst(proj.X, proj.Y, 6, 3.0, 2.5, color.RGBA{170, 110, 255, 255}, 25)...)
		}
	}
	g.projectiles = nil

	baseSouls := g.soulsCollected + convertedSouls + 15 // Bônus base pelo Soberano
	bonusSouls := baseSouls / 2                        // Bônus multiplicador 1.5x de vitória!
	totalSouls := baseSouls + bonusSouls

	progression.AddSoulFragments(totalSouls)
	save := progression.GetSave()
	save.TotalRuns++
	save.TotalKills += g.totalKills
	if TotalCampaignWaves > save.MaxWaveReached {
		save.MaxWaveReached = TotalCampaignWaves
	}
	_ = progression.Save()

	summary := RunSummary{
		Class:        g.player.ClassType,
		Score:        g.score,
		Kills:        g.totalKills,
		RelicsCount:  g.player.TotalRelics,
		DamageDealt:  g.totalDamageDealt,
		SurvivalSecs: g.survivalFrames / 60,
		BaseSouls:    baseSouls,
		BonusSouls:   bonusSouls,
		TotalSouls:   totalSouls,
	}

	ebiten.SetCursorMode(ebiten.CursorModeVisible)
	g.navigator.SwitchToVictory(summary)
}

func (g *GameplayScene) openRelicSelection() {
	g.isChoosingRelic = true
	g.relicChoices = entities.PickRandomRelics(3)
}

func (g *GameplayScene) findNearestTarget(x, y float64) (float64, float64, bool) {
	bestDist := 999999.0
	var tx, ty float64
	found := false

	if g.boss != nil && g.boss.Alive {
		bcx, bcy := g.boss.Center()
		d := math.Hypot(bcx-x, bcy-y)
		bestDist = d
		tx, ty = bcx, bcy
		found = true
	}

	for _, e := range g.enemies {
		if !e.Alive {
			continue
		}
		ex := e.X + e.Width/2
		ey := e.Y + e.Height/2
		d := math.Hypot(ex-x, ey-y)
		if d < bestDist {
			bestDist = d
			tx, ty = ex, ey
			found = true
		}
	}

	return tx, ty, found
}

func (g *GameplayScene) makePlayerProjectile(x, y, vx, vy, rad, dmg float64, clr color.RGBA, life int, kind entities.ProjectileKind) *entities.Projectile {
	proj := entities.NewProjectile(x, y, vx, vy, rad, dmg, clr, life, kind, false)
	if pierce := g.player.GetPierceChance(); pierce > 0 && rand.Float64() < pierce {
		proj.PierceLeft = 1
		proj.HitIDs = make(map[uint64]bool)
	}
	return proj
}

func (g *GameplayScene) handlePlayerAttack() {
	px, py := g.player.Center()
	mx, my := ebiten.CursorPosition()
	fmx, fmy := float64(mx), float64(my)

	switch g.player.ClassType {
	case entities.ClassWarrior:
		slash := entities.NewSlash(px, py, fmx, fmy, g.player.SlashRange, 68)
		g.slashes = append(g.slashes, slash)
		g.player.ResetAttackCD()
		audio.Play(audio.SFXSlash)

	case entities.ClassArcher:
		dx := fmx - px
		dy := fmy - py
		dist := math.Hypot(dx, dy)
		if dist > 0.001 {
			vx := (dx / dist) * g.player.ArrowSpeed
			vy := (dy / dist) * g.player.ArrowSpeed
			proj := g.makePlayerProjectile(px, py, vx, vy, 3.5, 28, color.RGBA{R: 120, G: 255, B: 120, A: 255}, 85, entities.ProjectileKindArrow)
			g.projectiles = append(g.projectiles, proj)
			g.player.ResetAttackCD()
			audio.Play(audio.SFXArrow)
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
				proj := g.makePlayerProjectile(px, py, vx, vy, 8.5, dmg, color.RGBA{R: 70, G: 160, B: 255, A: 255}, 90, entities.ProjectileKindMagic)
				g.projectiles = append(g.projectiles, proj)
			} else {
				dmg := 15.0 * g.player.SpellMultiplier
				proj := g.makePlayerProjectile(px, py, vx*0.9, vy*0.9, 4.0, dmg, color.RGBA{R: 160, G: 190, B: 255, A: 200}, 70, entities.ProjectileKindMagicBasic)
				g.projectiles = append(g.projectiles, proj)
			}
			g.player.ResetAttackCD()
			audio.Play(audio.SFXMagic)
		}

	case entities.ClassShotgunner:
		baseAngle := g.player.FacingAngle
		count := g.player.PelletCount
		spread := g.player.PelletSpread
		for i := 0; i < count; i++ {
			offset := (float64(i)/float64(count-1) - 0.5) * spread
			ang := baseAngle + offset + (rand.Float64()-0.5)*0.05
			spd := 10.0 + rand.Float64()*2.0
			vx := math.Cos(ang) * spd
			vy := math.Sin(ang) * spd
			proj := g.makePlayerProjectile(px, py, vx, vy, 3.2, g.player.PelletDamage, color.RGBA{255, 190, 70, 255}, 24, entities.ProjectileKindPellet)
			proj.Knockback = g.player.PelletKnockback
			g.projectiles = append(g.projectiles, proj)
		}
		g.player.MuzzleSmokeTimer = 18
		g.player.ResetAttackCD()
		g.addScreenShake(3.2)
		audio.Play(audio.SFXShotgun)

	case entities.ClassSheriff:
		if g.player.Ammo <= 0 || g.player.ReloadTimer > 0 {
			return
		}
		dx := fmx - px
		dy := fmy - py
		dist := math.Hypot(dx, dy)
		if dist > 0.001 {
			spd := 13.0
			vx := (dx / dist) * spd
			vy := (dy / dist) * spd

			dmg := 42.0
			isCrit := rand.Float64() < g.player.CritChance
			if isCrit {
				dmg *= g.player.CritMult
			}

			proj := g.makePlayerProjectile(px, py, vx, vy, 3.0, dmg, color.RGBA{240, 245, 255, 255}, 70, entities.ProjectileKindRevolver)
			proj.IsCrit = isCrit
			proj.Knockback = 2.5
			g.projectiles = append(g.projectiles, proj)

			g.player.Ammo--
			g.player.ResetAttackCD()
			audio.Play(audio.SFXRevolver)
			if g.player.Ammo <= 0 {
				g.player.StartReload()
			}
		}
	}
}

func (g *GameplayScene) handleSpecialSkill() {
	px, py := g.player.Center()
	mx, my := ebiten.CursorPosition()
	fmx, fmy := float64(mx), float64(my)

	switch g.player.ClassType {
	case entities.ClassWarrior:
		g.player.ActivateShield()
		audio.Play(audio.SFXShieldReflect)
		g.addFloatingText("POSTURA DE ESCUDO!", g.player.X-20, g.player.Y-25, color.RGBA{R: 255, G: 220, B: 80, A: 255})

	case entities.ClassArcher:
		rain := entities.NewArrowRain(fmx, fmy, g.player.ArrowRainRadius, g.player.ArrowRainDamage, 130)
		g.arrowRains = append(g.arrowRains, rain)
		g.player.SpecialCD = g.player.SpecialMaxCD
		audio.Play(audio.SFXArrow)
		g.addFloatingText("CHUVA DE FLECHAS!", fmx-40, fmy-20, color.RGBA{R: 120, G: 240, B: 140, A: 255})

	case entities.ClassMage:
		if g.player.Mana >= g.player.FirePillarManaCost {
			g.player.Mana -= g.player.FirePillarManaCost
			pillar := entities.NewFirePillar(fmx, fmy, g.player.FirePillarRadius, g.player.FirePillarDamage, g.player.FirePillarDuration)
			g.firePillars = append(g.firePillars, pillar)
			g.player.SpecialCD = g.player.SpecialMaxCD
			g.addScreenShake(6.5)
			audio.Play(audio.SFXExplosion)
			g.addFloatingText("PILAR DE FOGO!", fmx-35, fmy-25, color.RGBA{R: 255, G: 140, B: 40, A: 255})
		}

	case entities.ClassShotgunner:
		rad := g.player.ShockwaveRadius
		dmg := g.player.ShockwaveDamage * g.player.GetDamageMultiplier()
		g.shockwaves = append(g.shockwaves, entities.NewShockwave(px, py, rad, color.RGBA{255, 165, 40, 240}, 24, false))
		g.particles = append(g.particles, entities.SpawnBurst(px, py, 28, 5.5, 3.5, color.RGBA{255, 190, 70, 255}, 25)...)
		g.triggerHitstop(3)
		g.addScreenShake(10.0)
		audio.Play(audio.SFXExplosion)

		for _, e := range g.enemies {
			if !e.Alive {
				continue
			}
			ex, ey := e.X+e.Width/2, e.Y+e.Height/2
			dx, dy := ex-px, ey-py
			dist := math.Hypot(dx, dy)
			if dist <= rad {
				if dist < 1 {
					dist = 1
				}
				push := 14.0
				e.ApplyKnockback((dx/dist)*push, (dy/dist)*push)
				e.StunTimer = 90
				g.damageEnemy(e, dmg, false, 0, 0)
			}
			if len(g.enemies) == 0 {
				break
			}
		}

		if g.boss != nil && g.boss.Alive {
			bcx, bcy := g.boss.Center()
			if math.Hypot(bcx-px, bcy-py) <= rad+36 {
				g.damageBoss(dmg, false)
			}
		}

		g.player.SpecialCD = g.player.SpecialMaxCD
		g.addFloatingText("DISPARO SISMICO!", px-45, py-30, color.RGBA{255, 175, 40, 255})

	case entities.ClassSheriff:
		g.player.Ammo = g.player.MaxAmmo
		g.player.ReloadTimer = 0
		g.player.DeadeyeQueue = 6
		g.player.DeadeyeTimer = 0
		g.player.SpecialCD = g.player.SpecialMaxCD
		audio.Play(audio.SFXRelic)
		g.addFloatingText("HORA DA FORCA! (DEADEYE)", px-55, py-30, color.RGBA{255, 225, 50, 255})
	}
}

func (g *GameplayScene) fireDeadeyeShot() {
	px, py := g.player.Center()
	tx, ty, found := g.findNearestTarget(px, py)
	ang := g.player.FacingAngle
	if found {
		ang = math.Atan2(ty-py, tx-px)
	}

	spd := 14.5
	vx := math.Cos(ang) * spd
	vy := math.Sin(ang) * spd

	dmg := 48.0
	isCrit := rand.Float64() < (g.player.CritChance * 2.0)
	if isCrit {
		dmg *= g.player.CritMult
	}

	proj := g.makePlayerProjectile(px, py, vx, vy, 3.8, dmg, color.RGBA{255, 225, 60, 255}, 95, entities.ProjectileKindDeadeye)
	proj.IsHoming = true
	proj.IsCrit = isCrit
	proj.Knockback = 3.5
	g.projectiles = append(g.projectiles, proj)
	audio.Play(audio.SFXRevolver)
}

func (g *GameplayScene) spawnEnemy() {
	g.nextEnemyID++
	edge := rand.Intn(4)
	var sx, sy float64

	switch edge {
	case 0:
		sx = rand.Float64() * 800
		sy = 64
	case 1:
		sx = 775
		sy = 68 + rand.Float64()*515
	case 2:
		sx = rand.Float64() * 800
		sy = 575
	case 3:
		sx = 18
		sy = 68 + rand.Float64()*515
	}

	waveMult := 1.0 + float64(g.wave-1)*0.16
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

func (g *GameplayScene) triggerChaosAmbush(px, py float64) {
	g.shockwaves = append(g.shockwaves, entities.NewShockwave(px, py, 95, color.RGBA{220, 25, 45, 230}, 45, true))
	g.screenFlashClr = color.RGBA{140, 10, 25, 130}
	g.screenFlashTimer = 28
	g.addScreenShake(6.5)
	audio.Play(audio.SFXBossAlert)

	count := 5
	dist := 150.0
	for i := 0; i < count; i++ {
		g.nextEnemyID++
		angle := float64(i)*(math.Pi*2/float64(count)) + rand.Float64()*0.4
		sx := px + dist*math.Cos(angle)
		sy := py + dist*math.Sin(angle)
		if sx < 20 {
			sx = 20
		}
		if sx > 775 {
			sx = 775
		}
		if sy < 64 {
			sy = 64
		}
		if sy > 575 {
			sy = 575
		}

		waveMult := 1.0 + float64(g.wave-1)*0.15
		bat := entities.NewEnemy(g.nextEnemyID, entities.EnemyTypeBat, sx, sy, waveMult)
		g.enemies = append(g.enemies, bat)
		g.particles = append(g.particles, entities.SpawnBurst(sx, sy, 8, 2.5, 2.5, color.RGBA{210, 30, 50, 240}, 20)...)
	}
}

func (g *GameplayScene) onEnemyDefeated(e *entities.Enemy) {
	audio.Play(audio.SFXExplosion)
	g.score += e.ScoreValue * g.wave
	g.totalKills++
	g.killsInWave++

	if g.player.ClassType == entities.ClassSheriff || g.player.ClassType == entities.ClassShotgunner {
		g.recentFirearmKills = append(g.recentFirearmKills, g.survivalFrames)
		cutoff := g.survivalFrames - 300
		validIdx := 0
		for validIdx < len(g.recentFirearmKills) && g.recentFirearmKills[validIdx] < cutoff {
			validIdx++
		}
		if validIdx > 0 {
			g.recentFirearmKills = g.recentFirearmKills[validIdx:]
		}
		if len(g.recentFirearmKills) >= 15 {
			g.tryUnlockAchievement("quick_trigger")
		}
	}

	if g.weather.Active == world.WeatherBloodMoon {
		save := progression.GetSave()
		save.BloodMoonKills++
		if save.BloodMoonKills >= 25 {
			g.tryUnlockAchievement("blood_moon_hunter")
		}
	}

	soulDropChance := 0.18
	soulAmount := 1
	if e.Type == entities.EnemyTypeTank {
		soulDropChance = 1.0
		soulAmount = 2 + rand.Intn(2)
	} else if e.Type == entities.EnemyTypeShooter || e.Type == entities.EnemyTypeRunner {
		soulDropChance = 0.35
		soulAmount = 1
	}

	if g.weather.Active == world.WeatherBloodMoon {
		soulAmount *= 2
	}

	if rand.Float64() < soulDropChance {
		soulDrop := entities.NewSoulDrop(e.X+e.Width/2, e.Y+e.Height/2, soulAmount)
		g.drops = append(g.drops, soulDrop)
	}

	if rand.Float64() < 0.28 {
		g.rollRegularDrop(e.X+e.Width/2, e.Y+e.Height/2)
	}

	if (g.boss == nil || !g.boss.Alive) && g.killsInWave >= g.killsToWave && g.wave < TotalCampaignWaves {
		g.advanceToWave(g.wave + 1)
	}
}

func (g *GameplayScene) advanceToWave(nextWave int) {
	if nextWave > TotalCampaignWaves {
		nextWave = TotalCampaignWaves
	}
	g.wave = nextWave
	g.killsInWave = 0
	g.killsToWave += 10

	save := progression.GetSave()
	if g.wave > save.MaxWaveReached {
		save.MaxWaveReached = g.wave
		_ = progression.Save()
	}
	if g.wave >= 10 {
		g.tryUnlockAchievement("survivor")
	}

	if g.spawnDelay > 28 {
		g.spawnDelay -= 5
	}

	biomeIdx := ((g.wave - 1) / 5) % 3
	newBiome := world.BiomeType(biomeIdx)
	if newBiome != g.arena.CurrentBiome {
		g.arena.SetBiome(newBiome)
		g.addFloatingText(fmt.Sprintf("BIOMA: %s", newBiome.String()), 310, 220, color.RGBA{120, 240, 255, 255})
	}

	g.player.Heal(g.player.MaxHP * 0.25)
	if g.player.MaxMana > 0 {
		g.player.Mana = g.player.MaxMana
	}

	if g.wave%2 == 0 && g.weather.Active == world.WeatherNone && g.wave%5 != 0 {
		g.triggerWeatherEvent()
	}

	if g.wave == 3 || g.wave == 7 || g.wave == 12 {
		g.openRelicSelection()
	}

	if g.wave%5 == 0 {
		g.triggerBossWave()
	} else {
		g.addFloatingText(fmt.Sprintf("=== ONDA %d / %d ===", g.wave, TotalCampaignWaves), 330, 180, color.RGBA{R: 255, G: 230, B: 60, A: 255})
	}
}

func (g *GameplayScene) triggerWeatherEvent() {
	chosen := g.weather.TriggerRandom()
	g.weatherBannerText = "CLIMA: " + chosen.Name()
	g.weatherBannerTimer = 240
}

func (g *GameplayScene) triggerBossWave() {
	for _, e := range g.enemies {
		ex := e.X + e.Width/2
		ey := e.Y + e.Height/2
		g.particles = append(g.particles, entities.SpawnBurst(ex, ey, 12, 3.2, 3.0, color.RGBA{170, 110, 255, 240}, 30)...)
	}
	g.enemies = nil

	for i := len(g.projectiles) - 1; i >= 0; i-- {
		if g.projectiles[i].IsEnemy {
			g.projectiles = append(g.projectiles[:i], g.projectiles[i+1:]...)
		}
	}

	g.boss = entities.NewBoss(g.wave)
	g.bossDamagedPlayer = false
	g.bossAlertText = fmt.Sprintf("⚠ ALERTA: %s DESPERTOU! ⚠", g.boss.Name)
	g.bossAlertTimer = 180

	g.screenFlashClr = color.RGBA{255, 80, 40, 140}
	g.screenFlashTimer = 30
	g.addScreenShake(11.0)
	audio.Play(audio.SFXBossAlert)
	audio.PlayBGM(audio.TrackBoss)
}

func (g *GameplayScene) tryUnlockAchievement(id string) {
	if def, unlockedNow := progression.UnlockAchievement(id); unlockedNow {
		g.achieveBannerTxt = fmt.Sprintf("🏆 CONQUISTA DESBLOQUEADA: %s!", def.Title)
		g.achieveBannerTimer = 240
		audio.Play(audio.SFXRelic)
	}
}

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
	bonusSouls := g.score / 1200
	totalSouls := g.soulsCollected + bonusSouls

	progression.AddSoulFragments(totalSouls)
	save := progression.GetSave()
	save.TotalRuns++
	save.TotalKills += g.totalKills
	if g.wave > save.MaxWaveReached {
		save.MaxWaveReached = g.wave
	}
	_ = progression.Save()

	ebiten.SetCursorMode(ebiten.CursorModeVisible)
	g.navigator.SwitchToGameOver(g.score, g.wave, g.player.ClassType, totalSouls)
}

func (g *GameplayScene) addFloatingText(text string, x, y float64, clr color.RGBA) {
	g.floatingTexts = append(g.floatingTexts, &FloatingText{
		Text:     text,
		X:        x,
		Y:        y,
		Color:    clr,
		Lifetime: 50,
	})
}

// Draw renderiza o mundo na worldCanvas com Screen Shake e o HUD/Modais fixos por cima.
func (g *GameplayScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 14, 255})

	// 1. Renderiza o mundo de jogo na worldCanvas
	g.worldCanvas.Clear()
	g.arena.Draw(g.worldCanvas)
	g.weather.DrawGroundEffects(g.worldCanvas)

	for _, sw := range g.shockwaves {
		sw.Draw(g.worldCanvas)
	}
	for _, ar := range g.arrowRains {
		ar.Draw(g.worldCanvas)
	}
	for _, fp := range g.firePillars {
		fp.Draw(g.worldCanvas)
	}
	for _, d := range g.drops {
		d.Draw(g.worldCanvas)
	}
	for _, e := range g.enemies {
		e.Draw(g.worldCanvas)
	}
	if g.boss != nil && g.boss.Alive {
		g.boss.Draw(g.worldCanvas)
	}
	for _, p := range g.projectiles {
		p.Draw(g.worldCanvas)
	}
	for _, s := range g.slashes {
		s.Draw(g.worldCanvas)
	}

	g.player.Draw(g.worldCanvas)

	for _, arc := range g.lightningArcs {
		vector.StrokeLine(g.worldCanvas, float32(arc.X1), float32(arc.Y1), float32(arc.X2), float32(arc.Y2), 3.0, color.RGBA{120, 235, 255, 230}, true)
		vector.StrokeLine(g.worldCanvas, float32(arc.X1), float32(arc.Y1), float32(arc.X2), float32(arc.Y2), 1.5, color.RGBA{255, 255, 210, 255}, true)
	}
	for _, pt := range g.particles {
		pt.Draw(g.worldCanvas)
	}

	px, py := g.player.Center()
	g.weather.DrawOverlayEffects(g.worldCanvas, px, py)

	for _, ft := range g.floatingTexts {
		ebitenutil.DebugPrintAt(g.worldCanvas, ft.Text, int(ft.X), int(ft.Y))
	}

	// 2. Desenha worldCanvas na tela principal aplicando o deslocamento de Screen Shake
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.camOffsetX, g.camOffsetY)
	screen.DrawImage(g.worldCanvas, op)

	// 3. Flash de tela colorido
	if g.screenFlashTimer > 0 {
		clr := g.screenFlashClr
		clr.A = uint8(float64(g.screenFlashClr.A) * (float64(g.screenFlashTimer) / 30.0))
		vector.DrawFilledRect(screen, 0, 52, 800, 548, clr, false)
	}

	// 4. HUD Fixo (não treme durante Screen Shake)
	g.drawHUD(screen)

	// 5. Banners de Alerta
	if g.bossAlertTimer > 0 && (g.bossAlertTimer/10)%2 == 0 {
		vector.DrawFilledRect(screen, 190, 95, 420, 32, color.RGBA{120, 15, 25, 220}, false)
		vector.StrokeRect(screen, 190, 95, 420, 32, 2, color.RGBA{255, 200, 60, 255}, false)
		ebitenutil.DebugPrintAt(screen, g.bossAlertText, 240, 105)
	}
	if g.weatherBannerTimer > 0 {
		vector.DrawFilledRect(screen, 200, 132, 400, 26, color.RGBA{20, 38, 42, 220}, false)
		vector.StrokeRect(screen, 200, 132, 400, 26, 1.5, color.RGBA{100, 220, 180, 255}, false)
		ebitenutil.DebugPrintAt(screen, g.weatherBannerText, 235, 139)
	}
	if g.achieveBannerTimer > 0 {
		vector.DrawFilledRect(screen, 210, 548, 380, 32, color.RGBA{52, 38, 18, 240}, false)
		vector.StrokeRect(screen, 210, 548, 380, 32, 2, color.RGBA{255, 215, 60, 255}, false)
		ebitenutil.DebugPrintAt(screen, g.achieveBannerTxt, 235, 558)
	}

	// 6. Modais (Relíquia ou Pausa com Ficha da Run)
	if g.isChoosingRelic {
		g.drawRelicSelectionModal(screen)
	} else if g.isPaused {
		if g.inPauseOpts {
			g.drawPauseOptionsModal(screen)
		} else {
			g.drawPauseMenu(screen)
		}
	}
}

func (g *GameplayScene) drawHUD(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, 800, 52, color.RGBA{R: 14, G: 16, B: 24, A: 252}, false)
	vector.StrokeLine(screen, 0, 52, 800, 52, 2, color.RGBA{R: 65, G: 80, B: 115, A: 255}, false)

	g.hpBar.Draw(screen)

	if g.manaBar != nil {
		g.manaBar.Draw(screen)
	} else if g.player.ClassType == entities.ClassSheriff {
		g.drawSheriffCylinderHUD(screen)
	}

	g.specialBar.Draw(screen)
	if g.player.SpecialCD <= 0 {
		ebitenutil.DebugPrintAt(screen, "OK", 322, 10)
	} else {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.1fs", float64(g.player.SpecialCD)/60.0), 312, 10)
	}

	g.dashBar.Draw(screen)
	if g.player.DashCD <= 0 {
		ebitenutil.DebugPrintAt(screen, "OK", 322, 29)
	}

	var waveInfo string
	if g.boss != nil && g.boss.Alive {
		waveInfo = fmt.Sprintf("ONDA %d/%d (CHEFE)", g.wave, TotalCampaignWaves)
	} else {
		waveInfo = fmt.Sprintf("ONDA %d/%d (%d/%d)", g.wave, TotalCampaignWaves, g.killsInWave, g.killsToWave)
	}
	ebitenutil.DebugPrintAt(screen, waveInfo, 350, 10)

	if g.soulIcon != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1.4, 1.4)
		op.GeoM.Translate(488, 6)
		screen.DrawImage(g.soulIcon, op)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("+%d", g.soulsCollected), 510, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Reliq:%d", g.player.TotalRelics), 555, 10)

	totalSeconds := g.survivalFrames / 60
	timeStr := fmt.Sprintf("%02d:%02d", totalSeconds/60, totalSeconds%60)
	scoreStr := fmt.Sprintf("Pts:%d", g.score)
	ebitenutil.DebugPrintAt(screen, timeStr, 630, 10)
	ebitenutil.DebugPrintAt(screen, scoreStr, 690, 10)

	effects := g.player.GetActiveEffects()
	if len(effects) > 0 {
		effX := 350
		for _, eff := range effects {
			secsLeft := (eff.Remaining + 59) / 60
			label := fmt.Sprintf("%s (%ds)", eff.Name, secsLeft)
			badgeW := float32(len(label)*6 + 10)
			vector.DrawFilledRect(screen, float32(effX), 30, badgeW, 16, color.RGBA{R: 28, G: 32, B: 45, A: 220}, false)
			vector.StrokeRect(screen, float32(effX), 30, badgeW, 16, 1, eff.Color, false)
			ebitenutil.DebugPrintAt(screen, label, effX+5, 32)
			effX += int(badgeW) + 6
		}
	} else if g.weather.Active != world.WeatherNone {
		secsLeft := (g.weather.Timer + 59) / 60
		wStr := fmt.Sprintf("Clima (%ds): %s", secsLeft, g.weather.Active.Name())
		ebitenutil.DebugPrintAt(screen, wStr, 350, 32)
	} else {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[%s] | [ESC/P] Pausa", g.arena.CurrentBiome.String()), 350, 32)
	}

	if g.boss != nil && g.boss.Alive {
		g.bossBar.Draw(screen)
	}
}

func (g *GameplayScene) drawSheriffCylinderHUD(dst *ebiten.Image) {
	baseX := float32(12)
	baseY := float32(30)

	if g.player.ReloadTimer > 0 {
		ebitenutil.DebugPrintAt(dst, "RECARREGANDO...", int(baseX), int(baseY))
		return
	}

	for i := 0; i < g.player.MaxAmmo; i++ {
		cx := baseX + float32(i*14) + 6
		cy := baseY + 6
		if i < g.player.Ammo {
			vector.DrawFilledCircle(dst, cx, cy, 5, color.RGBA{235, 185, 50, 255}, true)
			vector.DrawFilledCircle(dst, cx, cy, 2, color.RGBA{255, 245, 200, 255}, true)
		} else {
			vector.StrokeCircle(dst, cx, cy, 4.5, 1.2, color.RGBA{100, 105, 120, 200}, true)
		}
	}
}

func (g *GameplayScene) drawRelicSelectionModal(dst *ebiten.Image) {
	vector.DrawFilledRect(dst, 0, 0, 800, 600, color.RGBA{R: 10, G: 8, B: 14, A: 215}, false)

	vector.DrawFilledRect(dst, 150, 65, 500, 72, color.RGBA{R: 48, G: 32, B: 22, A: 255}, false)
	vector.StrokeRect(dst, 150, 65, 500, 72, 2.5, color.RGBA{R: 245, G: 195, B: 60, A: 255}, false)
	ebitenutil.DebugPrintAt(dst, "=== ⚜ ESCOLHA UMA RELIQUIA SAGRADA ⚜ ===", 265, 80)
	ebitenutil.DebugPrintAt(dst, "Clique em um dos 3 artefatos para absorver seu poder nesta jornada", 195, 104)

	mx, my := ebiten.CursorPosition()
	for idx, rel := range g.relicChoices {
		if idx >= 3 {
			break
		}
		b := g.relicCardBounds[idx]
		x, y, w, h := float32(b.X), float32(b.Y), float32(b.W), float32(b.H)
		hovered := b.Contains(float64(mx), float64(my))

		bg := color.RGBA{R: 38, G: 28, B: 24, A: 250}
		border := rel.Color
		stroke := float32(2.0)
		if hovered {
			bg = color.RGBA{R: 62, G: 46, B: 36, A: 255}
			stroke = 3.5
		}

		vector.DrawFilledRect(dst, x, y, w, h, bg, false)
		vector.StrokeRect(dst, x, y, w, h, stroke, border, false)

		iconX, iconY := x+w/2, y+55
		// Moldura ornamental de pedestal medieval para o artefato
		vector.DrawFilledRect(dst, iconX-28, iconY-28, 56, 56, color.RGBA{18, 14, 20, 240}, false)
		vector.StrokeRect(dst, iconX-28, iconY-28, 56, 56, 2, rel.Color, false)
		vector.StrokeRect(dst, iconX-25, iconY-25, 50, 50, 1, color.RGBA{255, 215, 80, 140}, false)

		if sprite := rel.GetSprite(); sprite != nil {
			op := &ebiten.DrawImageOptions{}
			scale := 2.5
			op.GeoM.Scale(scale, scale)
			sw := float64(sprite.Bounds().Dx()) * scale
			sh := float64(sprite.Bounds().Dy()) * scale
			op.GeoM.Translate(float64(iconX)-sw/2, float64(iconY)-sh/2)
			dst.DrawImage(sprite, op)
		}

		titleX := int(x) + (int(w)-len(rel.Name)*6)/2
		ebitenutil.DebugPrintAt(dst, rel.Name, titleX, int(y)+95)

		stacks := g.player.Relics[rel.ID]
		if stacks > 0 {
			ebitenutil.DebugPrintAt(dst, fmt.Sprintf("(Possui: %d)", stacks), int(x)+70, int(y)+113)
		}

		g.drawWrappedText(dst, rel.EffectDesc, int(x)+16, int(y)+140, 29)

		if hovered {
			ebitenutil.DebugPrintAt(dst, ">> CLICAR PARA ESCOLHER <<", int(x)+26, int(y+h)-28)
		}
	}
}

func (g *GameplayScene) drawWrappedText(dst *ebiten.Image, text string, x, y, maxChars int) {
	line := 0
	for len(text) > maxChars {
		cut := maxChars
		for cut > 10 && text[cut] != ' ' {
			cut--
		}
		ebitenutil.DebugPrintAt(dst, text[:cut], x, y+line*16)
		text = text[cut+1:]
		line++
	}
	if len(text) > 0 {
		ebitenutil.DebugPrintAt(dst, text, x, y+line*16)
	}
}

// drawPauseMenu exibe o menu de pausa à esquerda e a Ficha Completa da Run à direita.
func (g *GameplayScene) drawPauseMenu(dst *ebiten.Image) {
	vector.DrawFilledRect(dst, 0, 0, 800, 600, color.RGBA{R: 10, G: 8, B: 6, A: 215}, false)

	// Painel Principal amplo (700x380)
	panelW := float32(700)
	panelH := float32(380)
	px := float32((800 - panelW) / 2)
	py := float32((600 - panelH) / 2)

	vector.DrawFilledRect(dst, px+5, py+6, panelW, panelH, color.RGBA{R: 5, G: 4, B: 3, A: 230}, false)
	vector.DrawFilledRect(dst, px, py, panelW, panelH, color.RGBA{R: 44, G: 30, B: 20, A: 255}, false)
	vector.StrokeRect(dst, px, py, panelW, panelH, 3.0, color.RGBA{R: 190, G: 150, B: 75, A: 255}, false)
	vector.StrokeRect(dst, px+5, py+5, panelW-10, panelH-10, 1.2, color.RGBA{R: 90, G: 70, B: 40, A: 200}, false)

	drawRivet(dst, px+12, py+12, color.RGBA{R: 245, G: 205, B: 100, A: 255})
	drawRivet(dst, px+panelW-12, py+12, color.RGBA{R: 245, G: 205, B: 100, A: 255})
	drawRivet(dst, px+12, py+panelH-12, color.RGBA{R: 245, G: 205, B: 100, A: 255})
	drawRivet(dst, px+panelW-12, py+panelH-12, color.RGBA{R: 245, G: 205, B: 100, A: 255})

	// Título da Pausa
	ebitenutil.DebugPrintAt(dst, "=== ⚜ BATALHA PAUSADA ⚜ ===", int(px)+100, int(py)+25)
	ebitenutil.DebugPrintAt(dst, "Pressione [ESC] ou [P] para retornar ao combate", int(px)+55, int(py)+46)

	// Linha divisória vertical entre Botões (Esquerda) e Ficha da Run (Direita)
	vector.StrokeLine(dst, px+315, py+70, px+315, py+panelH-25, 1.5, color.RGBA{R: 120, G: 90, B: 50, A: 200}, false)

	// Botões à esquerda
	g.btnResume.Draw(dst)
	g.btnPauseOpts.Draw(dst)
	g.btnRestart.Draw(dst)
	g.btnMainMenu.Draw(dst)

	// Ficha da Run Atual à direita
	fx := int(px) + 335
	fy := int(py) + 25
	ebitenutil.DebugPrintAt(dst, "[ FICHA DA EXPEDICAO ATUAL ]", fx+45, fy)

	totalSeconds := g.survivalFrames / 60
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Heroi: %s   |   Onda: %d / %d", g.player.ClassType, g.wave, TotalCampaignWaves), fx, fy+28)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Tempo: %02d:%02d   |   Abates: %d", totalSeconds/60, totalSeconds%60, g.totalKills), fx, fy+48)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Dano Total: %.0f   |   Almas: %d", g.totalDamageDealt, g.soulsCollected), fx, fy+68)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Mult. Dano: %.2fx   |   Vel. Base: %.1f", g.player.GetDamageMultiplier(), g.player.BaseSpeed), fx, fy+88)

	vector.StrokeLine(dst, float32(fx), float32(fy+112), px+panelW-25, float32(fy+112), 1, color.RGBA{120, 90, 50, 180}, false)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Reliquias Passivas Ativas (%d):", g.player.TotalRelics), fx, fy+122)

	ry := fy + 144
	hasAnyRelic := false
	for _, rel := range entities.AllRelics {
		stacks := g.player.Relics[rel.ID]
		if stacks > 0 {
			hasAnyRelic = true
			if sprite := rel.GetSprite(); sprite != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(1.0, 1.0)
				op.GeoM.Translate(float64(fx), float64(ry-2))
				dst.DrawImage(sprite, op)
			} else {
				vector.DrawFilledCircle(dst, float32(fx+6), float32(ry+5), 4, rel.Color, true)
			}
			ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%s (x%d)", rel.Name, stacks), fx+22, ry)
			ry += 20
		}
	}
	if !hasAnyRelic {
		ebitenutil.DebugPrintAt(dst, "Nenhuma reliquia obtida ainda.", fx, ry)
		ebitenutil.DebugPrintAt(dst, "Derrote Chefes ou alcance marcos de onda!", fx, ry+18)
	}
}

func (g *GameplayScene) drawPauseOptionsModal(dst *ebiten.Image) {
	vector.DrawFilledRect(dst, 0, 0, 800, 600, color.RGBA{R: 10, G: 8, B: 6, A: 225}, false)

	panelW := float32(420)
	panelH := float32(340)
	px := float32((800 - panelW) / 2)
	py := float32((600 - panelH) / 2)

	vector.DrawFilledRect(dst, px, py, panelW, panelH, color.RGBA{R: 44, G: 30, B: 20, A: 255}, false)
	vector.StrokeRect(dst, px, py, panelW, panelH, 3.0, color.RGBA{R: 190, G: 150, B: 75, A: 255}, false)
	ebitenutil.DebugPrintAt(dst, "=== ⚙ OPCOES E AUDIO ⚙ ===", int(px)+130, int(py)+22)

	g.btnPauseFS.Draw(dst)

	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Volume Mestre: %d%%", g.navigator.GetMasterVolume()), 340, 232)
	g.btnPauseMasterDown.Draw(dst)
	g.btnPauseMasterUp.Draw(dst)

	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Efeitos (SFX): %d%%", g.navigator.GetSFXVolume()), 340, 282)
	g.btnPauseSFXDown.Draw(dst)
	g.btnPauseSFXUp.Draw(dst)

	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Musica (BGM):  %d%%", g.navigator.GetMusicVolume()), 340, 332)
	g.btnPauseBGMDown.Draw(dst)
	g.btnPauseBGMUp.Draw(dst)

	g.btnPauseBack.Draw(dst)
}
