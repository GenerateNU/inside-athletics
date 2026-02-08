package routeTests

import (
	h "inside-athletics/internal/handlers/post_like"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"testing"

	"github.com/google/uuid"
)

// newPostLikeTestUser returns a User for testing
func newPostLikeTestUser(id uuid.UUID, unique string) models.User {
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

func TestCreatePostLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	userID := uuid.New()
	postID := uuid.New()
	body := h.CreatePostLikeRequest{UserID: userID, PostID: postID}

	resp := api.Post("/api/v1/user/", "Authorization: Bearer mock-token", body)

	var u h.CreatePostLikeResponse
	DecodeTo(&u, resp)

	if u.ID == uuid.Nil {
		t.Fatalf("Expected like ID to be created, got nil")
	}

	// Verify the like exists and has correct user_id and post_id
	var like models.PostLike
	if err := testDB.DB.Where("id = ?", u.ID).First(&like).Error; err != nil {
		t.Fatalf("Expected like to exist in DB: %s", err.Error())
	}
	if like.UserID != userID {
		t.Fatalf("Unexpected user_id: got %s, expected %s", like.UserID.String(), userID.String())
	}
	if like.PostID != postID {
		t.Fatalf("Unexpected post_id: got %s, expected %s", like.PostID.String(), postID.String())
	}
}

func TestGetPostLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user := newPostLikeTestUser(uuid.New(), "get-like")
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}

	postID := uuid.New()
	postLike := models.PostLike{
		UserID: user.ID,
		PostID: postID,
	}
	likeResp := testDB.DB.Create(&postLike)
	_, err = utils.HandleDBError(&postLike, likeResp.Error)
	if err != nil {
		t.Fatalf("Unable to add post like to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/user/"+postLike.ID.String(), "Authorization: Bearer mock-token")

	var u h.GetPostLikeResponse
	DecodeTo(&u, resp)

	if u.UserID != user.ID {
		t.Fatalf("Unexpected user_id: got %s, expected %s", u.UserID.String(), user.ID.String())
	}
	if u.PostID != postID {
		t.Fatalf("Unexpected post_id: got %s, expected %s", u.PostID.String(), postID.String())
	}
}

func TestGetLikeCount(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	postID := uuid.New()
	// Making likes
	for i := 0; i < 3; i++ {
		user := newPostLikeTestUser(uuid.New(), uuid.New().String()[:8])
		userResp := testDB.DB.Create(&user)
		_, err := utils.HandleDBError(&user, userResp.Error)
		if err != nil {
			t.Fatalf("Unable to add user to table: %s", err.Error())
		}
		postLike := models.PostLike{UserID: user.ID, PostID: postID}
		likeResp := testDB.DB.Create(&postLike)
		_, err = utils.HandleDBError(&postLike, likeResp.Error)
		if err != nil {
			t.Fatalf("Unable to add post like to table: %s", err.Error())
		}
	}

	resp := api.Get("/api/v1/user/post/"+postID.String()+"/like-count", "Authorization: Bearer mock-token")

	var u h.GetLikeCountResponse
	DecodeTo(&u, resp)

	if u.Total != 3 {
		t.Fatalf("Unexpected total: got %d, expected 3", u.Total)
	}
}

func TestCheckUserLikedPost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user := newPostLikeTestUser(uuid.New(), "check-liked")
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}

	postID := uuid.New()
	postLike := models.PostLike{UserID: user.ID, PostID: postID}
	likeResp := testDB.DB.Create(&postLike)
	_, err = utils.HandleDBError(&postLike, likeResp.Error)
	if err != nil {
		t.Fatalf("Unable to add post like to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/user/post/"+postID.String()+"/check-like?user_id="+user.ID.String(), "Authorization: Bearer mock-token")

	var u h.CheckUserLikedPostResponse
	DecodeTo(&u, resp)

	if !u.Liked {
		t.Fatalf("Expected liked true, got false")
	}

	// Check that different user didn't like the post
	otherUserID := uuid.New()
	resp2 := api.Get("/api/v1/user/post/"+postID.String()+"/check-like?user_id="+otherUserID.String(), "Authorization: Bearer mock-token")
	var u2 h.CheckUserLikedPostResponse
	DecodeTo(&u2, resp2)
	if u2.Liked {
		t.Fatalf("Expected liked false for other user, got true")
	}
}

func TestDeletePostLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user := newPostLikeTestUser(uuid.New(), "delete-like")
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}

	postLike := models.PostLike{UserID: user.ID, PostID: uuid.New()}
	likeResp := testDB.DB.Create(&postLike)
	_, err = utils.HandleDBError(&postLike, likeResp.Error)
	if err != nil {
		t.Fatalf("Unable to add post like to table: %s", err.Error())
	}

	resp := api.Delete("/api/v1/user/"+postLike.ID.String(), "Authorization: Bearer mock-token")

	var u h.DeletePostLikeResponse
	DecodeTo(&u, resp)

	if u.Message != "Like was deleted successfully" {
		t.Fatalf("Unexpected message: got %s", u.Message)
	}
}
