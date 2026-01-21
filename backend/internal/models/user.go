package models

import (
	"time"

	"github.com/google/uuid"
)


type User struct {
	ID         uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `sql:"index" json:"deleted_at"`
	Name string `json:"name" example:"Suli" doc:"The name of a user" gorm:"type:varchar(100);not null"`
}
