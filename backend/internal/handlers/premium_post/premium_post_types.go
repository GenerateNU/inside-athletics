package premiumpost

import (
	"inside-athletics/internal/models"

	"github.com/google/uuid"
)

// get all posts - limit int for how many posts they want
// create
// update - only update s3 key, title, and content
// delete
// get all posts by sport - limit int for how many posts they want
// get all posts by college - limit int for how many posts they want
// get all posts by tag - limit int for how many posts they want i dont think we need this

// an attachmenttype is a
type AttachmentType string

// for optional pdf, image or video
const (
	AttachmentTypePDF   AttachmentType = "pdf"
	AttachmentTypeImage AttachmentType = "image"
	AttachmentTypeVideo AttachmentType = "video"
)

type GetPremiumPostsBySportParams struct {
}

type GetPremiumPostsByCollegeParams struct {
}

type GetPremiumPostsByTagParams struct {
}

type CreatePremiumPostParams struct {
	SportID        *uuid.UUID             `json:"sport_id" gorm:"type:uuid;default:null"`
	CollegeID      *uuid.UUID             `json:"college_id" gorm:"type:uuid;default:null"`
	Tags           []uuid.UUID            `json:"tag" type:"tag"`
	Title          string                 `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content        string                 `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	AttachmentKey  *string                `json:"attachment_key,omitempty" example:"abc123.pdf" gorm:"type:text;default:null"`
	AttachmentType *models.AttachmentType `json:"attachment_type,omitempty" example:"video" gorm:"type:varchar(10);default:null" validate:"omitempty,oneof=pdf image video"`
}

type CreatePremiumPostResponse struct {
	ID             uuid.UUID              `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthorID       *uuid.UUID             `json:"author" type:"uuid"`
	SportID        *uuid.UUID             `json:"sport_id" gorm:"type:uuid;default:null"`
	CollegeID      *uuid.UUID             `json:"college_id" gorm:"type:uuid;default:null"`
	Tags           []models.Tag           `json:"tag" type:"tag"`
	Title          string                 `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content        string                 `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	AttachmentKey  *string                `json:"attachment_key,omitempty" example:"abc123.pdf" gorm:"type:text;default:null"`
	AttachmentType *models.AttachmentType `json:"attachment_type,omitempty" example:"video" gorm:"type:varchar(10);default:null" validate:"omitempty,oneof=pdf image video"`
}

// ToPremiumPostResponse converts a PremiumPost model to a premiumPostResponse
func ToCreatePremiumPostResponse(post *models.PremiumPost, id uuid.UUID) *CreatePremiumPostResponse {
	var userId *uuid.UUID
	if id == post.AuthorID {
		uid := post.AuthorID
		userId = &uid
	}
	return &CreatePremiumPostResponse{
		ID:        post.ID,
		AuthorID:  userId,
		SportID:   post.SportID,
		CollegeID: post.CollegeID,
		Tags:      post.Tags,
		Title:     post.Title,
		Content:   post.Content,
	}
}

type UpdatePremiumPostRequest struct {
}

type DeletePremiumPostRequest struct {
}
