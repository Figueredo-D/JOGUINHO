package entities

import (
	"testing"
)

func TestBossBuffScaling(t *testing.T) {
	// Onda 5
	b5 := NewBoss(5)
	expectedHP5 := 1500.0 * 2.5
	expectedDmg5 := 28.0 * 2.0
	if b5.MaxHP != expectedHP5 {
		t.Errorf("Boss Onda 5 MaxHP: esperado %v, obteve %v", expectedHP5, b5.MaxHP)
	}
	if b5.Damage != expectedDmg5 {
		t.Errorf("Boss Onda 5 Dano: esperado %v, obteve %v", expectedDmg5, b5.Damage)
	}

	// Onda 10
	b10 := NewBoss(10)
	expectedHP10 := 2700.0 * 2.5
	expectedDmg10 := 38.0 * 2.0
	if b10.MaxHP != expectedHP10 {
		t.Errorf("Boss Onda 10 MaxHP: esperado %v, obteve %v", expectedHP10, b10.MaxHP)
	}
	if b10.Damage != expectedDmg10 {
		t.Errorf("Boss Onda 10 Dano: esperado %v, obteve %v", expectedDmg10, b10.Damage)
	}

	// Onda 15 (Chefe Final)
	b15 := NewBoss(15)
	expectedHP15 := 4200.0 * 2.5
	expectedDmg15 := 42.0 * 2.0
	if b15.MaxHP != expectedHP15 {
		t.Errorf("Boss Onda 15 MaxHP: esperado %v, obteve %v", expectedHP15, b15.MaxHP)
	}
	if b15.Damage != expectedDmg15 {
		t.Errorf("Boss Onda 15 Dano: esperado %v, obteve %v", expectedDmg15, b15.Damage)
	}
	if !b15.IsFinal {
		t.Errorf("Boss Onda 15 deveria ter IsFinal = true")
	}

	// Onda 20 (Suporte a ondas superiores)
	b20 := NewBoss(20)
	if b20.MaxHP <= b15.MaxHP {
		t.Errorf("Boss Onda 20 deveria ter mais vida que onda 15")
	}
	if b20.Damage <= b15.Damage {
		t.Errorf("Boss Onda 20 deveria ter mais dano que onda 15")
	}
}

func TestBossAggressionBelowHalfHP(t *testing.T) {
	b := NewBoss(5)
	// Reduz vida para menos de 50%
	b.HP = b.MaxHP * 0.40
	b.State = BossStateChase
	b.FrameCount = 45 // Update incrementa para 46 (46 % 46 == 0 dispara rajada)

	out := b.Update(400, 500)
	if len(out.Projectiles) < 5 {
		t.Errorf("Boss com < 50%% HP deveria disparar rajada em leque (esperado >= 5 projéteis, obteve %d)", len(out.Projectiles))
	}
}
