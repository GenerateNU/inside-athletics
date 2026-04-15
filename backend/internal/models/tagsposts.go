package models

import (
	"time"

	"github.com/google/uuid"
)

type TagPost struct {
	ID           uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TagID        uuid.UUID  `json:"tag_id" gorm:"type:uuid;not null"`
	Tag          Tag        `json:"-" gorm:"foreignKey:TagID;references:ID;constraint:OnDelete:CASCADE"`
	PostableID   uuid.UUID  `json:"postable_id" gorm:"type:uuid;not null"`
	PostableType string     `json:"postable_type" gorm:"type:varchar(20);not null" validate:"required,oneof=post premium_post"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `sql:"index" json:"deleted_at"`
}
