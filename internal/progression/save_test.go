package progression

import (
	"os"
	"testing"
)

func TestProgressionSaveAndLoad(t *testing.T) {
	testFile := "test_save.json"
	defer os.Remove(testFile)

	save := DefaultSaveData()
	save.SoulFragments = 150
	save.Warrior.MaxHPLevel = 2
	save.Archer.SpeedLevel = 3

	err := saveToPath(save, testFile)
	if err != nil {
		t.Fatalf("Erro ao salvar: %v", err)
	}

	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Erro ao ler save: %v", err)
	}

	if len(data) == 0 {
		t.Fatalf("Arquivo de save vazio")
	}
}

func TestUpgradeCost(t *testing.T) {
	cost0 := UpgradeCost(0)
	cost1 := UpgradeCost(1)
	costMax := UpgradeCost(MaxLevel)

	if cost0 <= 0 || cost1 <= cost0 {
		t.Errorf("Custos devem ser crescentes: nvl 0=%d, nvl 1=%d", cost0, cost1)
	}

	if costMax != 0 {
		t.Errorf("Custo para nível máximo deveria ser 0, obteve %d", costMax)
	}
}
