package unitTests

import (
	"inside-athletics/internal/handlers/comment"
	"inside-athletics/internal/models"
	"testing"

	"github.com/google/uuid"
)

// Verifies ToCommentResponse omits user_id for anonymous comments when caller is not super user.
func TestToCommentResponse_WhenAnonymousAndNotSuperUser_HidesUserID(t *testing.T) {
	userID := uuid.New()
	c := &models.Comment{
		ID:          uuid.New(),
		UserID:      userID,
		IsAnonymous: true,
		PostID:      uuid.New(),
		Description: "Anonymous",
	}
	resp := comment.ToCommentResponse(c, false)
	if resp.UserID != nil {
		t.Errorf("expected user_id nil for anonymous when not super user, got %v", resp.UserID)
	}
	if resp.IsAnonymous != true {
		t.Error("expected is_anonymous true")
	}
}

// Verifies ToCommentResponse includes user_id for anonymous comments when caller is super user.
func TestToCommentResponse_WhenAnonymousAndSuperUser_ShowsUserID(t *testing.T) {
	userID := uuid.New()
	c := &models.Comment{
		ID:          uuid.New(),
		UserID:      userID,
		IsAnonymous: true,
		PostID:      uuid.New(),
		Description: "Anonymous",
	}
	resp := comment.ToCommentResponse(c, true)
	if resp.UserID == nil || *resp.UserID != userID {
		t.Errorf("expected user_id %s for super user, got %v", userID, resp.UserID)
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
	resp := comment.ToCommentResponse(c, false)
	if resp.UserID == nil || *resp.UserID != userID {
		t.Errorf("expected user_id %s when not anonymous, got %v", userID, resp.UserID)
	}
}
