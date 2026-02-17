package tag

import (
	"context"
	"inside-athletics/internal/handlers/tagpost"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type TagService struct {
	tagDB     *TagDB
	tagPostDB *tagpost.TagPostDB
}

func (u *TagService) GetTagByName(ctx context.Context, input *GetTagByNameParams) (*utils.ResponseBody[GetTagResponse], error) {
	name := input.Name
	tag, err := u.tagDB.GetTagByName(name)
	respBody := &utils.ResponseBody[GetTagResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetTagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}

	return &utils.ResponseBody[GetTagResponse]{
		Body: response,
	}, err
}

func (u *TagService) GetTagById(ctx context.Context, input *GetTagByIDParams) (*utils.ResponseBody[GetTagResponse], error) {
	id := input.ID
	tag, err := u.tagDB.GetTagByID(id)
	respBody := &utils.ResponseBody[GetTagResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetTagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}

	return &utils.ResponseBody[GetTagResponse]{
		Body: response,
	}, err
}

// GetPostsByTag retrieves post IDs for a tag
func (u *TagService) GetPostsByTag(ctx context.Context, input *GetPostsByTagParams) (*utils.ResponseBody[GetPostsByTagResponse], error) {
	posts, err := u.tagPostDB.GetPostsByTag(input.TagID)
	respBody := &utils.ResponseBody[GetPostsByTagResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetPostsByTagResponse{
		TagID:   input.TagID,
		PostIDs: *posts,
	}

	return &utils.ResponseBody[GetPostsByTagResponse]{
		Body: response,
	}, err
}

func (u *TagService) CreateTag(ctx context.Context, input *CreateTagInput) (*utils.ResponseBody[CreateTagResponse], error) {
	respBody := &utils.ResponseBody[CreateTagResponse]{}

	tag := &models.Tag{
		Name: input.Body.Name,
	}

	createdTag, err := u.tagDB.CreateTag(tag)

	if err != nil {
		return respBody, err
	}

	response := &CreateTagResponse{
		ID:   createdTag.ID,
		Name: createdTag.Name,
	}

	return &utils.ResponseBody[CreateTagResponse]{
		Body: response,
	}, err
}

func (u *TagService) UpdateTag(cts context.Context, input *UpdateTagInput) (*utils.ResponseBody[UpdateTagResponse], error) {
	respBody := &utils.ResponseBody[UpdateTagResponse]{}

	updatedTag, err := u.tagDB.UpdateTag(input.ID, &input.Body)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &UpdateTagResponse{
		ID:   updatedTag.ID,
		Name: updatedTag.Name,
	}

	return respBody, nil
}

func (u *TagService) DeleteTag(ctx context.Context, input *GetTagByIDParams) (*utils.ResponseBody[DeleteTagResponse], error) {
	respBody := &utils.ResponseBody[DeleteTagResponse]{}

	err := u.tagDB.DeleteTag(input.ID)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &DeleteTagResponse{
		ID: input.ID,
	}

	return respBody, nil
}
