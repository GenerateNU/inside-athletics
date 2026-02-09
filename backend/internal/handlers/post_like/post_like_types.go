package post_like

import (
	"github.com/google/uuid"
)

type GetPostLikeParams struct {
	ID uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID of the like"`
}

type GetPostLikeResponse struct {
	UserID uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"UserID of the like"`
	PostID uuid.UUID `json:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"PostID of the like"`
}

type CreatePostLikeInput struct {
	Body CreatePostLikeBody
}

type CreatePostLikeBody struct {
	UserID uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"UserID of the like"`
	PostID uuid.UUID `json:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"PostID of the like"`
}

type CreatePostLikeResponse struct {
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the like"`
}

type DeletePostLikeParams struct {
	ID uuid.UUID `path:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the like"`
}

type DeletePostLikeResponse struct {
	Message string `json:"message" example:"Like was deleted successfully" doc:"Message to display"`
}

// For the single endpoint that returns like count and whether the user has liked.
type GetPostLikeInfoParams struct {
	PostID uuid.UUID `path:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Post to get like info for"`
	UserID uuid.UUID `query:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Optional; if provided, response includes whether this user has liked"`
}

// Retrieves like count and whether the current user has liked.
type GetPostLikeInfoResponse struct {
	Total int  `json:"total" example:"25" doc:"Total number of likes on the post"`
	Liked bool `json:"liked" example:"true" doc:"Whether the requested user has liked the post"`
}
