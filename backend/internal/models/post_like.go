package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// A PostLike represents a like on a post
type PostLike struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	UserID uuid.UUID `json:"user_id" gorm:"foreignKey;type:uuid"`
	PostID uuid.UUID `json:"post_id" gorm:"foreignKey;type:uuid"`
}
