package models

import "github.com/google/uuid"

type UserRole struct {
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`
	RoleID uuid.UUID `json:"role_id" gorm:"type:uuid;primaryKey"`
	User   User      `json:"user" gorm:"constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
	Role   Role      `json:"role" gorm:"constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
}
