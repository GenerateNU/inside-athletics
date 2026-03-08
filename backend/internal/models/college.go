package models

import (
	"time"

	"github.com/google/uuid"
)

type College struct {
	ID uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	// useful meta-data
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`

	Name         string   `json:"name" example:"Northeastern University" doc:"The name of a college" gorm:"type:varchar(200);not null"`
	State        string   `json:"state" example:"Massachusetts" doc:"The state of the college" gorm:"type:varchar(100);not null"`
	City         string   `json:"city" example:"Boston" doc:"The city of the college" gorm:"type:varchar(100);not null"`
	Website      string   `json:"website" example:"https://www.northeastern.edu" doc:"The website of the college" gorm:"type:varchar(500);not null"`
	AcademicRank *int16   `json:"academic_rank" example:"53" doc:"The academic rank of the college" gorm:"type:smallint"`
	DivisionRank Division `json:"division_rank" enum:"1,2,3" example:"1" doc:"NCAA division (1, 2, or 3)" gorm:"type:uint;not null"`

	// Stores the S3 Key of the image.
	Logo string `json:"logo" doc:"The S3 key for the logo of the college" gorm:"type:varchar(500)"`
}
