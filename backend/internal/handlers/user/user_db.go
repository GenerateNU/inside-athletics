package user

import (
	"context"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
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

func (u *UserDB) GetCurrentUserID(ctx context.Context) (uuid.UUID, error) {
	rawID := ctx.Value("user_id")
	if rawID == nil {
		return uuid.Nil, huma.Error401Unauthorized("User not authenticated")
	}

	userID, ok := rawID.(string)
	if !ok {
		return uuid.Nil, huma.Error500InternalServerError("Invalid user ID in context")
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, huma.Error400BadRequest("Invalid user ID", err)
	}

	return parsedID, nil
}

func (u *UserDB) CreateUser(user *models.User) (*models.User, error) {
	dbResponse := u.db.Create(user)
	return utils.HandleDBError(user, dbResponse.Error)
}

func (u *UserDB) UpdateUser(user *models.User) (*models.User, error) {
	dbResponse := u.db.Save(user)
	return utils.HandleDBError(user, dbResponse.Error)
}

func (u *UserDB) DeleteUser(id uuid.UUID) error {
	dbResponse := u.db.Delete(&models.User{}, "id = ?", id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.User{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
