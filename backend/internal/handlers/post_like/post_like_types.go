package like

import (
	"github.com/google/uuid"
)

type GetPostLikeParams struct {
	ID uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID to identify the user"`
}

type GetPostLikeResponse struct {
	UserID uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"UserID of the like"`
	PostID uuid.UUID `json:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"PostID of the like"`
}

type CreatePostLikeRequest struct {
	UserID uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"UserID of the like"`
	PostID uuid.UUID `json:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"PostID of the like"`
}

// Maybe we should also return the userID and postID of the like, not sure
type CreatePostLikeResponse struct {
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the like"`
}

type DeletePostLikeParams struct {
	ID uuid.UUID `path:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the like"`
}

type DeletePostLikeResponse struct {
	Message string `json:"message" example:"Like was deleted successfully" doc:"Message to display"`
}
