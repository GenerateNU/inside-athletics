package routeTests

import (
	"inside-athletics/internal/handlers/sportfollow"
	"inside-athletics/internal/models"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// seedUserAndSport creates a user and sport for testing
func seedUserAndSport(t *testing.T, testDB *TestDatabase, unique string) (models.User, models.Sport) {
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

	sport := models.Sport{
		ID:   uuid.New(),
		Name: "sport-" + unique,
	}

	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if err := testDB.DB.Create(&sport).Error; err != nil {
		t.Fatalf("failed to create sport: %v", err)
	}

	return user, sport
}

func TestGetSportFollowsByUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, sport := seedUserAndSport(t, testDB, "get-sport-follows-by-user")

	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sportfollow"},
	},
		user.ID,
	)

	createdSportFollow := models.SportFollow{
		ID:      uuid.New(),
		UserID:  user.ID,
		SportID: sport.ID,
	}

	if err := testDB.DB.Create(&createdSportFollow).Error; err != nil {
		t.Fatalf("Unable to add sport follow to table: %v", err)
	}

	resp := api.Get("/api/v1/user/sport/follows", authHeader)

	var response sportfollow.GetSportFollowsByUserResponse

	DecodeTo(&response, resp)

	sportIDs := []uuid.UUID{createdSportFollow.SportID}
	if !reflect.DeepEqual(response.SportIDs, sportIDs) {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestGetFollowingUsersBySport(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, sport := seedUserAndSport(t, testDB, "get-following-users-by-sport")

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sportfollow"},
	})

	createdSportFollow := models.SportFollow{
		ID:      uuid.New(),
		UserID:  user.ID,
		SportID: sport.ID,
	}

	if err := testDB.DB.Create(&createdSportFollow).Error; err != nil {
		t.Fatalf("Unable to add sport follow to table: %v", err)
	}

	resp := api.Get("/api/v1/user/sport/"+sport.ID.String()+"/users", authHeader)

	var response sportfollow.GetFollowingUsersBySportResponse

	DecodeTo(&response, resp)

	userIDs := []uuid.UUID{user.ID}
	if !reflect.DeepEqual(response.UserIDs, userIDs) {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestCreateSportFollow(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, sport := seedUserAndSport(t, testDB, "create-sport-follow")

	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sportfollow"},
	},
		user.ID,
	)

	reqBody := sportfollow.CreateSportFollowBody{
		SportID: sport.ID,
	}

	resp := api.Post("/api/v1/user/sport/", reqBody, authHeader)

	var response sportfollow.CreateSportFollowResponse
	DecodeTo(&response, resp)

	if response.UserID != user.ID {
		t.Fatalf("Unexpected UserID in response: %s", resp.Body.String())
	}
	if response.SportID != sport.ID {
		t.Fatalf("Unexpected SportID in response: %s", resp.Body.String())
	}

	var sportFollow models.SportFollow
	if err := testDB.DB.Where("user_id = ? AND sport_id = ?", user.ID, sport.ID).First(&sportFollow).Error; err != nil {
		t.Fatalf("SportFollow not found in DB after creation: %v", err)
	}
}

func TestDeleteSportFollow(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, sport := seedUserAndSport(t, testDB, "delete-sport-follow")

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionDeleteOwn, Resource: "sportfollow"},
	}, user.ID)

	createdSportFollow := models.SportFollow{
		ID:      uuid.New(),
		UserID:  user.ID,
		SportID: sport.ID,
	}

	if err := testDB.DB.Create(&createdSportFollow).Error; err != nil {
		t.Fatalf("Unable to add sport follow to table: %v", err)
	}

	resp := api.Delete("/api/v1/user/sport/"+createdSportFollow.ID.String(), authHeader)

	var response sportfollow.DeleteSportFollowResponse
	DecodeTo(&response, resp)

	if response.Message != "Sport follow was deleted successfully" {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}

	var sportFollow models.SportFollow
	err := testDB.DB.Where("id = ?", createdSportFollow.ID).First(&sportFollow).Error
	if err == nil {
		t.Fatalf("SportFollow still exists in DB after deletion")
	}
}

func TestDeleteSportFollow_ForbiddenForOtherUsersFollow(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	owner, sport := seedUserAndSport(t, testDB, "delete-other-sport-follow-owner")
	otherUser, _ := seedUserAndSport(t, testDB, "delete-other-sport-follow-actor")

	assignRoleToUser(t, testDB.DB, owner.ID, getRoleID(t, testDB.DB, models.RoleUser))
	assignRoleToUser(t, testDB.DB, otherUser.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUserForRole(t, testDB.DB, models.RoleUser, []permissionSpec{
		{Action: models.PermissionDeleteOwn, Resource: "sportfollow"},
	}, otherUser.ID)

	createdSportFollow := models.SportFollow{
		ID:      uuid.New(),
		UserID:  owner.ID,
		SportID: sport.ID,
	}

	if err := testDB.DB.Create(&createdSportFollow).Error; err != nil {
		t.Fatalf("Unable to add sport follow to table: %v", err)
	}

	resp := api.Delete("/api/v1/user/sport/"+createdSportFollow.ID.String(), authHeader)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 when deleting another user's sport follow, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestDeleteSportFollow_ForbiddenForModeratorOnOtherUsersFollow(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	owner, sport := seedUserAndSport(t, testDB, "delete-other-sport-follow-owner-moderator")
	moderatorUser, _ := seedUserAndSport(t, testDB, "delete-other-sport-follow-moderator")

	assignRoleToUser(t, testDB.DB, owner.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUserForRole(t, testDB.DB, models.RoleModerator, []permissionSpec{
		{Action: models.PermissionDeleteOwn, Resource: "sportfollow"},
	}, moderatorUser.ID)

	createdSportFollow := models.SportFollow{
		ID:      uuid.New(),
		UserID:  owner.ID,
		SportID: sport.ID,
	}

	if err := testDB.DB.Create(&createdSportFollow).Error; err != nil {
		t.Fatalf("Unable to add sport follow to table: %v", err)
	}

	resp := api.Delete("/api/v1/user/sport/"+createdSportFollow.ID.String(), authHeader)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403 when moderator deletes another user's sport follow, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestGetSportFollowsByUser_InvalidUUID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sportfollow"},
	})

	resp := api.Get("/api/v1/user/sport/not-a-valid-uuid/follows", authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid UUID, got %d: %s", resp.Code, resp.Body.String())
	}
	if resp.Code != http.StatusUnprocessableEntity {
		t.Logf("invalid UUID returned status %d (expected 422): %s", resp.Code, resp.Body.String())
	}
}

func TestGetFollowingUsersBySport_InvalidUUID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sportfollow"},
	})

	resp := api.Get("/api/v1/user/sport/invalid-uuid/users", authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid UUID, got %d: %s", resp.Code, resp.Body.String())
	}
	if resp.Code != http.StatusUnprocessableEntity {
		t.Logf("invalid UUID returned status %d (expected 422): %s", resp.Code, resp.Body.String())
	}
}

func TestCreateSportFollow_DuplicateReturns409(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, sport := seedUserAndSport(t, testDB, "dup-sportfollow")

	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sportfollow"},
	},
		user.ID,
	)

	reqBody := sportfollow.CreateSportFollowBody{
		SportID: sport.ID,
	}

	resp1 := api.Post("/api/v1/user/sport/", reqBody, authHeader)
	if resp1.Code != http.StatusOK {
		t.Fatalf("first create expected 200, got %d: %s", resp1.Code, resp1.Body.String())
	}

	resp2 := api.Post("/api/v1/user/sport/", reqBody, authHeader)
	if resp2.Code != http.StatusConflict {
		t.Errorf("expected 409 for duplicate sport follow, got %d: %s", resp2.Code, resp2.Body.String())
	}
}

func TestCreateSportFollow_NonExistentSportReturnsError(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	user, _ := seedUserAndSport(t, testDB, "create-nonexistent-sport")

	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissionsGivenUser(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "sportfollow"},
	},
		user.ID,
	)

	reqBody := sportfollow.CreateSportFollowBody{
		SportID: uuid.New(),
	}

	resp := api.Post("/api/v1/user/sport/", reqBody, authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected error for non-existent sport, got 200: %s", resp.Body.String())
	}
}

func TestDeleteSportFollow_NonExistentReturns404(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionDelete, Resource: "sportfollow"},
	})

	nonExistentID := uuid.New()
	resp := api.Delete("/api/v1/user/sport/"+nonExistentID.String(), authHeader)

	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent sport follow, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestDeleteSportFollow_InvalidUUID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionDelete, Resource: "sportfollow"},
	})

	resp := api.Delete("/api/v1/user/sport/not-a-valid-uuid", authHeader)

	if resp.Code == http.StatusOK {
		t.Fatalf("expected non-200 for invalid UUID, got %d: %s", resp.Code, resp.Body.String())
	}
	if resp.Code != http.StatusUnprocessableEntity {
		t.Logf("invalid UUID returned status %d (expected 422): %s", resp.Code, resp.Body.String())
	}
}
