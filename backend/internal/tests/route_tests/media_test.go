package routeTests

import (
	"inside-athletics/internal/handlers/media"
	"inside-athletics/internal/models"
	"net/http"
	"testing"
)

func TestCreateMedia(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Media{}); err != nil {
		t.Fatalf("failed to migrate media table: %v", err)
	}

	media.Route(testDB.API, testDB.DB)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "media"},
	})

	body := map[string]any{
		"s3key":      "test s3key",
		"title":      "test title",
		"media_type": "jpeg",
	}

	resp := api.Post("/api/v1/media/", body, authHeader)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result media.MediaResponse
	DecodeTo(&result, resp)

	if result.S3Key != "test s3key" {
		t.Errorf("expected s3key test s3key, got %s", result.S3Key)
	}

	if result.Title != "test title" {
		t.Errorf("expected title test title, got %v", result.Title)
	}

	if result.MediaType != "jpeg" {
		t.Errorf("expected media type jpeg, got %v", result.MediaType)
	}

}

func TestGetMedia(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Media{}); err != nil {
		t.Fatalf("failed to migrate media table: %v", err)
	}

	media.Route(testDB.API, testDB.DB)
	api := testDB.API
	mediaDB := media.NewMediaDB(testDB.DB)

	authHeader := authHeaderWithPermissions(t, testDB.DB, nil)

	vid := &models.Media{
		S3Key: "test s3key",
		Title: "test title",
	}

	createdMedia, err := mediaDB.CreateMedia(vid)
	if err != nil {
		t.Fatalf("failed to create title: %v", err)
	}

	resp := api.Get("/api/v1/media/"+createdMedia.ID.String(), authHeader)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result media.MediaResponse
	DecodeTo(&result, resp)

	if result.S3Key != "test s3key" {
		t.Errorf("expected s3key test s3key, got %s", result.S3Key)
	}

	if result.Title != "test title" {
		t.Errorf("expected title test title, got %v", result.Title)
	}
}

func TestDeleteMedia(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Media{}); err != nil {
		t.Fatalf("failed to migrate media table: %v", err)
	}

	media.Route(testDB.API, testDB.DB)
	api := testDB.API
	mediaDB := media.NewMediaDB(testDB.DB)

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionDelete, Resource: "media"},
	})

	vid := &models.Media{
		S3Key: "test s3key",
		Title: "test title",
	}
	createdMedia, err := mediaDB.CreateMedia(vid)
	if err != nil {
		t.Fatalf("failed to create media: %v", err)
	}

	resp := api.Delete("/api/v1/media/"+createdMedia.ID.String(), authHeader)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	// Verify the media is deleted
	getResp := api.Get("/api/v1/media/"+createdMedia.ID.String(), authHeader)
	if getResp.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getResp.Code)
	}
}
