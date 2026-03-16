package models

import (
	"time"

	"github.com/google/uuid"
)

// A PostLike represents a like on a post
type PostLike struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UserID 	  uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_post_likes_user_id_post_id"`
	User	  User           `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	PostID    uuid.UUID      `json:"post_id" gorm:"type:uuid;not null;uniqueIndex:idx_post_likes_user_id_post_id"`
	Post      Post           `json:"-" gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE"`
}
