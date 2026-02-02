package post

import (
	"context"
	"inside-athletics/internal/utils"
	"github.com/danielgtaylor/huma/v2"
	"gitlab.com/golang-utils/isnil"
	"gorm.io/gorm"
)

type PostService struct {
	postDB *postDB
}

// NewPostService creates a new PostService instance
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{
		PostDB: NewPostDB(db),
	}
}

func (s *PostService) CreatePost(ctx context.Context, input *struct{ Body CreatePostRequest }) (*utils.ResponseBody[PostResponse], error) {
	// Validate business rules
	if isnil.IsNil(input.Body.AuthorId) {
		return nil, huma.Error422UnprocessableEntity("Author ID cannot be null")
	}
	if isnil.IsNil(input.Body.SportId) {
		return nil, huma.Error422UnprocessableEntity("Sport ID cannot be null")
	}
	if input.Body.Title == "" {
		return nil, huma.Error422UnprocessableEntity("Title cannot be empty")
	}
	if input.Body.Content == "" {
		return nil, huma.Error422UnprocessableEntity("Content cannot be empty")
	}
	if input.Body.isAnonymous == nil {
		return nil, huma.Error422UnprocessableEntity("isAnonymous cannot be null")
	}

	post, err := utils.HandleDBError(s.postDB.CreatePost(input.Body.AuthorId, input.Body.SportId, input.Body.Title, input.Body.Content, input.Body.isAnonymous))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PostResponse]{
		Body: ToPostResponse(post),
	}, nil
}

func (s *PostService) GetPostByID(ctx context.Context, input *GetPostByIDParams) (*utils.ResponseBody[PostResponse], error) {
	post, err := utils.HandleDBError(s.postDB.GetPostByID((input.ID)))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PostResponse]{
		Body: ToPostResponse(post),
	}, nil
}

func (s *PostService) GetPostByAuthorID(ctx context.Context, input *GetPostByAuthorIDParams) (*utils.ResponseBody[PostResponse], error) {
	post, err := utils.HandleDBError(s.postDB.GetPostByAuthorID((input.ID)))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PostResponse]{
		Body: ToPostResponse(post),
	}, nil
}