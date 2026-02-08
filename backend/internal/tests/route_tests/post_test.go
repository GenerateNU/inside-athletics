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

	if err := testDB.DB.AutoMigrate(&models.Post{}); err != nil {
		t.Fatalf("failed to migrate posts table: %v", err)
	}

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

	if err := testDB.DB.AutoMigrate(&models.Post{}); err != nil {
		t.Fatalf("failed to migrate posts table: %v", err)
	}

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

func TestGetPostByAuthorId(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Post{}); err != nil {
		t.Fatalf("failed to migrate posts table: %v", err)
	}

	post.Route(testDB.API, testDB.DB)
	api := testDB.API
	postDB := post.NewPostDB(testDB.DB)

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

	// Generate a UUID for the user that will be used in the Authorization header
	userID := uuid.New()

	user := map[string]any{
		"first_name":              "Joe",
		"last_name":               "Bob",
		"email":                   "bobjoe123@email.com",
		"username":                "bjproathlete",
		"bio":                     "My name is Bob and I'm a pro athlete",
		"account_type":            true,
		"sport":                   []string{"hockey"},
		"expected_grad_year":      2027,
		"verified_athlete_status": "pending",
		"college":                 "Northeastern University",
		"division":                1,
	}

	resp_user := api.Post("/api/v1/user/", user, "Authorization: Bearer "+userID.String())
	if resp_user.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp_user.Code, resp_user.Body.String())
	}

	var createdUser models.User
	DecodeTo(&createdUser, resp_user)

	authorID := createdUser.ID

	createdPost, err := postDB.CreatePost(
		authorID,
		sportID,
		"Looking for thoughts on NEU Fencing!",
		"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?",
		true,
	)

	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	resp := api.Get("/api/v1/post/by-author/"+createdPost.AuthorId.String(), "Authorization: Bearer mock-token")

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
		t.Errorf("expected Likes 0, got %d", result.Likes)
	}
	if result.IsAnonymous != true {
		t.Errorf("expected IsAnonymous %v, got %v", true, result.IsAnonymous)
	}
}
