package sport

import "github.com/google/uuid"

// GET Sport
type GetSportByNameParams struct {
	Name string `path:"name" maxLength:"30" example:"Joe" doc:"Name to identify test data"`
}

type GetAllSportsParams struct {
    Limit  int `query:"limit" default:"50" example:"50" doc:"Number of sports to return"`
    Offset int `query:"offset" default:"0" example:"0" doc:"Number of sports to skip"`
}

type GetAllSportsResponse struct {
    Sports []SportSummary `json:"sports" doc:"List of sports"`
    Total  int            `json:"total" example:"25" doc:"Total number of sports"`
}

type SportSummary struct {
    ID         uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the sport"`
    Name       string    `json:"name" example:"Women's Soccer" doc:"Name of the sport"`
    Popularity int32     `json:"popularity" example:"20000" doc:"Number of players"`
}

type GetSportByIDParams struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the sport"`
}

type GetSportResponse struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the sport"`
	Name string    `json:"name" example:"Women's Soccer" doc:"Name of the sport"`
}

// POST Sport
type CreateSportRequest struct {
	Name string `path:"name" maxLength:"30" example:"Joe" doc:"Name to identify test data"`
}

// PATCH Sport
type UpdateSportRequest struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the sport"`
	Name string    `json:"name" example:"Women's Soccer" doc:"Name of the sport"`
}

// DELETE Sport
type DeleteSportRequest struct {
	ID   uuid.UUID `json:"id" example:"1" doc:"ID of the sport"`
}

type DeleteSportResponse struct {
    ID   uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the deleted sport"`
    Name string    `json:"name" example:"Women's Soccer" doc:"Name of the deleted sport"`
}