package premiumpost

import (
	"inside-athletics/internal/models"

	"github.com/google/uuid"
)

// get all posts - limit int for how many posts they want
// create
// update - title, content, and optional media_id
// delete
// get all posts by sport - limit int for how many posts they want
// get all posts by college - limit int for how many posts they want
// get all posts by tag - limit int for how many posts they want i dont think we need this

// Retrieve all posts
type GetAllPremiumPostsParams struct {
	Limit  int `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset int `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}

// Given an AuthorID, return all posts that the author has posted (with pagination)
type GetPremiumPostsByAuthorIDParams struct {
	AuthorID uuid.UUID `path:"author_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Author ID to filter posts"`
	Limit    int       `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset   int       `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}

// Given a SportID, return all posts related to the sport (with pagnination)
type GetPremiumPostsBySportIDParams struct {
	SportID uuid.UUID `path:"sport_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Sport ID to filter posts"`
	Limit   int       `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset  int       `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}

// Given a CollegeID, return all posts related to the college (with pagination)
type GetPremiumPostsByCollegeIDParams struct {
	CollegeID uuid.UUID `path:"college_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"College ID to filter posts"`
	Limit     int       `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset    int       `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}

// Given a TagID, return all posts related to the tag (with pagination)
type GetPremiumPostsByTagIDParams struct {
	TagID  uuid.UUID `path:"tag_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Tag ID to filter posts"`
	Limit  int       `query:"limit" default:"50" example:"50" doc:"Number of posts to return"`
	Offset int       `query:"offset" default:"0" example:"0" doc:"Number of posts to skip"`
}

type PremiumPostResponse struct {
	ID             uuid.UUID              `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Author         *models.User           `json:"author" type:"user"`
	Sport          *models.Sport          `json:"sport" type:"sport"`
	College        *models.College        `json:"college" type:"college"`
	Tags           []models.Tag           `json:"tags" type:"tag"`
	Title          string                 `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content        string                 `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	MediaID        *uuid.UUID             `json:"media_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	Media          *models.Media          `json:"media,omitempty"`
}

type GetAllPremiumPostsResponse struct {
	Posts []PremiumPostResponse `json:"posts" doc:"List of premium posts"`
	Total int                   `json:"total" example:"100" doc:"Total number of premium posts"`
}

type GetPremiumPostsByAuthorIDResponse struct {
	Posts []PremiumPostResponse `json:"posts" doc:"List of premium posts for the author"`
	Total int                   `json:"total" example:"25" doc:"Total number of premium posts for this author"`
}

type GetPremiumPostsBySportIDResponse struct {
	Posts []PremiumPostResponse `json:"posts" doc:"List of premium posts for the sport"`
	Total int                   `json:"total" example:"25" doc:"Total number of premium posts for this sport"`
}

type GetPremiumPostsByCollegeIDResponse struct {
	Posts []PremiumPostResponse `json:"posts" doc:"List of premium posts for the college"`
	Total int                   `json:"total" example:"25" doc:"Total number of premium posts for this college"`
}

type GetPremiumPostsByTagIDResponse struct {
	Posts []PremiumPostResponse `json:"posts" doc:"List of premium posts for the tag"`
	Total int                   `json:"total" example:"25" doc:"Total number of premium posts for this tag"`
}

type CreatePremiumPostParams struct {
	SportID        *uuid.UUID             `json:"sport_id" gorm:"type:uuid;default:null"`
	CollegeID      *uuid.UUID             `json:"college_id" gorm:"type:uuid;default:null"`
	Tags           []uuid.UUID            `json:"tag" type:"tag"`
	Title          string                 `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content        string                 `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	MediaID        *uuid.UUID             `json:"media_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type CreatePremiumPostResponse struct {
	ID             uuid.UUID              `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthorID       *uuid.UUID             `json:"author" type:"uuid"`
	SportID        *uuid.UUID             `json:"sport_id" gorm:"type:uuid;default:null"`
	CollegeID      *uuid.UUID             `json:"college_id" gorm:"type:uuid;default:null"`
	Tags           []models.Tag           `json:"tag" type:"tag"`
	Title          string                 `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content        string                 `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	MediaID        *uuid.UUID             `json:"media_id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	Media          *models.Media          `json:"media,omitempty"`
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
		MediaID:   post.MediaID,
		Media:     post.Media,
	}
}

func ToPremiumPostResponse(post *models.PremiumPost) *PremiumPostResponse {
	return &PremiumPostResponse{
		ID:      post.ID,
		Author:  &post.Author,
		Sport:   post.Sport,
		College: post.College,
		Tags:    post.Tags,
		Title:   post.Title,
		Content: post.Content,
		MediaID: post.MediaID,
		Media:   post.Media,
	}
}

type GetSearchPremiumPostParam struct {
	SearchStr string `query:"search_str" binding:"required" example:"Northeastern University" doc:"String to fuzzy search premium posts on"`
	Limit     int    `query:"limit" default:"20" example:"10" doc:"Cap on the number of posts to return"`
	Offset    int    `query:"offset" default:"0" example:"8" doc:"Number of entries to skip for pagination"`
}

type GetSearchPremiumPostResponse struct {
	Posts []PremiumPostResponse `json:"posts" doc:"List of premium post responses found for given search"`
	Count int64                 `json:"count" example:"5" doc:"Count of search results found for given search"`
}

type GetFilterPremiumPostsParams struct {
	CollegeIds string `query:"college_ids" default:"" example:"98d830a4-3ddd-441f-a8b8-12d99b597894,98d830a4-3ddd-441f-a8b8-12d99b597894" doc:"Comma seperated list of college_ids to filter by"`
	SportIds   string `query:"sport_ids" default:"" example:"98d830a4-3ddd-441f-a8b8-12d99b597894,98d830a4-3ddd-441f-a8b8-12d99b597894" doc:"Comma seperated list of sport_ids to filter by"`
	TagIds     string `query:"tag_ids" default:"" example:"98d830a4-3ddd-441f-a8b8-12d99b597894,98d830a4-3ddd-441f-a8b8-12d99b597894" doc:"Comma seperated list of tag_ids to filter by"`
	Limit      int    `query:"limit" default:"20" example:"20" doc:"Number of posts to return when filtering"`
	Offset     int    `query:"offset" default:"0" example:"8" doc:"Number of entries in the database to offset by"`
}

type GetFilterPremiumPostsResponse struct {
	Posts []PremiumPostResponse `json:"posts" doc:"List of filtered premium posts"`
	Total int                   `json:"total" example:"100" doc:"Total number of matching premium posts"`
}

type UpdatePremiumPostRequest struct {
	Title   *string     `json:"title,omitempty" minLength:"1" maxLength:"100"`
	Content *string     `json:"content,omitempty" minLength:"1" maxLength:"5000"`
	MediaID *uuid.UUID  `json:"media_id,omitempty"`
}

type DeletePremiumPostRequest struct {
	Message string    `json:"message" example:"Premium post deleted successfully" doc:"Success message"`
	ID      uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of the deleted premium post"`
}
