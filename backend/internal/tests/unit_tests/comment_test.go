package unitTests

import (
	"inside-athletics/internal/handlers/comment"
	"inside-athletics/internal/models"
	"testing"

	"github.com/google/uuid"
)

// Verifies ToCommentResponse omits user_id for anonymous comments when caller is not the user who made the comment.
func TestToCommentResponse_WhenAnonymousAndNotOwnUser_HidesUserID(t *testing.T) {
	userID := uuid.New()
	c := &models.Comment{
		ID:          uuid.New(),
		UserID:      userID,
		IsAnonymous: true,
		PostID:      uuid.New(),
		Description: "Anonymous",
	}
	resp := comment.ToCommentResponse(c, uuid.New())
	if resp.User != nil {
		t.Errorf("expected user_id nil for anonymous when not own user, got %v", resp.User.ID)
	}
	if resp.IsAnonymous != true {
		t.Error("expected is_anonymous true")
	}
}

// Verifies ToCommentResponse includes user_id for anonymous comments when caller is the user who made the comment.
func TestToCommentResponse_WhenAnonymousAndNotOwnUser_ShowsUserID(t *testing.T) {
	userID := uuid.New()
	c := &models.Comment{
		ID:          uuid.New(),
		UserID:      userID,
		IsAnonymous: true,
		PostID:      uuid.New(),
		Description: "Anonymous",
	}
	resp := comment.ToCreateCommentResponse(c, c.UserID)
	if resp.UserID == nil || *resp.UserID != userID {
		t.Errorf("expected user_id %s for own user, got %v", userID, resp.UserID)
	}
}

// Verifies ToCommentResponse includes user_id for non-anonymous comments.
func TestToCommentResponse_WhenNotAnonymous_ShowsUserID(t *testing.T) {
	userID := uuid.New()
	c := &models.Comment{
		ID:          uuid.New(),
		UserID:      userID,
		IsAnonymous: false,
		PostID:      uuid.New(),
		Description: "Not anonymous",
	}
	resp := comment.ToCreateCommentResponse(c, uuid.New())
	if resp.UserID == nil || *resp.UserID != userID {
		t.Errorf("expected user_id %s when not anonymous, got %v", userID, resp.UserID)
	}
}
