package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID        *uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	S3Key     string     `json:"s3key" example:"1" doc:"The S3 key of the media object" gorm:"type:varchar(200);not null"`
	Title     string     `json:"title" example:"Video of Athlete" doc:"The title of the media" gorm:"type:varchar(200);not null"`
	MediaType MediaType  `json:"media_type" example:"jpeg" doc:"The media type of the object" gorm:"type:varchar(200);not null"`
	CreatedAt time.Time  `json:"created_at"`
}

type MediaType string

const (
	JPEG MediaType = "jpeg"
	PNG  MediaType = "png"
	MP4  MediaType = "mp4"
	MOV  MediaType = "mov"
	WEBM MediaType = "webm"
	PDF  MediaType = "pdf"
)
