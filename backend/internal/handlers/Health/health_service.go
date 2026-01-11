package health

import (
	types "inside-athletics/internal/handlers/Health/types"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type HealthService struct {
	healthDB *HealthDB
}

func (h *HealthService) CheckHealth(c *fiber.Ctx) (*utils.ResponseBody[types.HealthcheckResponse], error) {
	response := &utils.ResponseBody[types.HealthcheckResponse]{}
	error := h.healthDB.Ping(id)
	if error != nil {
		return nil, huma.Error500InternalServerError("Database is unable to be reached") 
	}
	response.body = "Database reached successfully"
	return response, nil
}

func (h *HealthService) Health(c *fiber.Ctx) (*utils.ResponseBody[types.HealthcheckResponse], error) {
	response := &utils.ResponseBody[types.HealthcheckResponse]{}
	response.body = "Welcome to Inside Athletics API Version 1.0.0"
	return response, nil
}
