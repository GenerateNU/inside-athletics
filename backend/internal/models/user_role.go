package models

import "github.com/google/uuid"

type UserRole struct {
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`
    RoleID uuid.UUID `json:"role_id" gorm:"type:uuid;primaryKey"`
    User   User      `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
    Role   Role      `json:"-" gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE"`
}
