package routeTests

import (
	h "inside-athletics/internal/handlers/user"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func addCollegeAndSport(t *testing.T, testDB *TestDatabase) (models.College, models.Sport) {
	college := models.College{
		ID:           NortheasternID,
		Name:         "Northeastern University",
		State:        "Massachusetts",
		City:         "Boston",
		Website:      "https://www.northeastern.edu",
		DivisionRank: models.Division(1),
	}
	if err := testDB.DB.Create(&college).Error; err != nil {
		t.Fatalf("Unable to create college: %s", err.Error())
	}
	sport := models.Sport{
		ID:   SoccerID,
		Name: "Women's Soccer",
	}
	if err := testDB.DB.Create(&sport).Error; err != nil {
		t.Fatalf("Unable to create sport: %s", err.Error())
	}
	return college, sport
}

func TestGetUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	college, sport := addCollegeAndSport(t, testDB)

	user := models.User{
		ID:                      uuid.New(),
		FirstName:               "Suli",
		LastName:                "Test",
		Email:                   "suli@example.com",
		Username:                "suli",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusVerified,
		CollegeID:               &college.ID,
		SportID:                 &sport.ID,
		Division:                divisionPtr(models.DivisionI),
	}
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}
	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "user"},
	})

	resp := api.Get("/api/v1/user/"+user.ID.String(), authHeader)

	var u h.GetUserResponse
	DecodeTo(&u, resp)

	t.Log()
	if u.ID != user.ID ||
		u.FirstName != "Suli" ||
		u.LastName != "Test" ||
		u.Email != "suli@example.com" ||
		u.Username != "suli" ||
		u.AccountType != false ||
		u.College == nil ||
		u.College.Name != "Northeastern University" ||
		u.Sport == nil ||
		u.Sport.Name != "Women's Soccer" ||
		u.VerifiedAthleteStatus != models.VerifiedAthleteStatusVerified {
		t.Fatalf("Unexpected response: %+v", u)
	}
}

func TestGetUserWithCollegeAndSport(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	college, sport := addCollegeAndSport(t, testDB)

	user := models.User{
		ID:                      uuid.New(),
		FirstName:               "Suli",
		LastName:                "Test",
		Email:                   "suli@example.com",
		Username:                "suli",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusVerified,
		CollegeID:               &college.ID,
		SportID:                 &sport.ID,
		Division:                divisionPtr(models.DivisionI),
	}
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}
	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissions(t, testDB.DB, nil)

	resp := api.Get("/api/v1/user/"+user.ID.String(), authHeader)

	var u h.GetUserResponse
	DecodeTo(&u, resp)
	if u.ID != user.ID ||
		u.College == nil ||
		u.Sport == nil ||
		u.Division == nil ||
		*u.Division != models.DivisionI {
		t.Fatalf("Unexpected response: %+v", u)
	}
}

func TestGetUserNotFound(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	authHeader := authHeaderWithPermissions(t, testDB.DB, nil)

	resp := api.Get("/api/v1/user/"+NonExistentID.String(), authHeader)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("Expected 404, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestGetCurrentUserID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	userID := uuid.New()
	user := models.User{
		ID:                      userID,
		FirstName:               "Suli",
		LastName:                "Test",
		Email:                   "suli@example.com",
		Username:                "suli",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusNone,
	}
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}
	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	resp := api.Get("/api/v1/user/current", "Authorization: Bearer "+userID.String())

	var u h.GetUserResponse
	DecodeTo(&u, resp)

	if u.ID != user.ID ||
		u.FirstName != "Suli" ||
		u.LastName != "Test" ||
		u.Email != "suli@example.com" ||
		u.Username != "suli" ||
		u.AccountType != false ||
		u.VerifiedAthleteStatus != models.VerifiedAthleteStatusNone {
		t.Fatalf("Unexpected response: %+v", u)
	}
}

func TestCreateUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	college, sport := addCollegeAndSport(t, testDB)

	userID := uuid.NewString()
	payload := h.CreateUserBody{
		FirstName:             "Suli",
		LastName:              "Test",
		Email:                 "suli@example.com",
		Username:              "suli",
		Bio:                   strPtr("My bio"),
		AccountType:           true,
		SportID:               &sport.ID,
		ExpectedGradYear:      2027,
		VerifiedAthleteStatus: models.VerifiedAthleteStatusVerified,
		CollegeID:             &college.ID,
		Division:              divisionPtr(models.DivisionI),
	}

	resp := api.Post("/api/v1/user/", "Authorization: Bearer "+userID, payload)

	var u h.CreateUserResponse
	DecodeTo(&u, resp)
	if u.ID.String() != userID || u.Name != "Suli" {
		t.Fatalf("Unexpected response: %+v", u)
	}
}

func TestCreateUserWithNoneStatus(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	userID := uuid.NewString()
	payload := h.CreateUserBody{
		FirstName:             "Suli",
		LastName:              "Test",
		Email:                 "suli@example.com",
		Username:              "suli",
		Bio:                   strPtr("My bio"),
		AccountType:           true,
		ExpectedGradYear:      2027,
		VerifiedAthleteStatus: models.VerifiedAthleteStatusNone,
	}

	resp := api.Post("/api/v1/user/", "Authorization: Bearer "+userID, payload)

	var u h.CreateUserResponse
	DecodeTo(&u, resp)
	if u.ID.String() != userID || u.Name != "Suli" {
		t.Fatalf("Unexpected response: %+v", u)
	}
}

func TestUpdateUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user := models.User{
		ID:                      uuid.New(),
		FirstName:               "Suli",
		LastName:                "Test",
		Email:                   "suli@example.com",
		Username:                "suli",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}
	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionUpdate, Resource: "user"},
	})

	update := h.UpdateUserBody{
		FirstName: strPtr("Updated"),
	}

	resp := api.Patch("/api/v1/user/"+user.ID.String(), authHeader, update)

	var u h.UpdateUserResponse
	DecodeTo(&u, resp)
	if u.FirstName != "Updated" {
		t.Fatalf("Unexpected response: %+v", u)
	}
}

func TestDeleteUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user := models.User{
		ID:                      uuid.New(),
		FirstName:               "Suli",
		LastName:                "Test",
		Email:                   "suli@example.com",
		Username:                "suli",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}
	assignRoleToUser(t, testDB.DB, user.ID, getRoleID(t, testDB.DB, models.RoleUser))

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionDelete, Resource: "user"},
	})

	resp := api.Delete("/api/v1/user/"+user.ID.String(), authHeader)

	var u h.DeleteUserResponse
	DecodeTo(&u, resp)

	if u.ID != user.ID {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestAssignRoleToUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	adminRoleID := getRoleID(t, testDB.DB, models.RoleAdmin)
	adminUserID := uuid.New()
	adminUser := models.User{
		ID:                      adminUserID,
		FirstName:               "Admin",
		LastName:                "User",
		Email:                   "admin@example.com",
		Username:                "admin",
		Account_Type:            true,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&adminUser).Error; err != nil {
		t.Fatalf("failed to create admin user: %v", err)
	}
	assignRoleToUser(t, testDB.DB, adminUserID, adminRoleID)
	ensurePermissionForRole(t, testDB.DB, adminRoleID, models.PermissionCreate, "user")

	targetUserID := uuid.New()
	targetUser := models.User{
		ID:                      targetUserID,
		FirstName:               "Target",
		LastName:                "User",
		Email:                   "target@example.com",
		Username:                "target",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&targetUser).Error; err != nil {
		t.Fatalf("failed to create target user: %v", err)
	}

	moderatorRoleID := getRoleID(t, testDB.DB, models.RoleModerator)
	body := map[string]any{
		"role_id": moderatorRoleID,
	}

	resp := api.Post("/api/v1/user/"+targetUserID.String()+"/roles", body, "Authorization: Bearer "+adminUserID.String())
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var assigned h.AssignRoleResponse
	DecodeTo(&assigned, resp)
	if assigned.UserID != targetUserID || assigned.Role.ID != moderatorRoleID {
		t.Fatalf("unexpected role assignment response: %+v", assigned)
	}
}

func strPtr(v string) *string {
	return &v
}

func divisionPtr(v models.Division) *models.Division {
	return &v
}
