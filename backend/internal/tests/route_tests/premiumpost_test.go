package routeTests

import (
	premiumpost "inside-athletics/internal/handlers/premium_post"
	"inside-athletics/internal/models"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func premiumPostAdminAuthHeader(t *testing.T, testDB *TestDatabase) string {
	t.Helper()
	return authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
		{Action: models.PermissionCreate, Resource: "post"},
		{Action: models.PermissionCreate, Resource: "premiumpost"},
		{Action: models.PermissionUpdate, Resource: "premiumpost"},
		{Action: models.PermissionUpdateOwn, Resource: "premiumpost"},
		{Action: models.PermissionDelete, Resource: "premiumpost"},
		{Action: models.PermissionDeleteOwn, Resource: "premiumpost"},
	}, JohnID)
}

// newPremiumPostCreateBody returns a valid create body
func newPremiumPostCreateBody(t *testing.T, testDB *TestDatabase, title, content string) map[string]any {
	t.Helper()
	tag1 := models.Tag{Name: "recruiting"}
	tag2 := models.Tag{Name: "fencing"}
	testDB.DB.Create(&tag1)
	testDB.DB.Create(&tag2)
	college := models.College{
		ID:   NortheasternID,
		Name: "Northeastern University",
	}
	if err := testDB.DB.Create(&college).Error; err != nil {
		t.Fatalf("failed to create college: %v", err)
	}
	return map[string]any{
		"sport_id":   SoccerID,
		"college_id": NortheasternID,
		"tag": []string{
			tag1.ID.String(),
			tag2.ID.String(),
		},
		"title":   title,
		"content": content,
	}
}

// Test creating a premium post
func TestCreatePremiumPost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)

	tag1 := models.Tag{Name: "recruiting"}
	tag2 := models.Tag{Name: "fencing"}
	testDB.DB.Create(&tag1)
	testDB.DB.Create(&tag2)

	college := models.College{
		ID:   NortheasternID,
		Name: "Northeastern University",
	}
	if err := testDB.DB.Create(&college).Error; err != nil {
		t.Fatalf("failed to create college: %v", err)
	}

	body := map[string]any{
		"sport_id":   SoccerID,
		"college_id": NortheasternID,
		"tag": []string{
			tag1.ID.String(),
			tag2.ID.String(),
		},
		"title":   "Post with tags",
		"content": "Content content make content description.",
	}

	resp := api.Post("/api/v1/post/premium/", body, authHeader)
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

	if result.Title != "Post with tags" {
		t.Errorf("expected title %q, got %q", "Post with tags", result.Title)
	}

	if result.Content != "Content content make content description." {
		t.Errorf("expected content %q, got %q", "Content content make content description.", result.Content)
	}
}

// test that regular users cannot create a premium post
func TestCreatePremiumPost_UserForbidden(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)

	assignRoleToUser(t, testDB.DB, JohnID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUserForRole(t, testDB.DB, models.RoleUser, []permissionSpec{}, JohnID)

	tag1 := models.Tag{Name: "recruiting"}
	tag2 := models.Tag{Name: "fencing"}
	testDB.DB.Create(&tag1)
	testDB.DB.Create(&tag2)

	college := models.College{
		ID:   NortheasternID,
		Name: "Northeastern University",
	}
	if err := testDB.DB.Create(&college).Error; err != nil {
		t.Fatalf("failed to create college: %v", err)
	}

	body := map[string]any{
		"sport_id":   SoccerID,
		"college_id": NortheasternID,
		"tag": []string{
			tag1.ID.String(),
			tag2.ID.String(),
		},
		"title":   "Post with tags",
		"content": "Content content make content description.",
	}

	resp := api.Post("/api/v1/post/premium/", body, authHeader)

	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for regular user, got %d", resp.Code)
	}
}

// test retrieving premium posts by author
func TestGetPremiumPostByAuthorId(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	tag1 := models.Tag{Name: "recruiting"}
	tag2 := models.Tag{Name: "fencing"}
	testDB.DB.Create(&tag1)
	testDB.DB.Create(&tag2)

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

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	tag1 := models.Tag{Name: "recruiting"}
	tag2 := models.Tag{Name: "fencing"}
	testDB.DB.Create(&tag1)
	testDB.DB.Create(&tag2)

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

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

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

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	college := models.College{
		ID:   NortheasternID,
		Name: "Northeastern University",
	}
	if err := testDB.DB.Create(&college).Error; err != nil {
		t.Fatalf("failed to create college: %v", err)
	}

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

// test to expect a 422 when required JSON keys are missing
func TestCreatePremiumPostMissingRequiredJSONFields(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)

	body := map[string]any{
		"title":   "No sport/college/tag keys",
		"content": "Schema requires sport_id, college_id, and tag.",
	}

	resp := api.Post("/api/v1/post/premium/", body, authHeader)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d: %s", resp.Code, resp.Body.String())
	}
}

// test to assert an admin can create more than one premium post
func TestAdminCanCreateMultiplePremiumPosts(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)

	body := newPremiumPostCreateBody(t, testDB, "First premium", "First body.")
	resp1 := api.Post("/api/v1/post/premium/", body, authHeader)
	if resp1.Code != http.StatusOK {
		t.Fatalf("first post expected 200, got %d: %s", resp1.Code, resp1.Body.String())
	}

	body["title"] = "Second premium"
	body["content"] = "Second body."
	resp2 := api.Post("/api/v1/post/premium/", body, authHeader)
	if resp2.Code != http.StatusOK {
		t.Fatalf("second post expected 200, got %d: %s", resp2.Code, resp2.Body.String())
	}
}

// test returning 200 with zero posts for an author with no premium posts
func TestGetPremiumPostsByAuthorEmpty(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)

	emptyAuthor := uuid.New()
	resp := api.Get("/api/v1/posts/premium/by-author/"+emptyAuthor.String(), authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result premiumpost.GetPremiumPostsByAuthorIDResponse
	DecodeTo(&result, resp)

	if result.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Total)
	}
	if len(result.Posts) != 0 {
		t.Errorf("expected no posts, got %d", len(result.Posts))
	}
}

// test rejection on non-UUID path segments for premium list routes
func TestBadValidationPremiumPost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, nil)

	resp := api.Get("/api/v1/posts/premium/by-author/random string", authHeader)
	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid author id, got %d: %s", resp.Code, resp.Body.String())
	}

	resp = api.Get("/api/v1/posts/premium/by-sport/random string", authHeader)
	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid sport id, got %d: %s", resp.Code, resp.Body.String())
	}

	resp = api.Get("/api/v1/posts/premium/by-tag/random string", authHeader)
	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid tag id, got %d: %s", resp.Code, resp.Body.String())
	}

	resp = api.Get("/api/v1/posts/premium/by-college/random string", authHeader)
	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid college id, got %d: %s", resp.Code, resp.Body.String())
	}
}

// test checks limit/offset on GET all for pagination stuff
func TestGetAllPremiumPostsPagination(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)

	for i := 0; i < 3; i++ {
		_, err := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
			AuthorID: JohnID,
			SportID:  &SoccerID,
			Title:    "Pagination post",
			Content:  "Body",
		})
		if err != nil {
			t.Fatalf("failed to create post %d: %v", i, err)
		}
	}

	resp := api.Get("/api/v1/posts/premium/?limit=2&offset=0", authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result premiumpost.GetAllPremiumPostsResponse
	DecodeTo(&result, resp)

	if result.Total != 3 {
		t.Errorf("expected total 3, got %d", result.Total)
	}
	if len(result.Posts) != 2 {
		t.Errorf("expected 2 posts with limit=2, got %d", len(result.Posts))
	}

	resp2 := api.Get("/api/v1/posts/premium/?limit=2&offset=2", authHeader)
	if resp2.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp2.Code, resp2.Body.String())
	}

	var page2 premiumpost.GetAllPremiumPostsResponse
	DecodeTo(&page2, resp2)

	if page2.Total != 3 {
		t.Errorf("expected total 3 on page 2, got %d", page2.Total)
	}
	if len(page2.Posts) != 1 {
		t.Errorf("expected 1 post with offset=2, got %d", len(page2.Posts))
	}
}

// test updating a premium post
func TestUpdatePremiumPost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	post, err := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "Original Title",
		Content:  "Original content.",
	})
	if err != nil {
		t.Fatalf("failed to create premium post: %v", err)
	}

	updatedTitle := "Updated Title"
	updatedContent := "Updated content."
	body := map[string]any{
		"title":   updatedTitle,
		"content": updatedContent,
	}

	resp := api.Patch("/api/v1/posts/premium/"+post.ID.String(), authHeader, body)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result premiumpost.PremiumPostResponse
	DecodeTo(&result, resp)

	if result.Title != updatedTitle {
		t.Errorf("expected title %q, got %q", updatedTitle, result.Title)
	}
	if result.Content != updatedContent {
		t.Errorf("expected content %q, got %q", updatedContent, result.Content)
	}
	if result.ID != post.ID {
		t.Errorf("expected id %v, got %v", post.ID, result.ID)
	}
}

// test that a regular user cannot update a premium post
func TestUpdatePremiumPost_UserForbidden(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	post, err := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "Original Title",
		Content:  "Original content.",
	})
	if err != nil {
		t.Fatalf("failed to create premium post: %v", err)
	}

	authHeader := authHeaderWithPermissionsGivenUserForRole(t, testDB.DB, models.RoleUser, []permissionSpec{}, JohnID)

	body := map[string]any{
		"title": "Hacked Title",
	}

	resp := api.Patch("/api/v1/posts/premium/"+post.ID.String(), authHeader, body)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for regular user, got %d", resp.Code)
	}
}

// test deleting a premium post
func TestDeletePremiumPost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	post, err := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "Post to delete",
		Content:  "This post will be deleted.",
	})
	if err != nil {
		t.Fatalf("failed to create premium post: %v", err)
	}

	resp := api.Delete("/api/v1/posts/premium/"+post.ID.String(), authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result premiumpost.DeletePremiumPostRequest
	DecodeTo(&result, resp)

	if result.ID != post.ID {
		t.Errorf("expected deleted post id %v, got %v", post.ID, result.ID)
	}
}

// test that a regular user cannot delete a premium post
func TestDeletePremiumPost_UserForbidden(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	premiumpostDB := premiumpost.NewPremiumPostDB(testDB.DB)

	post, err := premiumpostDB.CreatePremiumPost(&models.PremiumPost{
		AuthorID: JohnID,
		SportID:  &SoccerID,
		Title:    "Post to delete",
		Content:  "This post will be deleted.",
	})
	if err != nil {
		t.Fatalf("failed to create premium post: %v", err)
	}

	authHeader := authHeaderWithPermissionsGivenUserForRole(t, testDB.DB, models.RoleUser, []permissionSpec{}, JohnID)

	resp := api.Delete("/api/v1/posts/premium/"+post.ID.String(), authHeader)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for regular user, got %d", resp.Code)
	}
}

// test deleting a non-existent premium post
func TestDeletePremiumPost_NotFound(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	premiumpost.Route(testDB.API, testDB.DB, nil)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	authHeader := premiumPostAdminAuthHeader(t, testDB)

	resp := api.Delete("/api/v1/posts/premium/"+uuid.New().String(), authHeader)
	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", resp.Code, resp.Body.String())
	}
}
