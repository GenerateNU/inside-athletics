package routeTests

import (
	"inside-athletics/internal/handlers/post"
	"inside-athletics/internal/models"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func CreateUserAndSport(testDB *TestDatabase, t *testing.T) {
	user := models.User{
		ID: JohnID,
		FirstName:               "Test",
		LastName:                "User",
		Email:                   "test-john@example.com",
		Username:                "testuser-john",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}

	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	popularity := int32(100)
	soccer := models.Sport {
		ID: SoccerID,
		Name: "Soccer",
		Popularity: &popularity,
	}

	if err := testDB.DB.Create(&soccer).Error; err != nil {
		t.Fatalf("failed to create sport: %v", err)
	}
}
func TestCreatePost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
		{Action: models.PermissionCreate, Resource: "post"},
	})

	CreateUserAndSport(testDB, t)
	
	body := map[string]any{
		"author_id":    JohnID,
		"sport_id":     SoccerID,
		"title":        "Looking for thoughts on NEU Fencing!",
		"content":      "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?",
		"is_anonymous": true,
		"tags":         []map[string]any{},
	}

	resp := api.Post("/api/v1/post/", body, authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post.CreatePostResponse
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

	if result.Likes != 0 {
		t.Errorf("expected Likes 0, got %d", result.Likes)
	}

	if result.IsAnonymous != true {
		t.Errorf("expected IsAnonymous true, got %v", result.IsAnonymous)
	}
}

func TestCreatePostWithTags(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
		{Action: models.PermissionCreate, Resource: "post"},
	})

	CreateUserAndSport(testDB, t)

	tag1 := models.Tag{Name: "recruiting"}
	tag2 := models.Tag{Name: "fencing"}
	testDB.DB.Create(&tag1)
	testDB.DB.Create(&tag2)

	body := map[string]any{
		"author_id":  	JohnID,
		"sport_id":     SoccerID,
		"title":        "Post with tags",
		"content":      "Testing that tags are associated with this post correctly.",
		"is_anonymous": false,
		"tags": []map[string]any{
			{"id": tag1.ID},
			{"id": tag2.ID},
		},
	}

	resp := api.Post("/api/v1/post/", body, authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post.CreatePostResponse
	DecodeTo(&result, resp)

	if len(result.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(result.Tags))
	}
}

func TestGetPostById(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API
	postDB := post.NewPostDB(testDB.DB)

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
	})

	CreateUserAndSport(testDB, t)

	createdPost, err := postDB.CreatePost(&models.Post{
		AuthorID:    JohnID,
		SportID:     &SoccerID,
		Title:       "Looking for thoughts on NEU Fencing!",
		Content:     "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?",
		IsAnonymous: true,
	}, []post.TagRequest{})
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	resp := api.Get("/api/v1/post/"+createdPost.ID.String(), authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post.PostResponse
	DecodeTo(&result, resp)

	if result.Author == nil || result.Author.ID != JohnID {
		t.Errorf("expected authorID %v, got %v", JohnID, result.Author)
	}

	if result.Sport == nil || result.Sport.ID != SoccerID {
		t.Errorf("expected sportID %v, got %v", SoccerID, result.Sport)
	}

	if result.Title != "Looking for thoughts on NEU Fencing!" {
		t.Errorf("expected title %q, got %q", "Looking for thoughts on NEU Fencing!", result.Title)
	}

	if result.Content != "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" {
		t.Errorf("expected content %q, got %q", "My name is Bob Joe...", result.Content)
	}

	if result.Likes != 0 {
		t.Errorf("expected Likes 0, got %d", result.Likes)
	}

	if result.IsAnonymous != true {
		t.Errorf("expected IsAnonymous true, got %v", result.IsAnonymous)
	}
}

func TestGetPostByIdNotFound(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, nil)

	resp := api.Get("/api/v1/post/"+uuid.New().String(), authHeader)
	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.Code)
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

	authHeader := authHeaderWithPermissions(t, testDB.DB, nil)

	resp := api.Get("/api/v1/post/"+"random string", authHeader)
	if resp.Code == http.StatusOK {
		t.Fatalf("expected status 422, got %d: %s", resp.Code, resp.Body.String())
	}

	resp = api.Get("/api/v1/posts/by-sport/"+"random string", authHeader)
	if resp.Code == http.StatusOK {
		t.Fatalf("expected status 422, got %d: %s", resp.Code, resp.Body.String())
	}

	resp = api.Get("/api/v1/posts/by-author/"+"random string", authHeader)
	if resp.Code == http.StatusOK {
		t.Fatalf("expected status 422, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestGetPostByAuthorId(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API
	postDB := post.NewPostDB(testDB.DB)

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
	})

	CreateUserAndSport(testDB, t)

	_, err1 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "First Post About Fencing", Content: "This is the first post content", IsAnonymous: false,
	},[]post.TagRequest{})
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	_, err2 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "Second Post About Basketball", Content: "This is the second post content", IsAnonymous: true,
	},[]post.TagRequest{})
	if err2 != nil {
		t.Fatalf("failed to create post 2: %v", err2)
	}

	resp := api.Get("/api/v1/posts/by-author/"+JohnID.String(), authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result post.GetPostsByAuthorIDResponse
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

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
	})

	CreateUserAndSport(testDB, t)

	_, err1 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "First Post About Fencing", Content: "This is the first post content", IsAnonymous: false,
	},[]post.TagRequest{})
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	_, err2 := postDB.CreatePost(&models.Post{
		AuthorID:JohnID, SportID: &SoccerID,
		Title: "Second Post About Basketball", Content: "This is the second post content", IsAnonymous: true,
	},[]post.TagRequest{})
	if err2 != nil {
		t.Fatalf("failed to create post 2: %v", err2)
	}

	resp := api.Get("/api/v1/posts/by-sport/"+SoccerID.String(), authHeader)
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

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
	})

	CreateUserAndSport(testDB, t)

	_, err1 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "First Post About Fencing", Content: "This is the first post content", IsAnonymous: false,
	},[]post.TagRequest{})
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	_, err2 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID, 
		Title: "Second Post About Basketball", Content: "This is the second post content", IsAnonymous: true,
	},[]post.TagRequest{})
	if err2 != nil {
		t.Fatalf("failed to create post 2: %v", err2)
	}

	resp := api.Get("/api/v1/posts/", authHeader)
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

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
		{Action: models.PermissionUpdate, Resource: "post"},
	})

	CreateUserAndSport(testDB, t)

	createdPost, err := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "Original Title", Content: "Original content", IsAnonymous: false,
	},[]post.TagRequest{})
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	updateBody := map[string]any{
		"title":   "Updated Title",
		"content": "Updated content about the program",
	}

	resp := api.Patch("/api/v1/post/"+createdPost.ID.String(), updateBody, authHeader)
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

func TestUpdatePostNotFound(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionUpdate, Resource: "post"},
	})

	updateBody := map[string]any{"title": "Doesn't Matter"}
	resp := api.Patch("/api/v1/post/"+uuid.New().String(), updateBody, authHeader)
	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.Code)
	}
}

func TestDeletePost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API
	postDB := post.NewPostDB(testDB.DB)

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
		{Action: models.PermissionDelete, Resource: "post"},
	})

	CreateUserAndSport(testDB, t)

	createdPost, err := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "Post to Delete", Content: "This post will be deleted", IsAnonymous: false,
	},[]post.TagRequest{})
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	resp := api.Delete("/api/v1/post/"+createdPost.ID.String(), authHeader)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 204, got %d: %s", resp.Code, resp.Body.String())
	}

	getResp := api.Get("/api/v1/post/"+createdPost.ID.String(), authHeader)
	if getResp.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getResp.Code)
	}
}

func TestDeletePostNotFound(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionDelete, Resource: "post"},
	})

	resp := api.Delete("/api/v1/post/"+uuid.New().String(), authHeader)
	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.Code)
	}
}

func uuidDereference(v *uuid.UUID) uuid.UUID{
	return *v
}