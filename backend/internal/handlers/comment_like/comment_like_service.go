package like

import (
	"context"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type CommentLikeService struct {
	commentLikeDB *CommentLikeDB
}

// GetCommentLike retrieves a like given an id
func (u *CommentLikeService) GetCommentLike(ctx context.Context, input *GetCommentLikeParams) (*utils.ResponseBody[GetCommentLikeResponse], error) {
	like, err := u.commentLikeDB.GetCommentLike(input.ID)
	respBody := &utils.ResponseBody[GetCommentLikeResponse]{}
	if err != nil {
		return respBody, err
	}
	return &utils.ResponseBody[GetCommentLikeResponse]{
		Body: &GetCommentLikeResponse{
			UserID:    like.UserID,
			CommentID: like.CommentID,
		},
	}, nil
}

// CreateCommentLike creates a like on a comment
func (u *CommentLikeService) CreateCommentLike(ctx context.Context, input *CreateCommentLikeRequest) (*utils.ResponseBody[CreateCommentLikeResponse], error) {
	commentLike := &models.CommentLike{
		UserID:    input.UserID,
		CommentID: input.CommentID,
	}
	created, err := u.commentLikeDB.CreateCommentLike(commentLike)
	respBody := &utils.ResponseBody[CreateCommentLikeResponse]{}
	if err != nil {
		return respBody, err
	}
	return &utils.ResponseBody[CreateCommentLikeResponse]{
		Body: &CreateCommentLikeResponse{ID: created.ID},
	}, nil
}

// DeleteCommentLike deletes a like on a comment
func (u *CommentLikeService) DeleteCommentLike(ctx context.Context, input *DeleteCommentLikeParams) (*utils.ResponseBody[DeleteCommentLikeResponse], error) {
	err := u.commentLikeDB.DeleteCommentLike(input.ID)
	respBody := &utils.ResponseBody[DeleteCommentLikeResponse]{}
	if err != nil {
		return respBody, err
	}
	return &utils.ResponseBody[DeleteCommentLikeResponse]{
		Body: &DeleteCommentLikeResponse{Message: "Like was deleted successfully"},
	}, nil
}

// GetLikeCount returns the total number of likes for a comment
func (u *CommentLikeService) GetLikeCount(ctx context.Context, input *GetLikeCountParams) (*utils.ResponseBody[GetLikeCountResponse], error) {
	count, err := u.commentLikeDB.GetLikeCount(input.CommentID)
	respBody := &utils.ResponseBody[GetLikeCountResponse]{}
	if err != nil {
		return respBody, err
	}
	return &utils.ResponseBody[GetLikeCountResponse]{
		Body: &GetLikeCountResponse{Total: int(count)},
	}, nil
}

// CheckUserLikedComment returns whether the given user has liked the comment
func (u *CommentLikeService) CheckUserLikedComment(ctx context.Context, input *CheckUserLikedCommentParams) (*utils.ResponseBody[CheckUserLikedCommentResponse], error) {
	liked, err := u.commentLikeDB.CheckUserLikedComment(input.UserID, input.CommentID)
	respBody := &utils.ResponseBody[CheckUserLikedCommentResponse]{}
	if err != nil {
		return respBody, err
	}
	return &utils.ResponseBody[CheckUserLikedCommentResponse]{
		Body: &CheckUserLikedCommentResponse{Liked: liked},
	}, nil
}
