package routeTests

import (
	"inside-athletics/internal/handlers/tagfollow"
	"inside-athletics/internal/models"

	"net/http"
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

	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "tagfollow"},
	},
		user.ID,
	)

	createdTagFollow := models.TagFollow{
		ID:     uuid.New(),
		UserID: user.ID,
		TagID:  tag.ID,
	}

	if err := testDB.DB.Create(&createdTagFollow).Error; err != nil {
		t.Fatalf("Unable to add tag follow to table: %v", err)
	}

	resp := api.Get("/api/v1/user/tag/follows", authHeader)

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

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "tagfollow"},
	})

	createdTagFollow := models.TagFollow{
		ID:     uuid.New(),
		UserID: user.ID,
		TagID:  tag.ID,
	}

	if err := testDB.DB.Create(&createdTagFollow).Error; err != nil {
		t.Fatalf("Unable to add tag follow to table: %v", err)
	}

	resp := api.Get("/api/v1/user/tag/"+tag.ID.String()+"/users", authHeader)

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

	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "tagfollow"},
	},
		user.ID,
	)

	reqBody := tagfollow.CreateTagFollowBody{
		TagID: tag.ID,
	}

	resp := api.Post("/api/v1/user/tag/", reqBody, authHeader)

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

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionDeleteOwn, Resource: "tagfollow"},
	}, user.ID)

	createdTagFollow := models.TagFollow{
		ID:     uuid.New(),
		UserID: user.ID,
		TagID:  tag.ID,
	}

	if err := testDB.DB.Create(&createdTagFollow).Error; err != nil {
		t.Fatalf("Unable to add tag follow to table: %v", err)
	}

	resp := api.Delete("/api/v1/user/tag/"+createdTagFollow.ID.String(), authHeader)

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

func TestDeleteTagFollow_ForbiddenForOtherUsersFollow(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	owner, tag := seedUserAndTag(t, testDB, "delete-other-tag-follow-owner")
	otherUser, _ := seedUserAndTag(t, testDB, "delete-other-tag-follow-actor")

	assignRoleToUser(t, testDB.DB, owner.ID, getRoleID(t, testDB.DB, models.RoleUser))
	assignRoleToUser(t, testDB.DB, otherUser.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUserForRole(t, testDB.DB, models.RoleUser, []permissionSpec{
		{Action: models.PermissionDeleteOwn, Resource: "tagfollow"},
	}, otherUser.ID)

	createdTagFollow := models.TagFollow{
		ID:     uuid.New(),
		UserID: owner.ID,
		TagID:  tag.ID,
	}

	if err := testDB.DB.Create(&createdTagFollow).Error; err != nil {
		t.Fatalf("Unable to add tag follow to table: %v", err)
	}

	resp := api.Delete("/api/v1/user/tag/"+createdTagFollow.ID.String(), authHeader)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 when deleting another user's tag follow, got %d: %s", resp.Code, resp.Body.String())
	}
}

// Retrieving list of tags by user with an invalid tag UUID should throw errors
func TestGetTagFollowsByUser_InvalidUUID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "tagfollow"},
	})

	resp := api.Get("/api/v1/user/tag/not-a-valid-uuid/follows", authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid UUID, got %d: %s", resp.Code, resp.Body.String())
	}
	if resp.Code != http.StatusUnprocessableEntity {
		t.Logf("invalid UUID returned status %d (expected 422): %s", resp.Code, resp.Body.String())
	}
}

// Retrieving list of users by tag with an invalid tag UUID should throw an error
func TestGetFollowingUsersByTag_InvalidUUID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "tagfollow"},
	})

	resp := api.Get("/api/v1/user/tag/invalid-uuid/users", authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid UUID, got %d: %s", resp.Code, resp.Body.String())
	}
	if resp.Code != http.StatusUnprocessableEntity {
		t.Logf("invalid UUID returned status %d (expected 422): %s", resp.Code, resp.Body.String())
	}
}

// Creating a tag follow when it already exists should throw an error
func TestCreateTagFollow_DuplicateReturns409(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, tag := seedUserAndTag(t, testDB, "dup-tagfollow")

	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "tagfollow"},
	},
		user.ID,
	)

	reqBody := tagfollow.CreateTagFollowBody{
		TagID: tag.ID,
	}

	resp1 := api.Post("/api/v1/user/tag/", reqBody, authHeader)
	if resp1.Code != http.StatusOK {
		t.Fatalf("first create expected 200, got %d: %s", resp1.Code, resp1.Body.String())
	}

	resp2 := api.Post("/api/v1/user/tag/", reqBody, authHeader)
	if resp2.Code != http.StatusConflict {
		t.Errorf("expected 409 for duplicate tag follow, got %d: %s", resp2.Code, resp2.Body.String())
	}
}

// Creating tag follow with a nonexistent tag should throw an error
func TestCreateTagFollow_NonExistentTagReturnsError(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, _ := seedUserAndTag(t, testDB, "create-nonexistent-tag")

	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "tagfollow"},
	},
		user.ID,
	)

	reqBody := tagfollow.CreateTagFollowBody{
		TagID: uuid.New(),
	}

	resp := api.Post("/api/v1/user/tag/", reqBody, authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected error for non-existent tag, got 200: %s", resp.Body.String())
	}
}

// Creating a tag follow with invalid body should return error
func TestCreateTagFollow_InvalidBodyReturnsError(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "tagfollow"},
	})

	// Invalid body: invalid random UUID strings
	reqBody := map[string]any{
		"user_id": "not-a-valid-uuid",
		"tag_id":  "also-invalid",
	}
	resp := api.Post("/api/v1/user/tag/", reqBody, authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected error for invalid UUIDs in body, got 200: %s", resp.Body.String())
	}
}

// Deleting a nonexistent tag follow should throw an error
func TestDeleteTagFollow_NonExistentReturns404(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionDelete, Resource: "tagfollow"},
	})

	nonExistentID := uuid.New()
	resp := api.Delete("/api/v1/user/tag/"+nonExistentID.String(), authHeader)

	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent tag follow, got %d: %s", resp.Code, resp.Body.String())
	}
}

// Invalid UUID should throw error when deleting tag follow
func TestDeleteTagFollow_InvalidUUID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionDelete, Resource: "tagfollow"},
	})

	resp := api.Delete("/api/v1/user/tag/not-a-valid-uuid", authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid UUID, got %d: %s", resp.Code, resp.Body.String())
	}
	if resp.Code != http.StatusUnprocessableEntity {
		t.Logf("invalid UUID returned status %d (expected 422): %s", resp.Code, resp.Body.String())
	}
}
