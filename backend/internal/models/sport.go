package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Sport represents a sport entity in the system
type Sport struct {
	ID         uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Name       string         `json:"name" example:"Women's Soccer" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Popularity *int32         `json:"popularity,omitempty" example:"20000" gorm:"type:int"`
}
