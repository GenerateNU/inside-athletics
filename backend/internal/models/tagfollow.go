package models

import (
	"time"

	"github.com/google/uuid"
)

type TagFollow struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`

	TagID  uuid.UUID `json:"tag_id" doc:"TagID of tag the user followed" gorm:"foreignKey:TagID;references:ID;constraint:OnDelete:CASCADE"`
	UserID uuid.UUID `json:"user_id" doc:"UserID of user that followed tag" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}
