package routeTests

import (
	"inside-athletics/internal/handlers/comment"
	"inside-athletics/internal/models"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestCreateComment(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	postID := uuid.New()
	post := models.Post{ID: postID, AuthorId: uuid.New(), Title: "Test Post", Content: "Test content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	userID := uuid.New()

	body := map[string]any{
		"user_id":      userID.String(),
		"post_id":      postID.String(),
		"description":  "A test comment",
		"is_anonymous": false,
	}

	resp := api.Post("/api/v1/comment/", body, "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result comment.CommentResponse
	DecodeTo(&result, resp)
	if result.Description != "A test comment" {
		t.Errorf("expected description A test comment, got %s", result.Description)
	}
	if result.PostID != postID {
		t.Errorf("expected post_id %s, got %s", postID, result.PostID)
	}
	if result.UserID == nil || *result.UserID != userID {
		t.Errorf("expected user_id for non-anonymous comment, got %v", result.UserID)
	}
}

// Asserts anonymous comments hide user_id when caller is not a super user (default context).
func TestCreateCommentAnonymous(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	postID := uuid.New()
	post := models.Post{ID: postID, AuthorId: uuid.New(), Title: "Test Post", Content: "Test content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	userID := uuid.New()

	body := map[string]any{
		"user_id":      userID.String(),
		"post_id":      postID.String(),
		"description":  "Anonymous comment",
		"is_anonymous": true,
	}

	resp := api.Post("/api/v1/comment/", body, "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result comment.CommentResponse
	DecodeTo(&result, resp)
	if result.IsAnonymous != true {
		t.Errorf("expected is_anonymous true, got %v", result.IsAnonymous)
	}
	if result.UserID != nil {
		t.Errorf("expected user_id omitted for anonymous when not super user, got %v", result.UserID)
	}
}

func TestGetComment(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	postID := uuid.New()
	post := models.Post{ID: postID, AuthorId: uuid.New(), Title: "Test Post", Content: "Test content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	commentDB := comment.NewCommentDB(testDB.DB)
	userID := uuid.New()
	c := &models.Comment{UserID: userID, PostID: postID, Description: "Get me"}
	created, err := commentDB.CreateComment(c)
	if err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	resp := api.Get("/api/v1/comment/"+created.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result comment.CommentResponse
	DecodeTo(&result, resp)
	if result.ID != created.ID || result.Description != "Get me" {
		t.Errorf("expected same comment, got %+v", result)
	}
}

func TestGetCommentsByPost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	postID := uuid.New()
	post := models.Post{ID: postID, AuthorId: uuid.New(), Title: "Test Post", Content: "Test content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	commentDB := comment.NewCommentDB(testDB.DB)
	userID := uuid.New()
	for _, desc := range []string{"First", "Second"} {
		c := &models.Comment{UserID: userID, PostID: postID, Description: desc}
		if _, err := commentDB.CreateComment(c); err != nil {
			t.Fatalf("failed to create comment: %v", err)
		}
	}

	resp := api.Get("/api/v1/post/"+postID.String()+"/comments", "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result []comment.CommentResponse
	DecodeTo(&result, resp)
	if len(result) < 2 {
		t.Errorf("expected at least 2 comments, got %d", len(result))
	}
}

func TestGetReplies(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	postID := uuid.New()
	post := models.Post{ID: postID, AuthorId: uuid.New(), Title: "Test Post", Content: "Test content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	commentDB := comment.NewCommentDB(testDB.DB)
	userID := uuid.New()
	parent := &models.Comment{UserID: userID, PostID: postID, Description: "Parent"}
	createdParent, err := commentDB.CreateComment(parent)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}
	reply := &models.Comment{UserID: userID, PostID: postID, ParentCommentID: &createdParent.ID, Description: "Reply"}
	if _, err := commentDB.CreateComment(reply); err != nil {
		t.Fatalf("failed to create reply: %v", err)
	}

	resp := api.Get("/api/v1/comment/"+createdParent.ID.String()+"/replies", "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result []comment.CommentResponse
	DecodeTo(&result, resp)
	if len(result) != 1 {
		t.Errorf("expected 1 reply, got %d", len(result))
	}
	if len(result) > 0 && result[0].Description != "Reply" {
		t.Errorf("expected reply description Reply, got %s", result[0].Description)
	}
}

func TestUpdateComment(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	postID := uuid.New()
	post := models.Post{ID: postID, AuthorId: uuid.New(), Title: "Test Post", Content: "Test content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	commentDB := comment.NewCommentDB(testDB.DB)
	userID := uuid.New()
	c := &models.Comment{UserID: userID, PostID: postID, Description: "Original"}
	created, err := commentDB.CreateComment(c)
	if err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	updateBody := map[string]any{"description": "Updated"}
	resp := api.Patch("/api/v1/comment/"+created.ID.String(), updateBody, "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result comment.CommentResponse
	DecodeTo(&result, resp)
	if result.Description != "Updated" {
		t.Errorf("expected description Updated, got %s", result.Description)
	}
}

func TestDeleteComment(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	postID := uuid.New()
	post := models.Post{ID: postID, AuthorId: uuid.New(), Title: "Test Post", Content: "Test content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	commentDB := comment.NewCommentDB(testDB.DB)
	userID := uuid.New()
	c := &models.Comment{UserID: userID, PostID: postID, Description: "To delete"}
	created, err := commentDB.CreateComment(c)
	if err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	resp := api.Delete("/api/v1/comment/"+created.ID.String(), "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	getResp := api.Get("/api/v1/comment/"+created.ID.String(), "Authorization: Bearer mock-token")
	if getResp.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getResp.Code)
	}
}

func TestCreateReplyToReplyReturns400(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	postID := uuid.New()
	post := models.Post{ID: postID, AuthorId: uuid.New(), Title: "Test Post", Content: "Test content"}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	commentDB := comment.NewCommentDB(testDB.DB)
	userID := uuid.New()
	parent := &models.Comment{UserID: userID, PostID: postID, Description: "Parent"}
	createdParent, err := commentDB.CreateComment(parent)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}
	reply := &models.Comment{UserID: userID, PostID: postID, ParentCommentID: &createdParent.ID, Description: "Reply"}
	createdReply, err := commentDB.CreateComment(reply)
	if err != nil {
		t.Fatalf("failed to create reply: %v", err)
	}

	body := map[string]any{
		"user_id":           userID.String(),
		"post_id":           postID.String(),
		"parent_comment_id": createdReply.ID.String(),
		"description":       "Reply to reply",
		"is_anonymous":      false,
	}
	resp := api.Post("/api/v1/comment/", body, "Authorization: Bearer mock-token")
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for reply-to-reply (one layer only), got %d: %s", resp.Code, resp.Body.String())
	}
}
