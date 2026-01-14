package health

import (
	"context"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
)

type HealthService struct {
	healthDB *HealthDB
}

func (h *HealthService) CheckHealth(ctx context.Context, input *utils.EmptyInput) (*utils.ResponseBody[HealthResponse], error) {
	err := h.healthDB.Ping()
	resp := &utils.ResponseBody[HealthResponse]{}
	if err != nil {
		return resp, huma.Error500InternalServerError("Database is unable to be reached")
	}

	resp.Body = &HealthResponse{Message: "Database is reachable"}

	return resp, nil
}

func (h *HealthService) Health(ctx context.Context, input *utils.EmptyInput) (*utils.ResponseBody[HealthResponse], error) {
	return &utils.ResponseBody[HealthResponse]{
		Body: &HealthResponse{
			Message: "Welcome to Inside Athletics API Version 1.0.0",
		},
	}, nil
}
