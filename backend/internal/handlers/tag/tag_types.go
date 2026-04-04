package tag

import (
	"inside-athletics/internal/handlers/post"
	"inside-athletics/internal/models"

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

type GetTagByTypeParams struct {
	Type models.TagType `json:"type" example:"sports" doc:"The type of the tag" gorm:"type:varchar(50);not null"`
}

type GetPostsByTagParams struct {
	TagID uuid.UUID `path:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Tag"`
}

type GetTagResponse struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the tag"`
	Name string    `json:"name" example:"Hockey" doc:"The name of the tag"`
	Type models.TagType `json:"type" example:"sports" doc:"The type of the tag" gorm:"type:varchar(50);not null"`
}

type CreateTagInput struct {
	Body CreateTagBody
}

type CreateTagBody struct {
	Name string `json:"name" example:"Hockey" doc:"The name of the tag to create"`
	Type models.TagType `json:"type" example:"sports" doc:"The type of the tag" gorm:"type:varchar(50);not null"`
}

type CreateTagResponse struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the tag created"`
	Name string    `json:"name" example:"Hockey" doc:"The name of the tag created"`
	Type models.TagType `json:"type" example:"sports" doc:"The type of the tag" gorm:"type:varchar(50);not null"`
}

type UpdateTagInput struct {
	ID   uuid.UUID `path:"id" example:"1" doc:"ID of the tag to update"`
	Body UpdateTagBody
}

type UpdateTagBody struct {
	Name string `json:"name" example:"Hockey" doc:"The new name to update"`
	Type models.TagType `json:"type" example:"sports" doc:"The type of the tag" gorm:"type:varchar(50);not null"`
}

type UpdateTagResponse struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the tag updated"`
	Name string    `json:"name" example:"Northeastern Hockey" doc:"The updated name of the tag"`
	Type models.TagType `json:"type" example:"sports" doc:"The type of the tag" gorm:"type:varchar(50);not null"`
}

type DeleteTagResponse struct {
	ID uuid.UUID `json:"id" example:"1" doc:"ID of the deleted tag"`
}
