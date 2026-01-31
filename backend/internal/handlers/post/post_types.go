package post

import (
	"time"

	"github.com/google/uuid"
)

type GetPostBySportIdParams struct {
	ID uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID to identify the post"`
}

type PostResponse struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthorId    uuid.UUID `json:"author_id" type:"uuid" default:"gen_random_uuid()"`
	SportId     uuid.UUID `json:"sport_id" type:"uuid" default:"gen_random_uuid()"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content     string    `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	UpVotes     int64     `json:"numUpVotes,omitempty" example:"20000" gorm:"type:int"`
	DownVotes   int64     `json:"numDownVotes,omitempty" example:"20000" gorm:"type:int"`
	IsAnonymous bool      `json:"isAnanymous"`
}
