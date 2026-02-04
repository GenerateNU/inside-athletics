package like

import (
	"context"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type PostLikeService struct {
	postLikeDB *PostLikeDB
}

// GetPostLike retrieves a like given an id
func (u *PostLikeService) GetPostLike(ctx context.Context, input *GetPostLikeParams) (*utils.ResponseBody[GetPostLikeResponse], error) {
	like, err := u.postLikeDB.GetPostLike(input.ID)
	respBody := &utils.ResponseBody[GetPostLikeResponse]{}
	if err != nil {
		return respBody, err
	}
	return &utils.ResponseBody[GetPostLikeResponse]{
		Body: &GetPostLikeResponse{
			UserID:    like.UserID,
			PostID: like.PostID,
		},
	}, nil
}

// CreatePostLike creates a like on a post
func (u *PostLikeService) CreatePostLike(ctx context.Context, input *CreatePostLikeRequest) (*utils.ResponseBody[CreatePostLikeResponse], error) {
	postLike := &models.PostLike{
		UserID:    input.UserID,
		PostID: input.PostID,
	}
	created, err := u.postLikeDB.CreatePostLike(postLike)
	respBody := &utils.ResponseBody[CreatePostLikeResponse]{}
	if err != nil {
		return respBody, err
	}
	return &utils.ResponseBody[CreatePostLikeResponse]{
		Body: &CreatePostLikeResponse{ID: created.ID},
	}, nil
}

// DeletePostLike deletes a like on a post
func (u *PostLikeService) DeletePostLike(ctx context.Context, input *DeletePostLikeParams) (*utils.ResponseBody[DeletePostLikeResponse], error) {
	err := u.postLikeDB.DeletePostLike(input.ID)
	respBody := &utils.ResponseBody[DeletePostLikeResponse]{}
	if err != nil {
		return respBody, err
	}
	return &utils.ResponseBody[DeletePostLikeResponse]{
		Body: &DeletePostLikeResponse{Message: "Like was deleted successfully"},
	}, nil
}
