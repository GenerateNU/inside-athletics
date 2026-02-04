package like

import (
	"github.com/google/uuid"
)

type GetCommentLikeParams struct {
	ID uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID to identify the user"`
}

type GetCommentLikeResponse struct {
	UserID uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"UserID of the like"`
	CommentID uuid.UUID `json:"comment_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"CommentID of the like"`
}

// CreateCommentLikeRequest is the request body for POST (no path params)
type CreateCommentLikeRequest struct {
	UserID    uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"UserID of the like"`
	CommentID uuid.UUID `json:"comment_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"CommentID of the like"`
}

// Do we want to return the id?
type CreateCommentLikeResponse struct {
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the like"`
}

type DeleteCommentLikeParams struct {
	ID uuid.UUID `path:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the like"`
}

type DeleteCommentLikeResponse struct {
	Message string `json:"message" example:"Like was deleted successfully" doc:"Message to display"`
}