package unitTests

import (
	"context"
	"inside-athletics/internal/s3"
	"time"
)

// MockS3Client implements s3.Client for tests; returns deterministic URLs and records calls.
type MockS3Client struct {
	UploadURLPrefix    string // prefix for presigned upload URLs (default "https://mock-upload.s3.local/")
	DownloadURLPrefix  string // prefix for presigned download URLs (default "https://mock-download.s3.local/")
	UploadCalls        []MockS3UploadCall
	DownloadCalls      []MockS3DownloadCall
	HeadObjectCalls    []string
	DeleteObjectCalls  []string
	// Configures what HeadObject returns (size, metadata, or Err).
	HeadObjectResponse struct {
		Size     int64
		Metadata map[string]string
		Err      error
	}
}

// Records one PresignedUploadURL call for assertions.
type MockS3UploadCall struct {
	Key         string
	ContentType string
	ExpiresIn   time.Duration
}

// Records one PresignedDownloadURL call for assertions.
type MockS3DownloadCall struct {
	Key       string
	ExpiresIn time.Duration
}

// Returns a deterministic upload URL and records the call.
func (m *MockS3Client) PresignedUploadURL(ctx context.Context, key, contentType string, expiresIn time.Duration, metadata map[string]string) (string, error) {
	m.UploadCalls = append(m.UploadCalls, MockS3UploadCall{Key: key, ContentType: contentType, ExpiresIn: expiresIn})
	prefix := m.UploadURLPrefix
	if prefix == "" {
		prefix = "https://mock-upload.s3.local/"
	}
	return prefix + key, nil
}

// Returns a deterministic download URL and records the call.
func (m *MockS3Client) PresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error) {
	m.DownloadCalls = append(m.DownloadCalls, MockS3DownloadCall{Key: key, ExpiresIn: expiresIn})
	prefix := m.DownloadURLPrefix
	if prefix == "" {
		prefix = "https://mock-download.s3.local/"
	}
	return prefix + key, nil
}

// Records the key and returns configured size/metadata or error.
func (m *MockS3Client) HeadObject(ctx context.Context, key string) (size int64, metadata map[string]string, err error) {
	m.HeadObjectCalls = append(m.HeadObjectCalls, key)
	if m.HeadObjectResponse.Err != nil {
		return 0, nil, m.HeadObjectResponse.Err
	}
	meta := m.HeadObjectResponse.Metadata
	if meta == nil {
		meta = make(map[string]string)
	}
	return m.HeadObjectResponse.Size, meta, nil
}

// Records the key and returns nil.
func (m *MockS3Client) DeleteObject(ctx context.Context, key string) error {
	m.DeleteObjectCalls = append(m.DeleteObjectCalls, key)
	return nil
}

// NewMockS3Client returns a MockS3Client that records calls for assertions.
func NewMockS3Client() *MockS3Client {
	return &MockS3Client{
		UploadCalls:   make([]MockS3UploadCall, 0),
		DownloadCalls: make([]MockS3DownloadCall, 0),
	}
}

// Ensure MockS3Client implements s3.Client at compile time.
var _ s3.Client = (*MockS3Client)(nil)
