package sportfollow

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type SportFollowService struct {
	sportfollowDB *SportFollowDB
}

func (u *SportFollowService) getCurrentUserID(ctx context.Context) (uuid.UUID, error) {
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

// Given a UserID, get all the sports that they follow
func (u *SportFollowService) GetSportFollowsByUser(ctx context.Context, input *GetSportFollowsByUserParams) (*utils.ResponseBody[GetSportFollowsByUserResponse], error) {
	userID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
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
	userID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	
	sportfollow := &models.SportFollow{
		SportID: input.Body.SportID,
		UserID:    userID,
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
