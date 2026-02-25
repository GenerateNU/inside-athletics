package content

import (
	"context"
	"inside-athletics/internal/s3"
	"inside-athletics/internal/utils"
)

// Holds S3 service for premium content upload/download URLs.
type ContentService struct {
	s3 *s3.Service
}

// Returns a ContentService that uses the given S3 service.
func NewContentService(s3Service *s3.Service) *ContentService {
	return &ContentService{s3: s3Service}
}

// Returns a presigned upload URL and key/expiry for the request body.
func (c *ContentService) GetUploadURL(ctx context.Context, input *GetUploadURLInput) (*utils.ResponseBody[s3.GetUploadURLResponse], error) {
	resp, err := c.s3.GetUploadURL(ctx, s3.GetUploadURLInput{
		FileName:    input.Body.FileName,
		FileType:    input.Body.FileType,
		ContentKind: input.Body.ContentKind,
		ContentID:   input.Body.ContentID,
		UserID:      input.Body.UserID,
	})
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[s3.GetUploadURLResponse]{Body: resp}, nil
}

// Returns a presigned download URL for the given key.
func (c *ContentService) GetDownloadURL(ctx context.Context, input *GetDownloadURLParams) (*utils.ResponseBody[s3.GetDownloadURLResponse], error) {
	resp, err := c.s3.GetDownloadURL(ctx, input.Key)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[s3.GetDownloadURLResponse]{Body: resp}, nil
}

// Confirms upload via HeadObject and returns download URL and size/metadata.
func (c *ContentService) ConfirmUpload(ctx context.Context, input *ConfirmUploadInput) (*utils.ResponseBody[s3.ConfirmUploadResponse], error) {
	resp, err := c.s3.ConfirmUpload(ctx, input.Body.Key)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[s3.ConfirmUploadResponse]{Body: resp}, nil
}

// Deletes the object at key from S3.
func (c *ContentService) DeleteContent(ctx context.Context, input *DeleteContentParams) (*utils.ResponseBody[DeleteContentResponse], error) {
	if err := c.s3.DeleteObject(ctx, input.Key); err != nil {
		return nil, err
	}
	return &utils.ResponseBody[DeleteContentResponse]{Body: &DeleteContentResponse{Message: "deleted"}}, nil
}
