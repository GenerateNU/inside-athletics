package tag

import (
	"context"
	"inside-athletics/internal/handlers/post"
	"inside-athletics/internal/handlers/tagpost"
	"inside-athletics/internal/models"
	"inside-athletics/internal/s3"
	"inside-athletics/internal/utils"
	"net/url"
)

type TagService struct {
	tagDB     *TagDB
	tagPostDB *tagpost.TagPostDB
	s3        *s3.Service
}

func (u *TagService) GetTagByName(ctx context.Context, input *GetTagByNameParams) (*utils.ResponseBody[GetTagResponse], error) {
	name, decodeErr := url.PathUnescape(input.Name)
	if decodeErr != nil {
		name = input.Name
	}
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
		p := &(*posts)[i]
		if url := s3.ResolveKey(ctx, u.s3, p.Author.ProfilePicture); url != "" {
			p.Author.ProfilePicture = url
		}
		if p.College != nil {
			if url := s3.ResolveKey(ctx, u.s3, p.College.Logo); url != "" {
				p.College.Logo = url
			}
		}
		postResponses = append(postResponses, *post.ToPostResponse(p, userID))
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

func (u *TagService) FuzzySearchFor(ctx context.Context, input *utils.SearchParam) (*utils.ResponseBody[utils.SearchResults[*GetTagResponse]], error) {
	return utils.FuzzySearchService(input, models.Tag{}, GetTagResponse{}, "name", u.tagDB.db, getTagResponse)
}

func getTagResponse(tag *models.Tag) *GetTagResponse {
	return &GetTagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}
}
