package premiumpost

import (
	"context"
	"fmt"
	"inside-athletics/internal/models"
	"inside-athletics/internal/s3"
	"inside-athletics/internal/utils"
	"regexp"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- reminders
// add permission to only allow moderators and admins for access to get, create, update, delete
// add permission so that only paid users can view the posts

type PremiumPostService struct {
	premiumPostDB *PremiumPostDB
	s3            *s3.Service
}

func NewPremiumPostService(db *gorm.DB, s3Svc *s3.Service) *PremiumPostService {
	return &PremiumPostService{
		premiumPostDB: NewPremiumPostDB(db),
		s3:            s3Svc,
	}
}

// resolveMediaKey replaces post.Media.S3Key and post.Author.ProfilePicture with presigned download URLs in-place.
func (s *PremiumPostService) resolveMediaKey(ctx context.Context, post *models.PremiumPost) {
	if post.Media != nil {
		if url := s3.ResolveKey(ctx, s.s3, post.Media.S3Key); url != "" {
			post.Media.S3Key = url
		}
	}
	if post.Author.ProfilePicture != "" {
		if url := s3.ResolveKey(ctx, s.s3, post.Author.ProfilePicture); url != "" {
			post.Author.ProfilePicture = url
		}
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
		AuthorID:  id,
		SportID:   input.Body.SportID,
		CollegeID: input.Body.CollegeID,
		Title:     input.Body.Title,
		Content:   input.Body.Content,
		MediaID:   input.Body.MediaID,
	}

	createdPost, err := utils.HandleDBError(
		s.premiumPostDB.CreatePremiumPost(premiumPost),
	)

	if err != nil {
		return nil, err
	}

	s.resolveMediaKey(ctx, createdPost)

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
		s.resolveMediaKey(ctx, &posts[i])
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
		s.resolveMediaKey(ctx, &posts[i])
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
		s.resolveMediaKey(ctx, &posts[i])
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
		s.resolveMediaKey(ctx, &posts[i])
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
		s.resolveMediaKey(ctx, &posts[i])
		postResponses = append(postResponses, *ToPremiumPostResponse(&posts[i]))
	}

	return &utils.ResponseBody[GetPremiumPostsByTagIDResponse]{
		Body: &GetPremiumPostsByTagIDResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}

// FuzzySearchForPremiumPost searches premium posts by title using fuzzy matching
func (s *PremiumPostService) FuzzySearchForPremiumPost(ctx context.Context, input *GetSearchPremiumPostParam) (*utils.ResponseBody[GetSearchPremiumPostResponse], error) {
	posts, total, err := s.premiumPostDB.FuzzySearchForPremiumPost(input.SearchStr, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PremiumPostResponse, 0, len(posts))
	for i := range posts {
		s.resolveMediaKey(ctx, &posts[i])
		postResponses = append(postResponses, *ToPremiumPostResponse(&posts[i]))
	}

	return &utils.ResponseBody[GetSearchPremiumPostResponse]{
		Body: &GetSearchPremiumPostResponse{
			Posts: postResponses,
			Count: total,
		},
	}, nil
}

// FilterPremiumPosts filters premium posts by college, sport, and tag IDs
func (s *PremiumPostService) FilterPremiumPosts(ctx context.Context, input *GetFilterPremiumPostsParams) (*utils.ResponseBody[GetFilterPremiumPostsResponse], error) {
	uuidPattern := `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	fullPattern := fmt.Sprintf(`^%s(?:,%s)*$`, uuidPattern, uuidPattern)
	re := regexp.MustCompile(fullPattern)

	if input.CollegeIds != "" && !re.MatchString(input.CollegeIds) {
		return nil, huma.Error400BadRequest("Expected comma separated list of uuids with no spaces for college input uuid,uuid")
	}
	if input.SportIds != "" && !re.MatchString(input.SportIds) {
		return nil, huma.Error400BadRequest("Expected comma separated list of uuids with no spaces for sport input uuid,uuid")
	}
	if input.TagIds != "" && !re.MatchString(input.TagIds) {
		return nil, huma.Error400BadRequest("Expected comma separated list of uuids with no spaces for tag input uuid,uuid")
	}

	mapUUID := func(id string) uuid.UUID {
		parsedId, _ := uuid.Parse(id)
		return parsedId
	}

	collegeIds := []uuid.UUID{}
	sportIds := []uuid.UUID{}
	tagIds := []uuid.UUID{}
	if input.CollegeIds != "" {
		collegeIds = utils.MapList(strings.Split(input.CollegeIds, ","), mapUUID)
	}
	if input.SportIds != "" {
		sportIds = utils.MapList(strings.Split(input.SportIds, ","), mapUUID)
	}
	if input.TagIds != "" {
		tagIds = utils.MapList(strings.Split(input.TagIds, ","), mapUUID)
	}

	posts, total, err := s.premiumPostDB.FilterPremiumPosts(collegeIds, sportIds, tagIds, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PremiumPostResponse, 0, len(posts))
	for i := range posts {
		s.resolveMediaKey(ctx, &posts[i])
		postResponses = append(postResponses, *ToPremiumPostResponse(&posts[i]))
	}

	return &utils.ResponseBody[GetFilterPremiumPostsResponse]{
		Body: &GetFilterPremiumPostsResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}

// UpdatePremiumPost updates an existing premium post
func (s *PremiumPostService) UpdatePremiumPost(ctx context.Context, input *struct {
	ID   uuid.UUID `path:"id"`
	Body UpdatePremiumPostRequest
}) (*utils.ResponseBody[PremiumPostResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}

	updatedPost, err := s.premiumPostDB.UpdatePremiumPost(input.ID, input.Body, userID)
	if err != nil {
		return nil, err
	}

	s.resolveMediaKey(ctx, updatedPost)

	return &utils.ResponseBody[PremiumPostResponse]{
		Body: ToPremiumPostResponse(updatedPost),
	}, nil
}

// DeletePremiumPost soft deletes a premium post by ID
func (s *PremiumPostService) DeletePremiumPost(ctx context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*utils.ResponseBody[DeletePremiumPostRequest], error) {
	id := input.ID
	err := s.premiumPostDB.DeletePremiumPost(id)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[DeletePremiumPostRequest]{
		Body: &DeletePremiumPostRequest{
			Message: fmt.Sprintf("Premium post %s deleted successfully", id.String()),
			ID:      id,
		},
	}, nil
}
