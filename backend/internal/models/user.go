package models

import (
	"gorm.io/gorm"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	gorm.model 
	Name string `json:"name" maxLength:"30" example:"Joe" doc:"test name for example data"`
}


func (user *User) Validate() error {
	return validation.ValidateStruct(&user,
		// Name cannot be empty and can be between 5 and 50 characters
		validation.Field(&a.Name, validation.Required, validation.Length(5, 50))
	)

}