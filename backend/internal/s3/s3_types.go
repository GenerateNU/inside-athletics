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
