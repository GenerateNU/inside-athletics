package models

import (
	"errors"
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

var (
	ErrPermissionSpecInvalid    = errors.New("permission spec must include action and resource")
	ErrPermissionActionInvalid  = errors.New("permission action is invalid")
	ErrPermissionResourceInvalid = errors.New("permission resource is invalid")
)

func IsValidPermissionAction(action PermissionAction) bool {
	switch action {
	case PermissionCreate,
		PermissionUpdate,
		PermissionDelete,
		PermissionUpdateOwn,
		PermissionDeleteOwn:
		return true
	default:
		return false
	}
}

func ValidatePermissionSpec(action PermissionAction, resource string) error {
	if action == "" || resource == "" {
		return ErrPermissionSpecInvalid
	}
	if !IsValidPermissionAction(action) {
		return ErrPermissionActionInvalid
	}
	return nil
}

func ValidatePermissionAction(action PermissionAction) error {
	if action == "" || !IsValidPermissionAction(action) {
		return ErrPermissionActionInvalid
	}
	return nil
}

func ValidatePermissionResource(resource string) error {
	if resource == "" {
		return ErrPermissionResourceInvalid
	}
	return nil
}

type Permission struct {
	ID              uuid.UUID        `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `json:"deleted_at,omitempty" gorm:"index"`
	Action          PermissionAction `json:"action" gorm:"type:varchar(50);not null"`
	Resource        string           `json:"resource" gorm:"type:varchar(50);not null"`
}
