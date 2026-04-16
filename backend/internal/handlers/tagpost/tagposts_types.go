package tagpost

import (
	"github.com/google/uuid"
	models "inside-athletics/internal/models"
)

// REQUEST/RESPONSE FOR ADDS TAGS TO POST
type PostTagsForPostRequest struct {
	PostId uuid.UUID
	TagIds []uuid.UUIDs
}

type PostTagsForPostResponseBody struct {
	Tags	[]models.Tag
}

type PostTagsForPostResponse struct {
	Body PostTagsForPostResponseBody
}


// PARAM/RESPONSE GETTING ALL TAGS FOR A POST
type GetTagsByPostParam struct {
	PostID uuid.UUID `path:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post"`
}

type GetTagsByPostResponse struct {
	PostID uuid.UUID   `json:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post"`
	Tag   []models.Tag `json:"tag_ids" example:"[\"123e4567-e89b-12d3-a456-426614174000\",\"123e4567-e89b-12d3-a456-426614174001\"]" doc:"The tag ids associated with a post"`
}



// GETTING SINGLE TAG POST
type GetTagPostByIdParam struct {
	ID uuid.UUID `path:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the TagsPosts item"`
}

type GetTagPostByIDResponse struct {
	ID     uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"the id of the item"`
	TagID  uuid.UUID `json:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"the tag id"`
	PostID uuid.UUID `json:"post_ids" example:"123e4567-e89b-12d3-a456-426614174000" doc:"The post id"`
}

// CREATING SINGLE TAG POST
type CreateTagPostInput struct {
	Body CreateTagPostBody
}

type CreateTagPostBody struct {
	PostID uuid.UUID `json:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post"`
	TagID  uuid.UUID `json:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the tag"`
}

type CreateTagPostsResponse struct {
	ID     uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the tagpost item created"`
	PostID uuid.UUID `json:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post"`
	TagID  uuid.UUID `json:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the tag"`
}


// UPDATING TAG/POST
type UpdateTagPostInput struct {
	ID   uuid.UUID `path:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the tagpost to update"`
	Body UpdateTagPostBody
}

type UpdateTagPostBody struct {
	PostID uuid.UUID `json:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"the post id to update to"`
	TagID  uuid.UUID `json:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"the tag id to update to"`
}

type UpdateTagPostResponse struct {
	ID     uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the tagpost updated"`
	PostID uuid.UUID `json:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"the updated post id"`
	TagID  uuid.UUID `json:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"the updated tag id"`
}


// DELETING TAG/POST
type DeleteTagPostResponse struct {
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the deleted tagpost"`
}
