package premiumpost

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

// --- reminders
// add permission to only allow moderators and admins for access to get, create, update, delete
// add permission so that only paid users can view the posts

type PremiumPostService struct {
	premiumPostDB *PremiumPostDB
}

func NewPremiumPostService(db *gorm.DB) *PremiumPostService {
	return &PremiumPostService{
		premiumPostDB: NewPremiumPostDB(db),
	}
}

// CreatePremiumPost creates a new post in the db
func (s *PremiumPostService) CreatePremiumPost(ctx context.Context, input *struct{ Body CreatePremiumPostParams }) (*utils.ResponseBody[CreatePremiumPostResponse], error) {
	id, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if len(input.Body.Tags) == 0 && input.Body.SportID == nil && input.Body.CollegeID == nil {
		return nil, huma.Error400BadRequest("Need to have at least a single tag on a post")
	}
	premiumPost := &models.PremiumPost{
		AuthorID:       id,
		SportID:        input.Body.SportID,
		CollegeID:      input.Body.CollegeID,
		Title:          input.Body.Title,
		Content:        input.Body.Content,
		AttachmentKey:  input.Body.AttachmentKey,
		AttachmentType: input.Body.AttachmentType,
	}

	createdPost, err := utils.HandleDBError(
		s.premiumPostDB.CreatePremiumPost(premiumPost),
	)

	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[CreatePremiumPostResponse]{
		Body: ToCreatePremiumPostResponse(createdPost, id),
	}, nil
}

// GetAllPremiumPosts returns all premium posts
func (s *PremiumPostService) GetAllPremiumPosts(ctx context.Context, input *GetAllPremiumPostsParams) (*utils.ResponseBody[GetAllPremiumPostsResponse], error) {
	posts, total, err := s.premiumPostDB.GetAllPremiumPosts(input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PremiumPostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPremiumPostResponse(&posts[i]))
	}

	return &utils.ResponseBody[GetAllPremiumPostsResponse]{
		Body: &GetAllPremiumPostsResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}

// GetPremiumPostsByAuthorID returns all premium posts related to a given author
func (s *PremiumPostService) GetPremiumPostsByAuthorID(ctx context.Context, input *GetPremiumPostsByAuthorIDParams) (*utils.ResponseBody[GetPremiumPostsByAuthorIDResponse], error) {
	posts, total, err := s.premiumPostDB.GetPremiumPostsByAuthorID(input.Limit, input.Offset, input.AuthorID)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PremiumPostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPremiumPostResponse(&posts[i]))
	}

	return &utils.ResponseBody[GetPremiumPostsByAuthorIDResponse]{
		Body: &GetPremiumPostsByAuthorIDResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}

// GetPremiumPostsBySportID returns all premium posts related to a given sport
func (s *PremiumPostService) GetPremiumPostsBySportID(ctx context.Context, input *GetPremiumPostsBySportIDParams) (*utils.ResponseBody[GetPremiumPostsBySportIDResponse], error) {
	posts, total, err := s.premiumPostDB.GetPremiumPostsBySportID(input.Limit, input.Offset, input.SportID)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PremiumPostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPremiumPostResponse(&posts[i]))
	}

	return &utils.ResponseBody[GetPremiumPostsBySportIDResponse]{
		Body: &GetPremiumPostsBySportIDResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}

// GetPremiumPostsByCollegeID returns all premium posts related to a given college
func (s *PremiumPostService) GetPremiumPostsByCollegeID(ctx context.Context, input *GetPremiumPostsByCollegeIDParams) (*utils.ResponseBody[GetPremiumPostsByCollegeIDResponse], error) {
	posts, total, err := s.premiumPostDB.GetPremiumPostsByCollegeID(input.Limit, input.Offset, input.CollegeID)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PremiumPostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPremiumPostResponse(&posts[i]))
	}

	return &utils.ResponseBody[GetPremiumPostsByCollegeIDResponse]{
		Body: &GetPremiumPostsByCollegeIDResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}

// GetPremiumPostsByTagID returns all premium posts related to a given tag
func (s *PremiumPostService) GetPremiumPostsByTagID(ctx context.Context, input *GetPremiumPostsByTagIDParams) (*utils.ResponseBody[GetPremiumPostsByTagIDResponse], error) {
	posts, total, err := s.premiumPostDB.GetPremiumPostsByTagID(input.Limit, input.Offset, input.TagID)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PremiumPostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPremiumPostResponse(&posts[i]))
	}

	return &utils.ResponseBody[GetPremiumPostsByTagIDResponse]{
		Body: &GetPremiumPostsByTagIDResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}
