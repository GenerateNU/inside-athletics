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

	UserID          uuid.UUID  `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"UserID the comment belongs to (always stored; hidden for anonymous when not super user)" gorm:"foreignKey;type:uuid;not null"`
	IsAnonymous     bool       `json:"is_anonymous" doc:"If true, user_id is hidden from regular users; super user always sees user_id" gorm:"default:false;not null"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"CommentID the comment belongs to" gorm:"foreignKey;type:uuid"`
	PostID          uuid.UUID  `json:"post_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"PostID the comment belongs to" gorm:"foreignKey;type:uuid;not null"`
	Description     string     `json:"description" example:"This is a helpful thread" maxLength:"1500" doc:"Content of the comment" gorm:"type:varchar(3000);not null"`
}
