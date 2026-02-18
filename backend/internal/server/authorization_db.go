package server

import (
	"inside-athletics/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthorizationDB wraps authorization-related queries for readability and reuse.
type AuthorizationDB struct {
	db *gorm.DB
}

func NewAuthorizationDB(db *gorm.DB) *AuthorizationDB {
	return &AuthorizationDB{db: db}
}

func (a *AuthorizationDB) UserExists(id uuid.UUID) error {
	var user models.User
	return a.db.Select("id").First(&user, "id = ?", id).Error
}

func (a *AuthorizationDB) UserHasPermission(userID uuid.UUID, action models.PermissionAction, resource string) (bool, error) {
	var count int64
	err := a.db.Table("user_roles").
		Joins("JOIN role_permissions rp ON rp.role_id = user_roles.role_id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("user_roles.user_id = ? AND p.action = ? AND p.resource = ?", userID, action, resource).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
