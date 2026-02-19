package models

import (
	"time"

	"github.com/google/uuid"
)

// A PostLike represents a like on a post
type PostLike struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`

	UserID uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"UserID the like belongs to" gorm:"type:uuid;not null"`
	User User `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`

	PostID   uuid.UUID `json:"post_id" doc:"PostID the like belongs to" gorm:"type:uuid;not null"`
	Post     Post      `json:"-" gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE"`
}
