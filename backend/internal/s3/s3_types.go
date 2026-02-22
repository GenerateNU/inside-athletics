package s3

import (
	"context"
	"time"
)

// Defines S3 operations for the premium content hub.
type Client interface {
	// Returns a time-limited PUT URL for the given key (expires after expiresIn).
	PresignedUploadURL(ctx context.Context, key, contentType string, expiresIn time.Duration, metadata map[string]string) (string, error)
	// Returns a time-limited GET URL for the given key.
	PresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)
	// Returns object size and metadata to confirm upload before persisting.
	HeadObject(ctx context.Context, key string) (size int64, metadata map[string]string, err error)
	// Removes the object at key.
	DeleteObject(ctx context.Context, key string) error
}

// Holds S3 bucket, region, and presigned URL expiry.
type Config struct {
	Bucket              string        // S3 bucket name
	Region              string        // AWS region, must match bucket
	PresignedURLExpiry  time.Duration // how long upload/download URLs stay valid
}

// Default presigned URL TTL (1 hour).
const DefaultPresignedURLExpiry = time.Hour

// Content kind for premium hub (image, video, pdf).
const (
	ContentKindImage = "image"
	ContentKindVideo = "video"
	ContentKindPDF   = "pdf"
)

// Input for requesting a presigned upload URL.
type GetUploadURLInput struct {
	FileName    string // original filename (used in key and as documentId for confirm)
	FileType    string // MIME type, e.g. image/jpeg, application/pdf
	ContentKind string // image, video, or pdf
	ContentID   string // optional; preferred for key path if set
	UserID      string // optional; used for key path if ContentID empty
}

// Response after generating a presigned upload URL.
type GetUploadURLResponse struct {
	UploadURL  string `json:"upload_url"`
	Key        string `json:"key"`
	DocumentID string `json:"document_id"` // pass back on confirm
	ExpiresIn  int    `json:"expires_in"` // seconds until URL expires
}

// Response after generating a presigned download URL (for PDF/image).
type GetDownloadURLResponse struct {
	DownloadURL string `json:"download_url"`
	ExpiresIn   int    `json:"expires_in"` // seconds until URL expires
}

// Response after confirming an upload (HeadObject + presigned download URL).
type ConfirmUploadResponse struct {
	Key         string            `json:"key"`
	DownloadURL string            `json:"download_url"`
	Size        int64             `json:"size"`
	Metadata    map[string]string  `json:"metadata,omitempty"`
}
