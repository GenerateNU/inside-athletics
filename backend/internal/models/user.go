package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `json:"name" example:"Suli" doc:"The name of a user" gorm:"type:varchar(100);not null"`
}
