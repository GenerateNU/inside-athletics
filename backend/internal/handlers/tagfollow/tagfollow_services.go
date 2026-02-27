package tagfollow

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type TagFollowService struct {
	tagfollowDB *TagFollowDB
}

// Given a UserID, get all the tags that they follow
func (u *TagFollowService) GetTagFollowsByUser(ctx context.Context, input *GetTagFollowsByUserParams) (*utils.ResponseBody[GetTagFollowsByUserResponse], error) {
	userID := input.UserID
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

// Given a tag and a user, creates a tag follow
func (u *TagFollowService) CreateTagFollow(ctx context.Context, input *CreateTagFollowBody) (*utils.ResponseBody[CreateTagFollowResponse], error) {
	tagfollow := &models.TagFollow{
		TagID:  input.TagID,
		UserID: input.UserID,
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
