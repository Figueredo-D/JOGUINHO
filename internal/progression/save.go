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

// SaveData representa os dados persistidos da conta / meta-progressão.
type SaveData struct {
	SoulFragments int             `json:"soul_fragments"`
	TotalRuns     int             `json:"total_runs"`
	TotalKills    int             `json:"total_kills"`
	Warrior       WarriorUpgrades `json:"warrior"`
	Archer        ArcherUpgrades  `json:"archer"`
	Mage          MageUpgrades    `json:"mage"`
}

var (
	currentSave *SaveData
	saveMutex   sync.Mutex
)

// DefaultSaveData retorna os dados iniciais padrão para novos jogadores.
func DefaultSaveData() *SaveData {
	return &SaveData{
		SoulFragments: 0,
		TotalRuns:     0,
		TotalKills:    0,
		Warrior:       WarriorUpgrades{},
		Archer:        ArcherUpgrades{},
		Mage:          MageUpgrades{},
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

	return &s
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
