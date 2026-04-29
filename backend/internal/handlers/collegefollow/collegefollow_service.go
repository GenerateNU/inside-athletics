package collegefollow

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type CollegeFollowService struct {
	collegefollowDB *CollegeFollowDB
}

func (u *CollegeFollowService) getCurrentUserID(ctx context.Context) (uuid.UUID, error) {
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

// Given a UserID, get all colleges that they follow
func (u *CollegeFollowService) GetCollegeFollowsByUser(ctx context.Context, input *GetCollegeFollowsByUserParams) (*utils.ResponseBody[GetCollegeFollowsByUserResponse], error) {
	userID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}

	colleges, err := u.collegefollowDB.GetCollegeFollowsByUser(userID)
	if err != nil {
		return nil, err
	}

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
	userID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}

	collegefollow := &models.CollegeFollow{
		CollegeID: input.Body.CollegeID,
		UserID:    userID,
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

	userID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}

	err = u.collegefollowDB.DeleteCollegeFollow(userID, input.CollegeID)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &DeleteCollegeFollowResponse{
		Message: "College follow was deleted successfully",
	}

	return respBody, nil
}
