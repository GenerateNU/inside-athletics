package comment

import (
	models "inside-athletics/internal/models"

	"github.com/google/uuid"
)

// Defines parameters for getting a comment by ID
type GetCommentParams struct {
	ID uuid.UUID `path:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of comment"`
}

// Defines parameters for getting top-level comments for a post
type GetCommentsByPostParams struct {
	PostID uuid.UUID `path:"post_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of the post"`
}

// Defines parameters for getting replies to a comment
type GetReplyParams struct {
	ID uuid.UUID `path:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of parent comment to get replies for"`
}

// Defines the response structure for a comment
type CommentResponse struct {
	ID                uuid.UUID    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of comment"`
	User              *models.User `json:"user,omitempty" doc:"The user who commented; omitted for anonymous comments"`
	IsAnonymous       bool         `json:"is_anonymous" doc:"True if posted as anonymous; frontend can show 'Anonymous' when user_id is omitted"`
	ParentCommentID   *uuid.UUID   `json:"parent_comment_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"CommentID this comment is in response to"`
	PostID            uuid.UUID    `json:"post_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"PostID of the post this comment is under"`
	Description       string       `json:"description" example:"This is a helpful thread" maxLength:"1500" doc:"Content of the comment"`
	LikeCount         int64        `json:"like_count" doc:"Number of total likes on comment" example:"20000" gorm:"type:int"`
	IsLiked           bool         `json:"is_liked" doc:"If current user has liked this comment" example:"true" gorm:"type:bool"`
	IsVerifiedAthlete bool         `json:"is_verified_athlete" doc:"If commenter is a verified athlete" example:"true" gorm:"type:bool"`
}

type CreateCommentResponse struct {
	ID              uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of comment"`
	UserID          uuid.UUID  `json:"user,omitempty" doc:"The user who commented; omitted for anonymous comments"`
	IsAnonymous     bool       `json:"is_anonymous" doc:"True if posted as anonymous; frontend can show 'Anonymous' when user_id is omitted"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"CommentID the comment belongs to"`
	PostID          uuid.UUID  `json:"post_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"PostID the comment belongs to"`
	Description     string     `json:"description" example:"This is a helpful thread" maxLength:"1500" doc:"Content of the comment"`
}

// The full input for creating a comment
type CreateCommentInput struct {
	Body CreateCommentBody
}

// Defines the request body for creating a new comment
type CreateCommentBody struct {
	IsAnonymous     bool       `json:"is_anonymous" doc:"If true, user_id is hidden from regular users; super user always sees it"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000" doc:"CommentID the comment belongs to"`
	PostID          uuid.UUID  `json:"post_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"PostID the comment belongs to"`
	Description     string     `json:"description" example:"This is a helpful thread" maxLength:"300" doc:"Content of the comment"`
}

// The full input for updating a comment
type UpdateCommentInput struct {
	ID   uuid.UUID `path:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of comment to update"`
	Body UpdateCommentBody
}

// Defines the request body for updating a comment
type UpdateCommentBody struct {
	Description string `json:"description" example:"This is a helpful thread" maxLength:"300" doc:"Updated comment text"`
}

// Defines the request for deleting a comment
type DeleteCommentRequest struct {
	ID uuid.UUID `path:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of comment to delete"`
}

// Converts a Comment model to a CommentResponse, sets UserID to nil for anonymous comments when caller is not super user.
func ToCommentResponse(c *models.Comment, id uuid.UUID) *CommentResponse {
	var user *models.User
	if (!c.IsAnonymous) || (id == c.UserID) {
		u := c.User
		user = &u
	}
	return &CommentResponse{
		ID:                c.ID,
		User:              user,
		IsAnonymous:       c.IsAnonymous,
		ParentCommentID:   c.ParentCommentID,
		PostID:            c.PostID,
		Description:       c.Description,
		LikeCount:         c.LikeCount,
		IsLiked:           c.IsLiked,
		IsVerifiedAthlete: c.User.Verified_Athlete_Status == models.VerifiedAthleteStatusVerified,
	}
}

// Converts a Comment model to a CreateCommentResponse, sets UserID to nil for anonymous comments when caller is not super user.
func ToCreateCommentResponse(c *models.Comment, id uuid.UUID) *CreateCommentResponse {
	return &CreateCommentResponse{
		ID:              c.ID,
		UserID:          c.UserID,
		IsAnonymous:     c.IsAnonymous,
		ParentCommentID: c.ParentCommentID,
		PostID:          c.PostID,
		Description:     c.Description,
	}
}
