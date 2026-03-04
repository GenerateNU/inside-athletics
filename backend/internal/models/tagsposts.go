package models

import (
	"time"

	"github.com/google/uuid"
)

type TagPost struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PostID    uuid.UUID  `json:"post_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Post 	  Post		 `json:"-" gorm:"foreignKey:PostId;references:ID;constraint:OnDelete:CASCADE"`
	TagID     uuid.UUID  `json:"tag_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Tag 	  Tag        `json:"-" gorm:"foreignKey:TagId;references:ID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
