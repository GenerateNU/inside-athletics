package routeTests

import (
	h "inside-athletics/internal/handlers/comment_like"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"testing"

	"github.com/google/uuid"
)

func TestCreateCommentLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	// Use from post_like_tests
	user := newPostLikeTestUser(uuid.New(), "create-comment-like")
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}

	commentID := uuid.New()
	body := h.CreateCommentLikeRequest{UserID: user.ID, CommentID: commentID}

	resp := api.Post("/api/v1/user/", "Authorization: Bearer mock-token", body)

	var u h.CreateCommentLikeResponse
	DecodeTo(&u, resp)

	if u.ID == uuid.Nil {
		t.Fatalf("Expected like ID to be created, got nil")
	}

	var like models.CommentLike
	if err := testDB.DB.Where("id = ?", u.ID).First(&like).Error; err != nil {
		t.Fatalf("Expected like to exist in DB: %s", err.Error())
	}
	if like.UserID != user.ID {
		t.Fatalf("Unexpected user_id: got %s, expected %s", like.UserID.String(), user.ID.String())
	}
	if like.CommentID != commentID {
		t.Fatalf("Unexpected comment_id: got %s, expected %s", like.CommentID.String(), commentID.String())
	}
}

func TestGetCommentLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user := newPostLikeTestUser(uuid.New(), "get-comment-like")
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}

	commentID := uuid.New()
	commentLike := models.CommentLike{UserID: user.ID, CommentID: commentID}
	likeResp := testDB.DB.Create(&commentLike)
	_, err = utils.HandleDBError(&commentLike, likeResp.Error)
	if err != nil {
		t.Fatalf("Unable to add comment like to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/user/"+commentLike.ID.String(), "Authorization: Bearer mock-token")

	var u h.GetCommentLikeResponse
	DecodeTo(&u, resp)

	if u.UserID != user.ID {
		t.Fatalf("Unexpected user_id: got %s, expected %s", u.UserID.String(), user.ID.String())
	}
	if u.CommentID != commentID {
		t.Fatalf("Unexpected comment_id: got %s, expected %s", u.CommentID.String(), commentID.String())
	}
}

func TestGetCommentLikeCount(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	commentID := uuid.New()
	// Make likes
	for i := 0; i < 3; i++ {
		user := newPostLikeTestUser(uuid.New(), uuid.New().String()[:8])
		userResp := testDB.DB.Create(&user)
		_, err := utils.HandleDBError(&user, userResp.Error)
		if err != nil {
			t.Fatalf("Unable to add user to table: %s", err.Error())
		}
		commentLike := models.CommentLike{UserID: user.ID, CommentID: commentID}
		likeResp := testDB.DB.Create(&commentLike)
		_, err = utils.HandleDBError(&commentLike, likeResp.Error)
		if err != nil {
			t.Fatalf("Unable to add comment like to table: %s", err.Error())
		}
	}

	resp := api.Get("/api/v1/user/comment/"+commentID.String()+"/like-count", "Authorization: Bearer mock-token")

	var u h.GetLikeCountResponse
	DecodeTo(&u, resp)

	if u.Total != 3 {
		t.Fatalf("Unexpected total: got %d, expected 3", u.Total)
	}
}

func TestCheckUserLikedComment(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user := newPostLikeTestUser(uuid.New(), "check-comment-liked")
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}

	commentID := uuid.New()
	commentLike := models.CommentLike{UserID: user.ID, CommentID: commentID}
	likeResp := testDB.DB.Create(&commentLike)
	_, err = utils.HandleDBError(&commentLike, likeResp.Error)
	if err != nil {
		t.Fatalf("Unable to add comment like to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/user/comment/"+commentID.String()+"/check-like?user_id="+user.ID.String(), "Authorization: Bearer mock-token")

	var u h.CheckUserLikedCommentResponse
	DecodeTo(&u, resp)

	if !u.Liked {
		t.Fatalf("Expected liked true, got false")
	}

	// Check that different user didn't like the comment
	otherUserID := uuid.New()
	resp2 := api.Get("/api/v1/user/comment/"+commentID.String()+"/check-like?user_id="+otherUserID.String(), "Authorization: Bearer mock-token")
	var u2 h.CheckUserLikedCommentResponse
	DecodeTo(&u2, resp2)
	if u2.Liked {
		t.Fatalf("Expected liked false for other user, got true")
	}
}

func TestDeleteCommentLike(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user := newPostLikeTestUser(uuid.New(), "delete-comment-like")
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}

	commentLike := models.CommentLike{UserID: user.ID, CommentID: uuid.New()}
	likeResp := testDB.DB.Create(&commentLike)
	_, err = utils.HandleDBError(&commentLike, likeResp.Error)
	if err != nil {
		t.Fatalf("Unable to add comment like to table: %s", err.Error())
	}

	resp := api.Delete("/api/v1/user/"+commentLike.ID.String(), "Authorization: Bearer mock-token")

	var u h.DeleteCommentLikeResponse
	DecodeTo(&u, resp)

	if u.Message != "Like was deleted successfully" {
		t.Fatalf("Unexpected message: got %s", u.Message)
	}
}
