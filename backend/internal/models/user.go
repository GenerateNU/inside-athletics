package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.model 
	Name string `gorm:"type:varchar(100);not null"`
}