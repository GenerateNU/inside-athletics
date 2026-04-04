package tag

import (
	"context"
	"inside-athletics/internal/handlers/post"
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
		Type: tag.Type,
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
		Type: tag.Type,
	}

	return &utils.ResponseBody[GetTagResponse]{
		Body: response,
	}, err
}

func (u *TagService) GetTagsByType(ctx context.Context, input *GetTagsByTypeParams) (*utils.ResponseBody[[]GetTagResponse], error) {
	tagType := input.Type
	tags, err := u.tagDB.GetTagsByType(tagType)
	respBody := &utils.ResponseBody[[]GetTagResponse]{}

	if err != nil {
		return respBody, err
	}

	responses := make([]GetTagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = GetTagResponse{
			ID:   tag.ID,
			Name: tag.Name,
			Type: tag.Type,
		}
	}

	return &utils.ResponseBody[[]GetTagResponse]{
		Body: &responses,
	}, nil
}

// Returns an array of post ids that are tagged with a unique tag, determined by the tag id.
func (u *TagService) GetPostsByTag(ctx context.Context, input *GetPostsByTagParam) (*utils.ResponseBody[GetPostsByTagResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	posts, err := u.tagDB.GetPostsByTag(input.TagID, input.Limit, input.Offset, userID)
	respBody := &utils.ResponseBody[GetPostsByTagResponse]{}

	if err != nil {
		return respBody, err
	}

	postResponses := make([]post.PostResponse, 0, len(*posts))
	for i := range *posts {
		postResponses = append(postResponses, *post.ToPostResponse(&((*posts)[i]), userID))
	}

	response := &GetPostsByTagResponse{
		Posts: postResponses,
	}

	return &utils.ResponseBody[GetPostsByTagResponse]{
		Body: response,
	}, err
}

func (u *TagService) CreateTag(ctx context.Context, input *CreateTagInput) (*utils.ResponseBody[CreateTagResponse], error) {
	respBody := &utils.ResponseBody[CreateTagResponse]{}

	tag := &models.Tag{
		Name: input.Body.Name,
		Type: input.Body.Type,
	}

	createdTag, err := u.tagDB.CreateTag(tag)

	if err != nil {
		return respBody, err
	}

	response := &CreateTagResponse{
		ID:   createdTag.ID,
		Name: createdTag.Name,
		Type: createdTag.Type,
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
		Type: updatedTag.Type,
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
