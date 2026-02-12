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

	UserID uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"UserID the like belongs to" gorm:"type:uuid;not null"`
	User   User      `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`

	CommentID uuid.UUID `json:"comment_id" doc:"CommentID the like belongs to" gorm:"type:uuid;not null"`
	Comment   Comment   `json:"-" gorm:"foreignKey:CommentID;references:ID;constraint:OnDelete:CASCADE"`
}
