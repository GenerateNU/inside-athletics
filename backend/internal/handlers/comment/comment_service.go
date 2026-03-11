package comment

import (
	"context"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
)

type CommentService struct {
	commentDB *CommentDB
}

// Creates a new comment.
func (s *CommentService) CreateComment(ctx context.Context, input *CreateCommentInput) (*utils.ResponseBody[CreateCommentResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	// Enforce one layer of replies: parent must be top-level (parent_comment_id IS NULL)
	if input.Body.ParentCommentID != nil {
		parent, err := s.commentDB.GetCommentByID(*input.Body.ParentCommentID, userID)
		if err != nil {
			return nil, err
		}
		if parent.ParentCommentID != nil {
			return nil, huma.Error400BadRequest("Replies only allowed to top-level comments; one layer of replies")
		}
	}

	// Create the comment model
	comment := &models.Comment{
		UserID:          userID,
		IsAnonymous:     input.Body.IsAnonymous,
		ParentCommentID: input.Body.ParentCommentID,
		PostID:          input.Body.PostID,
		Description:     input.Body.Description,
	}

	// Create the comment in the database
	created, err := s.commentDB.CreateComment(comment)
	if err != nil {
		return nil, err
	}

	// Convert the comment to a response
	return &utils.ResponseBody[CreateCommentResponse]{
		Body: ToCreateCommentResponse(created, userID),
	}, nil
}

// Retrieves a single comment by ID.
func (s *CommentService) GetComment(ctx context.Context, input *GetCommentParams) (*utils.ResponseBody[CommentResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	// Get the comment from the database
	comment, err := s.commentDB.GetCommentByID(input.ID, userID)
	if err != nil {
		return nil, err
	}

	// Convert the comment to a response
	return &utils.ResponseBody[CommentResponse]{
		Body: ToCommentResponse(comment, userID),
	}, nil
}

// Retrieves top-level comments for a post.
func (s *CommentService) GetCommentsByPost(ctx context.Context, input *GetCommentsByPostParams) (*utils.ResponseBody[[]CommentResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	// Get the comments from the database
	comments, err := s.commentDB.GetCommentsByPost(input.PostID, userID)
	if err != nil {
		_, humaErr := utils.HandleDBError[[]CommentResponse](nil, err)
		return nil, humaErr
	}

	// Convert the comments to responses
	responses := make([]CommentResponse, len(comments))
	for i := range comments {
		responses[i] = *ToCommentResponse(&comments[i], userID)
	}

	return &utils.ResponseBody[[]CommentResponse]{Body: &responses}, nil
}

// Retrieves replies to a comment.
func (s *CommentService) GetReplies(ctx context.Context, input *GetReplyParams) (*utils.ResponseBody[[]CommentResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	// Get the replies from the database
	comments, err := s.commentDB.GetReplies(input.ID, userID)
	if err != nil {
		_, humaErr := utils.HandleDBError[[]CommentResponse](nil, err)
		return nil, humaErr
	}

	// Convert the replies to responses
	responses := make([]CommentResponse, len(comments))
	for i := range comments {
		responses[i] = *ToCommentResponse(&comments[i], userID)
	}

	return &utils.ResponseBody[[]CommentResponse]{Body: &responses}, nil
}

// Updates a comment's description by ID.
func (s *CommentService) UpdateComment(ctx context.Context, input *UpdateCommentInput) (*utils.ResponseBody[CommentResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	// Update the comment in the database
	updated, err := s.commentDB.UpdateComment(input.ID, input.Body, userID)
	if err != nil {
		return nil, err
	}

	// Convert the comment to a response
	return &utils.ResponseBody[CommentResponse]{
		Body: ToCommentResponse(updated, userID),
	}, nil
}

// Soft-deletes a comment by ID.
func (s *CommentService) DeleteComment(ctx context.Context, input *DeleteCommentRequest) (*utils.ResponseBody[CommentResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	// Get the comment from the database
	comment, err := s.commentDB.GetCommentByID(input.ID, userID)
	if err != nil {
		return nil, err
	}

	// Delete the comment from the database
	err = s.commentDB.DeleteComment(input.ID)
	if err != nil {
		return nil, err
	}

	// Convert the comment to a response
	return &utils.ResponseBody[CommentResponse]{
		Body: ToCommentResponse(comment, userID),
	}, nil
}
