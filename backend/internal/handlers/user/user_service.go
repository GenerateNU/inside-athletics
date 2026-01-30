package user

import (
	"context"
	"encoding/json"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type UserService struct {
	userDB *UserDB
}

/*
*
Huma automatically validates the params based on which type you have passed it based on the struct tags (so nice!!)
IF there is an error it will automatically send the correct response back with the right status and message about the validation errors

Here we are mapping to a GetUserResponse so that we can control what the return type is. It's important to make seperate
return types so that we can control what information we are sending back instead of just the entire model
*/
func (u *UserService) GetUser(ctx context.Context, input *GetUserParams) (*utils.ResponseBody[GetUserResponse], error) {
	id := input.ID
	user, err := u.userDB.GetUser(id)
	respBody := &utils.ResponseBody[GetUserResponse]{}

	if err != nil {
		return respBody, err
	}

	// mapping to correct response type
	// we do this so we can control what values are
	// returned by the API
	response := &GetUserResponse{
		ID:                    user.ID,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		Email:                 user.Email,
		Username:              user.Username,
		Bio:                   user.Bio,
		AccountType:           user.Account_Type,
		Sport:                 user.Sport,
		ExpectedGradYear:      user.Expected_Grad_Year,
		VerifiedAthleteStatus: user.Verified_Athlete_Status,
		College:               user.College,
		Division:              user.Division,
	}

	return &utils.ResponseBody[GetUserResponse]{
		Body: response,
	}, err
}

func (u *UserService) GetCurrentUserID(ctx context.Context, input *utils.EmptyInput) (*utils.ResponseBody[GetCurrentUserIDResponse], error) {
	respBody := &utils.ResponseBody[GetCurrentUserIDResponse]{}

	userID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return respBody, err
	}

	user, err := u.userDB.GetUser(userID)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &GetCurrentUserIDResponse{
		ID:                    user.ID,
		FirstName:             user.FirstName,
		LastName:              user.LastName,
		Email:                 user.Email,
		Username:              user.Username,
		Bio:                   user.Bio,
		AccountType:           user.Account_Type,
		Sport:                 user.Sport,
		ExpectedGradYear:      user.Expected_Grad_Year,
		VerifiedAthleteStatus: user.Verified_Athlete_Status,
		College:               user.College,
		Division:              user.Division,
	}

	return respBody, nil
}

func (u *UserService) CreateUser(ctx context.Context, input *CreateUserInput) (*utils.ResponseBody[CreateUserResponse], error) {
	respBody := &utils.ResponseBody[CreateUserResponse]{}

	currentUserID, err := u.getCurrentUserID(ctx)
	if err != nil {
		return respBody, err
	}

	sportJSON, err := marshalSport(input.Body.Sport)
	if err != nil {
		return respBody, err
	}

	user := &models.User{
		ID:                      currentUserID,
		FirstName:               input.Body.FirstName,
		LastName:                input.Body.LastName,
		Email:                   input.Body.Email,
		Username:                input.Body.Username,
		Bio:                     input.Body.Bio,
		Account_Type:            input.Body.AccountType,
		Sport:                   sportJSON,
		Expected_Grad_Year:      input.Body.ExpectedGradYear,
		Verified_Athlete_Status: input.Body.VerifiedAthleteStatus,
		College:                 input.Body.College,
		Division:                input.Body.Division,
	}

	createdUser, err := u.userDB.CreateUser(user)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &CreateUserResponse{
		ID:   createdUser.ID,
		Name: createdUser.FirstName,
	}

	return respBody, nil
}

func (u *UserService) UpdateUser(ctx context.Context, input *UpdateUserInput) (*utils.ResponseBody[UpdateUserResponse], error) {
	respBody := &utils.ResponseBody[UpdateUserResponse]{}

	user, err := u.userDB.GetUser(input.ID)
	if err != nil {
		return respBody, err
	}

	body := input.Body
	if body.FirstName != nil {
		user.FirstName = *body.FirstName
	}
	if body.LastName != nil {
		user.LastName = *body.LastName
	}
	if body.Email != nil {
		user.Email = *body.Email
	}
	if body.Username != nil {
		user.Username = *body.Username
	}
	if body.Bio != nil {
		user.Bio = body.Bio
	}
	if body.AccountType != nil {
		user.Account_Type = *body.AccountType
	}
	if body.Sport != nil {
		sportJSON, err := marshalSport(*body.Sport)
		if err != nil {
			return respBody, err
		}
		user.Sport = sportJSON
	}
	if body.ExpectedGradYear != nil {
		user.Expected_Grad_Year = *body.ExpectedGradYear
	}
	if body.VerifiedAthleteStatus != nil {
		user.Verified_Athlete_Status = *body.VerifiedAthleteStatus
	}
	if body.College != nil {
		user.College = body.College
	}
	if body.Division != nil {
		user.Division = body.Division
	}

	updatedUser, err := u.userDB.UpdateUser(user)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &UpdateUserResponse{
		ID:   updatedUser.ID,
		Name: updatedUser.FirstName,
	}

	return respBody, nil
}

func (u *UserService) DeleteUser(ctx context.Context, input *GetUserParams) (*utils.ResponseBody[DeleteUserResponse], error) {
	respBody := &utils.ResponseBody[DeleteUserResponse]{}

	err := u.userDB.DeleteUser(input.ID)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &DeleteUserResponse{
		ID: input.ID,
	}

	return respBody, nil
}

func marshalSport(sport []string) ([]byte, error) {
	if sport == nil {
		return nil, nil
	}
	sportJSON, err := json.Marshal(sport)
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid sport values", err)
	}
	return sportJSON, nil
}

func (u *UserService) getCurrentUserID(ctx context.Context) (uuid.UUID, error) {
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
