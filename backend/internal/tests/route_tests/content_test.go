package routeTests

import (
	"bytes"
	"encoding/json"
	"inside-athletics/internal/s3"
	"net/url"
	"strings"
	"testing"
)


// Tests POST /api/v1/content/confirm-upload with mock S3.
func TestConfirmUpload(t *testing.T) {
	api := RegisterContentTestAPI(t)

	body := map[string]string{"key": "premium/image/user-1/photo.jpg"}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}
	resp := api.Post("/api/v1/content/confirm-upload", "Authorization: Bearer mock-token", "Content-Type: application/json", bytes.NewReader(jsonBody))
	if resp.Code != 200 {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}
	var out s3.ConfirmUploadResponse
	DecodeTo(&out, resp)
	if out.Key != "premium/image/user-1/photo.jpg" {
		t.Errorf("key: want premium/image/user-1/photo.jpg, got %s", out.Key)
	}
	if out.DownloadURL == "" || !strings.HasPrefix(out.DownloadURL, "https://mock-download.s3.local/") {
		t.Errorf("download_url: want mock URL, got %s", out.DownloadURL)
	}
	if out.Size != 1024 {
		t.Errorf("size: want 1024, got %d", out.Size)
	}
	if out.Metadata["filename"] != "photo.jpg" {
		t.Errorf("metadata filename: want photo.jpg, got %s", out.Metadata["filename"])
	}
}

// Tests DELETE /api/v1/content?key=... with mock S3.
func TestDeleteContent(t *testing.T) {
	api := RegisterContentTestAPI(t)

	key := "premium/pdf/user-1/doc.pdf"
	resp := api.Delete("/api/v1/content?key="+url.QueryEscape(key), "Authorization: Bearer mock-token")
	if resp.Code != 200 {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}
	var out struct {
		Message string `json:"message"`
	}
	DecodeTo(&out, resp)
	if out.Message != "deleted" {
		t.Errorf("message: want deleted, got %s", out.Message)
	}
}
