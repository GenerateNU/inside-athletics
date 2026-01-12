package health

import (
	"context"
	types "inside-athletics/internal/handlers/health/types"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
)

type HealthService struct {
	healthDB *HealthDB
}

func (h *HealthService) CheckHealth(ctx context.Context, input *types.EmptyInput) (*utils.ResponseBody[types.HealthResponse], error) {
	err := h.healthDB.Ping()
	resp := &utils.ResponseBody[types.HealthResponse]{}
	if err != nil {
		return resp, huma.Error500InternalServerError("Database is unable to be reached")
	}

	resp.Body = &types.HealthResponse{Message: "Database is reachable"}

	return resp, nil
}

func (h *HealthService) Health(ctx context.Context, input *types.EmptyInput) (*utils.ResponseBody[types.HealthResponse], error) {
	return &utils.ResponseBody[types.HealthResponse]{
		Body: &types.HealthResponse{
			Message: "Welcome to Inside Athletics API Version 1.0.0",
		},
	}, nil
}
