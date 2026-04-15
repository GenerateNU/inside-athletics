package tag

import (
	"inside-athletics/internal/handlers/post"

	"github.com/google/uuid"
)

// GETTING ALL POSTS FROM A TAG
type GetPostsByTagParam struct {
	TagID  uuid.UUID `path:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Tag"`
	Limit  int       `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset int       `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}
type GetPostsByTagResponse struct {
	Posts []post.PostResponse `json:"post_ids" doc:"The post ids associated with a tag"`
}

type GetTagByIDParams struct {
	ID uuid.UUID `path:"id" example:"1" doc:"ID to identify tag"`
}

type GetTagByNameParams struct {
	Name string `path:"name" example:"Hockey" doc:"Name to identify tag"`
}

type ListTagsParams struct {
	Limit  int `query:"limit" default:"100" example:"100" doc:"Number of tags to return"`
	Offset int `query:"offset" default:"0" example:"0" doc:"Number of tags to skip"`
}

type GetPostsByTagParams struct {
	TagID uuid.UUID `path:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Tag"`
}

type GetTagResponse struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the tag"`
	Name string    `json:"name" example:"Hockey" doc:"The name of the tag"`
}

type ListTagsResponse struct {
	Tags []GetTagResponse `json:"tags" doc:"List of tags"`
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
