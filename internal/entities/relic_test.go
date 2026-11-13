package entities

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"joguinho/internal/sprites"
)

func TestRelicSprites(t *testing.T) {
	for _, rel := range AllRelics {
		sprite := rel.GetSprite()
		if sprite == nil {
			t.Fatalf("Relíquia %s retornou sprite nulo", rel.Name)
		}
		if sprite.Bounds().Dx() == 0 || sprite.Bounds().Dy() == 0 {
			t.Fatalf("Relíquia %s retornou sprite com dimensões zeradas", rel.Name)
		}
	}
}

func TestPotions(t *testing.T) {
	hpPot := sprites.GetHealthPotionSprite()
	if hpPot == nil || hpPot.Bounds().Dx() == 0 {
		t.Errorf("GetHealthPotionSprite retornou sprite inválido")
	}

	manaPot := sprites.GetManaPotionSprite()
	if manaPot == nil || manaPot.Bounds().Dx() == 0 {
		t.Errorf("GetManaPotionSprite retornou sprite inválido")
	}
}
