package health

import (
	types "inside-athletics/internal/handlers/Health/types"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type UserService struct {
	userDB *UserDB
}

/**
Huma automatically validates the params based on which type you have passed it based on the struct tags (so nice!!)
IF there is an error it will automatically send the correct response back with the right status and message about the validation errors

Here we are mapping to a GetUserResponse so that we can control what the return type is. It's important to make seperate 
return types so that we can control what information we are sending back instead of just the entire model
*/
func (u *UserService) GetUser(c *fiber.Ctx, input *paramTypes.GetUserParams) (*utils.ResponseBody[types.GetUserResponse], error) {
	id := input.Name
	user, err := u.healthDB.GetUser(id)

	// mapping to correct response type
	// we do this so we can control what values are 
	// returned by the API
	response := types.GetUserResponse{
		ID:        user.ID,
        Name:      user.Name,
    }

	return &utils.ResponseBody[types.GetUserResponse]{
		Body: user
	}, err
}
