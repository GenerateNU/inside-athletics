package user

import "github.com/google/uuid"

type GetUserParams struct {
	ID uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID to identify the user"`
}

type GetUserResponse struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the user"`
	Name string    `json:"name" example:"Joe" doc:"Name of the user"`
}

type GetCurrentUserIDResponse struct {
	ID uuid.UUID `json:"id" example:"1" doc:"ID of the currently authenticated user"`
}
