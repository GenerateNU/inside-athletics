package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Survey represents an athlete's survey response for a sport program at a college
type Survey struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Foreign keys
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index" validate:"required"`
	CollegeID uuid.UUID `json:"college_id" gorm:"type:uuid;not null;index" validate:"required"`
	SportID   uuid.UUID `json:"sport_id" gorm:"type:uuid;not null;index" validate:"required"`

	// Associations
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	College College `json:"college,omitempty" gorm:"foreignKey:CollegeID"`
	Sport   Sport   `json:"sport,omitempty" gorm:"foreignKey:SportID"`

	// Ratings (1–5)
	PlayerDev                  int32 `json:"player_dev" gorm:"type:smallint;not null" validate:"required,min=1,max=5"`
	AcademicsAthleticsPriority int32 `json:"academics_athletics_priority" gorm:"type:smallint;not null" validate:"required,min=1,max=5"`
	AcademicCareerResources    int32 `json:"academic_career_resources" gorm:"type:smallint;not null" validate:"required,min=1,max=5"`
	MentalHealthPriority       int32 `json:"mental_health_priority" gorm:"type:smallint;not null" validate:"required,min=1,max=5"`
	Environment                int32 `json:"environment" gorm:"type:smallint;not null" validate:"required,min=1,max=5"`
	Culture                    int32 `json:"culture" gorm:"type:smallint;not null" validate:"required,min=1,max=5"`
	Transparency               int32 `json:"transparency" gorm:"type:smallint;not null" validate:"required,min=1,max=5"`
}