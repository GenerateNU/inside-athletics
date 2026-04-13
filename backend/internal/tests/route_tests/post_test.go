package routeTests

import (
	"errors"
	"inside-athletics/internal/handlers/post"
	"inside-athletics/internal/models"
	"net/http"
	"strconv"
	"strings"
	"sync"
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

	resp := api.Post("/api/v1/post/", authHeader, body)
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

// TestFreeUserCannotCreateSecondPost asserts free users get 403 when creating a second post.
func TestFreeUserCannotCreateSecondPost(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API

	CreateUserAndSport(testDB, t)

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
		{Action: models.PermissionCreate, Resource: "post"},
	}, JohnID)

	body := map[string]any{
		"sport_id": SoccerID, "title": "First post", "content": "Content.", "is_anonymous": false, "tags": []map[string]any{},
	}
	resp1 := api.Post("/api/v1/post/", body, authHeader)
	if resp1.Code != http.StatusOK {
		t.Fatalf("first post expected 200, got %d: %s", resp1.Code, resp1.Body.String())
	}

	body["title"] = "Second post"
	resp2 := api.Post("/api/v1/post/", body, authHeader)
	if resp2.Code != http.StatusForbidden {
		t.Fatalf("second post (free user) expected 403, got %d: %s", resp2.Code, resp2.Body.String())
	}
	if resp2.Body.String() != "" && !strings.Contains(resp2.Body.String(), "free post views") {
		t.Logf("response body (expected free limit message): %s", resp2.Body.String())
	}
}

func TestFreeUserConcurrentPostCreationEnforcesLimit(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	CreateUserAndSport(testDB, t)
	postDB := post.NewPostDB(testDB.DB)

	errs := make(chan error, 2)
	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, err := postDB.CreatePostWithAuthorLimit(&models.Post{
				AuthorID:    JohnID,
				SportID:     &SoccerID,
				Title:       "Concurrent post " + strconv.Itoa(i),
				Content:     "Content",
				IsAnonymous: false,
			}, []post.TagRequest{}, post.FreeUserMaxPosts)
			errs <- err
		}(i)
	}

	close(start)
	wg.Wait()
	close(errs)

	var successCount int
	var limitErrorCount int
	for err := range errs {
		switch {
		case err == nil:
			successCount++
		case errors.Is(err, post.ErrFreePostCreationLimitReached):
			limitErrorCount++
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if successCount != 1 {
		t.Fatalf("expected exactly one successful post creation, got %d", successCount)
	}
	if limitErrorCount != 1 {
		t.Fatalf("expected exactly one free-tier limit error, got %d", limitErrorCount)
	}

	count, err := postDB.CountPostsByAuthor(JohnID)
	if err != nil {
		t.Fatalf("failed to count posts by author: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly one persisted post, got %d", count)
	}
}

// TestFreeUserGetPostReturns403AfterMaxViews asserts free users get 403 when viewing more than FreeUserMaxPostViews distinct posts.
func TestFreeUserGetPostReturns403AfterMaxViews(t *testing.T) {
	t.Parallel()
	const maxViews = 5
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API
	postDB := post.NewPostDB(testDB.DB)

	CreateUserAndSport(testDB, t)

	freeUserID := uuid.New()
	freeUser := models.User{
		ID:                      freeUserID,
		FirstName:               "Free",
		LastName:                "User",
		Email:                   "free@example.com",
		Username:                "freeuser",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&freeUser).Error; err != nil {
		t.Fatalf("failed to create free user: %v", err)
	}
	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
	}, freeUserID)

	var postIDs []uuid.UUID
	for i := 0; i < maxViews+1; i++ {
		p, err := postDB.CreatePost(&models.Post{
			AuthorID: JohnID, SportID: &SoccerID,
			Title: "Post " + strconv.Itoa(i), Content: "Content", IsAnonymous: false,
		}, []post.TagRequest{})
		if err != nil {
			t.Fatalf("failed to create post %d: %v", i, err)
		}
		postIDs = append(postIDs, p.ID)
	}

	for i := 0; i < maxViews; i++ {
		resp := api.Get("/api/v1/post/"+postIDs[i].String(), authHeader)
		if resp.Code != http.StatusOK {
			t.Fatalf("view %d (post %s) expected 200, got %d: %s", i, postIDs[i], resp.Code, resp.Body.String())
		}
	}

	resp := api.Get("/api/v1/post/"+postIDs[maxViews].String(), authHeader)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("view %d (over limit) expected 403, got %d: %s", maxViews, resp.Code, resp.Body.String())
	}

	respAgain := api.Get("/api/v1/post/"+postIDs[0].String(), authHeader)
	if respAgain.Code != http.StatusOK {
		t.Fatalf("re-viewing first post should still be 200, got %d: %s", respAgain.Code, respAgain.Body.String())
	}
}

func TestFreeUserConcurrentDistinctPostViewsEnforceLimit(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	postDB := post.NewPostDB(testDB.DB)
	CreateUserAndSport(testDB, t)

	freeUserID := uuid.New()
	freeUser := models.User{
		ID:                      freeUserID,
		FirstName:               "Free",
		LastName:                "User",
		Email:                   "free-concurrent@example.com",
		Username:                "freeuser-concurrent",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&freeUser).Error; err != nil {
		t.Fatalf("failed to create free user: %v", err)
	}

	postIDs := make([]uuid.UUID, 0, post.FreeUserMaxPostViews+1)
	for i := 0; i < post.FreeUserMaxPostViews+1; i++ {
		createdPost, err := postDB.CreatePost(&models.Post{
			AuthorID:    JohnID,
			SportID:     &SoccerID,
			Title:       "Concurrent view target " + strconv.Itoa(i),
			Content:     "Content",
			IsAnonymous: false,
		}, []post.TagRequest{})
		if err != nil {
			t.Fatalf("failed to create post %d: %v", i, err)
		}
		postIDs = append(postIDs, createdPost.ID)
	}

	errs := make(chan error, len(postIDs))
	start := make(chan struct{})
	var wg sync.WaitGroup

	for _, postID := range postIDs {
		wg.Add(1)
		go func(postID uuid.UUID) {
			defer wg.Done()
			<-start
			errs <- postDB.RecordPostViewIfAllowed(freeUserID, postID, post.FreeUserMaxPostViews)
		}(postID)
	}

	close(start)
	wg.Wait()
	close(errs)

	var successCount int
	var limitErrorCount int
	for err := range errs {
		switch {
		case err == nil:
			successCount++
		case errors.Is(err, post.ErrFreePostViewLimitReached):
			limitErrorCount++
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if successCount != post.FreeUserMaxPostViews {
		t.Fatalf("expected %d successful view records, got %d", post.FreeUserMaxPostViews, successCount)
	}
	if limitErrorCount != 1 {
		t.Fatalf("expected exactly one free-tier view limit error, got %d", limitErrorCount)
	}

	count, err := postDB.CountViewedPostsByUser(freeUserID)
	if err != nil {
		t.Fatalf("failed to count viewed posts: %v", err)
	}
	if count != post.FreeUserMaxPostViews {
		t.Fatalf("expected %d persisted viewed posts, got %d", post.FreeUserMaxPostViews, count)
	}
}

// TestPremiumUserCanCreateMultiplePosts asserts premium users can create more than one post.
func TestPremiumUserCanCreateMultiplePosts(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	post.Route(testDB.API, testDB.DB)
	api := testDB.API

	CreateUserAndSport(testDB, t)
	if err := testDB.DB.Model(&models.User{}).Where("id = ?", JohnID).Update("Account_Type", true).Error; err != nil {
		t.Fatalf("failed to set premium: %v", err)
	}

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sport"},
		{Action: models.PermissionCreate, Resource: "post"},
	}, JohnID)

	body := map[string]any{
		"sport_id": SoccerID, "title": "First", "content": "Content.", "is_anonymous": false, "tags": []map[string]any{},
	}
	resp1 := api.Post("/api/v1/post/", body, authHeader)
	if resp1.Code != http.StatusOK {
		t.Fatalf("first post expected 200, got %d: %s", resp1.Code, resp1.Body.String())
	}
	body["title"] = "Second"
	resp2 := api.Post("/api/v1/post/", body, authHeader)
	if resp2.Code != http.StatusOK {
		t.Fatalf("second post (premium) expected 200, got %d: %s", resp2.Code, resp2.Body.String())
	}
}

func uuidDereference(v *uuid.UUID) uuid.UUID {
	return *v
}
