package sportfollow

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type SportFollowService struct {
	sportfollowDB *SportFollowDB
}

// Given a UserID, get all the sports that they follow
func (u *SportFollowService) GetSportFollowsByUser(ctx context.Context, input *GetSportFollowsByUserParams) (*utils.ResponseBody[GetSportFollowsByUserResponse], error) {
	userID := input.UserID
	sports, err := u.sportfollowDB.GetSportFollowsByUser(userID)
	respBody := &utils.ResponseBody[GetSportFollowsByUserResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetSportFollowsByUserResponse{
		UserID:   userID,
		SportIDs: *sports,
	}

	return &utils.ResponseBody[GetSportFollowsByUserResponse]{
		Body: response,
	}, err
}

// Given a SportID, get all the users that follow this sport
func (u *SportFollowService) GetFollowingUsersBySport(ctx context.Context, input *GetFollowingUsersBySportParams) (*utils.ResponseBody[GetFollowingUsersBySportResponse], error) {
	sportID := input.SportID
	users, err := u.sportfollowDB.GetFollowingUsersBySport(sportID)
	respBody := &utils.ResponseBody[GetFollowingUsersBySportResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetFollowingUsersBySportResponse{
		SportID: sportID,
		UserIDs: *users,
	}

	return &utils.ResponseBody[GetFollowingUsersBySportResponse]{
		Body: response,
	}, err
}

// Given a sport and a user, creates a sport follow if doesn't already exist
func (u *SportFollowService) CreateSportFollow(ctx context.Context, input *CreateSportFollowInput) (*utils.ResponseBody[CreateSportFollowResponse], error) {
	sportfollow := &models.SportFollow{
		SportID: input.Body.SportID,
		UserID:  input.Body.UserID,
	}

	createdSportFollow, err := u.sportfollowDB.CreateSportFollow(sportfollow)

	respBody := &utils.ResponseBody[CreateSportFollowResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &CreateSportFollowResponse{
		ID:      createdSportFollow.ID,
		UserID:  createdSportFollow.UserID,
		SportID: createdSportFollow.SportID,
	}

	return &utils.ResponseBody[CreateSportFollowResponse]{
		Body: response,
	}, err
}

func (u *SportFollowService) DeleteSportFollow(ctx context.Context, input *DeleteSportFollowParams) (*utils.ResponseBody[DeleteSportFollowResponse], error) {
	respBody := &utils.ResponseBody[DeleteSportFollowResponse]{}

	err := u.sportfollowDB.DeleteSportFollow(input.ID)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &DeleteSportFollowResponse{
		Message: "Sport follow was deleted successfully",
	}

	return respBody, nil
}
