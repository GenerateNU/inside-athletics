package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// A CommentLike represents a like on a comment in the system
type CommentLike struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	UserID uuid.UUID `json:"user_id" gorm:"foreignKey;type:uuid"`
	CommentID uuid.UUID `json:"comment_id" gorm:"foreignKey;type:uuid"`
}
