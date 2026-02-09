package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleName string

const (
	RoleUser      RoleName = "user"
	RoleAdmin     RoleName = "admin"
	RoleModerator RoleName = "moderator"
)

type Role struct {
	ID              uuid.UUID        `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `json:"deleted_at,omitempty" gorm:"index"`
	Name            RoleName         `json:"name" gorm:"type:varchar(50);not null;unique"`
	RolePermissions []RolePermission `json:"role_permissions,omitempty" gorm:"foreignKey:RoleID"`
}

type RolePermission struct {
	RoleID       uuid.UUID  `json:"role_id" gorm:"type:uuid;primaryKey"`
	PermissionID uuid.UUID  `json:"permission_id" gorm:"type:uuid;primaryKey"`
	Role         Role       `json:"role" gorm:"constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
	Permission   Permission `json:"permission" gorm:"constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
}
