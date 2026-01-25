package user

import (
	"context"
	"inside-athletics/internal/utils"
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
		ID:   user.ID,
		Name: user.FirstName,
	}

	return &utils.ResponseBody[GetUserResponse]{
		Body: response,
	}, err
}

func (u *UserService) GetCurrentUserID(ctx context.Context, input *utils.EmptyInput) (*utils.ResponseBody[GetCurrentUserIDResponse], error) {
	respBody := &utils.ResponseBody[GetCurrentUserIDResponse]{}

	userID, err := u.userDB.GetCurrentUserID(ctx)
	if err != nil {
		return respBody, err
	}

	respBody.Body = &GetCurrentUserIDResponse{
		ID: userID,
	}

	return respBody, nil
}
