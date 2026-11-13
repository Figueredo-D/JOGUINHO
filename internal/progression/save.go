package progression

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

const (
	SaveFileName = "save.json"
	MaxLevel     = 5
)

// WarriorUpgrades contém os níveis de melhoria do Guerreiro.
type WarriorUpgrades struct {
	MaxHPLevel  int `json:"max_hp_level"`  // +25 HP por nível
	DamageLevel int `json:"damage_level"`  // +12% dano e +6 alcance por nível
	ShieldLevel int `json:"shield_level"`  // -1.0s de recarga do escudo e +0.25s de duração por nível
}

// ArcherUpgrades contém os níveis de melhoria do Arqueiro.
type ArcherUpgrades struct {
	SpeedLevel     int `json:"speed_level"`      // +0.35 vel. movimento por nível
	AttackSpeedLvl int `json:"attack_speed_lvl"` // +10% vel. ataque e +1.2 vel. projétil por nível
	ArrowRainLevel int `json:"arrow_rain_level"` // +15% raio de impacto e +20% dano por nível
}

// MageUpgrades contém os níveis de melhoria do Mago.
type MageUpgrades struct {
	ManaLevel       int `json:"mana_level"`        // +20 MaxMana e +3.5 regen/s por nível
	SpellPowerLevel int `json:"spell_power_level"` // +15% dano base de magia por nível
	FirePillarLevel int `json:"fire_pillar_level"` // +15% raio de erupção e +20 frames de fogo por nível
}

// ShotgunnerUpgrades contém os níveis de melhoria do Atirador Pesado.
type ShotgunnerUpgrades struct {
	SpreadPelletsLevel int `json:"spread_pellets_level"` // +1 balote e menor dispersão por nível
	ShrapnelDmgLevel   int `json:"shrapnel_dmg_level"`   // +15% dano dos balotes e +knockback por nível
	ShockwaveLevel     int `json:"shockwave_level"`      // +20px alcance e +20% dano do Choque Sísmico
}

// SheriffUpgrades contém os níveis de melhoria do Xerife.
type SheriffUpgrades struct {
	ReloadSpeedLevel int `json:"reload_speed_level"` // -0.15s tempo de recarga do tambor por nível
	CritDamageLevel  int `json:"crit_damage_level"`  // +10% chance crítica e +20% dano crítico por nível
	ExtraAmmoLevel   int `json:"extra_ammo_level"`   // +1 bala adicional no tambor por nível
}

// AchievementDef define os metadados de uma conquista do jogo.
type AchievementDef struct {
	ID          string
	Title       string
	Description string
	Goal        int
}

// AllAchievements lista todas as conquistas disponíveis.
var AllAchievements = []AchievementDef{
	{
		ID:          "survivor",
		Title:       "Sobrevivente",
		Description: "Sobreviva ate a Onda 10 em uma unica jornada.",
		Goal:        10,
	},
	{
		ID:          "giant_slayer",
		Title:       "Mata-Gigantes",
		Description: "Derrote um Chefe (Boss) sem sofrer dano durante o combate contra ele.",
		Goal:        1,
	},
	{
		ID:          "quick_trigger",
		Title:       "Gatilho Rapido",
		Description: "Elimine 15 inimigos em menos de 5s usando Xerife ou Atirador Pesado.",
		Goal:        15,
	},
	{
		ID:          "iron_wall",
		Title:       "Muralha de Ferro",
		Description: "Reflita 30 projeteis inimigos usando o Escudo do Guerreiro.",
		Goal:        30,
	},
	{
		ID:          "relic_hoarder",
		Title:       "Colecionador de Reliquias",
		Description: "Acumule 4 ou mais Reliquias Passivas em uma mesma partida.",
		Goal:        4,
	},
	{
		ID:          "blood_moon_hunter",
		Title:       "Cacador da Lua Carmesim",
		Description: "Elimine 25 inimigos durante o evento climatico Lua Carmesim.",
		Goal:        25,
	},
}

// GameSettings armazena preferências de áudio e vídeo persistidas em save.json.
type GameSettings struct {
	Initialized  bool `json:"initialized"`
	MasterVolume int  `json:"master_volume"`
	SFXVolume    int  `json:"sfx_volume"`
	MusicVolume  int  `json:"music_volume"`
	Fullscreen   bool `json:"fullscreen"`
}

// DefaultSettings retorna as configurações padrão da Versão 1.0.
func DefaultSettings() GameSettings {
	return GameSettings{
		Initialized:  true,
		MasterVolume: 80,
		SFXVolume:    80,
		MusicVolume:  60,
		Fullscreen:   false,
	}
}

// SaveData representa os dados persistidos da conta / meta-progressão.
type SaveData struct {
	SoulFragments        int                `json:"soul_fragments"`
	TotalRuns            int                `json:"total_runs"`
	TotalKills           int                `json:"total_kills"`
	MaxWaveReached       int                `json:"max_wave_reached"`
	ShieldReflectCount   int                `json:"shield_reflect_count"`
	BloodMoonKills       int                `json:"blood_moon_kills"`
	UnlockedAchievements map[string]bool    `json:"unlocked_achievements"`
	Settings             GameSettings       `json:"settings"`
	Warrior              WarriorUpgrades    `json:"warrior"`
	Archer               ArcherUpgrades     `json:"archer"`
	Mage                 MageUpgrades       `json:"mage"`
	Shotgunner           ShotgunnerUpgrades `json:"shotgunner"`
	Sheriff              SheriffUpgrades    `json:"sheriff"`
}

var (
	currentSave *SaveData
	saveMutex   sync.Mutex
)

// DefaultSaveData retorna os dados iniciais padrão para novos jogadores.
func DefaultSaveData() *SaveData {
	return &SaveData{
		SoulFragments:        0,
		TotalRuns:            0,
		TotalKills:           0,
		UnlockedAchievements: make(map[string]bool),
		Settings:             DefaultSettings(),
		Warrior:              WarriorUpgrades{},
		Archer:               ArcherUpgrades{},
		Mage:                 MageUpgrades{},
		Shotgunner:           ShotgunnerUpgrades{},
		Sheriff:              SheriffUpgrades{},
	}
}

// GetSave retorna o save atual em memória (carregando se necessário).
func GetSave() *SaveData {
	saveMutex.Lock()
	defer saveMutex.Unlock()

	if currentSave == nil {
		currentSave = loadSaveFromDisk()
	}
	return currentSave
}

// Save grava o estado atual no disco.
func Save() error {
	saveMutex.Lock()
	defer saveMutex.Unlock()

	if currentSave == nil {
		currentSave = DefaultSaveData()
	}

	data, err := json.MarshalIndent(currentSave, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(SaveFileName, data, 0644)
}

func loadSaveFromDisk() *SaveData {
	// Verifica se arquivo existe
	if _, err := os.Stat(SaveFileName); os.IsNotExist(err) {
		s := DefaultSaveData()
		_ = saveToPath(s, SaveFileName)
		return s
	}

	data, err := os.ReadFile(SaveFileName)
	if err != nil {
		return DefaultSaveData()
	}

	var s SaveData
	if err := json.Unmarshal(data, &s); err != nil {
		return DefaultSaveData()
	}
	if s.UnlockedAchievements == nil {
		s.UnlockedAchievements = make(map[string]bool)
	}
	if !s.Settings.Initialized {
		s.Settings = DefaultSettings()
	}

	return &s
}

// UnlockAchievement marca uma conquista como concluída caso ainda não esteja e persiste no save.
// Retorna a definição da conquista e true se foi desbloqueada agora pela primeira vez.
func UnlockAchievement(id string) (AchievementDef, bool) {
	s := GetSave()
	saveMutex.Lock()
	if s.UnlockedAchievements == nil {
		s.UnlockedAchievements = make(map[string]bool)
	}
	already := s.UnlockedAchievements[id]
	if already {
		saveMutex.Unlock()
		return AchievementDef{}, false
	}
	s.UnlockedAchievements[id] = true
	saveMutex.Unlock()

	_ = Save()

	for _, a := range AllAchievements {
		if a.ID == id {
			return a, true
		}
	}
	return AchievementDef{ID: id, Title: id}, true
}

func saveToPath(s *SaveData, path string) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// AddSoulFragments adiciona fragmentos e persiste.
func AddSoulFragments(amount int) {
	s := GetSave()
	saveMutex.Lock()
	s.SoulFragments += amount
	saveMutex.Unlock()
	_ = Save()
}

// UpgradeCost calcula o custo em fragmentos para comprar o próximo nível.
func UpgradeCost(currentLevel int) int {
	if currentLevel >= MaxLevel {
		return 0
	}
	// Nível 0 -> 1: 15
	// Nível 1 -> 2: 30
	// Nível 2 -> 3: 55
	// Nível 3 -> 4: 90
	// Nível 4 -> 5: 140
	costs := []int{15, 30, 55, 90, 140}
	if currentLevel < len(costs) {
		return costs[currentLevel]
	}
	return 150
}
