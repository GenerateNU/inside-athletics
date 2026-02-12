package tag

import (
	"github.com/google/uuid"
)

type GetTagByIDParams struct {
	ID uuid.UUID `path:"id" example:"1" doc:"ID to identify tag"`
}

type GetTagByNameParams struct {
	Name string `path:"name" example:"Hockey" doc:"Name to identify tag"`
}

type GetTagResponse struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the tag"`
	Name string    `json:"name" example:"Hockey" doc:"The name of the tag"`
}

type CreateTagInput struct {
	Body CreateTagBody
}

type CreateTagBody struct {
	Name string `json:"name" example:"Hockey" doc:"The name of the tag to create"`
}

type CreateTagResponse struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the tag created"`
	Name string    `json:"name" example:"Hockey" doc:"The name of the tag created"`
}

type UpdateTagInput struct {
	ID   uuid.UUID `path:"id" example:"1" doc:"ID of the tag to update"`
	Body UpdateTagBody
}

type UpdateTagBody struct {
	Name string `json:"name" example:"Hockey" doc:"The new name to update"`
}

type UpdateTagResponse struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the tag updated"`
	Name string    `json:"name" example:"Northeastern Hockey" doc:"The updated name of the tag"`
}

type DeleteTagResponse struct {
	ID uuid.UUID `json:"id" example:"1" doc:"ID of the deleted tag"`
}
