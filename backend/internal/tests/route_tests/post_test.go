package routeTests

import (
	"fmt"
	"inside-athletics/internal/handlers/post"
	"inside-athletics/internal/handlers/sport"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func CreateUserAndSport(testDB *TestDatabase, t *testing.T) {
	user := models.User{
		ID:                      JohnID,
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
	soccer := models.Sport{
		ID:         SoccerID,
		Name:       "Soccer",
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

	CreateUserAndSport(testDB, t)

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
		{Action: models.PermissionCreate, Resource: "post"},
	},
		JohnID,
	)

	body := map[string]any{
		"sport_id":     SoccerID,
		"title":        "Looking for thoughts on NEU Fencing!",
		"content":      "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?",
		"is_anonymous": false,
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

	if result.IsAnonymous != false {
		t.Errorf("expected IsAnonymous false, got %v", result.IsAnonymous)
	}
}

func TestCreatePostWithoutTagsThrowsError(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API

	CreateUserAndSport(testDB, t)

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
		{Action: models.PermissionCreate, Resource: "post"},
	},
		JohnID,
	)

	body := map[string]any{
		"title":        "Looking for thoughts on NEU Fencing!",
		"content":      "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?",
		"is_anonymous": false,
		"tags":         []map[string]any{},
	}

	resp := api.Post("/api/v1/post/", body, authHeader)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestCreatePostWithTags(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API

	CreateUserAndSport(testDB, t)

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
		{Action: models.PermissionCreate, Resource: "post"},
	},
		JohnID,
	)

	tag1 := models.Tag{Name: "recruiting"}
	tag2 := models.Tag{Name: "fencing"}
	testDB.DB.Create(&tag1)
	testDB.DB.Create(&tag2)

	body := map[string]any{
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
		IsAnonymous: false,
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

	if result.LikeCount != 0 {
		t.Errorf("expected Likes 0, got %d", result.LikeCount)
	}

	if result.IsAnonymous != false {
		t.Errorf("expected IsAnonymous false, got %v", result.IsAnonymous)
	}
}

func TestGetPostByIdWithLikes(t *testing.T) {
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
		IsAnonymous: false,
	}, []post.TagRequest{})
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// seed a like
	like := models.PostLike{
		UserID: JohnID,
		PostID: createdPost.ID,
	}
	if err := testDB.DB.Create(&like).Error; err != nil {
		t.Fatalf("failed to create like: %v", err)
	}

	// Seed a top-level comment
	commentID := uuid.New()
	comment := models.Comment{
		ID:          commentID,
		UserID:      JohnID,
		PostID:      createdPost.ID,
		Description: "Test comment",
		IsAnonymous: false,
	}
	if err := testDB.DB.Create(&comment).Error; err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	// Seed a reply (comment with a parent)
	reply := models.Comment{
		UserID:          JohnID,
		PostID:          createdPost.ID,
		ParentCommentID: &commentID,
		Description:     "Test reply",
		IsAnonymous:     false,
	}
	if err := testDB.DB.Create(&reply).Error; err != nil {
		t.Fatalf("failed to create reply: %v", err)
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

	if result.LikeCount != 1 {
		t.Errorf("expected Likes 1, got %d", result.LikeCount)
	}

	if result.CommentCount != 2 {
		t.Errorf("expected Comments 2 , got %d", result.CommentCount)
	}

	if result.IsAnonymous != false {
		t.Errorf("expected IsAnonymous false, got %v", result.IsAnonymous)
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

	post1, err1 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "First Post About Fencing", Content: "This is the first post content", IsAnonymous: false,
	}, []post.TagRequest{})
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	like := models.PostLike{
		UserID: JohnID,
		PostID: post1.ID,
	}
	if err := testDB.DB.Create(&like).Error; err != nil {
		t.Fatalf("failed to create like: %v", err)
	}

	_, err2 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "Second Post About Basketball", Content: "This is the second post content", IsAnonymous: false,
	}, []post.TagRequest{})
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

	if result.Posts[0].LikeCount != 1 {
		t.Errorf("expected 1 like in first post")
	}
	if result.Posts[1].LikeCount != 0 {
		t.Errorf("expected 0 likes in first post")
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
	}, []post.TagRequest{})
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	_, err2 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "Second Post About Basketball", Content: "This is the second post content", IsAnonymous: false,
	}, []post.TagRequest{})
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
	}, []post.TagRequest{})
	if err1 != nil {
		t.Fatalf("failed to create post 1: %v", err1)
	}

	_, err2 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "Second Post About Basketball", Content: "This is the second post content", IsAnonymous: false,
	}, []post.TagRequest{})
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
	}, []post.TagRequest{})
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
	}, []post.TagRequest{})
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

func TestPostSearch(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	api := testDB.API

	postDB := post.NewPostDB(testDB.DB)

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "post"},
	})

	CreateUserAndSport(testDB, t)

	post1, err := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "Northeastern University Soccer Is LIT!", Content: "Wow I love NEU Soccer", IsAnonymous: false,
	}, []post.TagRequest{})
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	_, err1 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "Northwestern University Soccer Is LIT!", Content: "Wow I love NWU Soccer", IsAnonymous: false,
	}, []post.TagRequest{})
	if err1 != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Test searching for substring of title, expecting post1 to have higher similarity score
	resp := api.Get("/api/v1/posts/search?search_str=NorthE", authHeader)
	if resp.Code != http.StatusOK {
		t.Errorf("Expected 204 got %d %s", resp.Code, resp.Body.String())
	}
	var searchResp post.GetSearchResponse
	DecodeTo(&searchResp, resp)

	if searchResp.Count != 2 {
		t.Errorf("Expected 2 entries but got %d", searchResp.Count)
	}
	if searchResp.Posts[0].Title != post1.Title {
		t.Error("Expected post1 to have higher similarity score to search string")
	}
}

func TestTypoStillReturns(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	api := testDB.API

	postDB := post.NewPostDB(testDB.DB)

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "post"},
	})

	CreateUserAndSport(testDB, t)

	post1, err := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "Northeastern University Soccer Is LIT!", Content: "Wow I love NEU Soccer", IsAnonymous: false,
	}, []post.TagRequest{})
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	_, err1 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "I farted", Content: "Wow it smells", IsAnonymous: false,
	}, []post.TagRequest{})
	if err1 != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	resp := api.Get("/api/v1/posts/search?search_str=northeusternTest", authHeader)
	if resp.Code != http.StatusOK {
		t.Errorf("Expected 204 got %d %s", resp.Code, resp.Body.String())
	}
	var searchResp post.GetSearchResponse
	DecodeTo(&searchResp, resp)

	if searchResp.Count != 1 {
		t.Errorf("Expected 1 entries but got %d", searchResp.Count)
	}
	if searchResp.Posts[0].Title != post1.Title {
		t.Error("Expected to retrieve post 1")
	}
}

func TestSearchLimit(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	api := testDB.API

	postDB := post.NewPostDB(testDB.DB)

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "post"},
	})

	CreateUserAndSport(testDB, t)

	_, err := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "Northeastern University Soccer Is LIT!", Content: "Wow I love NEU Soccer", IsAnonymous: false,
	}, []post.TagRequest{})
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	_, err1 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID,
		Title: "I farted north", Content: "Wow it smells", IsAnonymous: false,
	}, []post.TagRequest{})
	if err1 != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	resp := api.Get("/api/v1/posts/search?search_str=northTest&limit=1", authHeader)
	if resp.Code != http.StatusOK {
		t.Errorf("Expected 204 got %d %s", resp.Code, resp.Body.String())
	}
	var searchResp post.GetSearchResponse
	DecodeTo(&searchResp, resp)

	if searchResp.Count != 2 {
		t.Errorf("Expected 2 total entries but got %d", searchResp.Count)
	}
	if len(searchResp.Posts) != 1 {
		t.Errorf("Expected only 1 entry to be returned")
	}
}

func TestFilterPosts(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	api := testDB.API

	postDB := post.NewPostDB(testDB.DB)

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "post"},
	})

	CreateUserAndSport(testDB, t)

	collegeId := uuid.New()
	neu := models.College{
		ID:           collegeId,
		Name:         "Northeastern University",
		State:        "Massachusetts",
		City:         "Boston",
		Website:      "https://www.northeastern.edu",
		DivisionRank: models.DivisionI,
	}
	collegeResp := testDB.DB.Create(&neu)
	_, errCollege := utils.HandleDBError(&neu, collegeResp.Error)

	if errCollege != nil {
		t.Fatalf("Unable to add college to table: %s", errCollege.Error())
	}

	sportDB := sport.NewSportDB(testDB.DB)
	popularity := int32(1)
	ermSport, errSport := sportDB.CreateSport("Erm Sport", &popularity)
	if errSport != nil {
		t.Fatal("Unable to create sport erm sport")
	}

	tags := make([]models.Tag, 0)
	for i := range 3 {
		tagId := uuid.New()
		tag := models.Tag{
			ID:   tagId,
			Name: fmt.Sprintf("Tag%d", i),
		}
		tags = append(tags, tag)
		testDB.DB.Create(&tag)
	}

	mapTagRequests := func(t models.Tag) post.TagRequest { return post.TagRequest{ID: t.ID} }

	_, err := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &SoccerID, CollegeID: &collegeId,
		Title: "Northeastern University Soccer Is LIT!", Content: "Wow I love NEU Soccer", IsAnonymous: false,
	}, utils.MapList(tags[0:2], mapTagRequests))

	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	_, err1 := postDB.CreatePost(&models.Post{
		AuthorID: JohnID, SportID: &ermSport.ID, CollegeID: &collegeId,
		Title: "I farted north", Content: "Wow it smells", IsAnonymous: false,
	}, utils.MapList(tags[1:3], mapTagRequests))
	if err1 != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	filterCollegesResp := api.Get("/api/v1/posts/filter?college_ids="+collegeId.String(), authHeader)
	if filterCollegesResp.Code != http.StatusOK {
		t.Fatalf("Expected a 200 but got %d", filterCollegesResp.Code)
	}

	var filteredColleges post.GetAllPostsResponse

	DecodeTo(&filteredColleges, filterCollegesResp)

	if filteredColleges.Total != 2 {
		t.Fatalf("Expected 2 posts in response for college filter got %d", filteredColleges.Total)
	}

	filterSportsResp := api.Get("/api/v1/posts/filter?sport_ids="+ermSport.ID.String(), authHeader)
	if filterSportsResp.Code != http.StatusOK {
		t.Fatalf("Expected to get status code 200 but got %d", filterSportsResp.Code)
	}

	var filteredSports post.GetAllPostsResponse

	DecodeTo(&filteredSports, filterSportsResp)

	if filteredSports.Total != 1 {
		t.Fatalf("Expected to get 1 filtered post got %d", filteredSports.Total)
	}

	filterSportsResp = api.Get("/api/v1/posts/filter?sport_ids="+ermSport.ID.String()+","+SoccerID.String(), authHeader)
	if filterSportsResp.Code != http.StatusOK {
		t.Fatalf("Expected to get status code 200 but got %d", filterSportsResp.Code)
	}

	DecodeTo(&filteredSports, filterSportsResp)

	if filteredSports.Total != 2 {
		t.Fatalf("Expected to get 2 filtered posts based on 2 sports got %d", filteredSports.Total)
	}

	tagIds := utils.MapList(tags, func(t models.Tag) string {
		return t.ID.String()
	})

	// test tag filtering
	filterTagsResp := api.Get("/api/v1/posts/filter?tag_ids="+strings.Join(tagIds, ","), authHeader)
	if filterCollegesResp.Code != http.StatusOK {
		t.Fatalf("Expected status code 200 but got %d", filterTagsResp.Code)
	}

	var filteredTags post.GetAllPostsResponse

	DecodeTo(&filteredTags, filterTagsResp)

	if filteredTags.Total != 2 {
		t.Fatalf("Expected 2 filtered posts but got %d", filteredTags.Total)
	}

	filterTagsResp = api.Get("/api/v1/posts/filter?tag_ids="+tagIds[0], authHeader)
	if filterCollegesResp.Code != http.StatusOK {
		t.Fatalf("Expected status code 200 but got %d", filterTagsResp.Code)
	}

	DecodeTo(&filteredTags, filterTagsResp)

	if filteredTags.Total != 1 {
		t.Fatalf("Expected 1 filtered posts but got %d", filteredTags.Total)
	}

	filterSportAndTagResp := api.Get(fmt.Sprintf("/api/v1/posts/filter?sport_ids=%s&tag_ids=%s", SoccerID.String(), tagIds[2]), authHeader)
	if filterSportAndTagResp.Code != http.StatusOK {
		t.Fatalf("Expected status code 200 but got %d", filterSportAndTagResp.Code)
	}

	var filterSportAndTag post.GetAllPostsResponse

	DecodeTo(&filterSportAndTag, filterSportAndTagResp)

	if filterSportAndTag.Total != 2 {
		t.Fatalf("Expected filter to return 2 posts but got %d", filterSportAndTag.Total)
	}
}

func uuidDereference(v *uuid.UUID) uuid.UUID {
	return *v
}
