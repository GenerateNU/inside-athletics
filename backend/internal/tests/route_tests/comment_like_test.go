package routeTests

import (
	"inside-athletics/internal/handlers/comment_like"
	"inside-athletics/internal/models"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

// seedUserPostAndComment reuses seedUserAndPost from comment_test.go, then creates a Comment so we can test
func seedUserPostAndComment(t *testing.T, testDB *TestDatabase, unique string) (models.User, models.Post, models.Comment) {
	t.Helper()
	user, post := seedUserAndPost(t, testDB, unique)
	comment := models.Comment{
		UserID:      user.ID,
		PostID:      post.ID,
		Description: "A comment",
	}
	if err := testDB.DB.Create(&comment).Error; err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}
	return user, post, comment
}

func TestCreateCommentLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, _, comment := seedUserPostAndComment(t, testDB, "create-like")

	body := map[string]any{
		"user_id":    user.ID.String(),
		"comment_id": comment.ID.String(),
	}
	resp := api.Post("/api/v1/comment-like/", body, "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result comment_like.CreateCommentLikeResponse
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

func TestGetCommentLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, _, comment := seedUserPostAndComment(t, testDB, "get-like")
	like := models.CommentLike{UserID: user.ID, CommentID: comment.ID}
	if err := testDB.DB.Create(&like).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Get("/api/v1/comment-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result comment_like.GetCommentLikeResponse
	DecodeTo(&result, resp)
	if result.CommentID != comment.ID || result.UserID != user.ID {
		t.Errorf("expected comment_id %s user_id %s, got %+v", comment.ID, user.ID, result)
	}
}

func TestGetCommentLikeInfo(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, _, comment := seedUserPostAndComment(t, testDB, "like-info")
	if err := testDB.DB.Create(&models.CommentLike{UserID: user.ID, CommentID: comment.ID}).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Get("/api/v1/comment-like/comment/"+comment.ID.String()+"/likes?user_id="+user.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result comment_like.GetCommentLikeInfoResponse
	DecodeTo(&result, resp)
	if result.Total != 1 {
		t.Errorf("expected total 1, got %v", result)
	}
	if !result.Liked {
		t.Error("expected liked true")
	}
}

func TestDeleteCommentLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, _, comment := seedUserPostAndComment(t, testDB, "delete-like")
	like := models.CommentLike{UserID: user.ID, CommentID: comment.ID}
	if err := testDB.DB.Create(&like).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Delete("/api/v1/comment-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result comment_like.DeleteCommentLikeResponse
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

	getResp := api.Get("/api/v1/comment-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if getResp.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getResp.Code)
	}
}

func TestCreateCommentLikeDuplicateReturns409(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, _, comment := seedUserPostAndComment(t, testDB, "dup-like")

	body := map[string]any{
		"user_id":    user.ID.String(),
		"comment_id": comment.ID.String(),
	}
	resp1 := api.Post("/api/v1/comment-like/", body, "Authorization: Bearer mock-token")
	if resp1.Code != http.StatusOK {
		t.Fatalf("first create expected 200, got %d: %s", resp1.Code, resp1.Body.String())
	}

	resp2 := api.Post("/api/v1/comment-like/", body, "Authorization: Bearer mock-token")
	if resp2.Code != http.StatusConflict {
		t.Errorf("expected 409 for duplicate like, got %d: %s", resp2.Code, resp2.Body.String())
	}
}
