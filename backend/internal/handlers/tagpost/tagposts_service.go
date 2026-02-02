package tagpost

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"reflect"
	"strings"
)

type TagPostService struct {
	tagpostDB *TagPostDB
}

func (u *TagPostService) GetTagsByPost(ctx context.Context, input *GetTagsByPostParam) (*utils.ResponseBody[GetTagsByPostResponse], error) {
	postID := input.PostID
	tags, err := u.tagpostDB.GetTagsByPost(postID)
	respBody := &utils.ResponseBody[GetTagsByPostResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetTagsByPostResponse{
		PostID: postID,
		TagIDs: *tags,
	}

	return &utils.ResponseBody[GetTagsByPostResponse]{
		Body: response,
	}, err
}

func (u *TagPostService) GetPostsByTag(ctx context.Context, input *GetPostsByTagParam) (*utils.ResponseBody[GetPostsbyTagResponse], error) {
	tagID := input.TagID
	posts, err := u.tagpostDB.GetPostsByTag(tagID)
	respBody := &utils.ResponseBody[GetPostsbyTagResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetPostsbyTagResponse{
		TagID:   tagID,
		PostIDs: *posts,
	}

	return &utils.ResponseBody[GetPostsbyTagResponse]{
		Body: response,
	}, err
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

	updates, err := buildTagUpdates(input.Body)
	if err != nil {
		return respBody, err
	}

	updatedTagPost, err := u.tagpostDB.UpdateTagPost(input.ID, updates)
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

func buildTagUpdates(body UpdateTagPostBody) (map[string]interface{}, error) {
	updates := make(map[string]interface{})
	val := reflect.ValueOf(body)
	typ := reflect.TypeOf(body)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name := strings.Split(tag, ",")[0]
		if name == "" {
			continue
		}

		fieldVal := val.Field(i)
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				continue
			}
			updates[name] = fieldVal.Elem().Interface()
			continue
		}

		updates[name] = fieldVal.Interface()
	}

	return updates, nil
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
