package routeTests

import (
	"inside-athletics/internal/handlers/comment_like"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestCreateCommentLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	if err := testDB.DB.AutoMigrate(&models.Post{}, &models.Comment{}, &models.CommentLike{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	api := testDB.API

	postID := uuid.New()
	post := models.Post{ID: postID, UserID: uuid.New(), Title: "Test", Content: "Content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	commentID := uuid.New()
	comment := models.Comment{ID: commentID, UserID: uuid.New(), PostID: postID, Description: "A comment"}
	if err := testDB.DB.Create(&comment).Error; err != nil {
		t.Fatalf("create comment: %v", err)
	}
	userID := uuid.New()

	body := map[string]any{
		"user_id":    userID.String(),
		"comment_id": commentID.String(),
	}
	resp := api.Post("/api/v1/comment-like/", body, "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result utils.ResponseBody[comment_like.CreateCommentLikeResponse]
	DecodeTo(&result, resp)
	if result.Body == nil || result.Body.ID == uuid.Nil {
		t.Error("expected non-zero like ID")
	}
}

func TestGetCommentLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	if err := testDB.DB.AutoMigrate(&models.Post{}, &models.Comment{}, &models.CommentLike{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	api := testDB.API

	postID := uuid.New()
	post := models.Post{ID: postID, UserID: uuid.New(), Title: "Test", Content: "Content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	commentID := uuid.New()
	comment := models.Comment{ID: commentID, UserID: uuid.New(), PostID: postID, Description: "A comment"}
	if err := testDB.DB.Create(&comment).Error; err != nil {
		t.Fatalf("create comment: %v", err)
	}
	userID := uuid.New()
	like := models.CommentLike{ID: uuid.New(), UserID: userID, CommentID: commentID}
	if err := testDB.DB.Create(&like).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Get("/api/v1/comment-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result utils.ResponseBody[comment_like.GetCommentLikeResponse]
	DecodeTo(&result, resp)
	if result.Body == nil || result.Body.CommentID != commentID || result.Body.UserID != userID {
		t.Errorf("expected comment_id %s user_id %s, got %+v", commentID, userID, result.Body)
	}
}

func TestGetCommentLikeInfo(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	if err := testDB.DB.AutoMigrate(&models.Post{}, &models.Comment{}, &models.CommentLike{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	api := testDB.API

	postID := uuid.New()
	post := models.Post{ID: postID, UserID: uuid.New(), Title: "Test", Content: "Content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	commentID := uuid.New()
	comment := models.Comment{ID: commentID, UserID: uuid.New(), PostID: postID, Description: "A comment"}
	if err := testDB.DB.Create(&comment).Error; err != nil {
		t.Fatalf("create comment: %v", err)
	}
	userID := uuid.New()
	if err := testDB.DB.Create(&models.CommentLike{UserID: userID, CommentID: commentID}).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Get("/api/v1/comment-like/comment/"+commentID.String()+"/likes?user_id="+userID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result utils.ResponseBody[comment_like.GetCommentLikeInfoResponse]
	DecodeTo(&result, resp)
	if result.Body == nil || result.Body.Total != 1 {
		t.Errorf("expected total 1, got %v", result.Body)
	}
	if result.Body != nil && !result.Body.Liked {
		t.Error("expected liked true")
	}
}

func TestDeleteCommentLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	if err := testDB.DB.AutoMigrate(&models.Post{}, &models.Comment{}, &models.CommentLike{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	api := testDB.API

	postID := uuid.New()
	post := models.Post{ID: postID, UserID: uuid.New(), Title: "Test", Content: "Content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	commentID := uuid.New()
	comment := models.Comment{ID: commentID, UserID: uuid.New(), PostID: postID, Description: "A comment"}
	if err := testDB.DB.Create(&comment).Error; err != nil {
		t.Fatalf("create comment: %v", err)
	}
	like := models.CommentLike{ID: uuid.New(), UserID: uuid.New(), CommentID: commentID}
	if err := testDB.DB.Create(&like).Error; err != nil {
		t.Fatalf("create like: %v", err)
	}

	resp := api.Delete("/api/v1/comment-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", resp.Code, resp.Body.String())
	}

	getResp := api.Get("/api/v1/comment-like/"+like.ID.String(), "Authorization: Bearer mock-token")
	if getResp.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getResp.Code)
	}
}

func TestCreateCommentLikeDuplicateReturns409(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	if err := testDB.DB.AutoMigrate(&models.Post{}, &models.Comment{}, &models.CommentLike{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	api := testDB.API

	postID := uuid.New()
	post := models.Post{ID: postID, UserID: uuid.New(), Title: "Test", Content: "Content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("create post: %v", err)
	}
	commentID := uuid.New()
	comment := models.Comment{ID: commentID, UserID: uuid.New(), PostID: postID, Description: "A comment"}
	if err := testDB.DB.Create(&comment).Error; err != nil {
		t.Fatalf("create comment: %v", err)
	}
	userID := uuid.New()

	body := map[string]any{
		"user_id":    userID.String(),
		"comment_id": commentID.String(),
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
