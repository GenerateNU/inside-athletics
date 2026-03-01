package models

import (
	"time"

	"github.com/google/uuid"
)

// A CommentLike represents a like on a comment in the system
type CommentLike struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_comment_likes_user_id_comment_id"`
	User      User           `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	CommentID uuid.UUID      `json:"comment_id" gorm:"type:uuid;not null;uniqueIndex:idx_comment_likes_user_id_comment_id"`
	Comment   Comment        `json:"-" gorm:"foreignKey:CommentID;references:ID;constraint:OnDelete:CASCADE"`
}
