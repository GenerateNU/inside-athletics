package routeTests

import (
	"inside-athletics/internal/handlers/collegefollow"
	"inside-athletics/internal/models"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// seedUserAndCollege creates a user and college for testing
func seedUserAndCollege(t *testing.T, testDB *TestDatabase, unique string) (models.User, models.College) {
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

	college := models.College{
		ID:            uuid.New(),
		Name:          "College-" + unique,
		State:         "Massachusetts",
		City:          "Boston",
		Website:       "https://example.com",
		DivisionRank:  1,
	}

	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if err := testDB.DB.Create(&college).Error; err != nil {
		t.Fatalf("failed to create college: %v", err)
	}

	return user, college
}

func TestGetCollegeFollowsByUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, college := seedUserAndCollege(t, testDB, "get-college-follows-by-user")

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
		{Action: models.PermissionCreate, Resource: "college"},
	})

	createdCollegeFollow := models.CollegeFollow{
		ID:        uuid.New(),
		UserID:    user.ID,
		CollegeID: college.ID,
	}

	if err := testDB.DB.Create(&createdCollegeFollow).Error; err != nil {
		t.Fatalf("Unable to add college follow to table: %v", err)
	}

	resp := api.Get("/api/v1/user/college/"+user.ID.String()+"/follows", authHeader)

	var response collegefollow.GetCollegeFollowsByUserResponse

	DecodeTo(&response, resp)

	collegeIDs := []uuid.UUID{createdCollegeFollow.CollegeID}
	if !reflect.DeepEqual(response.CollegeIDs, collegeIDs) {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestGetFollowingUsersByCollege(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, college := seedUserAndCollege(t, testDB, "get-following-users-by-college")

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
		{Action: models.PermissionCreate, Resource: "college"},
	})

	createdCollegeFollow := models.CollegeFollow{
		ID:        uuid.New(),
		UserID:    user.ID,
		CollegeID: college.ID,
	}

	if err := testDB.DB.Create(&createdCollegeFollow).Error; err != nil {
		t.Fatalf("Unable to add college follow to table: %v", err)
	}

	resp := api.Get("/api/v1/user/college/"+college.ID.String()+"/users", authHeader)

	var response collegefollow.GetFollowingUsersByCollegeResponse

	DecodeTo(&response, resp)

	userIDs := []uuid.UUID{user.ID}
	if !reflect.DeepEqual(response.UserIDs, userIDs) {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestCreateCollegeFollow(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, college := seedUserAndCollege(t, testDB, "create-college-follow")

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
		{Action: models.PermissionCreate, Resource: "college"},
	})

	reqBody := collegefollow.CreateCollegeFollowBody{
		UserID:    user.ID,
		CollegeID: college.ID,
	}

	resp := api.Post("/api/v1/user/college/", reqBody, authHeader)

	var response collegefollow.CreateCollegeFollowResponse
	DecodeTo(&response, resp)

	if response.UserID != user.ID {
		t.Fatalf("Unexpected UserID in response: %s", resp.Body.String())
	}
	if response.CollegeID != college.ID {
		t.Fatalf("Unexpected CollegeID in response: %s", resp.Body.String())
	}

	var collegeFollow models.CollegeFollow
	if err := testDB.DB.Where("user_id = ? AND college_id = ?", user.ID, college.ID).First(&collegeFollow).Error; err != nil {
		t.Fatalf("CollegeFollow not found in DB after creation: %v", err)
	}
}

func TestDeleteCollegeFollow(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, college := seedUserAndCollege(t, testDB, "delete-college-follow")

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
		{Action: models.PermissionCreate, Resource: "college"},
		{Action: models.PermissionDelete, Resource: "user"},
		{Action: models.PermissionDelete, Resource: "college"},
	})

	createdCollegeFollow := models.CollegeFollow{
		ID:        uuid.New(),
		UserID:    user.ID,
		CollegeID: college.ID,
	}

	if err := testDB.DB.Create(&createdCollegeFollow).Error; err != nil {
		t.Fatalf("Unable to add college follow to table: %v", err)
	}

	resp := api.Delete("/api/v1/user/college/"+createdCollegeFollow.ID.String(), authHeader)

	var response collegefollow.DeleteCollegeFollowResponse
	DecodeTo(&response, resp)

	if response.Message != "College follow was deleted successfully" {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}

	var collegeFollow models.CollegeFollow
	err := testDB.DB.Where("id = ?", createdCollegeFollow.ID).First(&collegeFollow).Error
	if err == nil {
		t.Fatalf("CollegeFollow still exists in DB after deletion")
	}
}

func TestGetCollegeFollowsByUser_InvalidUUID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
		{Action: models.PermissionCreate, Resource: "college"},
	})

	resp := api.Get("/api/v1/user/college/not-a-valid-uuid/follows", authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid UUID, got %d: %s", resp.Code, resp.Body.String())
	}
	if resp.Code != http.StatusUnprocessableEntity {
		t.Logf("invalid UUID returned status %d (expected 422): %s", resp.Code, resp.Body.String())
	}
}

func TestGetFollowingUsersByCollege_InvalidUUID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
		{Action: models.PermissionCreate, Resource: "college"},
	})

	resp := api.Get("/api/v1/user/college/invalid-uuid/users", authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid UUID, got %d: %s", resp.Code, resp.Body.String())
	}
	if resp.Code != http.StatusUnprocessableEntity {
		t.Logf("invalid UUID returned status %d (expected 422): %s", resp.Code, resp.Body.String())
	}
}

func TestCreateCollegeFollow_DuplicateReturns409(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, college := seedUserAndCollege(t, testDB, "dup-collegefollow")

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
		{Action: models.PermissionCreate, Resource: "college"},
	})

	reqBody := collegefollow.CreateCollegeFollowBody{
		UserID:    user.ID,
		CollegeID: college.ID,
	}

	resp1 := api.Post("/api/v1/user/college/", reqBody, authHeader)
	if resp1.Code != http.StatusOK {
		t.Fatalf("first create expected 200, got %d: %s", resp1.Code, resp1.Body.String())
	}

	resp2 := api.Post("/api/v1/user/college/", reqBody, authHeader)
	if resp2.Code != http.StatusConflict {
		t.Errorf("expected 409 for duplicate college follow, got %d: %s", resp2.Code, resp2.Body.String())
	}
}

func TestCreateCollegeFollow_NonExistentUserReturnsError(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	_, college := seedUserAndCollege(t, testDB, "create-nonexistent-user-college")

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
		{Action: models.PermissionCreate, Resource: "college"},
	})

	reqBody := collegefollow.CreateCollegeFollowBody{
		UserID:    uuid.New(),
		CollegeID: college.ID,
	}

	resp := api.Post("/api/v1/user/college/", reqBody, authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected error for non-existent user, got 200: %s", resp.Body.String())
	}
}

func TestCreateCollegeFollow_NonExistentCollegeReturnsError(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, _ := seedUserAndCollege(t, testDB, "create-nonexistent-college")

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
		{Action: models.PermissionCreate, Resource: "college"},
	})

	reqBody := collegefollow.CreateCollegeFollowBody{
		UserID:    user.ID,
		CollegeID: uuid.New(),
	}

	resp := api.Post("/api/v1/user/college/", reqBody, authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected error for non-existent college, got 200: %s", resp.Body.String())
	}
}

func TestDeleteCollegeFollow_NonExistentReturns404(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
		{Action: models.PermissionCreate, Resource: "college"},
		{Action: models.PermissionDelete, Resource: "user"},
		{Action: models.PermissionDelete, Resource: "college"},
	})

	nonExistentID := uuid.New()
	resp := api.Delete("/api/v1/user/college/"+nonExistentID.String(), authHeader)

	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent college follow, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestDeleteCollegeFollow_InvalidUUID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
		{Action: models.PermissionCreate, Resource: "college"},
		{Action: models.PermissionDelete, Resource: "user"},
		{Action: models.PermissionDelete, Resource: "college"},
	})

	resp := api.Delete("/api/v1/user/college/not-a-valid-uuid", authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid UUID, got %d: %s", resp.Code, resp.Body.String())
	}
	if resp.Code != http.StatusUnprocessableEntity {
		t.Logf("invalid UUID returned status %d (expected 422): %s", resp.Code, resp.Body.String())
	}
}
