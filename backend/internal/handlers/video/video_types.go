package video

import (
	models "inside-athletics/internal/models"

	"github.com/google/uuid"
)

// CreateVideoRequest defines the request body for creating a new video
type CreateVideoRequest struct {
	ID          *uuid.UUID   `json:"id,omitempty"`
	S3Key       string       `json:"s3key" example:"example_s3_key"`
	Title       string       `json:"title" example:"Suli doing a backflip !!"`
}

// VideoResponse defines the response structure for a video
type VideoResponse struct {
	ID          *uuid.UUID        `json:"id,omitempty"`
	S3Key       string       `json:"s3key" example:"example_s3_key"`
	Title       string       `json:"title" example:"Suli doing a backflip !!"`
}

// GetVideoParams defines parameters for getting a video
type GetVideoParams struct {
	ID uuid.UUID `path:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the video"`
}

// ToCreateVideoResponse converts a Video model to a videoResponse
func ToVideoResponse(video *models.Video) *VideoResponse {
	return &VideoResponse{
		ID:   video.ID,
		S3Key: video.S3Key,
		Title: video.Title,
	}
}