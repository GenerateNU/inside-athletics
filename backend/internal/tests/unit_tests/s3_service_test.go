package unitTests

import (
	"context"
	"errors"
	"inside-athletics/internal/s3"
	"strings"
	"testing"
	"time"
)

var errHeadNotFound = errors.New("head object: not found")

// Tests GetUploadURL key format, response shape, validation, and contentID/userID behavior.
func TestS3Service_GetUploadURL(t *testing.T) {
	ctx := context.Background()
	mock := NewMockS3Client()
	cfg := s3.Config{Bucket: "test", Region: "us-east-1", PresignedURLExpiry: time.Hour}
	svc := s3.NewService(mock, cfg)

	t.Run("returns upload URL and correct key format for contentID", func(t *testing.T) {
		resp, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
			FileName:    "photo.jpg",
			FileType:    "image/jpeg",
			ContentKind: s3.ContentKindImage,
			ContentID:   "content-123",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.UploadURL == "" {
			t.Error("expected non-empty upload_url")
		}
		if !strings.HasPrefix(resp.UploadURL, "https://mock-upload.s3.local/") {
			t.Errorf("upload_url should start with mock prefix, got %s", resp.UploadURL)
		}
		expectedKey := "premium/image/content-123/photo.jpg"
		if resp.Key != expectedKey {
			t.Errorf("key: want %s, got %s", expectedKey, resp.Key)
		}
		if resp.DocumentID != "photo.jpg" {
			t.Errorf("document_id: want photo.jpg, got %s", resp.DocumentID)
		}
		if resp.ExpiresIn != 3600 {
			t.Errorf("expires_in: want 3600, got %d", resp.ExpiresIn)
		}
		if len(mock.UploadCalls) != 1 {
			t.Fatalf("expected 1 upload call, got %d", len(mock.UploadCalls))
		}
		if mock.UploadCalls[0].Key != expectedKey {
			t.Errorf("mock upload key: want %s, got %s", expectedKey, mock.UploadCalls[0].Key)
		}
	})

	t.Run("uses userID when contentID empty", func(t *testing.T) {
		mock := NewMockS3Client()
		svc := s3.NewService(mock, cfg)
		resp, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
			FileName:    "doc.pdf",
			FileType:    "application/pdf",
			ContentKind: s3.ContentKindPDF,
			UserID:      "user-456",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expectedKey := "premium/pdf/user-456/doc.pdf"
		if resp.Key != expectedKey {
			t.Errorf("key: want %s, got %s", expectedKey, resp.Key)
		}
	})

	t.Run("prefers contentID over userID", func(t *testing.T) {
		mock := NewMockS3Client()
		svc := s3.NewService(mock, cfg)
		resp, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
			FileName:    "vid.mp4",
			FileType:    "video/mp4",
			ContentKind: s3.ContentKindVideo,
			ContentID:   "content-1",
			UserID:      "user-1",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(resp.Key, "content-1") {
			t.Errorf("key should contain contentID, got %s", resp.Key)
		}
	})

	t.Run("missing fileName returns error", func(t *testing.T) {
		_, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
			FileType:    "image/jpeg",
			ContentKind: s3.ContentKindImage,
			ContentID:   "c1",
		})
		if err == nil {
			t.Fatal("expected error for missing fileName")
		}
		if !strings.Contains(err.Error(), "fileName") {
			t.Errorf("error should mention fileName, got %v", err)
		}
	})

	t.Run("missing contentID and userID returns error", func(t *testing.T) {
		_, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
			FileName:    "x.jpg",
			FileType:    "image/jpeg",
			ContentKind: s3.ContentKindImage,
		})
		if err == nil {
			t.Fatal("expected error when contentID and userID both empty")
		}
		if !strings.Contains(err.Error(), "contentID") && !strings.Contains(err.Error(), "userID") {
			t.Errorf("error should mention contentID or userID, got %v", err)
		}
	})
}

// Tests GetDownloadURL URL and expiry, and empty key error.
func TestS3Service_GetDownloadURL(t *testing.T) {
	ctx := context.Background()
	mock := NewMockS3Client()
	cfg := s3.Config{Bucket: "test", Region: "us-east-1", PresignedURLExpiry: time.Hour}
	svc := s3.NewService(mock, cfg)

	t.Run("returns download URL and expires_in", func(t *testing.T) {
		key := "premium/image/user-1/photo.jpg"
		resp, err := svc.GetDownloadURL(ctx, key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.DownloadURL == "" {
			t.Error("expected non-empty download_url")
		}
		if !strings.HasPrefix(resp.DownloadURL, "https://mock-download.s3.local/") {
			t.Errorf("download_url should start with mock prefix, got %s", resp.DownloadURL)
		}
		if !strings.HasSuffix(resp.DownloadURL, key) {
			t.Errorf("download_url should end with key, got %s", resp.DownloadURL)
		}
		if resp.ExpiresIn != 3600 {
			t.Errorf("expires_in: want 3600, got %d", resp.ExpiresIn)
		}
		if len(mock.DownloadCalls) != 1 || mock.DownloadCalls[0].Key != key {
			t.Errorf("mock download key: want %s, got %v", key, mock.DownloadCalls)
		}
	})

	t.Run("empty key returns error", func(t *testing.T) {
		_, err := svc.GetDownloadURL(ctx, "")
		if err == nil {
			t.Fatal("expected error for empty key")
		}
		if !strings.Contains(err.Error(), "key") {
			t.Errorf("error should mention key, got %v", err)
		}
	})
}

// Tests that zero config expiry falls back to default 1 hour.
func TestS3Service_GetUploadURL_defaultExpiry(t *testing.T) {
	ctx := context.Background()
	mock := NewMockS3Client()
	cfg := s3.Config{Bucket: "test", Region: "us-east-1"} // PresignedURLExpiry zero
	svc := s3.NewService(mock, cfg)
	resp, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
		FileName: "x.pdf", FileType: "application/pdf", ContentKind: s3.ContentKindPDF, ContentID: "c1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ExpiresIn != 3600 {
		t.Errorf("default expires_in should be 3600, got %d", resp.ExpiresIn)
	}
}

// Tests ConfirmUpload: HeadObject then PresignedDownloadURL, response shape, and empty key / HeadObject error.
func TestS3Service_ConfirmUpload(t *testing.T) {
	ctx := context.Background()
	mock := NewMockS3Client()
	mock.HeadObjectResponse.Size = 1024
	mock.HeadObjectResponse.Metadata = map[string]string{"content-kind": "image", "filename": "photo.jpg"}
	cfg := s3.Config{Bucket: "test", Region: "us-east-1", PresignedURLExpiry: time.Hour}
	svc := s3.NewService(mock, cfg)

	t.Run("returns key, download URL, size, and metadata", func(t *testing.T) {
		key := "premium/image/user-1/photo.jpg"
		resp, err := svc.ConfirmUpload(ctx, key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Key != key {
			t.Errorf("key: want %s, got %s", key, resp.Key)
		}
		if resp.DownloadURL == "" || !strings.HasPrefix(resp.DownloadURL, "https://mock-download.s3.local/") {
			t.Errorf("download_url should be mock URL, got %s", resp.DownloadURL)
		}
		if resp.Size != 1024 {
			t.Errorf("size: want 1024, got %d", resp.Size)
		}
		if resp.Metadata["filename"] != "photo.jpg" {
			t.Errorf("metadata filename: want photo.jpg, got %s", resp.Metadata["filename"])
		}
		if len(mock.HeadObjectCalls) != 1 || mock.HeadObjectCalls[0] != key {
			t.Errorf("HeadObject should be called with key, got %v", mock.HeadObjectCalls)
		}
		if len(mock.DownloadCalls) != 1 || mock.DownloadCalls[0].Key != key {
			t.Errorf("PresignedDownloadURL should be called with key, got %v", mock.DownloadCalls)
		}
	})

	t.Run("empty key returns error", func(t *testing.T) {
		_, err := svc.ConfirmUpload(ctx, "")
		if err == nil {
			t.Fatal("expected error for empty key")
		}
		if !strings.Contains(err.Error(), "key") {
			t.Errorf("error should mention key, got %v", err)
		}
	})

	t.Run("HeadObject error is returned", func(t *testing.T) {
		mock2 := NewMockS3Client()
		mock2.HeadObjectResponse.Err = errHeadNotFound
		svc2 := s3.NewService(mock2, cfg)
		_, err := svc2.ConfirmUpload(ctx, "premium/pdf/x/doc.pdf")
		if err == nil {
			t.Fatal("expected error when HeadObject fails")
		}
		if err != errHeadNotFound {
			t.Errorf("error: want errHeadNotFound, got %v", err)
		}
	})
}

// Tests DeleteObject calls client and validates key.
func TestS3Service_DeleteObject(t *testing.T) {
	ctx := context.Background()
	mock := NewMockS3Client()
	cfg := s3.Config{Bucket: "test", Region: "us-east-1"}
	svc := s3.NewService(mock, cfg)

	t.Run("deletes and records key", func(t *testing.T) {
		key := "premium/image/user-1/photo.jpg"
		err := svc.DeleteObject(ctx, key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(mock.DeleteObjectCalls) != 1 || mock.DeleteObjectCalls[0] != key {
			t.Errorf("DeleteObject should be called with key, got %v", mock.DeleteObjectCalls)
		}
	})

	t.Run("empty key returns error", func(t *testing.T) {
		err := svc.DeleteObject(ctx, "")
		if err == nil {
			t.Fatal("expected error for empty key")
		}
		if !strings.Contains(err.Error(), "key") {
			t.Errorf("error should mention key, got %v", err)
		}
	})
}
