package routeTests

import (
	h "inside-athletics/internal/handlers/user"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"testing"

	"github.com/google/uuid"
)

func TestGetUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	// insert directly into DB to test
	// Don't use another endpoint to test this one - harder to tell which one is
	// incorrect if the test fails
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

	// Need to authenticate each request by passing an authorization header like this
	// when we start making endpoints that require a user-id you should add the user-id
	// you need here
	resp := api.Get("/api/v1/user/"+user.ID.String(), "Authorization: Bearer "+uuid.NewString())

	var u h.GetUserResponse

	DecodeTo(&u, resp)
	if u.ID != user.ID ||
		u.FirstName != "Suli" ||
		u.LastName != "Test" ||
		u.Email != "suli@example.com" ||
		u.Username != "suli" ||
		u.AccountType != false ||
		u.VerifiedAthleteStatus != models.VerifiedAthleteStatusPending {
		t.Fatalf("Unexpected response: %+v", u)
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
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)
	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/user/current", "Authorization: Bearer "+userID.String())

	var u h.GetCurrentUserIDResponse
	DecodeTo(&u, resp)

	if u.ID != user.ID ||
		u.FirstName != "Suli" ||
		u.LastName != "Test" ||
		u.Email != "suli@example.com" ||
		u.Username != "suli" ||
		u.AccountType != false ||
		u.VerifiedAthleteStatus != models.VerifiedAthleteStatusPending {
		t.Fatalf("Unexpected response: %+v", u)
	}
}

func TestCreateUser(t *testing.T) {
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
		Sport:                 []string{"hockey"},
		ExpectedGradYear:      2027,
		VerifiedAthleteStatus: models.VerifiedAthleteStatusPending,
		College:               strPtr("Northeastern University"),
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
		Sport:                 []string{"hockey"},
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

	update := h.UpdateUserBody{
		FirstName: strPtr("Updated"),
	}

	resp := api.Patch("/api/v1/user/"+user.ID.String(), "Authorization: Bearer "+uuid.NewString(), update)

	var u h.UpdateUserResponse
	DecodeTo(&u, resp)
	if u.Name != "Updated" {
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

	resp := api.Delete("/api/v1/user/"+user.ID.String(), "Authorization: Bearer "+uuid.NewString())

	var u h.DeleteUserResponse
	DecodeTo(&u, resp)

	if u.ID != user.ID {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func strPtr(v string) *string {
	return &v
}

func divisionPtr(v models.Division) *models.Division {
	return &v
}
