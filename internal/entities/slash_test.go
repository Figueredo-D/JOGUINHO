package entities

import (
	"testing"

	"joguinho/internal/collision"
)

func TestSlashLifecycleAndSparks(t *testing.T) {
	s := NewSlash(100, 100, 150, 100, 60, 50)
	if !s.Active {
		t.Fatalf("Novo slash deveria estar ativo")
	}

	// Executa frames de atualização
	s.Update()
	if len(s.Sparks) == 0 {
		t.Errorf("Slash deveria gerar faíscas dinâmicas durante o arco")
	}

	// Alvo na frente do corte
	hit := s.CollidesWithEnemy(1, collision.NewRect(130, 95, 20, 20))
	if !hit {
		t.Errorf("Inimigo diretamente no alcance angular deveria ser atingido")
	}

	// Inimigo atrás do guerreiro
	behind := s.CollidesWithEnemy(2, collision.NewRect(50, 100, 20, 20))
	if behind {
		t.Errorf("Inimigo atrás do ângulo do golpe não deveria ser atingido")
	}
}
