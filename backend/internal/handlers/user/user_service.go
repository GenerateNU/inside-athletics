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
	response := &utils.ResponseBody[types.GetUserResponse]{}
	healthModel, dbResponse := h.healthDB.GetFromDB(id)

	response.Body = healthModel
	return response, handleDbErrors(dbResponse)
}
