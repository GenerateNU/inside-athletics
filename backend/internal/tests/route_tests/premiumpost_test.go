package routeTests

import (
	premiumpost "inside-athletics/internal/handlers/premium_post"
	"inside-athletics/internal/models"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

// Test creating a premium post
func TestCreatePremiumPost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB)
	api := testDB.API

	CreateUserAndSport(testDB, t)

	// assigning admin user to JohnID
	assignRoleToUser(t, testDB.DB, JohnID, getRoleID(t, testDB.DB, models.RoleAdmin))

	body := map[string]any{
		"sport_id":     SoccerID,
		"title":        "Looking for thoughts on NEU Fencing!",
		"content":      "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?",
		"is_anonymous": false,
		"tags":         []map[string]any{},
	}

	resp := api.Post("/api/v1/premiumpost/", body, authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result premiumpost.CreatePremiumPostResponse
	DecodeTo(&result, resp)

	if result.AuthorID == nil || uuidDereference(result.AuthorID) != JohnID {
		t.Errorf("expected authorID %v, got %v", JohnID, result.AuthorID)
	}

	if result.SportID == nil || uuidDereference(result.SportID) != SoccerID {
		t.Errorf("expected sportID %v, got %v", SoccerID, result.SportID)
	}

	if result.Title != "Looking for thoughts on NEU Fencing!" {
		t.Errorf("expected title %q, got %q", "Looking for thoughts on NEU Fencing!", result.Title)
	}

	if result.Content != "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" {
		t.Errorf("expected content %q, got %q", "My name is Bob Joe...", result.Content)
	}
}

// test that regular users cannot create a premium post
func TestCreatePremiumPost_UserForbidden(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB)
	api := testDB.API

	CreateUserAndSport(testDB, t)

	assignRoleToUser(t, testDB.DB, JohnID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{}, JohnID)

	body := map[string]any{
		"sport_id": SoccerID,
		"title":    "Premium Post",
		"content":  "Testing permissions this shouldnt work!!",
	}

	resp := api.Post("/api/v1/premiumpost/", body, authHeader)

	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for regular user, got %d", resp.Code)
	}
}

// test retrieving premium posts by author
func TestGetPremiumPostByAuthorId(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB)
	api := testDB.API
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	CreateUserAndSport(testDB, t)

	assignRoleToUser(t, testDB.DB, JohnID, getRoleID(t, testDB.DB, models.RoleAdmin))

	post1, err1 := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "First Post About Fencing",
		Content:  "This is the first post content",
	})
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	_, err2 := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "Second Post About Basketball",
		Content:  "This is the second post content",
	})
	if err2 != nil {
		t.Fatalf("failed to create post 2: %v", err2)
	}

	resp := api.Get("/api/v1/posts/premium/by-author/"+JohnID.String(), authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result premiumpost.GetPremiumPostsByAuthorIDResponse
	DecodeTo(&result, resp)

	if result.Total < 2 {
		t.Errorf("expected at least 2 posts, got %d", result.Total)
	}

	if len(result.Posts) < 2 {
		t.Errorf("expected at least 2 posts in response, got %d", len(result.Posts))
	}

	if result.Posts[0].ID != post1.ID && result.Posts[1].ID != post1.ID {
		t.Errorf("expected first created post %v to appear in results", post1.ID)
	}
}

// test retrieving premium posts by sport
func TestGetPremiumPostBySportId(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB)
	api := testDB.API
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	CreateUserAndSport(testDB, t)
	assignRoleToUser(t, testDB.DB, JohnID, getRoleID(t, testDB.DB, models.RoleAdmin))

	_, err1 := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "Sport-filtered Post 1",
		Content:  "First post for sport filter",
	})
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	_, err2 := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "Sport-filtered Post 2",
		Content:  "Second post for sport filter",
	})
	if err2 != nil {
		t.Fatalf("failed to create post 2: %v", err2)
	}

	resp := api.Get("/api/v1/posts/premium/by-sport/"+SoccerID.String(), authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result premiumpost.GetPremiumPostsBySportIDResponse
	DecodeTo(&result, resp)

	if result.Total < 2 {
		t.Errorf("expected at least 2 posts, got %d", result.Total)
	}

	if len(result.Posts) < 2 {
		t.Errorf("expected at least 2 posts in response, got %d", len(result.Posts))
	}
}

// test retrieving premium post by tag
func TestGetPremiumPostByTagId(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB)
	api := testDB.API
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	CreateUserAndSport(testDB, t)
	assignRoleToUser(t, testDB.DB, JohnID, getRoleID(t, testDB.DB, models.RoleAdmin))

	// creating health and wellness tag
	tag := models.Tag{
		ID:   HealthAndWellnessID,
		Name: "Health And Wellness",
	}
	if err := testDB.DB.Create(&tag).Error; err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	post1, err1 := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "Tag-filtered Post 1",
		Content:  "First post for tag filter",
	})
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	if err := testDB.DB.Create(&models.TagPost{
		ID:           uuid.New(),
		PostableID:   post1.ID,
		PostableType: "premium_post",
		TagID:        HealthAndWellnessID,
	}).Error; err != nil {
		t.Fatalf("failed to create tag_post for post 1: %v", err)
	}

	post2, err2 := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "Tag-filtered Post 2",
		Content:  "Second post for tag filter",
	})
	if err2 != nil {
		t.Fatalf("failed to create post 2: %v", err2)
	}

	if err := testDB.DB.Create(&models.TagPost{
		ID:           uuid.New(),
		PostableID:   post2.ID,
		PostableType: "premium_post",
		TagID:        HealthAndWellnessID,
	}).Error; err != nil {
		t.Fatalf("failed to create tag_post for post 2: %v", err)
	}

	resp := api.Get("/api/v1/posts/premium/by-tag/"+HealthAndWellnessID.String(), authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result premiumpost.GetPremiumPostsByTagIDResponse
	DecodeTo(&result, resp)

	if result.Total < 2 {
		t.Errorf("expected at least 2 posts, got %d", result.Total)
	}

	if len(result.Posts) < 2 {
		t.Errorf("expected at least 2 posts in response, got %d", len(result.Posts))
	}
}

// test retrieving all premium posts
func TestGetAllPremiumPosts(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB)
	api := testDB.API
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	CreateUserAndSport(testDB, t)
	assignRoleToUser(t, testDB.DB, JohnID, getRoleID(t, testDB.DB, models.RoleAdmin))

	_, err1 := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "All premium posts — first",
		Content:  "First body",
	})
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	_, err2 := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "All premium posts — second",
		Content:  "Second body",
	})
	if err2 != nil {
		t.Fatalf("failed to create post 2: %v", err2)
	}

	resp := api.Get("/api/v1/posts/premium/", authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result premiumpost.GetAllPremiumPostsResponse
	DecodeTo(&result, resp)

	if result.Total < 2 {
		t.Errorf("expected at least 2 posts, got %d", result.Total)
	}
	if len(result.Posts) < 2 {
		t.Errorf("expected at least 2 posts in response, got %d", len(result.Posts))
	}
}
