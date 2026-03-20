package models

import (
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID *uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	S3Key string 
	Title string
	CreatedAt time.Time
}

