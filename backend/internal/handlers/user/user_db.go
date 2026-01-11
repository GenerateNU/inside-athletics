package user

import (
	"context"
	"fmt"
	models "inside-athletics/internal/models"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"gorm.io/gorm"
)

type UserDB struct {
	db *gorm.DB
}

/**
Here we are using GORM to interact with the database. This is an ORM (Object Relational Mapping)
which allows us to interact with the database without having to write raw SQL queries
*/
func (u *UserDB) GetUser(name string) (*models.User, error) {
	var user models.User
	dbResponse := u.db.Where("name = ?", name).First(&user)
	return handleDBError(&user, dbReponse.error) // helper function that maps GORM errors to Huma errors
}

