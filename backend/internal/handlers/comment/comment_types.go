package comment

import "github.com/google/uuid"

type GetCommentParams struct {
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of comment"`
}

type GetReplyParams struct {
	ParentCommentID *uuid.UUID `json:"parent_comment_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"Comment the comment belongs to"`
}

type CommentResponse struct {
	ID              uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of comment"`
	UserID          uuid.UUID  `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"UserID the comment belongs to"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"CommentID the comment belongs to"`
	PostID          uuid.UUID  `json:"post_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"PostID the comment belongs to"`
	Description     string     `json:"description" example:"This is a helpful thread" maxLength:"1500" doc:"Content of the comment"`
}

type CreateCommentRequest struct {
	UserID          uuid.UUID  `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"UserID the comment belongs to"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"CommentID the comment belongs to"`
	PostID          uuid.UUID  `json:"post_id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"PostID the comment belongs to"`
	Description     string     `json:"description" example:"This is a helpful thread" maxLength:"300" doc:"Content of the comment"`
}

type UpdateCommentRequest struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of comment"`
	Description string    `json:"description" example:"This is a helpful thread" maxLength:"300" doc:"Content of the comment"`
}

type UpdateCommentResponse struct {
	Description string `json:"description" example:"This is a helpful thread" maxLength:"300" doc:"Content of the comment"`
}

type DeleteCommentParams struct {
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of comment to delete"`
}

type DeleteCommentResponse struct {
	Message string    `json:"message" example:"Comment deleted successfully" doc:"Success message"`
	ID      uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of deleted comment"`
}
