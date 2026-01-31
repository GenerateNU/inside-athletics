package sport

import (
	"context"
	"errors"
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

	sport, err := s.sportDB.CreateSport(input.Body.Name, input.Body.Popularity)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[SportResponse]{
		Body: ToSportResponse(sport),
	}, nil
}

func (s *SportService) GetSportByName(ctx context.Context, input *GetSportByNameParams) (*utils.ResponseBody[SportResponse], error) {
	sport, err := s.sportDB.GetSportByName(input.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("sport not found")
		}
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
	sport, err := s.sportDB.GetSportByID(input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("sport not found")
		}
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
	sport, err := s.sportDB.GetSportByID(input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("sport not found")
		}
		return nil, err
	}

	// Apply partial updates
	if input.Body.Name != nil {
		if *input.Body.Name == "" {
			return nil, huma.Error422UnprocessableEntity("name cannot be empty")
		}
		sport.Name = *input.Body.Name
	}

	if input.Body.Popularity != nil {
		if *input.Body.Popularity < 0 {
			return nil, huma.Error422UnprocessableEntity("popularity cannot be negative")
		}
		sport.Popularity = input.Body.Popularity
	}

	updatedSport, err := s.sportDB.UpdateSport(sport)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[SportResponse]{
		Body: ToSportResponse(updatedSport),
	}, nil
}

func (s *SportService) DeleteSport(ctx context.Context, input *DeleteSportRequest) (*utils.ResponseBody[map[string]string], error) {
	err := s.sportDB.DeleteSport(input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.Error404NotFound("sport not found")
		}
		return nil, err
	}

	return &utils.ResponseBody[map[string]string]{
		Body: &map[string]string{
			"message": "Sport deleted successfully",
		},
	}, nil
}
