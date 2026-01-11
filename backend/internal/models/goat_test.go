package models

import "testing"

func TestGoat(t *testing.T) {
	goat := Goat{Id: 1, Name: "Joe", Age: 67}
	actual := goat.makeSomeNoise()
	expected := "BAAAAA"

	if actual != expected {
		t.Error("The goat did not make some noise...")
	}
}
