package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID                      uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               *time.Time     `sql:"index" json:"deleted_at"`
	FirstName               string         `json:"first_name" example:"Suli" doc:"The first name of a user" gorm:"type:varchar(100);not null"`
	LastName                string         `json:"last_name" example:"Suli" doc:"The last name of a user" gorm:"type:varchar(100);not null"`
	Email                   string         `json:"email" example:"suli123@email.com" doc:"The email of a user" gorm:"type:varchar(100);not null"`
	Username                string         `json:"username" example:"suliproathelete" doc:"The username of a user" gorm:"type:varchar(100);not null"`
	Bio                     *string        `json:"bio" example:"My name is Suli and I'm a pro athlete" doc:"The name of a user" gorm:"type:varchar(100);"` //nullable
	Account_Type            bool           `json:"account_type" example:"True" doc:"If the user has access to premium features" gorm:"type:bool;not null"`
	Sport                   datatypes.JSON `json:"sport" example:"hockey" doc:"The sport(s) the user is interested in" gorm:"type:jsonb;"` //nullable
	Expected_Grad_Year      uint           `json:"expected_grad_year" example:"2027" doc:"The user's grad year" gorm:"type:uint;"`         //nullable
	Verified_Athlete_Status VerifiedAthleteStatus `json:"verified_athelete_status" example:"pending" doc:"" gorm:"type:varchar(100);not null"`
	College                 *string        `json:"college" example:"Northeastern University" doc:"The college of a user" gorm:"type:varchar(100);"`
	Division                *Division      `json:"division" example:"1" doc:"The divison of their college" gorm:"type:uint;"`
}

type VerifiedAthleteStatus string

const (
	VerifiedAthleteStatusNone     VerifiedAthleteStatus = "none"
	VerifiedAthleteStatusPending  VerifiedAthleteStatus = "pending"
	VerifiedAthleteStatusVerified VerifiedAthleteStatus = "verified"
)

type Division uint

const (
	DivisionI   Division = 1
	DivisionII Division = 2
	DivisionIII Division = 3
)

// Validate the user
func (u *User) BeforeSave(tx *gorm.DB) error {
	// If verified athlete, both college and division are required
	if u.Verified_Athlete_Status == VerifiedAthleteStatusVerified {
		if u.College == nil || *u.College == "" {
			return errors.New("college is required when verified athlete status is 'verified'")
		}
		if u.Division == nil {
			return errors.New("division is required when verified athlete status is 'verified'")
		}
	}
	return nil
}
