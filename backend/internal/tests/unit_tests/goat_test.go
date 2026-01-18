package unitTests

import (
	"fmt"
	"inside-athletics/internal/models"
	"testing"
)

func TestGoat(t *testing.T) {
	goat := models.Goat{Name: "Joe", Age: 67}
	actual := goat.MakeSomeNoise()
	expected := "BAAAAA Joe"

	if actual != expected {
		t.Error("The goat did not make some noise...")
	}
}

func TestPointers(t *testing.T) {
	goat := models.Goat{Name: "Joe", Age: 67}
	goat.SetName("Suli")
	if goat.Name != "Suli" {
		t.Error("Set name didn't work...")
	}
	goat.SetNameCopy("Erm")
	if goat.Name != "Suli" {
		t.Error("I was wrong and it didn't use a copy")
	}

	copy := goat.MakeCopy()
	if &copy == &goat {
		t.Error("THEY ARE THE SAME OH MAN")
	}
	fmt.Printf("Copy pointer: %p, OG pointer: %p", &copy, &goat)
}
