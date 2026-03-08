package post

import (
	models "inside-athletics/internal/models"

	"github.com/google/uuid"
)

// CreatePostRequest defines the request body for creating a new post
type CreatePostRequest struct {
	AuthorId    uuid.UUID    `json:"author_id"`
	SportId     *uuid.UUID   `json:"sport_id,omitempty"`
	CollegeId   *uuid.UUID   `json:"college_id,omitempty"`
	Tags        []TagRequest `json:"tags,omitempty"`
	Title       string       `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content     string       `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	IsAnonymous bool         `json:"is_anonymous"`
}

type TagRequest struct {
	ID uuid.UUID `json:"id"`
}

// PostResponse defines the response structure for a post
type CreatePostResponse struct {
	ID          uuid.UUID    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthorID    *uuid.UUID   `json:"author" type:"uuid"`
	SportID     *uuid.UUID   `json:"sport" type:"uuid"`
	CollegeID   *uuid.UUID   `json:"college" type:"uuid"`
	Tags        []models.Tag `json:"tag" type:"tag"`
	Title       string       `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content     string       `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	IsAnonymous bool         `json:"is_anonymous"`
}

// PostResponse defines the response structure for a post
type PostResponse struct {
	ID           uuid.UUID       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Author       *models.User    `json:"author" type:"user"`
	Sport        *models.Sport   `json:"sport" type:"sport"`
	College      *models.College `json:"college" type:"college"`
	Tags         []models.Tag    `json:"tags" type:"tag"`
	Title        string          `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content      string          `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	LikeCount    int64           `json:"like_count,omitempty" example:"20000" gorm:"type:int"`
	CommentCount int64           `json:"comment_count,omitempty" example:"20" gorm:"type:int"`
	IsLiked      bool            `json:"is_liked,omitempty" example:"true" gorm:"type:bool"`
	IsAnonymous  bool            `json:"is_anonymous"`
	IsVerifiedAthlete bool       `json:"is_verified_athlete"`
}

// GetPostByIDParams defines parameters for getting a post by ID
type GetPostByIDParams struct {
	ID uuid.UUID `path:"id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the post"`
}

type GetPostByTagIDParams struct {
	TagId  uuid.UUID `path:"tag_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Sport ID to filter posts"`
	Limit  int       `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset int       `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}

type GetPostByTagIDResponse struct {
	Posts []PostResponse `json:"posts" doc:"List of posts for the sport"`
	Total int            `json:"total" example:"25" doc:"Total number of posts for this sport"`
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

type GetPostsBySportIDParams struct {
	SportId uuid.UUID `path:"sport_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Sport ID to filter posts"`
	Limit   int       `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset  int       `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}

type GetPostsBySportIDResponse struct {
	Posts []PostResponse `json:"posts" doc:"List of posts for the sport"`
	Total int            `json:"total" example:"25" doc:"Total number of posts for this sport"`
}

type GetPostsByCollegeIDParams struct {
	SportId uuid.UUID `path:"sport_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Sport ID to filter posts"`
	Limit   int       `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset  int       `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}

type GetPostsByCollegeIDResponse struct {
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
	Title       *string `json:"title,omitempty" minLength:"1" maxLength:"100"`
	Content     *string `json:"content,omitempty" minLength:"1" maxLength:"5000"`
	IsAnonymous *bool   `json:"is_anonymous,omitempty"`
}

// ToPostResponse converts a Post model to a postResponse
func ToPostResponse(post *models.Post, id uuid.UUID) *PostResponse {
	var author *models.User
	if !post.IsAnonymous || id == post.AuthorID {
		a := post.Author
		author = &a
	}
	return &PostResponse{
		ID:           post.ID,
		Author:       author,
		Sport:        post.Sport,
		College:      post.College,
		Tags:         post.Tags,
		Title:        post.Title,
		Content:      post.Content,
		IsAnonymous:  post.IsAnonymous,
		IsLiked:      post.IsLiked,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		IsVerifiedAthlete: post.Author.Verified_Athlete_Status == models.VerifiedAthleteStatusVerified,
	}
}

// ToPostResponse converts a Post model to a postResponse
func ToCreatePostResponse(post *models.Post, id uuid.UUID) *CreatePostResponse {
	var userId *uuid.UUID
	if (!post.IsAnonymous) || (id == post.AuthorID) {
		uid := post.AuthorID
		userId = &uid
	}
	return &CreatePostResponse{
		ID:                post.ID,
		AuthorID:          userId,
		SportID:           post.SportID,
		CollegeID:         post.CollegeID,
		Tags:              post.Tags,
		Title:             post.Title,
		Content:           post.Content,
		IsAnonymous:       post.IsAnonymous,
}
}

type DeletePostResponse struct {
	Message string    `json:"message" example:"College deleted successfully" doc:"Success message"`
	ID      uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of the deleted college"`
}
