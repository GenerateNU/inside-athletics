package models

import (
	"time"

	"github.com/google/uuid"
)

type TagPost struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PostID    uuid.UUID  `json:"post_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TagID     uuid.UUID  `json:"tag_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
