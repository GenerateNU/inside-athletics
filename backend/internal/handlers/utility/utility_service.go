package utility

import (
	"context"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type UtilityService struct {
	utilityDB *UtilityDB
}

func (s *UtilityService) GetAccessCheck(ctx context.Context, input *utils.EmptyInput) (*utils.ResponseBody[AccessCheckResponse], error) {
	rawID := ctx.Value("user_id")
	if rawID == nil {
		return nil, huma.Error401Unauthorized("User not authenticated")
	}

	userID, err := uuid.Parse(rawID.(string))
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid user ID", err)
	}

	hasPremium, err := s.utilityDB.UserHasPremium(userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to check premium access")
	}

	isAdmin, err := s.utilityDB.UserIsAdmin(userID)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to check admin access")
	}

	return &utils.ResponseBody[AccessCheckResponse]{
		Body: &AccessCheckResponse{
			HasPremium: hasPremium,
			IsAdmin:    isAdmin,
		},
	}, nil
}
