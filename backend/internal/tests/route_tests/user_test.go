package routeTests

import (
	h "inside-athletics/internal/handlers/user"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"testing"
)

func TestGetUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	// insert directly into DB to test
	// Don't use another endpoint to test this one - harder to tell which one is
	// incorrect if the test fails
	user := models.User{FirstName: "Suli"}
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)

	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}

	// Need to authenticate each request by passing an authorization header like this
	// when we start making endpoints that require a user-id you should add the user-id
	// you need here
	resp := api.Get("/api/v1/user/Suli", "Authorization: Bearer mock-token")

	var u h.GetUserResponse

	DecodeTo(&u, resp)

	if u.Name != "Suli" {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}
