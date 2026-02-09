package routeTests

import (
	"inside-athletics/internal/handlers/post_like"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestCreatePostLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	if err := testDB.DB.AutoMigrate(&models.Post{}, &models.PostLike{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	api := testDB.API

	postID := uuid.New()
	post := models.Post{ID: postID, UserID: uuid.New(), Title: "Test", Content: "Content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	userID := uuid.New()

	body := map[string]any{
		"user_id": userID.String(),
		"post_id": postID.String(),
	}
	resp := api.Post("/api/v1/post-like/", body, "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result utils.ResponseBody[post_like.CreatePostLikeResponse]
	DecodeTo(&result, resp)
	if result.Body == nil || result.Body.ID == uuid.Nil {
		t.Error("expected non-zero like ID")
	}
}

func TestGetPostLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	if err := testDB.DB.AutoMigrate(&models.Post{}, &models.PostLike{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	api := testDB.API

	postID := uuid.New()
	post := models.Post{ID: postID, UserID: uuid.New(), Title: "Test", Content: "Content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	userID := uuid.New()
	like := models.PostLike{ID: uuid.New(), UserID: userID, PostID: postID}
	if err := testDB.DB.Create(&like).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Get("/api/v1/post-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result utils.ResponseBody[post_like.GetPostLikeResponse]
	DecodeTo(&result, resp)
	if result.Body == nil || result.Body.PostID != postID || result.Body.UserID != userID {
		t.Errorf("expected post_id %s user_id %s, got %+v", postID, userID, result.Body)
	}
}

func TestGetPostLikeInfo(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	if err := testDB.DB.AutoMigrate(&models.Post{}, &models.PostLike{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	api := testDB.API

	postID := uuid.New()
	post := models.Post{ID: postID, UserID: uuid.New(), Title: "Test", Content: "Content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	userID := uuid.New()
	if err := testDB.DB.Create(&models.PostLike{UserID: userID, PostID: postID}).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Get("/api/v1/post-like/post/"+postID.String()+"/likes?user_id="+userID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result utils.ResponseBody[post_like.GetPostLikeInfoResponse]
	DecodeTo(&result, resp)
	if result.Body == nil || result.Body.Total != 1 {
		t.Errorf("expected total 1, got %v", result.Body)
	}
	if result.Body != nil && !result.Body.Liked {
		t.Error("expected liked true")
	}
}

func TestDeletePostLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	if err := testDB.DB.AutoMigrate(&models.Post{}, &models.PostLike{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	api := testDB.API

	postID := uuid.New()
	post := models.Post{ID: postID, UserID: uuid.New(), Title: "Test", Content: "Content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	like := models.PostLike{ID: uuid.New(), UserID: uuid.New(), PostID: postID}
	if err := testDB.DB.Create(&like).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Delete("/api/v1/post-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	getResp := api.Get("/api/v1/post-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if getResp.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getResp.Code)
	}
}

func TestCreatePostLikeDuplicateReturns409(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	if err := testDB.DB.AutoMigrate(&models.Post{}, &models.PostLike{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	api := testDB.API

	postID := uuid.New()
	post := models.Post{ID: postID, UserID: uuid.New(), Title: "Test", Content: "Content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	userID := uuid.New()

	body := map[string]any{
		"user_id": userID.String(),
		"post_id": postID.String(),
	}
	resp1 := api.Post("/api/v1/post-like/", body, "Authorization: Bearer mock-token")
	if resp1.Code != http.StatusOK {
		t.Fatalf("first create expected 200, got %d: %s", resp1.Code, resp1.Body.String())
	}

	resp2 := api.Post("/api/v1/post-like/", body, "Authorization: Bearer mock-token")
	if resp2.Code != http.StatusConflict {
		t.Errorf("expected 409 for duplicate like, got %d: %s", resp2.Code, resp2.Body.String())
	}
}
