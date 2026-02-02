package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Like represents a like entity in the system
type Like struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	PostID uuid.UUID `json:"post_id" gorm:"foreignKey;type:uuid"`
}
