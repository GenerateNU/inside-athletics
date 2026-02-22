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

// Builds an S3 key, gets a presigned PUT URL, and returns upload details.
func (s *Service) GetUploadURL(ctx context.Context, input GetUploadURLInput) (*GetUploadURLResponse, error) {
	if input.FileName == "" || input.FileType == "" || input.ContentKind == "" {
		return nil, fmt.Errorf("fileName, fileType, and contentKind are required")
	}
	owner := input.ContentID
	if owner == "" {
		owner = input.UserID
	}
	if owner == "" {
		return nil, fmt.Errorf("contentID or userID is required")
	}
	// Key format: premium/{kind}/{owner}/{fileName} so list/delete by content or user is easy.
	key := path.Join("premium", input.ContentKind, owner, input.FileName)
	expiresIn := s.cfg.PresignedURLExpiry
	if expiresIn == 0 {
		expiresIn = DefaultPresignedURLExpiry
	}
	metadata := map[string]string{
		"content-kind": input.ContentKind,
		"filename":     input.FileName,
	}
	uploadURL, err := s.client.PresignedUploadURL(ctx, key, input.FileType, expiresIn, metadata)
	if err != nil {
		return nil, err
	}
	return &GetUploadURLResponse{
		UploadURL:  uploadURL,
		Key:        key,
		DocumentID: input.FileName,
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
