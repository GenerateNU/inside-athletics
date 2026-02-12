package comment_like

import (
	"github.com/google/uuid"
)

type GetCommentLikeParams struct {
	ID uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID of the like"`
}

type GetCommentLikeResponse struct {
	UserID    uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"UserID of the like"`
	CommentID uuid.UUID `json:"comment_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"CommentID of the like"`
}

type CreateCommentLikeInput struct {
	Body CreateCommentLikeBody
}

type CreateCommentLikeBody struct {
	UserID    uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"UserID of the like"`
	CommentID uuid.UUID `json:"comment_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"CommentID of the like"`
}

type CreateCommentLikeResponse struct {
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the like"`
	Total int  `json:"total" example:"25" doc:"Total number of likes on the comment"`
	Liked bool `json:"liked" example:"true" doc:"Whether the requested user has liked the comment"`
}

type DeleteCommentLikeParams struct {
	ID uuid.UUID `path:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the like"`
	UserID    uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"UserID of the like"`
}

type DeleteCommentLikeResponse struct {
	Message string `json:"message" example:"Like was deleted successfully" doc:"Message to display"`
	Total   int    `json:"total" example:"25" doc:"Total number of likes on the comment"`
	Liked   bool   `json:"liked" example:"true" doc:"Whether the requested user has liked the comment"`
}

// Retrieves like count and whether the user has liked.
type GetCommentLikeInfoParams struct {
	CommentID uuid.UUID `path:"comment_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Comment to get like info for"`
	UserID    uuid.UUID `query:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Optional; if provided, response includes whether this user has liked"`
}

// Retrieves like count and whether the current user has liked.
type GetCommentLikeInfoResponse struct {
	Total int  `json:"total" example:"25" doc:"Total number of likes on the comment"`
	Liked bool `json:"liked" example:"true" doc:"Whether the requested user has liked the comment"`
}
