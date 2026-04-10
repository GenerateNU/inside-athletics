package post

import (
	"context"
	"fmt"
	"inside-athletics/internal/handlers/tagpost"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"regexp"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostService struct {
	postDB    *PostDB
	tagPostDB *tagpost.TagPostDB
}

// NewPostService creates a new PostService instance
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{
		postDB:    NewPostDB(db),
		tagPostDB: tagpost.NewTagPostDB(db),
	}
}

func (s *PostService) CreatePost(ctx context.Context, input *struct{ Body CreatePostRequest }) (*utils.ResponseBody[CreatePostResponse], error) {
	id, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if len(input.Body.Tags) == 0 && input.Body.SportId == nil && input.Body.CollegeId == nil {
		return nil, huma.Error400BadRequest("Need to have at least a single tag on a post")
	}
	post := &models.Post{
		AuthorID:    id,
		SportID:     input.Body.SportId,
		CollegeID:   input.Body.CollegeId,
		Title:       input.Body.Title,
		Content:     input.Body.Content,
		IsAnonymous: input.Body.IsAnonymous,
	}

	createdPost, err := utils.HandleDBError(
		s.postDB.CreatePost(post, input.Body.Tags),
	)

	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[CreatePostResponse]{
		Body: ToCreatePostResponse(createdPost, id),
	}, nil
}

// GetAllPosts retrieves all posts with pagination
func (s *PostService) GetAllPosts(ctx context.Context, input *GetAllPostsParams) (*utils.ResponseBody[GetAllPostsResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	posts, total, err := s.postDB.GetAllPosts(input.Limit, input.Offset, userID)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPostResponse(&posts[i], userID))
	}

	return &utils.ResponseBody[GetAllPostsResponse]{
		Body: &GetAllPostsResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}

func (s *PostService) GetPopularPosts(ctx context.Context, input *GetPopularPostsParams) (*utils.ResponseBody[GetPopularPostsResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}

	posts, total, err := s.postDB.GetPopularPosts(input.Limit, input.Offset, input.WindowHours, userID)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPostResponse(&posts[i], userID))
	}

	return &utils.ResponseBody[GetPopularPostsResponse]{
		Body: &GetPopularPostsResponse{
			Posts: postResponses,
			Total: total,
		},
	}, nil
}

// UpdatePost updates an existing post with partial updates
func (s *PostService) UpdatePost(ctx context.Context, input *struct {
	ID   uuid.UUID `path:"id"`
	Body UpdatePostRequest
}) (*utils.ResponseBody[PostResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	updatedPost, err := s.postDB.UpdatePost(input.ID, input.Body, userID)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PostResponse]{
		Body: ToPostResponse(updatedPost, userID),
	}, nil
}

func (s *PostService) GetPostByID(ctx context.Context, input *GetPostByIDParams) (*utils.ResponseBody[PostResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	post, err := s.postDB.GetPostByID((input.ID), userID)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PostResponse]{
		Body: ToPostResponse(post, userID),
	}, nil
}

func (s *PostService) GetPostBySportID(ctx context.Context, input *GetPostsBySportIDParams) (*utils.ResponseBody[GetPostsBySportIDResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	posts, total, err := s.postDB.GetPostsBySportID(input.Limit, input.Offset, input.SportId, userID)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPostResponse(&posts[i], userID))
	}

	return &utils.ResponseBody[GetPostsBySportIDResponse]{
		Body: &GetPostsBySportIDResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}

func (s *PostService) GetPostByAuthorID(ctx context.Context, input *GetPostsByAuthorIDParams) (*utils.ResponseBody[GetPostsByAuthorIDResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	posts, total, err := s.postDB.GetPostsByAuthorID(input.Limit, input.Offset, input.AuthorID, userID)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPostResponse(&posts[i], userID))
	}

	return &utils.ResponseBody[GetPostsByAuthorIDResponse]{
		Body: &GetPostsByAuthorIDResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}

// DeletePost soft deletes a post by ID
func (s *PostService) DeletePost(ctx context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*utils.ResponseBody[DeletePostResponse], error) {
	id := input.ID
	err := s.postDB.DeletePost(id)

	respBody := &utils.ResponseBody[DeletePostResponse]{}
	if err != nil {
		return respBody, err
	}

	response := &DeletePostResponse{
		Message: fmt.Sprintf("Post %s deleted successfully", id.String()),
		ID:      id,
	}

	return &utils.ResponseBody[DeletePostResponse]{
		Body: response,
	}, err
}

func (s *PostService) FuzzySearchForPost(ctx context.Context, input *GetSearchParam) (*utils.ResponseBody[GetSearchResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return utils.HandleDBError(&utils.ResponseBody[GetSearchResponse]{}, err)
	}
	searchStr := input.SearchStr
	posts, total, err := s.postDB.FuzzySearchForPost(userID, searchStr, input.Limit, input.Offset)
	if err != nil {
		return utils.HandleDBError(&utils.ResponseBody[GetSearchResponse]{}, err)
	}
	postResponses := make([]PostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPostResponse(&posts[i], userID))
	}

	return &utils.ResponseBody[GetSearchResponse]{
		Body: &GetSearchResponse{
			Posts: postResponses,
			Count: total,
		},
	}, nil
}

func (s *PostService) FilterPosts(ctx context.Context, input *GetFilterPostsParams) (*utils.ResponseBody[GetAllPostsResponse], error) {
	userID, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	uuidPattern := `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	fullPattern := fmt.Sprintf(`^%s(?:,%s)*$`, uuidPattern, uuidPattern)
	re := regexp.MustCompile(fullPattern)

	// input validation
	if input.CollegeIds != "" && !re.MatchString(input.CollegeIds) {
		return nil, huma.Error400BadRequest("Expected comma seperated list of uuids with no spaces for college input uuid,uuid")
	}
	if input.SportIds != "" && !re.MatchString(input.SportIds) {
		return nil, huma.Error400BadRequest("Expected comma seperated list of uuids with no spaces for sport input uuid,uuid")
	}
	if input.TagIds != "" && !re.MatchString(input.TagIds) {
		return nil, huma.Error400BadRequest("Expected comma seperated list of uuids with no spaces for tag input uuid,uuid")
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

	posts, total, err := s.postDB.FilterPosts(userID, collegeIds, sportIds, tagIds, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	postResponses := utils.MapList(posts, func(p models.Post) PostResponse {
		return *ToPostResponse(&p, p.ID)
	})
	return &utils.ResponseBody[GetAllPostsResponse]{
		Body: &GetAllPostsResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}
