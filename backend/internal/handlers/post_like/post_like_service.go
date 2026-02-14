package post_like

import (
	"context"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
)

type PostLikeService struct {
	postLikeDB *PostLikeDB
}

// Retrieves a like by ID.
func (u *PostLikeService) GetPostLike(ctx context.Context, input *GetPostLikeParams) (*utils.ResponseBody[GetPostLikeResponse], error) {
	like, err := u.postLikeDB.GetPostLike(input.ID)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[GetPostLikeResponse]{
		Body: &GetPostLikeResponse{
			UserID: like.UserID,
			PostID: like.PostID,
		},
	}, nil
}

// Creates a like on a post. Returns 409 if the user has already liked the post.
// Response includes total likes on the post and liked=true for the requesting user.
func (u *PostLikeService) CreatePostLike(ctx context.Context, input *CreatePostLikeInput) (*utils.ResponseBody[CreatePostLikeResponse], error) {
	_, liked, err := u.postLikeDB.GetPostLikeInfo(input.Body.PostID, input.Body.UserID)
	if err != nil {
		return nil, err
	}
	if liked {
		return nil, huma.Error409Conflict("User has already liked this post")
	}
	postLike := &models.PostLike{
		UserID: input.Body.UserID,
		PostID: input.Body.PostID,
	}
	created, err := u.postLikeDB.CreatePostLike(postLike)
	if err != nil {
		return nil, err
	}
	total, _, err := u.postLikeDB.GetPostLikeInfo(input.Body.PostID, input.Body.UserID)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[CreatePostLikeResponse]{
		Body: &CreatePostLikeResponse{
			ID:    created.ID,
			Total: int(total),
			Liked: true,
		},
	}, nil
}

// Deletes a like by ID. Response includes updated total likes on the post and if user liked post.
func (u *PostLikeService) DeletePostLike(ctx context.Context, input *DeletePostLikeParams) (*utils.ResponseBody[DeletePostLikeResponse], error) {
	like, err := u.postLikeDB.GetPostLike(input.ID)
	if err != nil {
		return nil, err
	}
	postID := like.PostID
	err = u.postLikeDB.DeletePostLike(input.ID)
	if err != nil {
		return nil, err
	}
	total, liked, err := u.postLikeDB.GetPostLikeInfo(postID, input.UserID)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[DeletePostLikeResponse]{
		Body: &DeletePostLikeResponse{
			Message: "Like was deleted successfully",
			Total:   int(total),
			Liked:   liked,
		},
	}, nil
}

// Retrieves like count and whether the user has liked the post.
func (u *PostLikeService) GetPostLikeInfo(ctx context.Context, input *GetPostLikeInfoParams) (*utils.ResponseBody[GetPostLikeInfoResponse], error) {
	count, liked, err := u.postLikeDB.GetPostLikeInfo(input.PostID, input.UserID)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[GetPostLikeInfoResponse]{
		Body: &GetPostLikeInfoResponse{Total: int(count), Liked: liked},
	}, nil
}
