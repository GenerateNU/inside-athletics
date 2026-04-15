package unitTests

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"inside-athletics/internal/s3"
	"io"
	"strings"
	"testing"
	"time"
)

var errHeadNotFound = errors.New("head object: not found")

// Tests GetUploadURL key format, response shape, validation, and contentID/userID behavior.
func TestS3Service_GetUploadURL(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mock := NewMockS3Client()
	cfg := s3.Config{Bucket: "test", Region: "us-east-1", PresignedURLExpiry: time.Hour}
	svc := s3.NewService(mock, cfg)

	t.Run("returns upload URL and correct key format for contentID", func(t *testing.T) {
		t.Parallel()
		mock := NewMockS3Client()
		svc := s3.NewService(mock, cfg)
		resp, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
			Key:      "premium/image/content-123/photo.jpg",
			FileType: "image/jpeg",
			FileName: "photo.jpg",
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
		t.Parallel()
		mock := NewMockS3Client()
		svc := s3.NewService(mock, cfg)
		resp, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
			Key:      "premium/pdf/user-456/doc.pdf",
			FileType: "application/pdf",
			FileName: "doc.pdf",
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
		t.Parallel()
		mock := NewMockS3Client()
		svc := s3.NewService(mock, cfg)
		resp, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
			Key:      "premium/video/content-1/vid.mp4",
			FileType: "video/mp4",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(resp.Key, "content-1") {
			t.Errorf("key should contain content-1, got %s", resp.Key)
		}
		if resp.DocumentID != "vid.mp4" {
			t.Errorf("document_id should default to path base, want vid.mp4, got %s", resp.DocumentID)
		}
	})

	t.Run("missing fileName defaults documentId to path base", func(t *testing.T) {
		t.Parallel()
		resp, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
			Key:      "premium/image/c1/photo.jpg",
			FileType: "image/jpeg",
		})
		if err != nil {
			t.Fatalf("fileName is optional: %v", err)
		}
		if resp.DocumentID != "photo.jpg" {
			t.Errorf("document_id should default to path base when fileName empty, want photo.jpg, got %s", resp.DocumentID)
		}
	})

	t.Run("missing key returns error", func(t *testing.T) {
		t.Parallel()
		_, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
			Key:      "",
			FileType: "image/jpeg",
			FileName: "x.jpg",
		})
		if err == nil {
			t.Fatal("expected error when key empty")
		}
		if !strings.Contains(err.Error(), "key") {
			t.Errorf("error should mention key, got %v", err)
		}
	})

	t.Run("missing fileType returns error", func(t *testing.T) {
		t.Parallel()
		_, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
			Key:      "premium/image/c1/x.jpg",
			FileType: "",
		})
		if err == nil {
			t.Fatal("expected error when fileType empty")
		}
		if !strings.Contains(err.Error(), "fileType") {
			t.Errorf("error should mention fileType, got %v", err)
		}
	})
}

// Tests GetDownloadURL URL and expiry, and empty key error.
func TestS3Service_GetDownloadURL(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mock := NewMockS3Client()
	cfg := s3.Config{Bucket: "test", Region: "us-east-1", PresignedURLExpiry: time.Hour}
	svc := s3.NewService(mock, cfg)

	t.Run("returns download URL and expires_in", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
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
	t.Parallel()
	ctx := context.Background()
	mock := NewMockS3Client()
	cfg := s3.Config{Bucket: "test", Region: "us-east-1"} // PresignedURLExpiry zero
	svc := s3.NewService(mock, cfg)
	resp, err := svc.GetUploadURL(ctx, s3.GetUploadURLInput{
		Key:      "premium/pdf/c1/x.pdf",
		FileType: "application/pdf",
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
	t.Parallel()
	ctx := context.Background()
	mock := NewMockS3Client()
	mock.HeadObjectResponse.Size = 1024
	mock.HeadObjectResponse.Metadata = map[string]string{"content-kind": "image", "filename": "photo.jpg"}
	cfg := s3.Config{Bucket: "test", Region: "us-east-1", PresignedURLExpiry: time.Hour}
	svc := s3.NewService(mock, cfg)

	t.Run("returns key, download URL, size, and metadata", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
		_, err := svc.ConfirmUpload(ctx, "")
		if err == nil {
			t.Fatal("expected error for empty key")
		}
		if !strings.Contains(err.Error(), "key") {
			t.Errorf("error should mention key, got %v", err)
		}
	})

	t.Run("HeadObject error is returned", func(t *testing.T) {
		t.Parallel()
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
	t.Parallel()
	ctx := context.Background()
	mock := NewMockS3Client()
	cfg := s3.Config{Bucket: "test", Region: "us-east-1"}
	svc := s3.NewService(mock, cfg)

	t.Run("deletes and records key", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
		err := svc.DeleteObject(ctx, "")
		if err == nil {
			t.Fatal("expected error for empty key")
		}
		if !strings.Contains(err.Error(), "key") {
			t.Errorf("error should mention key, got %v", err)
		}
	})
}

// --- CompressBytes tests (different media types) ---

func TestCompressBytes_EmptyInput(t *testing.T) {
	t.Parallel()
	out, err := s3.CompressBytes(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty gzip output for nil input")
	}
	dec, err := gzip.NewReader(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("output is not valid gzip: %v", err)
	}
	defer func() {
		_ = dec.Close()
	}()
	got, err := io.ReadAll(dec)
	if err != nil {
		t.Fatalf("decompress failed: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("round-trip: got %d bytes, want 0", len(got))
	}
}

func TestCompressBytes_PlainText(t *testing.T) {
	t.Parallel()
	src := []byte("Hello world. This is plain text that compresses well. " + "Repeated repeated repeated repeated repeated.")
	out, err := s3.CompressBytes(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// gzip overhead can make tiny inputs larger; we only assert valid round-trip.
	roundTripCompress(t, out, src)
}

func TestCompressBytes_JSON(t *testing.T) {
	t.Parallel()
	src := []byte(`{"name":"test","kind":"image","items":[1,2,3],"nested":{"a":true}}`)
	out, err := s3.CompressBytes(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	roundTripCompress(t, out, src)
}

func TestCompressBytes_HTML(t *testing.T) {
	t.Parallel()
	src := []byte(`<!DOCTYPE html><html><head><title>Test</title></head><body><p>Content</p></body></html>`)
	out, err := s3.CompressBytes(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	roundTripCompress(t, out, src)
}

func TestCompressBytes_ImageLike(t *testing.T) {
	t.Parallel()
	src := make([]byte, 256)
	src[0], src[1] = 0xFF, 0xD8 // JPEG magic
	for i := 2; i < len(src); i++ {
		src[i] = byte(i * 7)
	}
	out, err := s3.CompressBytes(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	roundTripCompress(t, out, src)
}

func TestCompressBytes_PDFLike(t *testing.T) {
	t.Parallel()
	src := []byte("%PDF-1.4\n%\xe2\xe3\xcf\xd3\n1 0 obj\n/Type /Catalog\nendobj\n" +
		"stream\nstream stream stream stream stream\nendstream\n")
	out, err := s3.CompressBytes(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	roundTripCompress(t, out, src)
}

func TestCompressBytes_VideoLike(t *testing.T) {
	t.Parallel()
	src := make([]byte, 1024)
	for i := range src {
		src[i] = byte(i*i + 1)
	}
	out, err := s3.CompressBytes(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	roundTripCompress(t, out, src)
}

func roundTripCompress(t *testing.T, compressed, want []byte) {
	t.Helper()
	dec, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatalf("output is not valid gzip: %v", err)
	}
	defer func() {
		_ = dec.Close()
	}()
	got, err := io.ReadAll(dec)
	if err != nil {
		t.Fatalf("decompress failed: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("round-trip mismatch: got %d bytes, want %d", len(got), len(want))
	}
}
