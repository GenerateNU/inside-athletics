package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Implements Client using AWS SDK.
type client struct {
	s3       *s3.Client
	presign  *s3.PresignClient
	cfg      Config
}

// Builds a Client from cfg using the default AWS credential chain.
func NewClient(ctx context.Context, cfg Config) (Client, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return nil, err
	}
	s3Client := s3.NewFromConfig(awsCfg)
	// PresignClient only generates signed URLs, s3 client does Head/Delete and is required by the presigner.
	presignClient := s3.NewPresignClient(s3Client)
	return &client{s3: s3Client, presign: presignClient, cfg: cfg}, nil
}

// Returns a time-limited PUT URL for the given key (expires after expiresIn).
func (c *client) PresignedUploadURL(ctx context.Context, key, contentType string, expiresIn time.Duration, metadata map[string]string) (string, error) {
	// SDK expects *string for most input fields.
	input := &s3.PutObjectInput{
		Bucket:      aws.String(c.cfg.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Metadata:    metadata,
	}
	resp, err := c.presign.PresignPutObject(ctx, input, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return "", err
	}
	return resp.URL, nil
}

// Returns a time-limited GET URL for the given key.
func (c *client) PresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
		Key:    aws.String(key),
	}
	resp, err := c.presign.PresignGetObject(ctx, input, s3.WithPresignExpires(expiresIn))
	if err != nil {
		return "", err
	}
	return resp.URL, nil
}

// Returns object size and metadata to confirm upload before persisting.
func (c *client) HeadObject(ctx context.Context, key string) (size int64, metadata map[string]string, err error) {
	out, err := c.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, nil, err
	}
	if out.Metadata != nil {
		metadata = out.Metadata
	} else {
		metadata = make(map[string]string) // avoid returning nil
	}
	size = 0
	// ContentLength is *int64, dereference only when non-nil.
	if out.ContentLength != nil {
		size = *out.ContentLength
	}
	return size, metadata, nil
}

// Removes the object at key.
func (c *client) DeleteObject(ctx context.Context, key string) error {
	_, err := c.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.cfg.Bucket),
		Key:    aws.String(key),
	})
	return err
}
