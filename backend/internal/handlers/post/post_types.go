package post

import (
	models "inside-athletics/internal/models"

	"github.com/google/uuid"
)

// CreatePostRequest defines the request body for creating a new post
type CreatePostRequest struct {
	AuthorId    uuid.UUID `json:"author_id"`
	SportId     uuid.UUID `json:"sport_id"`
	Title       string    `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content     string    `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	IsAnonymous bool      `json:"is_anonymous"`
}

// PostResponse defines the response structure for a post
type PostResponse struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthorId    uuid.UUID `json:"author_id" type:"uuid"`
	SportId     uuid.UUID `json:"sport_id" type:"uuid"`
	Title       string    `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content     string    `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	Likes       int64     `json:"Likes,omitempty" example:"20000" gorm:"type:int"`
	IsAnonymous bool      `json:"is_anonymous"`
}

// GetPostByIDParams defines parameters for getting a post by ID
type GetPostByIDParams struct {
	ID uuid.UUID `path:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post"`
}

type GetTagsByPostParams struct {
	PostID uuid.UUID `path:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post"`
}

type GetTagsByPostResponse struct {
	PostID uuid.UUID   `json:"post_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post"`
	TagIDs []uuid.UUID `json:"tag_ids" example:"[\"123e4567-e89b-12d3-a456-426614174000\",\"123e4567-e89b-12d3-a456-426614174001\"]" doc:"The tag ids associated with a post"`
}

// GetPostByAuthorIDParams defines parameters for getting a post by author ID
type GetPostsByAuthorIDParams struct {
	AuthorID uuid.UUID `path:"author_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the author"`
	Limit    int       `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset   int       `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}

type GetPostsByAuthorIDResponse struct {
	Posts []PostResponse `json:"posts" doc:"List of posts for the author"`
	Total int            `json:"total" example:"25" doc:"Total number of posts for this author"`
}

// ToPostResponse converts a Post model to a postResponse
func ToPostResponse(post *models.Post) *PostResponse {
	return &PostResponse{
		AuthorId:    post.AuthorId,
		SportId:     post.SportId,
		Title:       post.Title,
		Content:     post.Content,
		IsAnonymous: post.IsAnonymous,
	}
}

type GetPostsBySportIDParams struct {
	SportId uuid.UUID `path:"sport_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Sport ID to filter posts"`
	Limit   int       `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset  int       `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}

type GetPostsBySportIDResponse struct {
	Posts []PostResponse `json:"posts" doc:"List of posts for the sport"`
	Total int            `json:"total" example:"25" doc:"Total number of posts for this sport"`
}

// GetAllPostsParams defines query parameters for getting all posts
type GetAllPostsParams struct {
	Limit  int `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset int `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}

// GetAllPostsResponse defines the response for getting all posts
type GetAllPostsResponse struct {
	Posts []PostResponse `json:"posts" doc:"List of posts"`
	Total int            `json:"total" example:"100" doc:"Total number of posts"`
}

// UpdatePostRequest defines the request body for updating a post (all fields optional)
type UpdatePostRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=100" example:"Updated Title" doc:"Title of the post"`
	Content     *string `json:"content" binding:"omitempty,min=1,max=5000" example:"Updated content" doc:"Content of the post"`
	IsAnonymous *bool   `json:"is_anonymous,omitempty" required:"false" doc:"Whether the post is anonymous"`
}
