package tagpost

import (
	models "inside-athletics/internal/models"

	"github.com/google/uuid"
)

// REQUEST/RESPONSE FOR ADDING TAGS TO POST
type PostTagsForPostRequest struct {
	PostableID   uuid.UUID
	PostableType string
	TagIds       []uuid.UUID
}
type PostTagsForPostResponseBody struct {
	Tags []models.Tag
}
type PostTagsForPostResponse struct {
	Body PostTagsForPostResponseBody
}

// PARAM/RESPONSE GETTING ALL TAGS FOR A POST
type GetTagsByPostParam struct {
	PostableID   uuid.UUID `path:"postable_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post or premium post"`
	PostableType string    `path:"postable_type" example:"post" doc:"Type of the postable (post or premium_post)"`
}
type GetTagsByPostResponse struct {
	PostableID   uuid.UUID    `json:"postable_id"`
	PostableType string       `json:"postable_type"`
	Tags         []models.Tag `json:"tags"`
}

// GETTING SINGLE TAG POST
type GetTagPostByIdParam struct {
	ID uuid.UUID `path:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the TagsPosts item"`
}
type GetTagPostByIDResponse struct {
	ID           uuid.UUID `json:"id"`
	TagID        uuid.UUID `json:"tag_id"`
	PostableID   uuid.UUID `json:"postable_id"`
	PostableType string    `json:"postable_type"`
}

// CREATING SINGLE TAG POST
type CreateTagPostInput struct {
	Body CreateTagPostBody
}
type CreateTagPostBody struct {
	PostableID   uuid.UUID `json:"postable_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post or premium post"`
	PostableType string    `json:"postable_type" example:"post" doc:"Type: post or premium_post"`
	TagID        uuid.UUID `json:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the tag"`
}
type CreateTagPostsResponse struct {
	ID           uuid.UUID `json:"id"`
	PostableID   uuid.UUID `json:"postable_id"`
	PostableType string    `json:"postable_type"`
	TagID        uuid.UUID `json:"tag_id"`
}

// UPDATING TAG/POST
type UpdateTagPostInput struct {
	ID   uuid.UUID `path:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the tagpost to update"`
	Body UpdateTagPostBody
}
type UpdateTagPostBody struct {
	PostableID   uuid.UUID `json:"postable_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"the postable id to update to"`
	PostableType string    `json:"postable_type" example:"post" doc:"Type: post or premium_post"`
	TagID        uuid.UUID `json:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"the tag id to update to"`
}
type UpdateTagPostResponse struct {
	ID           uuid.UUID `json:"id"`
	PostableID   uuid.UUID `json:"postable_id"`
	PostableType string    `json:"postable_type"`
	TagID        uuid.UUID `json:"tag_id"`
}

// DELETING TAG/POST
type DeleteTagPostResponse struct {
	ID uuid.UUID `json:"id"`
}
