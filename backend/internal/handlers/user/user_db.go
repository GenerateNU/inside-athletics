package user

import (
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (u *UserDB) CreateUser(user *models.User) (*models.User, error) {
	dbResponse := u.db.Create(user)
	return utils.HandleDBError(user, dbResponse.Error)
}

func (u *UserDB) GetRoleIDByName(name models.RoleName) (uuid.UUID, error) {
	var role models.Role
	if err := u.db.Select("id").Where("name = ?", name).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return uuid.Nil, huma.Error500InternalServerError("Default role not found")
		}
		return uuid.Nil, huma.Error500InternalServerError("Database error", err)
	}
	return role.ID, nil
}

func (u *UserDB) UpdateUser(id uuid.UUID, updates UpdateUserBody) (*models.User, error) {
	var updatedUser models.User
	dbResponse := u.db.Model(&models.User{}).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updates).
		Scan(&updatedUser)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.User{}, dbResponse.Error)
		return nil, err
	}
	if dbResponse.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Resource not found")
	}
	return &updatedUser, nil
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
