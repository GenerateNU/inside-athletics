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

// This function creates a link between a user and a role in the user_roles table
// We use FirstOrCreate to avoid duplicate entries if the user already has the role
func (u *UserDB) AddUserRole(userID, roleID uuid.UUID) error {
	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	if err := u.db.Where("user_id = ? AND role_id = ?", userID, roleID).Find(&userRole).Error; err != nil {
		return huma.Error500InternalServerError("Failed to assign role to user", err)
	}
	return nil
}

func (u *UserDB) GetAllRolesForUser(userID uuid.UUID) (*[]models.Role, error) {
	var userRoles []models.Role
	err := u.db.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&userRoles).Error
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get user roles", err)
	}
	return &userRoles, nil
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
