package sport

import (
	"context"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SportService struct {
	sportDB *SportDB
}

// NewSportService creates a new SportService instance
func NewSportService(db *gorm.DB) *SportService {
	return &SportService{
		sportDB: NewSportDB(db),
	}
}

func (s *SportService) CreateSport(ctx context.Context, input *struct{ Body CreateSportRequest }) (*utils.ResponseBody[SportResponse], error) {
	// Validate business rules
	if input.Body.Name == "" {
		return nil, huma.Error422UnprocessableEntity("name cannot be empty")
	}
	if input.Body.Popularity != nil && *input.Body.Popularity < 0 {
		return nil, huma.Error422UnprocessableEntity("popularity cannot be negative")
	}

	sport, err := utils.HandleDBError(s.sportDB.CreateSport(input.Body.Name, input.Body.Popularity))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[SportResponse]{
		Body: ToSportResponse(sport),
	}, nil
}

func (s *SportService) GetSportByName(ctx context.Context, input *GetSportByNameParams) (*utils.ResponseBody[SportResponse], error) {
	sport, err := utils.HandleDBError(s.sportDB.GetSportByName(input.Name))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[SportResponse]{
		Body: ToSportResponse(sport),
	}, nil
}

func (s *SportService) GetAllSports(ctx context.Context, input *GetAllSportsParams) (*utils.ResponseBody[GetAllSportsResponse], error) {
	sports, total, err := s.sportDB.GetAllSports(input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	sportResponses := make([]SportResponse, 0, len(sports))
	for i := range sports {
		sportResponses = append(sportResponses, *ToSportResponse(&sports[i]))
	}

	return &utils.ResponseBody[GetAllSportsResponse]{
		Body: &GetAllSportsResponse{
			Sports: sportResponses,
			Total:  int(total),
		},
	}, nil
}

func (s *SportService) GetSportByID(ctx context.Context, input *GetSportByIDParams) (*utils.ResponseBody[SportResponse], error) {
	sport, err := utils.HandleDBError(s.sportDB.GetSportByID(input.ID))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[SportResponse]{
		Body: ToSportResponse(sport),
	}, nil
}

func (s *SportService) UpdateSport(ctx context.Context, input *struct {
	ID   uuid.UUID `path:"id"`
	Body UpdateSportRequest
}) (*utils.ResponseBody[SportResponse], error) {
	// Validate business rules for partial updates
	if input.Body.Name != nil && *input.Body.Name == "" {
		return nil, huma.Error422UnprocessableEntity("name cannot be empty")
	}
	if input.Body.Popularity != nil && *input.Body.Popularity < 0 {
		return nil, huma.Error422UnprocessableEntity("popularity cannot be negative")
	}

	updatedSport, err := utils.HandleDBError(s.sportDB.UpdateSport(input.ID, input.Body))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[SportResponse]{
		Body: ToSportResponse(updatedSport),
	}, nil
}

func (s *SportService) DeleteSport(ctx context.Context, input *DeleteSportRequest) (*utils.ResponseBody[SportResponse], error) {
	// First get the sport to return it
	sport, err := utils.HandleDBError(s.sportDB.GetSportByID(input.ID))
	if err != nil {
		return nil, err
	}

	// Actually delete it
	err = s.sportDB.DeleteSport(input.ID)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[SportResponse]{
		Body: ToSportResponse(sport),
	}, nil
}
