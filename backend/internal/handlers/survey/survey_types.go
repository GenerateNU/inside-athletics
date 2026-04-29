package survey

import (
	models "inside-athletics/internal/models"

	"github.com/google/uuid"
)

// CreateSurveyRequest defines the request body for submitting a survey
type CreateSurveyRequest struct {
	UserID                    uuid.UUID `json:"user_id" binding:"required" doc:"ID of the user submitting the survey"`
	CollegeID                 uuid.UUID `json:"college_id" binding:"required" doc:"ID of the college being rated"`
	SportID                   uuid.UUID `json:"sport_id" binding:"required" doc:"ID of the sport program being rated"`
	PlayerDev                 int32     `json:"player_dev" binding:"required,min=1,max=5" example:"4" doc:"Player development rating (1–5)"`
	AcademicsAthleticsPriority int32   `json:"academics_athletics_priority" binding:"required,min=1,max=5" example:"3" doc:"Academics vs athletics priority rating (1–5)"`
	AcademicCareerResources   int32     `json:"academic_career_resources" binding:"required,min=1,max=5" example:"4" doc:"Academic/career resources rating (1–5)"`
	MentalHealthPriority      int32     `json:"mental_health_priority" binding:"required,min=1,max=5" example:"3" doc:"Mental health priority rating (1–5)"`
	Environment               int32     `json:"environment" binding:"required,min=1,max=5" example:"5" doc:"Environment rating (1–5)"`
	Culture                   int32     `json:"culture" binding:"required,min=1,max=5" example:"4" doc:"Culture rating (1–5)"`
	Transparency              int32     `json:"transparency" binding:"required,min=1,max=5" example:"3" doc:"Transparency rating (1–5)"`
}

// DeleteSurveyRequest defines the path parameter for deleting a survey
type DeleteSurveyRequest struct {
	ID uuid.UUID `path:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the survey to delete"`
}

// GetSurveysByUserParams defines path + query params for listing a user's surveys
type GetSurveysByUserParams struct {
	UserID uuid.UUID `path:"user_id" binding:"required" doc:"ID of the user"`
	Limit  int       `query:"limit" default:"50" example:"50" doc:"Number of surveys to return"`
	Offset int       `query:"offset" default:"0" example:"0" doc:"Number of surveys to skip"`
}

// GetAverageRatingsParams defines optional query filters for the averages endpoint
type GetAverageRatingsParams struct {
	SportID   uuid.UUID `query:"sport_id" doc:"Filter by sport ID" required:"false"`
	CollegeID uuid.UUID `query:"college_id" doc:"Filter by college ID" required:"false"`
}

// SurveyResponse defines the response structure for a single survey
type SurveyResponse struct {
	ID                        uuid.UUID `json:"id" doc:"Survey ID"`
	UserID                    uuid.UUID `json:"user_id" doc:"User ID"`
	CollegeID                 uuid.UUID `json:"college_id" doc:"College ID"`
	SportID                   uuid.UUID `json:"sport_id" doc:"Sport ID"`
	PlayerDev                 int32     `json:"player_dev" doc:"Player development rating"`
	AcademicsAthleticsPriority int32   `json:"academics_athletics_priority" doc:"Academics vs athletics priority rating"`
	AcademicCareerResources   int32     `json:"academic_career_resources" doc:"Academic/career resources rating"`
	MentalHealthPriority      int32     `json:"mental_health_priority" doc:"Mental health priority rating"`
	Environment               int32     `json:"environment" doc:"Environment rating"`
	Culture                   int32     `json:"culture" doc:"Culture rating"`
	Transparency              int32     `json:"transparency" doc:"Transparency rating"`
}

// GetSurveysByUserResponse wraps a paginated list of the user's surveys
type GetSurveysByUserResponse struct {
	Surveys []SurveyResponse `json:"surveys" doc:"List of survey responses"`
	Total   int              `json:"total" doc:"Total number of surveys submitted by this user"`
}

// AverageRatingsRow is the raw DB scan target for the averages query
type AverageRatingsRow struct {
	SportID                    uuid.UUID `json:"sport_id"                    gorm:"column:sport_id"`
	CollegeID                  uuid.UUID `json:"college_id"                  gorm:"column:college_id"`
	PlayerDev                  float64   `json:"player_dev"                  gorm:"column:player_dev"`
	AcademicsAthleticsPriority float64   `json:"academics_athletics_priority" gorm:"column:academics_athletics_priority"`
	AcademicCareerResources    float64   `json:"academic_career_resources"   gorm:"column:academic_career_resources"`
	MentalHealthPriority       float64   `json:"mental_health_priority"      gorm:"column:mental_health_priority"`
	Environment                float64   `json:"environment"                 gorm:"column:environment"`
	Culture                    float64   `json:"culture"                     gorm:"column:culture"`
	Transparency               float64   `json:"transparency"                gorm:"column:transparency"`
	ResponseCount              int64     `json:"response_count"              gorm:"column:response_count"`
}

// AverageRatingsResponse wraps the list of grouped averages
type AverageRatingsResponse struct {
	Averages []AverageRatingsRow `json:"averages" doc:"Average ratings grouped by sport and college"`
}

// ToSurveyResponse converts a Survey model to a SurveyResponse
func ToSurveyResponse(m *models.Survey) *SurveyResponse {
	return &SurveyResponse{
		ID:                        m.ID,
		UserID:                    m.UserID,
		CollegeID:                 m.CollegeID,
		SportID:                   m.SportID,
		PlayerDev:                 m.PlayerDev,
		AcademicsAthleticsPriority: m.AcademicsAthleticsPriority,
		AcademicCareerResources:   m.AcademicCareerResources,
		MentalHealthPriority:      m.MentalHealthPriority,
		Environment:               m.Environment,
		Culture:                   m.Culture,
		Transparency:              m.Transparency,
	}
}