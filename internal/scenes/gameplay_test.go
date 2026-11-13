package scenes

import (
	"image/color"
	"testing"

	"joguinho/internal/entities"
	"joguinho/internal/world"
)

type mockNavigator struct {
	victoryCalled bool
	lastSummary   RunSummary
}

func (m *mockNavigator) SwitchToMenu()                                                       {}
func (m *mockNavigator) SwitchToClassSelect()                                                {}
func (m *mockNavigator) SwitchToSanctuary()                                                  {}
func (m *mockNavigator) SwitchToAchievements()                                               {}
func (m *mockNavigator) SwitchToHowToPlay()                                                  {}
func (m *mockNavigator) SwitchToGameplay(class entities.ClassType)                           {}
func (m *mockNavigator) SwitchToGameOver(score int, wave int, class entities.ClassType, s int) {}
func (m *mockNavigator) SwitchToVictory(summary RunSummary) {
	m.victoryCalled = true
	m.lastSummary = summary
}
func (m *mockNavigator) SwitchToOptions()      {}
func (m *mockNavigator) ToggleFullscreen()     {}
func (m *mockNavigator) IsFullscreen() bool    { return false }
func (m *mockNavigator) SetMasterVolume(v int) {}
func (m *mockNavigator) GetMasterVolume() int  { return 80 }
func (m *mockNavigator) SetSFXVolume(v int)    {}
func (m *mockNavigator) GetSFXVolume() int     { return 80 }
func (m *mockNavigator) SetMusicVolume(v int)  {}
func (m *mockNavigator) GetMusicVolume() int   { return 60 }
func (m *mockNavigator) QuitGame()             {}

func TestBossWaveTriggerDoesNotPanicOnEnemiesSlice(t *testing.T) {
	gs := NewGameplayScene(&mockNavigator{}, entities.ClassWarrior)
	gs.wave = 4
	gs.killsInWave = gs.killsToWave - 1

	for i := 0; i < 10; i++ {
		gs.enemies = append(gs.enemies, entities.NewEnemy(uint64(i+1), entities.EnemyTypeCommon, 100, 100, 1.0))
	}

	gs.slashes = append(gs.slashes, entities.NewSlash(100, 100, 120, 120, 80, 500))

	if err := gs.Update(); err != nil {
		t.Fatalf("Update retornou erro: %v", err)
	}

	if gs.boss == nil || !gs.boss.Alive {
		t.Errorf("Onda 5 deveria ter despertado o Chefe após o Board Clear")
	}
	if len(gs.enemies) != 0 {
		t.Errorf("Board Clear deveria ter zerado inimigos comuns, mas restaram %d", len(gs.enemies))
	}
}

func TestDashGrantsIFramesAndCooldown(t *testing.T) {
	p := entities.NewPlayer(entities.ClassArcher, 400, 300)
	if !p.TriggerDash() {
		t.Fatalf("Primeiro dash deveria ter sido ativado")
	}
	if p.InvincibleCD <= 0 {
		t.Errorf("Dash deveria conceder I-frames")
	}
	if p.TriggerDash() {
		t.Errorf("Segundo dash imediato não deveria ativar durante cooldown")
	}
}

func TestRelicBloodPactAndIronHeart(t *testing.T) {
	p := entities.NewPlayer(entities.ClassWarrior, 400, 300)
	baseHP := p.MaxHP
	baseDmg := p.GetDamageMultiplier()

	p.AddRelic(entities.RelicIronHeart)
	if p.MaxHP != baseHP+30 {
		t.Errorf("Coração de Ferro deveria adicionar +30 MaxHP, obteve %v", p.MaxHP)
	}

	p.AddRelic(entities.RelicBloodPact)
	if p.GetDamageMultiplier() <= baseDmg {
		t.Errorf("Pacto de Sangue deveria aumentar dano global")
	}

	p.HP = 100
	healed := p.Heal(40)
	if healed != 20 {
		t.Errorf("Pacto de Sangue deveria reduzir cura pela metade (esperado 20, obteve %v)", healed)
	}
}

func TestWeatherSystemTriggers(t *testing.T) {
	ws := world.NewWeatherSystem()
	if ws.Active != world.WeatherNone {
		t.Errorf("Clima inicial deveria ser WeatherNone")
	}
	chosen := ws.TriggerRandom()
	if chosen == world.WeatherNone || ws.Timer <= 0 {
		t.Errorf("TriggerRandom deveria ativar um evento climático com timer positivo")
	}
}

func TestFinalBossVictoryFlow(t *testing.T) {
	nav := &mockNavigator{}
	gs := NewGameplayScene(nav, entities.ClassMage)
	gs.wave = 15
	gs.soulsCollected = 10
	gs.boss = entities.NewBoss(15)
	if !gs.boss.IsFinal {
		t.Fatalf("Boss da onda 15 deveria ser o Soberano de Dorfieegu (IsFinal=true)")
	}

	// Adiciona 2 projéteis inimigos ativos que devem se converter em almas na vitória
	gs.projectiles = append(gs.projectiles,
		entities.NewProjectile(200, 200, 1, 1, 5, 10, color.RGBA{255, 50, 50, 255}, 50, entities.ProjectileKindBossOrb, true),
		entities.NewProjectile(250, 250, 1, 1, 5, 10, color.RGBA{255, 50, 50, 255}, 50, entities.ProjectileKindBossOrb, true),
	)

	// Mata o Soberano
	gs.damageBoss(99999, false)

	if !nav.victoryCalled {
		t.Fatalf("Derrotar o Soberano de Dorfieegu deveria chamar SwitchToVictory")
	}
	if nav.lastSummary.TotalSouls <= nav.lastSummary.BaseSouls {
		t.Errorf("Vitória deveria aplicar bônus multiplicador 1.5x nas almas (base=%d, total=%d)", nav.lastSummary.BaseSouls, nav.lastSummary.TotalSouls)
	}
	if nav.lastSummary.DamageDealt <= 0 {
		t.Errorf("Relatório de vitória deveria registrar dano total causado")
	}
}

func TestHitstopAndScreenShake(t *testing.T) {
	gs := NewGameplayScene(&mockNavigator{}, entities.ClassShotgunner)
	gs.triggerHitstop(3)
	gs.addScreenShake(10.0)

	prevFrames := gs.survivalFrames
	if err := gs.Update(); err != nil {
		t.Fatalf("Update falhou: %v", err)
	}

	if gs.survivalFrames != prevFrames {
		t.Errorf("Durante hitstop, survivalFrames não deveria avançar")
	}
	if gs.shakeIntensity >= 10.0 {
		t.Errorf("Shake deveria decair exponencialmente a cada frame")
	}
}
