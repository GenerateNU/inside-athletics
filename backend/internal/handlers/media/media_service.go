package media

import (
	"context"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaService struct {
	mediaDB *MediaDB
}

// NewMediaService creates a new MediaService instance
func NewMediaService(db *gorm.DB) *MediaService {
	return &MediaService{
		mediaDB: NewMediaDB(db),
	}
}

func (s *MediaService) CreateMedia(ctx context.Context, input *struct{ Body CreateMediaRequest }) (*utils.ResponseBody[MediaResponse], error) {
	media := &models.Media{
		S3Key:     input.Body.S3Key,
		Title:     input.Body.Title,
		MediaType: input.Body.MediaType,
	}

	createdMedia, err := utils.HandleDBError(
		s.mediaDB.CreateMedia(media),
	)

	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[MediaResponse]{
		Body: ToMediaResponse(createdMedia),
	}, nil
}

func (s *MediaService) GetMedia(ctx context.Context, input *GetMediaParams) (*utils.ResponseBody[MediaResponse], error) {
	media, err := s.mediaDB.GetMedia(&input.ID)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[MediaResponse]{
		Body: ToMediaResponse(media),
	}, nil
}

// DeleteMedia soft deletes a media by post ID
func (s *MediaService) DeleteMedia(ctx context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*utils.ResponseBody[MediaResponse], error) {
	id := &input.ID
	err := s.mediaDB.DeleteMedia(id)

	respBody := &utils.ResponseBody[MediaResponse]{}
	if err != nil {
		return respBody, err
	}

	response := &MediaResponse{
		ID: id,
	}

	return &utils.ResponseBody[MediaResponse]{
		Body: response,
	}, err
}
