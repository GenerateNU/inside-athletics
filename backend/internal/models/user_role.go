package models

import "github.com/google/uuid"

type UserRole struct {
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`
	RoleID uuid.UUID `json:"role_id" gorm:"type:uuid;primaryKey"`
}
