package media

import (
	models "inside-athletics/internal/models"

	"github.com/google/uuid"
)

// CreateMediaRequest defines the request body for creating a new media
type CreateMediaRequest struct {
	ID    *uuid.UUID `json:"id,omitempty"`
	S3Key string     `json:"s3key" example:"example_s3_key"`
	Title string     `json:"title" example:"Suli doing a backflip !!"`
}

// MediaResponse defines the response structure for a media
type MediaResponse struct {
	ID    *uuid.UUID `json:"id,omitempty"`
	S3Key string     `json:"s3key" example:"example_s3_key"`
	Title string     `json:"title" example:"Suli doing a backflip !!"`
}

// GetMediaParams defines parameters for getting a media
type GetMediaParams struct {
	ID uuid.UUID `path:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the media"`
}

// ToCreateMediaResponse converts a Media model to a mediaResponse
func ToMediaResponse(media *models.Media) *MediaResponse {
	return &MediaResponse{
		ID:    media.ID,
		S3Key: media.S3Key,
		Title: media.Title,
	}
}
