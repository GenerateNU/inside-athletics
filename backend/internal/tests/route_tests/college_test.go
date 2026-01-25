package routeTests

import (
	"encoding/json"
	h "inside-athletics/internal/handlers/college"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"testing"
)

func TestGetCollege(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	college := models.College{Name: "Northeastern University"}
	userResp := testDB.DB.Create(&college)
	_, err := utils.HandleDBError(&college, userResp.Error)

	if err != nil {
		t.Fatalf("Unable to add college to table: %s", err.Error())
	}

	// college's uuid
	resp := api.Get("/api/v1/college/"+college.ID.String(), "Authorization: Bearer mock-token")

	var u h.GetCollegeResponse

	DecodeTo(&u, resp)

	if u.Name != "Northeastern University" {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestCreateCollege(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	requestBody := h.CreateCollegeRequest{
		Name:         "Northeastern University",
		State:        "Massachusetts",
		City:         "Boston",
		DivisionRank: 1,
	}

	// converting to json string format
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Unable to marshal request body: %s", err.Error())
	}

	resp := api.Post("/api/v1/college", "Authorization: Bearer mock-token", "Content-Type: application/json", string(jsonBody))

	var response h.CreateCollegeResponse
	DecodeTo(&response, resp)

	if response.Name != "Northeastern University" {
		t.Fatalf("Unexpected name: got %s, expected Northeastern University", response.Name)
	}
}

func TestUpdateCollege(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	college := models.College{
		Name:         "Northeastern University",
		State:        "Massachusetts",
		City:         "Boston",
		DivisionRank: 1,
	}
	collegeResp := testDB.DB.Create(&college)
	_, err := utils.HandleDBError(&college, collegeResp.Error)
	if err != nil {
		t.Fatalf("Unable to add college to table: %s", err.Error())
	}

	// updating
	newName := "Northeastern University - Updated"
	newState := "MA"
	requestBody := h.UpdateCollegeRequest{
		Name:  &newName,
		State: &newState,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Unable to marshal request body: %s", err.Error())
	}

	resp := api.Put("/api/v1/college/"+college.ID.String(), "Authorization: Bearer mock-token", "Content-Type: application/json", string(jsonBody))

	var response h.UpdateCollegeResponse
	DecodeTo(&response, resp)

	if response.ID != college.ID {
		t.Fatalf("Unexpected ID: got %s, expected %s", response.ID.String(), college.ID.String())
	}
	if response.Name != "Northeastern University - Updated" {
		t.Fatalf("Unexpected name: got %s, expected Northeastern University - Updated", response.Name)
	}
	if response.State != "MA" {
		t.Fatalf("Unexpected state: got %s, expected MA", response.State)
	}
	if response.City != "Boston" {
		t.Fatalf("Unexpected city: got %s, expected Boston", response.City)
	}
	if response.DivisionRank != 1 {
		t.Fatalf("Unexpected division rank: got %d, expected 1", response.DivisionRank)
	}
}

func TestDeleteCollege(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	college := models.College{
		Name:         "Northeastern University",
		State:        "Massachusetts",
		City:         "Boston",
		DivisionRank: 1,
	}
	collegeResp := testDB.DB.Create(&college)
	_, err := utils.HandleDBError(&college, collegeResp.Error)
	if err != nil {
		t.Fatalf("Unable to add college to table: %s", err.Error())
	}

	resp := api.Delete("/api/v1/college/"+college.ID.String(), "Authorization: Bearer mock-token")

	var response h.DeleteCollegeResponse
	DecodeTo(&response, resp)

	if response.Message != "College deleted successfully" {
		t.Fatalf("Unexpected message: got %s, expected College deleted successfully", response.Message)
	}
}
