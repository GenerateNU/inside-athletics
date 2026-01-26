package sport

import (
	models "inside-athletics/internal/models"

	"github.com/google/uuid"
)

// GetSportByNameParams defines parameters for getting a sport by name
type GetSportByNameParams struct {
	Name string `path:"name" maxLength:"100" example:"Women's Soccer" doc:"Name of the sport to retrieve"`
}

// CreateSportRequest defines the request body for creating a new sport
type CreateSportRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=100" example:"Women's Soccer"`
	Popularity *int32 `json:"popularity" binding:"omitempty,gte=0" example:"20000"`
}

// GetAllSportsParams defines query parameters for getting all sports
type GetAllSportsParams struct {
	Limit  int `query:"limit" default:"50" example:"50" doc:"Number of sports to return"`
	Offset int `query:"offset" default:"0" example:"0" doc:"Number of sports to skip"`
}

// GetAllSportsResponse defines the response for getting all sports
type GetAllSportsResponse struct {
	Sports []SportResponse `json:"sports" doc:"List of sports"`
	Total  int             `json:"total" example:"25" doc:"Total number of sports"`
}

// GetSportByIDParams defines parameters for getting a sport by ID
type GetSportByIDParams struct {
	ID uuid.UUID `path:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the sport"`
}

// SportResponse defines the response structure for a sport
type SportResponse struct {
	ID         uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the sport"`
	Name       string    `json:"name" example:"Women's Soccer" doc:"Name of the sport"`
	Popularity *int32    `json:"popularity,omitempty" example:"20000" doc:"Number of players"`
}

// UpdateSportRequest defines the request body for updating a sport
type UpdateSportRequest struct {
	Name       *string `json:"name" binding:"omitempty,min=1,max=100" example:"Women's Soccer" doc:"Name of the sport"`
	Popularity *int32  `json:"popularity" binding:"omitempty,gte=0" example:"20000" doc:"Number of players"`
}

// DeleteSportRequest defines the request for deleting a sport
type DeleteSportRequest struct {
	ID uuid.UUID `path:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the sport"`
}

// ToSportResponse converts a Sport model to a SportResponse
func ToSportResponse(sport *models.Sport) *SportResponse {
	return &SportResponse{
		ID:         sport.ID,
		Name:       sport.Name,
		Popularity: sport.Popularity,
	}
}
