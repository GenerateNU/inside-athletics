package models

import (
	"time"

	"github.com/google/uuid"
)

type CollegeFollow struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`

	CollegeID uuid.UUID `json:"college_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"CollegeID of the college follow" gorm:"type:uuid;not null"`
	College   College   `json:"-" gorm:"foreignKey:TagID;references:ID;constraint:OnDelete:CASCADE"`

	UserID uuid.UUID `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"UserID of the college follow" gorm:"type:uuid;not null"`
	User   User
}
