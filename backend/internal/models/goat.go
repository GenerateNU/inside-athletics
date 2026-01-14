package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Goat struct {
	gorm.Model
	Name string `json:"name" example:"Suli" doc:"The name of a goat"`
	Age  int8   `json:"age" example:"67" doc:"The age of this goat"`
}

func (g *Goat) validateGoatStatus() (string, error) {
	/*
		TODO: Validate if this goat is really a GOAT.
		Heres the sitch:
		1) The goats id MUST be a positive integer (> 0)
		2) If the goats name doesn't start with an "s" are they really goated?
		3) They can't have a negative age... It also needs to be even with the exception of 67

		If the goat is really goated then you return "AYYY GOATED", nil (nil meaning there is no error)
		Else... return the empty string "", [INSERT_ERROR_HERE]
	*/
	return "", nil
}

func (g *Goat) MakeSomeNoise() string {
	return fmt.Sprintf("BAAAAA %s", g.Name)
}

func (g *Goat) SetName(name string) {
	g.Name = name
}

func (g Goat) SetNameCopy(name string) {
	g.Name = name
	fmt.Println(g.Name)
}

func (g Goat) MakeCopy() Goat {
	return g
}
