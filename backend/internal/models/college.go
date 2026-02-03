package models

import (
	"time"

	"github.com/google/uuid"
)

type College struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`

	Name         string `json:"name" example:"Northeastern University" doc:"The name of a college" gorm:"type:varchar(200);not null"`
	State        string `json:"state" example:"Massachusetts" doc:"The state of the college" gorm:"type:varchar(100);not null"`
	City         string `json:"city" example:"Boston" doc:"The city of the college" gorm:"type:varchar(100);not null"`
	Website      string `json:"website" example:"https://www.northeastern.edu" doc:"The website of the college" gorm:"type:varchar(500);not null"`
	AcademicRank *int16 `json:"academic_rank" example:"53" doc:"The academic rank of the college" gorm:"type:smallint"`
	DivisionRank int8   `json:"division_rank" example:"1" doc:"NCAA division (1, 2, or 3)" gorm:"type:int8;not null"`
	// Logo stores the URL/path to the college logo. Eventually this will point to a picture in S3.
	Logo string `json:"logo" example:"https://example.com/logo.png" doc:"The logo of the college" gorm:"type:varchar(500)"`
}
