package tagpost

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type TagPostService struct {
	tagpostDB *TagPostDB
}

func (u *TagPostService) GetTagPostById(ctx context.Context, input *GetTagPostByIdParam) (*utils.ResponseBody[GetTagPostByIDResponse], error) {
	id := input.ID
	tagpost, err := u.tagpostDB.GetTagPostById(id)
	respBody := &utils.ResponseBody[GetTagPostByIDResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetTagPostByIDResponse{
		ID:     id,
		TagID:  tagpost.TagID,
		PostID: tagpost.PostID,
	}

	return &utils.ResponseBody[GetTagPostByIDResponse]{
		Body: response,
	}, err
}

func (u *TagPostService) CreateTagPost(ctx context.Context, input *CreateTagPostInput) (*utils.ResponseBody[CreateTagPostsResponse], error) {
	respBody := &utils.ResponseBody[CreateTagPostsResponse]{}

	tagpost := &models.TagPost{
		PostID: input.Body.PostID,
		TagID:  input.Body.TagID,
	}

	createdTagPost, err := u.tagpostDB.CreateTagPost(tagpost)

	if err != nil {
		return respBody, err
	}

	response := &CreateTagPostsResponse{
		ID:     createdTagPost.ID,
		PostID: createdTagPost.PostID,
		TagID:  createdTagPost.TagID,
	}

	return &utils.ResponseBody[CreateTagPostsResponse]{
		Body: response,
	}, err
}

func (u *TagPostService) UpdateTagPost(cts context.Context, input *UpdateTagPostInput) (*utils.ResponseBody[UpdateTagPostResponse], error) {
	respBody := &utils.ResponseBody[UpdateTagPostResponse]{}

	updatedTagPost, err := u.tagpostDB.UpdateTagPost(input.ID, &input.Body)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &UpdateTagPostResponse{
		ID:     updatedTagPost.ID,
		PostID: updatedTagPost.PostID,
		TagID:  updatedTagPost.TagID,
	}

	return respBody, nil
}

func (u *TagPostService) DeleteTagPost(ctx context.Context, input *GetTagPostByIdParam) (*utils.ResponseBody[DeleteTagPostResponse], error) {
	respBody := &utils.ResponseBody[DeleteTagPostResponse]{}

	err := u.tagpostDB.DeleteTagPost(input.ID)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &DeleteTagPostResponse{
		ID: input.ID,
	}

	return respBody, nil
}
