package user

import (
	models "inside-athletics/internal/models"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type GetUserParams struct {
	ID uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID to identify the user"`
}

type GetUserResponse struct {
	ID                    uuid.UUID                    `json:"id" example:"1" doc:"ID of the user"`
	FirstName             string                       `json:"first_name" example:"Suli" doc:"The first name of a user"`
	LastName              string                       `json:"last_name" example:"Suli" doc:"The last name of a user"`
	Email                 string                       `json:"email" example:"suli123@email.com" doc:"The email of a user"`
	Username              string                       `json:"username" example:"suliproathlete" doc:"The username of a user"`
	Bio                   *string                      `json:"bio,omitempty" example:"My name is Suli and I'm a pro athlete" doc:"The bio of a user"`
	AccountType           bool                         `json:"account_type" example:"true" doc:"If the user has access to premium features"`
	Sport                 datatypes.JSON               `json:"sport,omitempty" example:"[\"hockey\",\"soccer\"]" doc:"The sport(s) the user is interested in"`
	ExpectedGradYear      uint                         `json:"expected_grad_year,omitempty" example:"2027" doc:"The user's grad year"`
	VerifiedAthleteStatus models.VerifiedAthleteStatus `json:"verified_athlete_status" example:"pending" doc:"Verification status for the athlete"`
	College               *string                      `json:"college,omitempty" example:"Northeastern University" doc:"The college of a user"`
	Division              *models.Division             `json:"division,omitempty" example:"1" doc:"The division of their college"`
	Roles                 []UserRoleResponse           `json:"roles,omitempty" doc:"Roles assigned to the user"`
}

type GetCurrentUserIDResponse struct {
	ID                    uuid.UUID                    `json:"id" example:"1" doc:"ID of the currently authenticated user"`
	FirstName             string                       `json:"first_name" example:"Suli" doc:"The first name of a user"`
	LastName              string                       `json:"last_name" example:"Suli" doc:"The last name of a user"`
	Email                 string                       `json:"email" example:"suli123@email.com" doc:"The email of a user"`
	Username              string                       `json:"username" example:"suliproathlete" doc:"The username of a user"`
	Bio                   *string                      `json:"bio,omitempty" example:"My name is Suli and I'm a pro athlete" doc:"The bio of a user"`
	AccountType           bool                         `json:"account_type" example:"true" doc:"If the user has access to premium features"`
	Sport                 datatypes.JSON               `json:"sport,omitempty" example:"[\"hockey\",\"soccer\"]" doc:"The sport(s) the user is interested in"`
	ExpectedGradYear      uint                         `json:"expected_grad_year,omitempty" example:"2027" doc:"The user's grad year"`
	VerifiedAthleteStatus models.VerifiedAthleteStatus `json:"verified_athlete_status" example:"pending" doc:"Verification status for the athlete"`
	College               *string                      `json:"college,omitempty" example:"Northeastern University" doc:"The college of a user"`
	Division              *models.Division             `json:"division,omitempty" example:"1" doc:"The division of their college"`
	Roles                 []UserRoleResponse           `json:"roles,omitempty" doc:"Roles assigned to the user"`
}

type UserRoleResponse struct {
	ID   uuid.UUID       `json:"id" example:"1" doc:"ID of the role"`
	Name models.RoleName `json:"name" example:"user" doc:"Name of the role"`
}

type CreateUserInput struct {
	Body CreateUserBody
}

type CreateUserBody struct {
	FirstName             string                       `json:"first_name" example:"Suli" doc:"The first name of a user"`
	LastName              string                       `json:"last_name" example:"Suli" doc:"The last name of a user"`
	Email                 string                       `json:"email" example:"suli123@email.com" doc:"The email of a user"`
	Username              string                       `json:"username" example:"suliproathlete" doc:"The username of a user"`
	Bio                   *string                      `json:"bio,omitempty" example:"My name is Suli and I'm a pro athlete" doc:"The bio of a user"`
	AccountType           bool                         `json:"account_type" example:"true" doc:"If the user has access to premium features"`
	Sport                 []string                     `json:"sport,omitempty" example:"[\"hockey\",\"soccer\"]" doc:"The sport(s) the user is interested in"`
	ExpectedGradYear      uint                         `json:"expected_grad_year,omitempty" example:"2027" doc:"The user's grad year"`
	VerifiedAthleteStatus models.VerifiedAthleteStatus `json:"verified_athlete_status" example:"pending" doc:"Verification status for the athlete"`
	College               *string                      `json:"college,omitempty" example:"Northeastern University" doc:"The college of a user"`
	Division              *models.Division             `json:"division,omitempty" example:"1" doc:"The division of their college"`
}

type CreateUserResponse struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the user"`
	Name string    `json:"name" example:"Joe" doc:"Name of the user"`
}

type UpdateUserInput struct {
	ID   uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID to identify the user"`
	Body UpdateUserBody
}

type UpdateUserBody struct {
	FirstName             *string                       `json:"first_name,omitempty" example:"Suli" doc:"The first name of a user"`
	LastName              *string                       `json:"last_name,omitempty" example:"Suli" doc:"The last name of a user"`
	Email                 *string                       `json:"email,omitempty" example:"suli123@email.com" doc:"The email of a user"`
	Username              *string                       `json:"username,omitempty" example:"suliproathlete" doc:"The username of a user"`
	Bio                   *string                       `json:"bio,omitempty" example:"My name is Suli and I'm a pro athlete" doc:"The bio of a user"`
	AccountType           *bool                         `json:"account_type,omitempty" example:"true" doc:"If the user has access to premium features"`
	Sport                 *[]string                     `json:"sport,omitempty" example:"[\"hockey\",\"soccer\"]" doc:"The sport(s) the user is interested in" gorm:"type:jsonb;serializer:json"`
	ExpectedGradYear      *uint                         `json:"expected_grad_year,omitempty" example:"2027" doc:"The user's grad year"`
	VerifiedAthleteStatus *models.VerifiedAthleteStatus `json:"verified_athlete_status,omitempty" example:"pending" doc:"Verification status for the athlete"`
	College               *string                       `json:"college,omitempty" example:"Northeastern University" doc:"The college of a user"`
	Division              *models.Division              `json:"division,omitempty" example:"1" doc:"The division of their college"`
}

type UpdateUserResponse struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the user"`
	Name string    `json:"name" example:"Joe" doc:"Name of the user"`
}

type DeleteUserResponse struct {
	ID uuid.UUID `json:"id" example:"1" doc:"ID of the deleted user"`
}
