package like

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

// Retrieves all likes from post as int
type GetLikeCountParams struct {
	PostID uuid.UUID `path:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"PostID to count likes for"`
}

// Response body for retrieving all likes from post as int
type GetLikeCountResponse struct {
	Total int `json:"total" example:"25" doc:"Total number of likes on post"`
}

// Checks if given user has liked the post
type CheckUserLikedPostParams struct {
	UserID uuid.UUID `query:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"User to check"`
	PostID uuid.UUID `path:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Post to check"`
}

type CheckUserLikedPostResponse struct {
	Liked bool `json:"liked" example:"true" doc:"Whether user liked the post"`
}
