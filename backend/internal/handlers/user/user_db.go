package user

import (
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserDB struct {
	db *gorm.DB
}

/*
*
Here we are using GORM to interact with the database. This is an ORM (Object Relational Mapping)
which allows us to interact with the database without having to write raw SQL queries
*/
func (u *UserDB) GetUser(id uuid.UUID) (*models.User, error) {
	var user models.User
	dbResponse := u.db.Where("id = ?", id).First(&user)
	return utils.HandleDBError(&user, dbResponse.Error) // helper function that maps GORM errors to Huma errors
}
