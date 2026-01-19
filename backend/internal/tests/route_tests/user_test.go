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

	// Not the best way to add a user.. i got lazy and didn't wanna make an endpoint
	user := models.User{Name: "Suli"}
	userResp := testDB.DB.Create(&user)
	_, err := utils.HandleDBError(&user, userResp.Error)

	if err != nil {
		t.Fatalf("Unable to add user to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/user/Suli", "Authorization: Bearer mock-token",)

	var u h.GetUserResponse

	DecodeTo(&u, resp)

	if u.Name != "Suli" {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}
