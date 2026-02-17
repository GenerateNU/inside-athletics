package routeTests

import (
	"inside-athletics/internal/handlers/comment"
	"inside-athletics/internal/models"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

// newCommentTestUser returns a User with only required fields for comment tests
func newCommentTestUser(id uuid.UUID, unique string) models.User {
	return models.User{
		ID:                      id,
		FirstName:               "Test",
		LastName:                "User",
		Email:                   "test-" + unique + "@example.com",
		Username:                "testuser-" + unique,
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
}

// seedUserAndPost creates a User and a Post (with that user as author) for comment tests
func seedUserAndPost(t *testing.T, testDB *TestDatabase, unique string) (models.User, models.Post) {
	t.Helper()
	user := newCommentTestUser(uuid.New(), unique)
	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	post := models.Post{
		AuthorId: user.ID,
		SportId:  uuid.New(),
		Title:    "Test Post",
		Content:  "Test content",
	}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	return user, post
}

func TestCreateComment(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, post := seedUserAndPost(t, testDB, "create-comment")

	body := map[string]any{
		"user_id":      user.ID.String(),
		"post_id":      post.ID.String(),
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
		t.Errorf("expected description 'A test comment', got %s", result.Description)
	}
	if result.PostID != post.ID {
		t.Errorf("expected post_id %s, got %s", post.ID, result.PostID)
	}
	if result.UserID == nil || *result.UserID != user.ID {
		t.Errorf("expected user_id for non-anonymous comment, got %v", result.UserID)
	}
}

// Asserts anonymous comments hide user_id when caller is not a super user (default context).
func TestCreateCommentAnonymous(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, post := seedUserAndPost(t, testDB, "create-anon")

	body := map[string]any{
		"user_id":      user.ID.String(),
		"post_id":      post.ID.String(),
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
	user, post := seedUserAndPost(t, testDB, "get-comment")
	commentDB := comment.NewCommentDB(testDB.DB)
	c := &models.Comment{UserID: user.ID, PostID: post.ID, Description: "Get me"}
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
	user, post := seedUserAndPost(t, testDB, "get-by-post")
	commentDB := comment.NewCommentDB(testDB.DB)
	for _, desc := range []string{"First", "Second"} {
		c := &models.Comment{UserID: user.ID, PostID: post.ID, Description: desc}
		if _, err := commentDB.CreateComment(c); err != nil {
			t.Fatalf("failed to create comment: %v", err)
		}
	}

	resp := api.Get("/api/v1/post/"+post.ID.String()+"/comments", "Authorization: Bearer mock-token")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	// checking for 2 comments in post
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
	user, post := seedUserAndPost(t, testDB, "get-replies")
	commentDB := comment.NewCommentDB(testDB.DB)
	parent := &models.Comment{UserID: user.ID, PostID: post.ID, Description: "Parent"}
	createdParent, err := commentDB.CreateComment(parent)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}
	reply := &models.Comment{UserID: user.ID, PostID: post.ID, ParentCommentID: &createdParent.ID, Description: "Reply"}
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
	user, post := seedUserAndPost(t, testDB, "update-comment")
	commentDB := comment.NewCommentDB(testDB.DB)
	c := &models.Comment{UserID: user.ID, PostID: post.ID, Description: "Original"}
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
	user, post := seedUserAndPost(t, testDB, "delete-comment")
	commentDB := comment.NewCommentDB(testDB.DB)
	c := &models.Comment{UserID: user.ID, PostID: post.ID, Description: "To delete"}
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
	user, post := seedUserAndPost(t, testDB, "reply-to-reply")
	commentDB := comment.NewCommentDB(testDB.DB)
	parent := &models.Comment{UserID: user.ID, PostID: post.ID, Description: "Parent"}
	createdParent, err := commentDB.CreateComment(parent)
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}
	reply := &models.Comment{UserID: user.ID, PostID: post.ID, ParentCommentID: &createdParent.ID, Description: "Reply"}
	createdReply, err := commentDB.CreateComment(reply)
	if err != nil {
		t.Fatalf("failed to create reply: %v", err)
	}

	body := map[string]any{
		"user_id":           user.ID.String(),
		"post_id":           post.ID.String(),
		"parent_comment_id": createdReply.ID.String(),
		"description":       "Reply to reply",
		"is_anonymous":      false,
	}
	resp := api.Post("/api/v1/comment/", body, "Authorization: Bearer mock-token")
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for reply-to-reply (one layer only), got %d: %s", resp.Code, resp.Body.String())
	}
}
