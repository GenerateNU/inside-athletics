package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`

	// Comments belong to a user
	UserID      uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"UserID the comment belongs to (always stored; hidden for anonymous when not super user)" gorm:"type:uuid;not null"`
	User        User      `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	IsAnonymous bool      `json:"is_anonymous" doc:"If true, user_id is hidden from regular users; super user always sees user_id" gorm:"default:false;not null"`

	// Optional parent comment (for replies)
	ParentCommentID *uuid.UUID `json:"parent_comment_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"CommentID the comment belongs to" gorm:"type:uuid"`
	ParentComment   *Comment   `json:"-" gorm:"foreignKey:ParentCommentID;references:ID;constraint:OnDelete:SET NULL"`

	// Comments will always belong to a post
	PostID   uuid.UUID `json:"post_id" doc:"PostID the comment belongs to" gorm:"type:uuid;not null"`
	Post     Post      `json:"-" gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE"`

	Description string `json:"description" example:"This is a helpful thread" maxLength:"1500" doc:"Content of the comment" gorm:"type:varchar(3000);not null"`
}
