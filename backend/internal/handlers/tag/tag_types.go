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

type GetPostsByTagParams struct {
	TagID uuid.UUID `path:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Tag"`
}

type GetPostsByTagResponse struct {
	TagID   uuid.UUID   `json:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Tag"`
	PostIDs []uuid.UUID `json:"post_ids" example:"[\"123e4567-e89b-12d3-a456-426614174000\",\"123e4567-e89b-12d3-a456-426614174001\"]" doc:"The post ids associated with a tag"`
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
