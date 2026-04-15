package s3

import (
	"context"
	"fmt"
	"path"
)

// Wraps Client and Config to expose high-level flows.
type Service struct {
	client Client
	cfg    Config
}

// Returns a Service that uses the given Client and Config.
func NewService(client Client, cfg Config) *Service {
	return &Service{client: client, cfg: cfg}
}

// Builds a presigned PUT URL for the given key and returns upload details.
func (s *Service) GetUploadURL(ctx context.Context, input GetUploadURLInput) (*GetUploadURLResponse, error) {
	if input.Key == "" || input.FileType == "" {
		return nil, fmt.Errorf("key and fileType are required")
	}
	key := input.Key
	expiresIn := s.cfg.PresignedURLExpiry
	if expiresIn == 0 {
		expiresIn = DefaultPresignedURLExpiry
	}
	documentID := input.FileName
	if documentID == "" {
		documentID = path.Base(key)
	}
	metadata := map[string]string{
		"filename": documentID,
	}
	uploadURL, err := s.client.PresignedUploadURL(ctx, key, input.FileType, expiresIn, metadata)
	if err != nil {
		return nil, err
	}
	return &GetUploadURLResponse{
		UploadURL:  uploadURL,
		Key:        key,
		DocumentID: documentID,
		ExpiresIn:  int(expiresIn.Seconds()),
	}, nil
}

// Builds a presigned GET URL for the given key and returns download details.
func (s *Service) GetDownloadURL(ctx context.Context, key string) (*GetDownloadURLResponse, error) {
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}
	expiresIn := s.cfg.PresignedURLExpiry
	if expiresIn == 0 {
		expiresIn = DefaultPresignedURLExpiry
	}
	downloadURL, err := s.client.PresignedDownloadURL(ctx, key, expiresIn)
	if err != nil {
		return nil, err
	}
	return &GetDownloadURLResponse{
		DownloadURL: downloadURL,
		ExpiresIn:   int(expiresIn.Seconds()),
	}, nil
}

// Verifies object exists via HeadObject, then returns download URL and size/metadata.
func (s *Service) ConfirmUpload(ctx context.Context, key string) (*ConfirmUploadResponse, error) {
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}
	size, metadata, err := s.client.HeadObject(ctx, key)
	if err != nil {
		return nil, err
	}
	expiresIn := s.cfg.PresignedURLExpiry
	if expiresIn == 0 {
		expiresIn = DefaultPresignedURLExpiry
	}
	downloadURL, err := s.client.PresignedDownloadURL(ctx, key, expiresIn)
	if err != nil {
		return nil, err
	}
	return &ConfirmUploadResponse{
		Key:         key,
		DownloadURL: downloadURL,
		Size:        size,
		Metadata:    metadata,
	}, nil
}

// Removes the object at key from S3.
func (s *Service) DeleteObject(ctx context.Context, key string) error {
	if key == "" {
		return fmt.Errorf("key is required")
	}
	return s.client.DeleteObject(ctx, key)
}

// ResolveKey returns the presigned download URL for key.
// Returns "" if svc is nil, key is empty, or the request fails.
func ResolveKey(ctx context.Context, svc *Service, key string) string {
	if svc == nil || key == "" {
		return ""
	}
	resp, err := svc.GetDownloadURL(ctx, key)
	if err != nil {
		return ""
	}
	return resp.DownloadURL
}
