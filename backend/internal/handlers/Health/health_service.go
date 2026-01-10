package health

import (
	"context"
	paramTypes "inside-athletics/internal/handlers/Health/health_route_params"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type HealthService struct {
	healthDB *HealthDB
}

func (h *HealthService) CheckHealth(ctx context.Context, input *struct{}) (*utils.ResponseBody[models.HealthModel], error) {
	healthModel := &models.HealthModel{Id: 1, Name: "YIPPEEE SO HEALTHY"}
	resp := &utils.ResponseBody[models.HealthModel]{Body: healthModel}
	return resp, nil
}

func (h *HealthService) GetHealthEntry(ctx context.Context, input *paramTypes.GetHealthParams) (*utils.ResponseBody[models.HealthModel], error) {
	id := input.Name
	response := &utils.ResponseBody[models.HealthModel]{}

	healthModel, err := h.healthDB.GetFromDB(id)

	if err != nil {
		return response, err
	}

	response.Body = healthModel
	return response, nil
}
