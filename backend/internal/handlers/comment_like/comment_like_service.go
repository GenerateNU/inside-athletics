package comment_like

import (
	"context"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
)

type CommentLikeService struct {
	commentLikeDB *CommentLikeDB
}

// Retrieves a like by ID.
func (u *CommentLikeService) GetCommentLike(ctx context.Context, input *GetCommentLikeParams) (*utils.ResponseBody[GetCommentLikeResponse], error) {
	like, err := u.commentLikeDB.GetCommentLike(input.ID)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[GetCommentLikeResponse]{
		Body: &GetCommentLikeResponse{
			UserID:    like.UserID,
			CommentID: like.CommentID,
		},
	}, nil
}

// Creates a like on a comment. Returns 409 if the user has already liked the comment.
// Response includes total likes on the comment and liked=true for the requesting user.
func (u *CommentLikeService) CreateCommentLike(ctx context.Context, input *CreateCommentLikeInput) (*utils.ResponseBody[CreateCommentLikeResponse], error) {
	_, liked, err := u.commentLikeDB.GetCommentLikeInfo(input.Body.CommentID, input.Body.UserID)
	if err != nil {
		return nil, err
	}
	if liked {
		return nil, huma.Error409Conflict("User has already liked this comment")
	}
	commentLike := &models.CommentLike{
		UserID:    input.Body.UserID,
		CommentID: input.Body.CommentID,
	}
	created, err := u.commentLikeDB.CreateCommentLike(commentLike)
	if err != nil {
		return nil, err
	}
	total, _, err := u.commentLikeDB.GetCommentLikeInfo(input.Body.CommentID, input.Body.UserID)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[CreateCommentLikeResponse]{
		Body: &CreateCommentLikeResponse{
			ID:    created.ID,
			Total: int(total),
			Liked: true,
		},
	}, nil
}

// Deletes a like by ID. Response includes updated total likes on the comment and if user liked comment.
func (u *CommentLikeService) DeleteCommentLike(ctx context.Context, input *DeleteCommentLikeParams) (*utils.ResponseBody[DeleteCommentLikeResponse], error) {
	like, err := u.commentLikeDB.GetCommentLike(input.ID)
	if err != nil {
		return nil, err
	}
	commentID := like.CommentID
	err = u.commentLikeDB.DeleteCommentLike(input.ID)
	if err != nil {
		return nil, err
	}
	total, liked, err := u.commentLikeDB.GetCommentLikeInfo(commentID, input.UserID)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[DeleteCommentLikeResponse]{
		Body: &DeleteCommentLikeResponse{
			Message: "Like was deleted successfully",
			Total:   int(total),
			Liked:   liked,
		},
	}, nil
}

// Retrieves like count and whether the user has liked the comment.
func (u *CommentLikeService) GetCommentLikeInfo(ctx context.Context, input *GetCommentLikeInfoParams) (*utils.ResponseBody[GetCommentLikeInfoResponse], error) {
	count, liked, err := u.commentLikeDB.GetCommentLikeInfo(input.CommentID, input.UserID)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[GetCommentLikeInfoResponse]{
		Body: &GetCommentLikeInfoResponse{Total: int(count), Liked: liked},
	}, nil
}
