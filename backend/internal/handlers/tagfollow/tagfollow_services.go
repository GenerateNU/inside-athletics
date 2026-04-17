package tagfollow

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type TagFollowService struct {
	tagfollowDB *TagFollowDB
}

func (u *TagFollowService) getCurrentUserID(ctx context.Context) (uuid.UUID, error) {
	rawID := ctx.Value("user_id")
	if rawID == nil {
		return uuid.Nil, huma.Error401Unauthorized("User not authenticated")
	}

	userID, ok := rawID.(string)
	if !ok {
		return uuid.Nil, huma.Error500InternalServerError("Invalid user ID in context")
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, huma.Error400BadRequest("Invalid user ID", err)
	}

	return parsedID, nil
}

// Given a UserID, get all the tags that they follow
func (u *TagFollowService) GetTagFollowsByUser(ctx context.Context, input *GetTagFollowsByUserParams) (*utils.ResponseBody[GetTagFollowsByUserResponse], error) {
	userID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}

	tags, err := u.tagfollowDB.GetTagFollowsByUser(userID)
	respBody := &utils.ResponseBody[GetTagFollowsByUserResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetTagFollowsByUserResponse{
		UserID: userID,
		TagIDs: *tags,
	}

	return &utils.ResponseBody[GetTagFollowsByUserResponse]{
		Body: response,
	}, err
}

// Given a TagID, get all the users that follow this tag
func (u *TagFollowService) GetFollowingUsersByTag(ctx context.Context, input *GetFollowingUsersByTagParams) (*utils.ResponseBody[GetFollowingUsersByTagResponse], error) {
	tagID := input.TagID
	users, err := u.tagfollowDB.GetFollowingUsersByTag(tagID)
	respBody := &utils.ResponseBody[GetFollowingUsersByTagResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetFollowingUsersByTagResponse{
		TagID:   tagID,
		UserIDs: *users,
	}

	return &utils.ResponseBody[GetFollowingUsersByTagResponse]{
		Body: response,
	}, err
}

// Given a tag and a user, creates a tag follow if doesn't already exist
func (u *TagFollowService) CreateTagFollow(ctx context.Context, input *CreateTagFollowInput) (*utils.ResponseBody[CreateTagFollowResponse], error) {
	userID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}

	tagfollow := &models.TagFollow{
		TagID:  input.Body.TagID,
		UserID: userID,
	}

	createdTagFollow, err := u.tagfollowDB.CreateTagFollow(tagfollow)

	respBody := &utils.ResponseBody[CreateTagFollowResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &CreateTagFollowResponse{
		ID:     createdTagFollow.ID,
		UserID: createdTagFollow.UserID,
		TagID:  createdTagFollow.TagID,
	}

	return &utils.ResponseBody[CreateTagFollowResponse]{
		Body: response,
	}, err
}

// Soft deletes tag follow
func (u *TagFollowService) DeleteTagFollow(ctx context.Context, input *DeleteTagFollowParams) (*utils.ResponseBody[DeleteTagFollowResponse], error) {
	respBody := &utils.ResponseBody[DeleteTagFollowResponse]{}

	err := u.tagfollowDB.DeleteTagFollow(input.ID)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &DeleteTagFollowResponse{
		Message: "Tag follow was deleted successfully",
	}

	return respBody, nil
}

// Soft deletes tag follow by tag id for current user.
func (u *TagFollowService) DeleteTagFollowByTag(ctx context.Context, input *DeleteTagFollowByTagParams) (*utils.ResponseBody[DeleteTagFollowResponse], error) {
	respBody := &utils.ResponseBody[DeleteTagFollowResponse]{}
	userID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return respBody, err
	}

	err = u.tagfollowDB.DeleteTagFollowByUserAndTag(userID, input.TagID)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &DeleteTagFollowResponse{
		Message: "Tag follow was deleted successfully",
	}
	return respBody, nil
}
