package models

import (
	"time"

	"github.com/google/uuid"
)

type Athlete struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	// useful meta-data
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`

	Name      string     `json:"name" example:"Suli Rashidzada" doc:"The name of an athlete" gorm:"type:varchar(200);not null"`
	SportID   *uuid.UUID `json:"sport_id" gorm:"type:uuid;default:null"`
	Sport     *Sport     `json:"-" gorm:"foreignKey:SportID;references:ID;constraint:OnDelete:SET NULL;"`
	CollegeID *uuid.UUID `json:"college_id" gorm:"type:uuid;default:null"`
	College   *College   `json:"-" gorm:"foreignKey:CollegeID;references:ID;constraint:OnDelete:SET NULL;"`
}
