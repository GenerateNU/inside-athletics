package utility

import (
	"inside-athletics/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UtilityDB struct {
	db *gorm.DB
}

func NewUtilityDB(db *gorm.DB) *UtilityDB {
	return &UtilityDB{db: db}
}

func (u *UtilityDB) UserHasPremium(userID uuid.UUID) (bool, error) {
	var count int64
	err := u.db.Model(&models.User{}).
		Where("id = ? AND account_type = true", userID).
		Count(&count).Error
	return count > 0, err
}

func (u *UtilityDB) UserIsAdmin(userID uuid.UUID) (bool, error) {
	var count int64
	err := u.db.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.name = ?", userID, models.RoleAdmin).
		Count(&count).Error
	return count > 0, err
}
