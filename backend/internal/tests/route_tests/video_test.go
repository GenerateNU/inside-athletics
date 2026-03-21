package routeTests

import (
	"inside-athletics/internal/handlers/video"
	"inside-athletics/internal/models"
	"net/http"
	"testing"
)

func TestCreateVideo(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Video{}); err != nil {
		t.Fatalf("failed to migrate video table: %v", err)
	}

	video.Route(testDB.API, testDB.DB)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "video"},
	})

	body := map[string]any{
		"s3key": "test s3key",
		"title": "test title",
	}

	resp := api.Post("/api/v1/video/", body, authHeader)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result video.VideoResponse
	DecodeTo(&result, resp)

	if result.S3Key != "test s3key" {
		t.Errorf("expected s3key test s3key, got %s", result.S3Key)
	}

	if result.Title != "test title" {
		t.Errorf("expected title test title, got %v", result.Title)
	}
}

func TestGetVideo(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Video{}); err != nil {
		t.Fatalf("failed to migrate video table: %v", err)
	}

	video.Route(testDB.API, testDB.DB)
	api := testDB.API
	videoDB := video.NewVideoDB(testDB.DB)

	authHeader := authHeaderWithPermissions(t, testDB.DB, nil)

	vid := &models.Video{
		S3Key: "test s3key",
		Title: "test title",
	}

	createdVideo, err := videoDB.CreateVideo(vid)
	if err != nil {
		t.Fatalf("failed to create title: %v", err)
	}

	resp := api.Get("/api/v1/video/"+createdVideo.ID.String(), authHeader)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result video.VideoResponse
	DecodeTo(&result, resp)

	if result.S3Key != "test s3key" {
		t.Errorf("expected s3key test s3key, got %s", result.S3Key)
	}

	if result.Title != "test title" {
		t.Errorf("expected title test title, got %v", result.Title)
	}
}

func TestDeleteVideo(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Video{}); err != nil {
		t.Fatalf("failed to migrate video table: %v", err)
	}

	video.Route(testDB.API, testDB.DB)
	api := testDB.API
	videoDB := video.NewVideoDB(testDB.DB)

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionDelete, Resource: "video"},
	})

	vid := &models.Video{
		S3Key: "test s3key",
		Title: "test title",
	}
	createdVideo, err := videoDB.CreateVideo(vid)
	if err != nil {
		t.Fatalf("failed to create video: %v", err)
	}

	resp := api.Delete("/api/v1/video/"+createdVideo.ID.String(), authHeader)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	// Verify the video is deleted
	getResp := api.Get("/api/v1/video/"+createdVideo.ID.String(), authHeader)
	if getResp.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getResp.Code)
	}
}
