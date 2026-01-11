package health

import (
	"context"
	types "inside-athletics/internal/handlers/health/types"
	"github.com/danielgtaylor/huma/v2"
)

type HealthService struct {
	healthDB *HealthDB
}

func (h *HealthService) CheckHealth(ctx context.Context, input *types.EmptyInput) (*types.HealthResponse, error) {
	error := h.healthDB.Ping()
	if error != nil {
		return nil, huma.Error500InternalServerError("Database is unable to be reached")
	}
	return &types.HealthResponse{
		Message: "Database is reachable",
	}, nil
}

func (h *HealthService) Health(ctx context.Context, input *types.EmptyInput) (*types.HealthResponse, error) {
	return &types.HealthResponse{
		Message: "Welcome to Inside Athletics API Version 1.0.0",
	}, nil
}
