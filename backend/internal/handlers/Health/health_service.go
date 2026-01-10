package health

import (
	"context"
	paramTypes "inside-athletics/internal/handlers/Health/health_route_params"
	"inside-athletics/internal/models"
)

type HealthService struct {
	healthDB *HealthDB
}

func (h *HealthService) CheckHealth(ctx context.Context, input *struct{}) (*models.ResponseBody[models.HealthModel], error) {
	healthModel := &models.HealthModel{Id: 1, Name: "Stinky"}
	resp := &models.ResponseBody[models.HealthModel]{Body: healthModel}
	return resp, nil
}

func (h *HealthService) GetHealthEntry(ctx context.Context, input *paramTypes.GetHealthParams) (*models.HealthModel, error) {
	id := input.Name

	healthModel, err := h.healthDB.GetFromDB(id)

	if err != nil {
		return &models.HealthModel{}, err
	}

	return healthModel, nil
}
