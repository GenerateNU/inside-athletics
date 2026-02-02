package post

import (
	"context"
	"inside-athletics/internal/models"
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

// GetPostBySportId retrieves all posts for a specific sport
func (s *PostService) GetPostBySportId(ctx context.Context, input *GetPostBySportIdParams) (*utils.ResponseBody[GetPostBySportIdResponse], error) {
	posts, total, err := s.postDB.GetPostBySportId(input.SportId, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	postResponses := make([]PostResponse, 0, len(posts))
	for i := range posts {
		postResponses = append(postResponses, *ToPostResponse(&posts[i]))
	}

	return &utils.ResponseBody[GetPostBySportIdResponse]{
		Body: &GetPostBySportIdResponse{
			Posts: postResponses,
			Total: int(total),
		},
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

	updatedPost, err := utils.HandleDBError(s.postDB.UpdatePost(post))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PostResponse]{
		Body: ToPostResponse(updatedPost),
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

// ToPostResponse converts a Post model to a PostResponse
func ToPostResponse(post *models.Post) *PostResponse {
	return &PostResponse{
		ID:          post.ID,
		AuthorId:    post.AuthorId,
		SportId:     post.SportId,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		UpVotes:     post.UpVotes,
		DownVotes:   post.DownVotes,
		IsAnonymous: post.IsAnonymous,
	}
}
