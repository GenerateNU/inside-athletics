package models

import (
	"time"

	"github.com/google/uuid"
)

type SportFollow struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`

	SportID uuid.UUID `json:"sport_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"SportID of the sport follow" gorm:"type:uuid;not null;uniqueIndex:idx_user_sport"`
	Sport   Sport     `json:"-" gorm:"foreignKey:SportID;references:ID;constraint:OnDelete:CASCADE"`

	UserID uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"UserID of the sport follow" gorm:"type:uuid;not null;uniqueIndex:idx_user_sport"`
	User   User      `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}
