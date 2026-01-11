package tests

import (
	"testing"
	"inside-athletics/internal/models"
)

func TestGoat(t *testing.T) {
	goat := models.Goat{Id: 1, Name: "Joe", Age: 67}
	actual := goat.MakeSomeNoise()
	expected := "BAAAAA"

	if actual != expected {
		t.Error("The goat did not make some noise...")
	}
}
