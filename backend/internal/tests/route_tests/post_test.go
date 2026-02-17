package routeTests

import (
	"inside-athletics/internal/handlers/post"
	"inside-athletics/internal/models"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestCreatePost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API

	authorID := uuid.New()

	popularity := int32(100000)

	sport := map[string]any{
		"name":       "Women's Basketball",
		"popularity": popularity,
	}

	resp_sport := api.Post("/api/v1/sport/", sport, "Authorization: Bearer mock-token")
	if resp_sport.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp_sport.Code, resp_sport.Body.String())
	}

	var createdSport models.Sport
	DecodeTo(&createdSport, resp_sport)

	sportID := createdSport.ID

	body := map[string]any{
		"author_id":    authorID,
		"sport_id":     sportID,
		"title":        "Looking for thoughts on NEU Fencing!",
		"content":      "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?",
		"is_anonymous": true,
	}

	resp := api.Post("/api/v1/post/", body, "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post.PostResponse
	DecodeTo(&result, resp)

	if result.AuthorId != authorID {
		t.Errorf("expected authorID %v, got %v", authorID, result.AuthorId)
	}

	if result.SportId != sportID {
		t.Errorf("expected sportID %v, got %v", sportID, result.SportId)
	}

	if result.Title != "Looking for thoughts on NEU Fencing!" {
		t.Errorf("expected title %q, got %q", "Looking for thoughts on NEU Fencing!", result.Title)
	}

	if result.Content != "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" {
		t.Errorf("expected content %q, got %q", "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?", result.Content)
	}

	if result.Likes != 0 {
		t.Errorf("expected UpVotes 0, got %d", result.Likes)
	}

	if result.IsAnonymous != true {
		t.Errorf("expected IsAnonymous %v, got %v", true, result.IsAnonymous)
	}
}

func TestGetPostById(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API
	postDB := post.NewPostDB(testDB.DB)

	authorID := uuid.New()

	popularity := int32(100000)

	sport := map[string]any{
		"name":       "Women's Basketball",
		"popularity": popularity,
	}

	resp_sport := api.Post("/api/v1/sport/", sport, "Authorization: Bearer mock-token")
	if resp_sport.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp_sport.Code, resp_sport.Body.String())
	}

	var createdSport models.Sport
	DecodeTo(&createdSport, resp_sport)

	sportID := createdSport.ID

	createdPost, err := postDB.CreatePost(
		authorID,
		sportID,
		"Looking for thoughts on NEU Fencing!",
		"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?", // content
		true,
	)

	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	resp := api.Get("/api/v1/post/"+createdPost.ID.String(), "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post.PostResponse
	DecodeTo(&result, resp)

	if result.AuthorId != authorID {
		t.Errorf("expected authorID %v, got %v", authorID, result.AuthorId)
	}

	if result.SportId != sportID {
		t.Errorf("expected sportID %v, got %v", sportID, result.SportId)
	}

	if result.Title != "Looking for thoughts on NEU Fencing!" {
		t.Errorf("expected title %q, got %q", "Looking for thoughts on NEU Fencing!", result.Title)
	}

	if result.Content != "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" {
		t.Errorf("expected content %q, got %q", "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?", result.Content)
	}

	if result.Likes != 0 {
		t.Errorf("expected UpVotes 0, got %d", result.Likes)
	}

	if result.IsAnonymous != true {
		t.Errorf("expected IsAnonymous %v, got %v", true, result.IsAnonymous)
	}
}

func TestBadValidation(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Post{}); err != nil {
		t.Fatalf("failed to migrate posts table: %v", err)
	}

	post.Route(testDB.API, testDB.DB)
	api := testDB.API

	resp := api.Get("/api/v1/post/"+"random string", "Authorization: Bearer mock-token")

	if resp.Code == http.StatusOK {
		t.Fatalf("expected status 422, got %d: %s", resp.Code, resp.Body.String())
	}

	resp = api.Get("/api/v1/posts/by-sport/"+"random string", "Authorization: Bearer mock-token")

	if resp.Code == http.StatusOK {
		t.Fatalf("expected status 422, got %d: %s", resp.Code, resp.Body.String())
	}

	resp = api.Get("/api/v1/posts/by-author/"+"random string", "Authorization: Bearer mock-token")

	if resp.Code == http.StatusOK {
		t.Fatalf("expected status 422, got %d: %s", resp.Code, resp.Body.String())
	}

}
func TestGetPostByAuthorId(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Post{}); err != nil {
		t.Fatalf("failed to migrate posts table: %v", err)
	}

	post.Route(testDB.API, testDB.DB)
	api := testDB.API
	postDB := post.NewPostDB(testDB.DB)

	// Create a sport first
	popularity := int32(100000)
	sport := map[string]any{
		"name":       "Women's Basketball",
		"popularity": popularity,
	}
	resp_sport := api.Post("/api/v1/sport/", sport, "Authorization: Bearer mock-token")
	if resp_sport.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp_sport.Code, resp_sport.Body.String())
	}
	var createdSport models.Sport
	DecodeTo(&createdSport, resp_sport)

	// Create two posts
	authorID := uuid.New()
	sportID := uuid.New()

	_, err1 := postDB.CreatePost(
		authorID,
		sportID,
		"First Post About Fencing",
		"This is the first post content",
		false,
	)
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	_, err2 := postDB.CreatePost(
		authorID,
		sportID,
		"Second Post About Basketball",
		"This is the second post content",
		true,
	)
	if err2 != nil {
		t.Fatalf("failed to create post 2: %v", err2)
	}

	resp := api.Get("/api/v1/posts/by-author/" + authorID.String(), "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post.GetPostsBySportIDResponse
	DecodeTo(&result, resp)

	if result.Total < 2 {
		t.Errorf("expected at least 2 posts, got %d", result.Total)
	}

	if len(result.Posts) < 2 {
		t.Errorf("expected at least 2 posts in response, got %d", len(result.Posts))
	}
}

func TestGetPostsBySportId(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API
	postDB := post.NewPostDB(testDB.DB)

	// Create a sport first
	popularity := int32(100000)
	sport := map[string]any{
		"name":       "Women's Basketball",
		"popularity": popularity,
	}
	resp_sport := api.Post("/api/v1/sport/", sport, "Authorization: Bearer mock-token")
	if resp_sport.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp_sport.Code, resp_sport.Body.String())
	}
	var createdSport models.Sport
	DecodeTo(&createdSport, resp_sport)
	sportID := createdSport.ID

	// Create two posts
	authorID1 := uuid.New()
	authorID2 := uuid.New()

	_, err1 := postDB.CreatePost(
		authorID1,
		sportID,
		"First Post About Fencing",
		"This is the first post content",
		false,
	)
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	_, err2 := postDB.CreatePost(
		authorID2,
		sportID,
		"Second Post About Basketball",
		"This is the second post content",
		true,
	)
	if err2 != nil {
		t.Fatalf("failed to create post 2: %v", err2)
	}

	resp := api.Get("/api/v1/posts/by-sport/" + sportID.String(), "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post.GetPostsBySportIDResponse
	DecodeTo(&result, resp)

	if result.Total < 2 {
		t.Errorf("expected at least 2 posts, got %d", result.Total)
	}

	if len(result.Posts) < 2 {
		t.Errorf("expected at least 2 posts in response, got %d", len(result.Posts))
	}
}

func TestGetAllPosts(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API
	postDB := post.NewPostDB(testDB.DB)

	// Create a sport first
	popularity := int32(100000)
	sport := map[string]any{
		"name":       "Women's Basketball",
		"popularity": popularity,
	}
	resp_sport := api.Post("/api/v1/sport/", sport, "Authorization: Bearer mock-token")
	if resp_sport.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp_sport.Code, resp_sport.Body.String())
	}
	var createdSport models.Sport
	DecodeTo(&createdSport, resp_sport)
	sportID := createdSport.ID

	// Create two posts
	authorID1 := uuid.New()
	authorID2 := uuid.New()

	_, err1 := postDB.CreatePost(
		authorID1,
		sportID,
		"First Post About Fencing",
		"This is the first post content",
		false,
	)
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	_, err2 := postDB.CreatePost(
		authorID2,
		sportID,
		"Second Post About Basketball",
		"This is the second post content",
		true,
	)
	if err2 != nil {
		t.Fatalf("failed to create post 2: %v", err2)
	}

	resp := api.Get("/api/v1/posts/", "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post.GetAllPostsResponse
	DecodeTo(&result, resp)

	if result.Total < 2 {
		t.Errorf("expected at least 2 posts, got %d", result.Total)
	}

	if len(result.Posts) < 2 {
		t.Errorf("expected at least 2 posts in response, got %d", len(result.Posts))
	}
}

func TestUpdatePost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API
	postDB := post.NewPostDB(testDB.DB)

	// Create a sport first
	popularity := int32(100000)
	sport := map[string]any{
		"name":       "Women's Basketball",
		"popularity": popularity,
	}
	resp_sport := api.Post("/api/v1/sport/", sport, "Authorization: Bearer mock-token")
	if resp_sport.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp_sport.Code, resp_sport.Body.String())
	}
	var createdSport models.Sport
	DecodeTo(&createdSport, resp_sport)
	sportID := createdSport.ID

	// Create a post
	authorID := uuid.New()
	createdPost, err := postDB.CreatePost(
		authorID,
		sportID,
		"Original Title",
		"Original content",
		false,
	)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Update the post
	updateBody := map[string]any{
		"title":   "Updated Title",
		"content": "Updated content about the program",
	}

	resp := api.Patch("/api/v1/post/"+createdPost.ID.String(), updateBody, "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post.PostResponse
	DecodeTo(&result, resp)

	if result.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %s", result.Title)
	}

	if result.Content != "Updated content about the program" {
		t.Errorf("expected content 'Updated content about the program', got %s", result.Content)
	}

	if result.IsAnonymous != false {
		t.Errorf("expected IsAnonymous false, got %v", result.IsAnonymous)
	}
}

func TestDeletePost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API
	postDB := post.NewPostDB(testDB.DB)

	// Create a sport first
	popularity := int32(100000)
	sport := map[string]any{
		"name":       "Women's Basketball",
		"popularity": popularity,
	}
	resp_sport := api.Post("/api/v1/sport/", sport, "Authorization: Bearer mock-token")
	if resp_sport.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp_sport.Code, resp_sport.Body.String())
	}
	var createdSport models.Sport
	DecodeTo(&createdSport, resp_sport)
	sportID := createdSport.ID

	// Create a post
	authorID := uuid.New()
	createdPost, err := postDB.CreatePost(
		authorID,
		sportID,
		"Post to Delete",
		"This post will be deleted",
		false,
	)
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Delete the post
	resp := api.Delete("/api/v1/post/"+createdPost.ID.String(), "Authorization: Bearer mock-token")

	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d: %s", resp.Code, resp.Body.String())
	}

	// Verify the post is deleted by trying to get it
	getResp := api.Get("/api/v1/post/"+createdPost.ID.String(), "Authorization: Bearer mock-token")
	if getResp.Code != 500 {
		t.Errorf("expected 404 after delete, got %d", getResp.Code)
	}
}
