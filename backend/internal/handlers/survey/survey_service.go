package survey

import (
	"context"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SurveyService struct {
	surveyDB *SurveyDB
}

func NewSurveyService(db *gorm.DB) *SurveyService {
	return &SurveyService{
		surveyDB: NewSurveyDB(db),
	}
}

func (s *SurveyService) CreateSurvey(ctx context.Context, input *struct{ Body CreateSurveyRequest }) (*utils.ResponseBody[SurveyResponse], error) {
	b := input.Body
	if err := validateRatings(b.PlayerDev, b.AcademicsAthleticsPriority, b.AcademicCareerResources, b.MentalHealthPriority, b.Environment, b.Culture, b.Transparency); err != nil {
		return nil, err
	}

	survey, err := utils.HandleDBError(s.surveyDB.CreateSurvey(b))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[SurveyResponse]{
		Body: ToSurveyResponse(survey),
	}, nil
}

func (s *SurveyService) DeleteSurvey(ctx context.Context, input *DeleteSurveyRequest) (*utils.ResponseBody[SurveyResponse], error) {
	survey, err := utils.HandleDBError(s.surveyDB.GetSurveyByID(input.ID))
	if err != nil {
		return nil, err
	}

	if err := s.surveyDB.DeleteSurvey(input.ID); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[SurveyResponse]{
		Body: ToSurveyResponse(survey),
	}, nil
}

func (s *SurveyService) GetSurveysByUser(ctx context.Context, input *GetSurveysByUserParams) (*utils.ResponseBody[GetSurveysByUserResponse], error) {
	surveys, total, err := s.surveyDB.GetSurveysByUserID(input.UserID, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	responses := make([]SurveyResponse, 0, len(surveys))
	for i := range surveys {
		responses = append(responses, *ToSurveyResponse(&surveys[i]))
	}

	return &utils.ResponseBody[GetSurveysByUserResponse]{
		Body: &GetSurveysByUserResponse{
			Surveys: responses,
			Total:   int(total),
		},
	}, nil
}

func (s *SurveyService) GetAverageRatings(ctx context.Context, input *GetAverageRatingsParams) (*utils.ResponseBody[AverageRatingsResponse], error) {
	rows, err := s.surveyDB.GetAverageRatings(input.SportID, input.CollegeID)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[AverageRatingsResponse]{
		Body: &AverageRatingsResponse{
			Averages: rows,
		},
	}, nil
}

// validateRatings ensures all 7 rating fields are within the 1–5 range
func validateRatings(vals ...int32) error {
	names := []string{
		"player_dev",
		"academics_athletics_priority",
		"academic_career_resources",
		"mental_health_priority",
		"environment",
		"culture",
		"transparency",
	}
	for i, v := range vals {
		if v < 1 || v > 5 {
			return huma.Error422UnprocessableEntity(names[i] + " must be between 1 and 5")
		}
	}
	return nil
}

// ensure uuid import is used (used in types/db, kept here for clarity)
var _ = uuid.UUID{}