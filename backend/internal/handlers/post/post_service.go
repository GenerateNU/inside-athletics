package post

import (
	"context"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostService struct {
	postDB *PostDB
}

// NewPostService creates a new PostService instance
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{
		postDB: NewPostDB(db),
	}
}

func (s *PostService) CreatePost(ctx context.Context, input *struct{ Body CreatePostRequest }) (*utils.ResponseBody[PostResponse], error) {
	post, err := utils.HandleDBError(
		s.postDB.CreatePost(
			input.Body.AuthorId,
			input.Body.SportId,
			input.Body.Title,
			input.Body.Content,
			input.Body.IsAnonymous,
		),
	)

	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PostResponse]{
		Body: ToPostResponse(post),
	}, nil
}

// GetAllPosts retrieves all posts with pagination
func (s *PostService) GetAllPosts(ctx context.Context, input *GetAllPostsParams) (*utils.ResponseBody[GetAllPostsResponse], error) {
	posts, total, err := s.postDB.GetAllPosts(input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPostResponse(&posts[i]))
	}

	return &utils.ResponseBody[GetAllPostsResponse]{
		Body: &GetAllPostsResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}

// UpdatePost updates an existing post with partial updates
func (s *PostService) UpdatePost(ctx context.Context, input *struct {
	ID   uuid.UUID `path:"id"`
	Body UpdatePostRequest
}) (*utils.ResponseBody[PostResponse], error) {
	post, err := utils.HandleDBError(s.postDB.GetPostByID(input.ID))
	if err != nil {
		return nil, err
	}

	// Apply partial updates
	if input.Body.Title != nil {
		post.Title = *input.Body.Title
	}

	if input.Body.Content != nil {
		post.Content = *input.Body.Content
	}

	if input.Body.IsAnonymous != nil {
		post.IsAnonymous = *input.Body.IsAnonymous
	}

	updatedPost, err := s.postDB.UpdatePost(post)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PostResponse]{
		Body: ToPostResponse(updatedPost),
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

func (s *PostService) GetPostBySportID(ctx context.Context, input *GetPostsBySportIDParams) (*utils.ResponseBody[GetPostsBySportIDResponse], error) {
	posts, total, err := s.postDB.GetPostsBySportID(input.Limit, input.Offset, input.SportId)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPostResponse(&posts[i]))
	}

	return &utils.ResponseBody[GetPostsBySportIDResponse]{
		Body: &GetPostsBySportIDResponse{
			Posts: postResponses,
			Total: int(total),
		},
	}, nil
}

func (s *PostService) GetPostByAuthorID(ctx context.Context, input *GetPostsByAuthorIDParams) (*utils.ResponseBody[GetPostsByAuthorIDResponse], error) {
	posts, total, err := s.postDB.GetPostsByAuthorID(input.Limit, input.Offset, input.AuthorID)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPostResponse(&posts[i]))
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
}) (*struct{}, error) {
	if err := s.postDB.DeletePost(input.ID); err != nil {
		return nil, err
	}
	return &struct{}{}, nil
}
