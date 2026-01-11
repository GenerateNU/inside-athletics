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

func (u *UserDB) GetUser(name string) (*models.User, error) {
	var user models.User
	dbResponse := u.db.Where("name = ?", name).First(&user)
	return handleDBError(&user, dbReponse.error)
}

