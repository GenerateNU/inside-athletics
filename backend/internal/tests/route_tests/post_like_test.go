package routeTests

import (
	"inside-athletics/internal/handlers/post_like"
	"inside-athletics/internal/models"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestCreatePostLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	// seedUserAndPost is from comment_test
	user, post := seedUserAndPost(t, testDB, "create-like")

	body := map[string]any{
		"user_id": user.ID.String(),
		"post_id": post.ID.String(),
	}
	resp := api.Post("/api/v1/post-like/", body, "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post_like.CreatePostLikeResponse

	DecodeTo(&result, resp)
	if result.ID == uuid.Nil {
		t.Error("expected non-zero like ID")
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
	if !result.Liked {
		t.Error("expected liked true after create")
	}
}

func TestGetPostLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, post := seedUserAndPost(t, testDB, "get-like")
	like := models.PostLike{UserID: user.ID, PostID: post.ID}
	if err := testDB.DB.Create(&like).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Get("/api/v1/post-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post_like.GetPostLikeResponse

	DecodeTo(&result, resp)
	if result.PostID != post.ID || result.UserID != user.ID {
		t.Errorf("expected post_id %s user_id %s, got %+v", post.ID, user.ID, result)
	}
}

func TestGetPostLikeInfo(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, post := seedUserAndPost(t, testDB, "like-info")
	if err := testDB.DB.Create(&models.PostLike{UserID: user.ID, PostID: post.ID}).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Get("/api/v1/post-like/post/"+post.ID.String()+"/likes?user_id="+user.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post_like.GetPostLikeInfoResponse

	DecodeTo(&result, resp)
	if result.Total != 1 {
		t.Errorf("expected total 1, got %v", result)
	}
	if !result.Liked {
		t.Error("expected liked true")
	}
}

func TestDeletePostLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, post := seedUserAndPost(t, testDB, "delete-like")
	like := models.PostLike{UserID: user.ID, PostID: post.ID}
	if err := testDB.DB.Create(&like).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Delete("/api/v1/post-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post_like.DeletePostLikeResponse

	DecodeTo(&result, resp)
	if result.Message != "Like was deleted successfully" {
		t.Errorf("expected success message, got %+v", result)
	}
	if result.Total != 0 {
		t.Errorf("expected total 0 after delete, got %d", result.Total)
	}
	if result.Liked {
		t.Error("expected liked false after delete")
	}

	getResp := api.Get("/api/v1/post-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if getResp.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getResp.Code)
	}
}

func TestCreatePostLikeDuplicateReturns409(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, post := seedUserAndPost(t, testDB, "dup-like")

	body := map[string]any{
		"user_id": user.ID.String(),
		"post_id": post.ID.String(),
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
