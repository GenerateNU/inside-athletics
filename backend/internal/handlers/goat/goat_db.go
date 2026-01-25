package goat

import (
	"gorm.io/gorm"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type GoatDB struct {
	db *gorm.DB // our connection to Supabase
}

func (g *GoatDB) GetGoat(id uint) (*models.Goat, error) {
	var goat models.Goat
	dbResponse := g.db.Where("id = ?", id).First(&goat) // gets the entry with the given id and copies the data into goat
	return utils.HandleDBError(&goat, dbResponse.Error) // Error handling utility function!!!
}