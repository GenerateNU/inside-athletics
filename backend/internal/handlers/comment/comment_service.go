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

// Checks if the caller can see user_id on anonymous comments.
// Lowkenuinely not sure how we should do this, ask TLs in PR
// In prod, when real auth/roles exist, should derive from role/permission instead of context.
func (s *CommentService) forSuperUser(ctx context.Context) bool {
	raw := ctx.Value("for_super_user")
	if raw == nil {
		return false
	}
	b, ok := raw.(bool)
	return ok && b
}

// Creates a new comment.
func (s *CommentService) CreateComment(ctx context.Context, input *CreateCommentInput) (*utils.ResponseBody[CommentResponse], error) {
	// Enforce one layer of replies: parent must be top-level (parent_comment_id IS NULL)
	if input.Body.ParentCommentID != nil {
		parent, err := s.commentDB.GetCommentByID(*input.Body.ParentCommentID)
		if err != nil {
			return nil, err
		}
		if parent.ParentCommentID != nil {
			return nil, huma.Error400BadRequest("Replies only allowed to top-level comments; one layer of replies")
		}
	}

	// Create the comment model
	comment := &models.Comment{
		UserID:          input.Body.UserID,
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
	return &utils.ResponseBody[CommentResponse]{
		Body: ToCommentResponse(created, s.forSuperUser(ctx)),
	}, nil
}

// Retrieves a single comment by ID.
func (s *CommentService) GetComment(ctx context.Context, input *GetCommentParams) (*utils.ResponseBody[CommentResponse], error) {
	// Get the comment from the database
	comment, err := s.commentDB.GetCommentByID(input.ID)
	if err != nil {
		return nil, err
	}

	// Convert the comment to a response
	return &utils.ResponseBody[CommentResponse]{
		Body: ToCommentResponse(comment, s.forSuperUser(ctx)),
	}, nil
}

// Retrieves top-level comments for a post.
func (s *CommentService) GetCommentsByPost(ctx context.Context, input *GetCommentsByPostParams) (*utils.ResponseBody[[]CommentResponse], error) {
	// Get the comments from the database
	comments, err := s.commentDB.GetCommentsByPost(input.PostID)
	if err != nil {
		_, humaErr := utils.HandleDBError[[]CommentResponse](nil, err)
		return nil, humaErr
	}

	// Convert the comments to responses
	responses := make([]CommentResponse, len(comments))
	for i := range comments {
		responses[i] = *ToCommentResponse(&comments[i], s.forSuperUser(ctx))
	}

	return &utils.ResponseBody[[]CommentResponse]{Body: &responses}, nil
}

// Retrieves replies to a comment.
func (s *CommentService) GetReplies(ctx context.Context, input *GetReplyParams) (*utils.ResponseBody[[]CommentResponse], error) {
	// Get the replies from the database
	comments, err := s.commentDB.GetReplies(input.ID)
	if err != nil {
		_, humaErr := utils.HandleDBError[[]CommentResponse](nil, err)
		return nil, humaErr
	}

	// Convert the replies to responses
	responses := make([]CommentResponse, len(comments))
	for i := range comments {
		responses[i] = *ToCommentResponse(&comments[i], s.forSuperUser(ctx))
	}

	return &utils.ResponseBody[[]CommentResponse]{Body: &responses}, nil
}

// Updates a comment's description by ID.
func (s *CommentService) UpdateComment(ctx context.Context, input *UpdateCommentInput) (*utils.ResponseBody[CommentResponse], error) {
	// Get the comment from the database
	comment, err := s.commentDB.GetCommentByID(input.ID)
	if err != nil {
		return nil, err
	}

	comment.Description = input.Body.Description
	// Update the comment in the database
	updated, err := s.commentDB.UpdateComment(comment)
	if err != nil {
		return nil, err
	}

	// Convert the comment to a response
	return &utils.ResponseBody[CommentResponse]{
		Body: ToCommentResponse(updated, s.forSuperUser(ctx)),
	}, nil
}

// Soft-deletes a comment by ID.
func (s *CommentService) DeleteComment(ctx context.Context, input *DeleteCommentRequest) (*utils.ResponseBody[CommentResponse], error) {
	// Get the comment from the database
	comment, err := s.commentDB.GetCommentByID(input.ID)
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
		Body: ToCommentResponse(comment, s.forSuperUser(ctx)),
	}, nil
}
