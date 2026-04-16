package routeTests

import (
	"inside-athletics/internal/handlers/comment"
	"inside-athletics/internal/handlers/post"
	"inside-athletics/internal/models"
	"net/http"
	"strconv"
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
	popularity := int32(100)
	soccer := models.Sport{
		ID:         SoccerID,
		Name:       "Soccer",
		Popularity: &popularity,
	}
	post := models.Post{
		AuthorID: user.ID,
		SportID:  &SoccerID,
		Title:    "Test Post",
		Content:  "Test content",
	}

	if err := testDB.DB.FirstOrCreate(&soccer).Error; err != nil {
		t.Fatalf("failed to create sport: %v", err)
	}
	if err := testDB.DB.FirstOrCreate(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}
	return user, post
}

func TestCreateComment(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, post := seedUserAndPost(t, testDB, "create-comment")

	body := map[string]any{
		"post_id":      post.ID.String(),
		"description":  "A test comment",
		"is_anonymous": false,
	}

	resp := api.Post("/api/v1/comment/", body, "Authorization: Bearer "+user.ID.String())
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result comment.CreateCommentResponse
	DecodeTo(&result, resp)
	if result.Description != "A test comment" {
		t.Errorf("expected description 'A test comment', got %s", result.Description)
	}
	if result.PostID != post.ID {
		t.Errorf("expected post_id %s, got %s", post.ID, result.PostID)
	}
}

// Asserts anonymous comments hide user_id when caller is not the user who made the comment.
func TestCreateCommentAnonymous(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, post := seedUserAndPost(t, testDB, "create-anon")

	body := map[string]any{
		"post_id":      post.ID.String(),
		"description":  "Anonymous comment",
		"is_anonymous": true,
	}

	resp := api.Post("/api/v1/comment/", body, "Authorization: Bearer "+user.ID.String())
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result comment.CreateCommentResponse
	DecodeTo(&result, resp)
	if result.IsAnonymous != true {
		t.Errorf("expected is_anonymous true, got %v", result.IsAnonymous)
	}
}

// testing get anonymous comment, when user is not user who made comment, will not return user id
// also testing when user who made comment gets it, userId is shown
func TestGetComment(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, post := seedUserAndPost(t, testDB, "get-comment")
	commentDB := comment.NewCommentDB(testDB.DB)
	c := &models.Comment{UserID: user.ID, PostID: post.ID, Description: "Get me", IsAnonymous: true}
	created, err := commentDB.CreateComment(c)
	if err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	resp := api.Get("/api/v1/comment/"+created.ID.String(), "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}
	var result comment.CommentResponse
	DecodeTo(&result, resp)
	if result.ID != created.ID || result.Description != "Get me" {
		t.Errorf("expected same comment, got %+v", result)
	}
	if result.User != nil {
		t.Errorf("expected nil UserId for Anonymous")
	}

	// Free-tier users must have accessed the post before viewing its comments.
	_ = api.Get("/api/v1/post/"+post.ID.String(), "Authorization: Bearer "+user.ID.String())

	resp2 := api.Get("/api/v1/comment/"+created.ID.String(), "Authorization: Bearer "+user.ID.String())
	if resp2.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp2.Code, resp2.Body.String())
	}
	var result2 comment.CommentResponse
	DecodeTo(&result2, resp2)
	if result2.ID != created.ID || result2.Description != "Get me" {
		t.Errorf("expected same comment, got %+v", result2)
	}
	if result2.User.ID != user.ID {
		t.Errorf("expected UserId for Anonymous")
	}

}

func TestGetCommentWithLikes(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, post := seedUserAndPost(t, testDB, "get-comment-2")
	commentDB := comment.NewCommentDB(testDB.DB)
	c := &models.Comment{UserID: user.ID, PostID: post.ID, Description: "Get me", IsAnonymous: true}
	created, err := commentDB.CreateComment(c)
	if err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	l1 := &models.CommentLike{UserID: user.ID, CommentID: c.ID}
	testDB.DB.Create(&l1)

	resp := api.Get("/api/v1/comment/"+created.ID.String(), "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}
	var result comment.CommentResponse
	DecodeTo(&result, resp)
	if result.ID != created.ID || result.Description != "Get me" {
		t.Errorf("expected same comment, got %+v", result)
	}
	if result.User != nil {
		t.Errorf("expected nil UserId for Anonymous")
	}

	// Free-tier users must have accessed the post before viewing its comments.
	_ = api.Get("/api/v1/post/"+post.ID.String(), "Authorization: Bearer "+user.ID.String())

	resp2 := api.Get("/api/v1/comment/"+created.ID.String(), "Authorization: Bearer "+user.ID.String())
	if resp2.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp2.Code, resp2.Body.String())
	}
	var result2 comment.CommentResponse
	DecodeTo(&result2, resp2)
	if result2.ID != created.ID || result2.Description != "Get me" {
		t.Errorf("expected same comment, got %+v", result2)
	}
	if result2.User.ID != user.ID {
		t.Errorf("expected UserId for Anonymous")
	}

	if result2.LikeCount != 1 {
		t.Errorf("expected 1 like got, %+v", result2.LikeCount)
	}

	if result2.IsLiked != true {
		t.Errorf("expected comment to be liked by user got, %+v", result2.IsLiked)
	}

}

func TestGetCommentsByPost(t *testing.T) {
	t.Parallel()
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

	resp := api.Get("/api/v1/post/"+post.ID.String()+"/comments", "Authorization: Bearer "+mockUUID)
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
	t.Parallel()
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

	resp := api.Get("/api/v1/comment/"+createdParent.ID.String()+"/replies", "Authorization: Bearer "+mockUUID)
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
	t.Parallel()
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
	resp := api.Patch("/api/v1/comment/"+created.ID.String(), updateBody, "Authorization: Bearer "+mockUUID)
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
	t.Parallel()
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

	resp := api.Delete("/api/v1/comment/"+created.ID.String(), "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	getResp := api.Get("/api/v1/comment/"+created.ID.String(), "Authorization: Bearer "+mockUUID)
	if getResp.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getResp.Code)
	}
}

func TestCreateReplyToReplyReturns400(t *testing.T) {
	t.Parallel()
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
		"post_id":           post.ID.String(),
		"parent_comment_id": createdReply.ID.String(),
		"description":       "Reply to reply",
		"is_anonymous":      false,
	}
	resp := api.Post("/api/v1/comment/", body, "Authorization: Bearer "+user.ID.String())
	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for reply-to-reply (one layer only), got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestFreeUserGetCommentsByPostLimitedToFirstThreeViewedPosts(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	viewer := newCommentTestUser(uuid.New(), "viewer-comments-limit")
	if err := testDB.DB.Create(&viewer).Error; err != nil {
		t.Fatalf("failed to create viewer user: %v", err)
	}

	author := newCommentTestUser(uuid.New(), "author-comments-limit")
	if err := testDB.DB.Create(&author).Error; err != nil {
		t.Fatalf("failed to create author user: %v", err)
	}

	popularity := int32(100)
	soccer := models.Sport{ID: SoccerID, Name: "Soccer", Popularity: &popularity}
	if err := testDB.DB.FirstOrCreate(&soccer).Error; err != nil {
		t.Fatalf("failed to create sport: %v", err)
	}

	postDB := post.NewPostDB(testDB.DB)
	commentDB := comment.NewCommentDB(testDB.DB)
	var postIDs []uuid.UUID
	for i := 0; i < 4; i++ {
		p, err := postDB.CreatePost(&models.Post{
			AuthorID: author.ID,
			SportID:  &SoccerID,
			Title:    "Comments gate post " + strconv.Itoa(i),
			Content:  "content",
		}, []post.TagRequest{})
		if err != nil {
			t.Fatalf("failed to create post %d: %v", i, err)
		}
		postIDs = append(postIDs, p.ID)
		if _, err := commentDB.CreateComment(&models.Comment{
			UserID:      author.ID,
			PostID:      p.ID,
			Description: "comment " + strconv.Itoa(i),
		}); err != nil {
			t.Fatalf("failed to create comment for post %d: %v", i, err)
		}
	}

	authHeader := "Authorization: Bearer " + viewer.ID.String()
	// Record four viewed posts for the free user.
	for i := range 4 {
		_ = api.Get("/api/v1/post/"+postIDs[i].String(), authHeader)
	}

	for i := range 3 {
		resp := api.Get("/api/v1/post/"+postIDs[i].String()+"/comments", authHeader)
		if resp.Code != http.StatusOK {
			t.Fatalf("expected 200 for post %d comments, got %d: %s", i, resp.Code, resp.Body.String())
		}
	}

	resp := api.Get("/api/v1/post/"+postIDs[3].String()+"/comments", authHeader)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for 4th viewed post comments, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestFreeUserGetCommentBlockedOutsideFirstThreeViewedPosts(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	viewer := newCommentTestUser(uuid.New(), "viewer-comment-single")
	if err := testDB.DB.Create(&viewer).Error; err != nil {
		t.Fatalf("failed to create viewer user: %v", err)
	}

	author := newCommentTestUser(uuid.New(), "author-comment-single")
	if err := testDB.DB.Create(&author).Error; err != nil {
		t.Fatalf("failed to create author user: %v", err)
	}

	popularity := int32(100)
	soccer := models.Sport{ID: SoccerID, Name: "Soccer", Popularity: &popularity}
	if err := testDB.DB.FirstOrCreate(&soccer).Error; err != nil {
		t.Fatalf("failed to create sport: %v", err)
	}

	postDB := post.NewPostDB(testDB.DB)
	commentDB := comment.NewCommentDB(testDB.DB)
	var comments []models.Comment
	var postIDs []uuid.UUID
	for i := 0; i < 4; i++ {
		p, err := postDB.CreatePost(&models.Post{
			AuthorID: author.ID,
			SportID:  &SoccerID,
			Title:    "Single comment gate post " + strconv.Itoa(i),
			Content:  "content",
		}, []post.TagRequest{})
		if err != nil {
			t.Fatalf("failed to create post %d: %v", i, err)
		}
		postIDs = append(postIDs, p.ID)
		c, err := commentDB.CreateComment(&models.Comment{
			UserID:      author.ID,
			PostID:      p.ID,
			Description: "comment " + strconv.Itoa(i),
		})
		if err != nil {
			t.Fatalf("failed to create comment for post %d: %v", i, err)
		}
		comments = append(comments, *c)
	}

	authHeader := "Authorization: Bearer " + viewer.ID.String()
	for i := 0; i < 4; i++ {
		_ = api.Get("/api/v1/post/"+postIDs[i].String(), authHeader)
	}

	respAllowed := api.Get("/api/v1/comment/"+comments[0].ID.String(), authHeader)
	if respAllowed.Code != http.StatusOK {
		t.Fatalf("expected 200 for comment on first viewed post, got %d: %s", respAllowed.Code, respAllowed.Body.String())
	}

	respBlocked := api.Get("/api/v1/comment/"+comments[3].ID.String(), authHeader)
	if respBlocked.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for comment on 4th viewed post, got %d: %s", respBlocked.Code, respBlocked.Body.String())
	}
}

func TestPremiumUserCanViewCommentsOutsideFirstThreeViewedPosts(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	viewer := newCommentTestUser(uuid.New(), "viewer-premium-comments")
	viewer.Account_Type = true
	if err := testDB.DB.Create(&viewer).Error; err != nil {
		t.Fatalf("failed to create premium viewer user: %v", err)
	}

	author := newCommentTestUser(uuid.New(), "author-premium-comments")
	if err := testDB.DB.Create(&author).Error; err != nil {
		t.Fatalf("failed to create author user: %v", err)
	}

	popularity := int32(100)
	soccer := models.Sport{ID: SoccerID, Name: "Soccer", Popularity: &popularity}
	if err := testDB.DB.FirstOrCreate(&soccer).Error; err != nil {
		t.Fatalf("failed to create sport: %v", err)
	}

	postDB := post.NewPostDB(testDB.DB)
	commentDB := comment.NewCommentDB(testDB.DB)
	var comments []models.Comment
	var postIDs []uuid.UUID
	for i := 0; i < 4; i++ {
		p, err := postDB.CreatePost(&models.Post{
			AuthorID: author.ID,
			SportID:  &SoccerID,
			Title:    "Premium comments post " + strconv.Itoa(i),
			Content:  "content",
		}, []post.TagRequest{})
		if err != nil {
			t.Fatalf("failed to create post %d: %v", i, err)
		}
		postIDs = append(postIDs, p.ID)
		c, err := commentDB.CreateComment(&models.Comment{
			UserID:      author.ID,
			PostID:      p.ID,
			Description: "comment " + strconv.Itoa(i),
		})
		if err != nil {
			t.Fatalf("failed to create comment for post %d: %v", i, err)
		}
		comments = append(comments, *c)
	}

	authHeader := "Authorization: Bearer " + viewer.ID.String()
	// Even after viewing four posts, premium users should not be blocked.
	for i := range 4 {
		_ = api.Get("/api/v1/post/"+postIDs[i].String(), authHeader)
	}

	// Premium users should still be able to read comment lists and single comments on the 4th post.
	respList := api.Get("/api/v1/post/"+postIDs[3].String()+"/comments", authHeader)
	if respList.Code != http.StatusOK {
		t.Fatalf("expected 200 for premium user comments on 4th viewed post, got %d: %s", respList.Code, respList.Body.String())
	}

	respSingle := api.Get("/api/v1/comment/"+comments[3].ID.String(), authHeader)
	if respSingle.Code != http.StatusOK {
		t.Fatalf("expected 200 for premium user single comment on 4th viewed post, got %d: %s", respSingle.Code, respSingle.Body.String())
	}
}
