package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionAction string

const (
	PermissionCreate    PermissionAction = "create"
	PermissionUpdate    PermissionAction = "update"
	PermissionDelete    PermissionAction = "delete"
	PermissionUpdateOwn PermissionAction = "update_own"
	PermissionDeleteOwn PermissionAction = "delete_own"
)

type Permission struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Action    PermissionAction `json:"action" gorm:"type:varchar(50);not null"`
	Resource  string           `json:"resource" gorm:"type:varchar(50);not null"`
	Roles     []Role           `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}
