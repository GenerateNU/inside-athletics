package routeTests

import (
	"inside-athletics/internal/handlers/tagfollow"
	"inside-athletics/internal/models"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// seedUserAndTag creates a user and tag so we can test
func seedUserAndTag(t *testing.T, testDB *TestDatabase, unique string) (models.User, models.Tag) {
	t.Helper()

	user := models.User{
		ID:                      uuid.New(),
		FirstName:               "Test",
		LastName:                "User",
		Email:                   "test-" + unique + "@example.com",
		Username:                "testuser-" + unique,
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}

	tag := models.Tag{
		ID:   uuid.New(),
		Name: "tag-" + unique,
	}

	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if err := testDB.DB.Create(&tag).Error; err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	return user, tag
}

func TestGetTagFollowsByUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, tag := seedUserAndTag(t, testDB, "create-user-tag")

	createdTagFollow := models.TagFollow{
		ID:     uuid.New(),
		UserID: user.ID,
		TagID:  tag.ID,
	}

	if err := testDB.DB.Create(&createdTagFollow).Error; err != nil {
		t.Fatalf("Unable to add tag follow to table: %v", err)
	}

	resp := api.Get("/api/v1/user/tag/"+user.ID.String()+"/follows", "Authorization: Bearer mock-token")

	var response tagfollow.GetTagFollowsByUserResponse

	DecodeTo(&response, resp)

	tagIDs := []uuid.UUID{createdTagFollow.TagID}
	if !reflect.DeepEqual(response.TagIDs, tagIDs) {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestGetFollowingUsersByTag(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, tag := seedUserAndTag(t, testDB, "get-following-users-by-tag")

	createdTagFollow := models.TagFollow{
		ID:     uuid.New(),
		UserID: user.ID,
		TagID:  tag.ID,
	}

	if err := testDB.DB.Create(&createdTagFollow).Error; err != nil {
		t.Fatalf("Unable to add tag follow to table: %v", err)
	}

	resp := api.Get("/api/v1/user/tag/"+tag.ID.String()+"/users", "Authorization: Bearer mock-token")

	var response tagfollow.GetFollowingUsersByTagResponse

	DecodeTo(&response, resp)

	userIDs := []uuid.UUID{user.ID}
	if !reflect.DeepEqual(response.UserIDs, userIDs) {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestCreateTagFollow(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, tag := seedUserAndTag(t, testDB, "create-tag-follow")

	reqBody := tagfollow.CreateTagFollowBody{
		UserID: user.ID,
		TagID:  tag.ID,
	}

	resp := api.Post("/api/v1/user/tag/", reqBody, "Authorization: Bearer mock-token")

	var response tagfollow.CreateTagFollowResponse
	DecodeTo(&response, resp)

	if response.UserID != user.ID {
		t.Fatalf("Unexpected UserID in response: %s", resp.Body.String())
	}
	if response.TagID != tag.ID {
		t.Fatalf("Unexpected TagID in response: %s", resp.Body.String())
	}

	// verifying that it exists in DB
	var tagFollow models.TagFollow
	if err := testDB.DB.Where("user_id = ? AND tag_id = ?", user.ID, tag.ID).First(&tagFollow).Error; err != nil {
		t.Fatalf("TagFollow not found in DB after creation: %v", err)
	}
}

func TestDeleteTagFollow(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, tag := seedUserAndTag(t, testDB, "delete-tag-follow")

	createdTagFollow := models.TagFollow{
		ID:     uuid.New(),
		UserID: user.ID,
		TagID:  tag.ID,
	}

	if err := testDB.DB.Create(&createdTagFollow).Error; err != nil {
		t.Fatalf("Unable to add tag follow to table: %v", err)
	}

	resp := api.Delete("/api/v1/user/tag/"+createdTagFollow.ID.String(), "Authorization: Bearer mock-token")

	var response tagfollow.DeleteTagFollowResponse
	DecodeTo(&response, resp)

	if response.Message != "Tag follow was deleted successfully" {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}

	// verify it's soft deleted in DB
	var tagFollow models.TagFollow
	err := testDB.DB.Where("id = ?", createdTagFollow.ID).First(&tagFollow).Error
	if err == nil {
		t.Fatalf("TagFollow still exists in DB after deletion")
	}
}
