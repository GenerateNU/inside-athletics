package video

import (
	"context"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoService struct {
	videoDB *VideoDB
}

// NewVideoService creates a new VideoService instance
func NewVideoService(db *gorm.DB) *VideoService {
	return &VideoService{
		videoDB: NewVideoDB(db),
	}
}

func (s *VideoService) CreateVideo(ctx context.Context, input *struct{ Body CreateVideoRequest }) (*utils.ResponseBody[VideoResponse], error) {
	video := &models.Video{
		S3Key: input.Body.S3Key,
		Title: input.Body.Title,
	}

	createdVideo, err := utils.HandleDBError(
		s.videoDB.CreateVideo(video),
	)

	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[VideoResponse]{
		Body: ToVideoResponse(createdVideo),
	}, nil
}

func (s *VideoService) GetVideo(ctx context.Context, input *GetVideoParams) (*utils.ResponseBody[VideoResponse], error) {
	video, err := s.videoDB.GetVideo(&input.ID)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[VideoResponse]{
		Body: ToVideoResponse(video),
	}, nil
}

// DeleteVideo soft deletes a video by post ID
func (s *VideoService) DeleteVideo(ctx context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*utils.ResponseBody[VideoResponse], error) {
	id := &input.ID
	err := s.videoDB.DeleteVideo(id)

	respBody := &utils.ResponseBody[VideoResponse]{}
	if err != nil {
		return respBody, err
	}

	response := &VideoResponse{
		ID: id,
	}

	return &utils.ResponseBody[VideoResponse]{
		Body: response,
	}, err
}
