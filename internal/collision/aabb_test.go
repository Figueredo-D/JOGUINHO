package collision

import (
	"testing"
)

func TestRectIntersects(t *testing.T) {
	r1 := NewRect(10, 10, 20, 20)
	r2 := NewRect(15, 15, 20, 20) // Sobreposto
	r3 := NewRect(40, 40, 20, 20) // Longe

	if !r1.Intersects(r2) {
		t.Errorf("r1 e r2 deveriam colidir")
	}

	if r1.Intersects(r3) {
		t.Errorf("r1 e r3 não deveriam colidir")
	}
}

func TestRectContains(t *testing.T) {
	r := NewRect(10, 10, 50, 50)

	if !r.Contains(20, 20) {
		t.Errorf("Ponto (20,20) deveria estar dentro do retângulo")
	}

	if r.Contains(5, 5) {
		t.Errorf("Ponto (5,5) não deveria estar dentro do retângulo")
	}
}
