package health

import (
	types "inside-athletics/internal/handlers/Health/types"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type UserService struct {
	userDB *UserDB
}

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
