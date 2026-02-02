package post

import (
	models "inside-athletics/internal/models"

	"github.com/google/uuid"
)

// CreatePostRequest defines the request body for creating a new post
type CreatePostRequest struct {
	AuthorId    uuid.UUID      `json:"author_id" type:"uuid""`
	SportId     uuid.UUID      `json:"sport_id" type:"uuid"`
	Title       string         `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content     string         `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	IsAnonymous bool           `json:"is_anonymous"`
}

// PostResponse defines the response structure for a post
type PostResponse struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthorId    uuid.UUID      `json:"author_id" type:"uuid"`
	SportId     uuid.UUID      `json:"sport_id" type:"uuid"`
	Title       string         `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content     string         `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	Likes       int64          `json:"Likes,omitempty" example:"20000" gorm:"type:int"`
	IsAnonymous bool           `json:"is_anonymous"`
}


// GetPostByIDParams defines parameters for getting a post by ID
type GetPostByIDParams struct {
	ID uuid.UUID `path:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post"`
}

// GetPostByAuthorIDParams defines parameters for getting a post by author ID
type GetPostByAuthorIDParams struct {
	ID uuid.UUID `path:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post"`
}

// ToPostResponse converts a Post model to a postResponse
func ToPostResponse(post *models.Post) *PostResponse {
	return &PostResponse{
		AuthorId:   post.AuthorId,
		SportId:    post.SportId,
		Title:      post.Title,
		Content:    post.Content,
		IsAnonymous: post.IsAnonymous,
	}
}