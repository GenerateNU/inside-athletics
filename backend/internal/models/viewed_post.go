package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ViewedPost records that a user (free tier) has viewed a post for limit enforcement.
type ViewedPost struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_viewed_posts_user_post"`
	User      User           `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	PostID    uuid.UUID      `json:"post_id" gorm:"type:uuid;not null;uniqueIndex:idx_viewed_posts_user_post"`
	Post      Post           `json:"-" gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
