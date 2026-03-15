package collegefollow

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type CollegeFollowService struct {
	collegefollowDB *CollegeFollowDB
}

// Given a UserID, get all colleges that they follow
func (u *CollegeFollowService) GetCollegeFollowsByUser(ctx context.Context, input *GetCollegeFollowsByUserParams) (*utils.ResponseBody[GetCollegeFollowsByUserResponse], error) {
	userID := input.UserID
	colleges, err := u.collegefollowDB.GetCollegeFollowsByUser(userID)
	respBody := &utils.ResponseBody[GetCollegeFollowsByUserResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetCollegeFollowsByUserResponse{
		UserID:     userID,
		CollegeIDs: *colleges,
	}

	return &utils.ResponseBody[GetCollegeFollowsByUserResponse]{
		Body: response,
	}, err
}

// Given a CollegeID, get all the users that follow this college
func (u *CollegeFollowService) GetFollowingUsersByCollege(ctx context.Context, input *GetFollowingUsersByCollegeParams) (*utils.ResponseBody[GetFollowingUsersByCollegeResponse], error) {
	collegeID := input.CollegeID
	users, err := u.collegefollowDB.GetFollowingUsersByCollege(collegeID)
	respBody := &utils.ResponseBody[GetFollowingUsersByCollegeResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetFollowingUsersByCollegeResponse{
		CollegeID: collegeID,
		UserIDs:   *users,
	}

	return &utils.ResponseBody[GetFollowingUsersByCollegeResponse]{
		Body: response,
	}, err
}

// Given a college and a user, creates a college follow if doesn't already exist
func (u *CollegeFollowService) CreateCollegeFollow(ctx context.Context, input *CreateCollegeFollowInput) (*utils.ResponseBody[CreateCollegeFollowResponse], error) {
	collegefollow := &models.CollegeFollow{
		CollegeID: input.Body.CollegeID,
		UserID:    input.Body.UserID,
	}

	createdCollegeFollow, err := u.collegefollowDB.CreateCollegeFollow(collegefollow)

	respBody := &utils.ResponseBody[CreateCollegeFollowResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &CreateCollegeFollowResponse{
		ID:        createdCollegeFollow.ID,
		UserID:    createdCollegeFollow.UserID,
		CollegeID: createdCollegeFollow.CollegeID,
	}

	return &utils.ResponseBody[CreateCollegeFollowResponse]{
		Body: response,
	}, err
}

func (u *CollegeFollowService) DeleteCollegeFollow(ctx context.Context, input *DeleteCollegeFollowParams) (*utils.ResponseBody[DeleteCollegeFollowResponse], error) {
	respBody := &utils.ResponseBody[DeleteCollegeFollowResponse]{}

	err := u.collegefollowDB.DeleteCollegeFollow(input.ID)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &DeleteCollegeFollowResponse{
		Message: "College follow was deleted successfully",
	}

	return respBody, nil
}
