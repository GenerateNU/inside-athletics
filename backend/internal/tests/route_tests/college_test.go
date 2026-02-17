package routeTests

import (
	"bytes"
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

	college := models.College{
		Name:         "Northeastern University",
		State:        "Massachusetts",
		City:         "Boston",
		Website:      "https://www.northeastern.edu",
		DivisionRank: 1,
	}
	collegeResp := testDB.DB.Create(&college)
	_, err := utils.HandleDBError(&college, collegeResp.Error)

	// make sure college was added to db
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
		Website:      "https://www.northeastern.edu",
		DivisionRank: 1,
	}

	// converting to json string format
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Unable to marshal request body: %s", err.Error())
	}

	resp := api.Post("/api/v1/college", "Authorization: Bearer mock-token", "Content-Type: application/json",
		bytes.NewReader(jsonBody))

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
		Website:      "https://www.northeastern.edu",
		DivisionRank: 1,
	}
	collegeResp := testDB.DB.Create(&college)
	_, err := utils.HandleDBError(&college, collegeResp.Error)
	if err != nil {
		t.Fatalf("Unable to add college to table: %s", err.Error())
	}

	// updating - these are the only fields that should change
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

	resp := api.Put("/api/v1/college/"+college.ID.String(),
		"Authorization: Bearer mock-token",
		"Content-Type: application/json",
		bytes.NewReader(jsonBody))

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

	// also validate other fields did not change
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
		Website:      "https://www.northeastern.edu",
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

	expectedMessage := "College " + college.ID.String() + " deleted successfully"
	if response.Message != expectedMessage {
		t.Fatalf("Unexpected message: got %s, expected %s", response.Message, expectedMessage)
	}
	if response.ID != college.ID {
		t.Fatalf("Unexpected ID: got %s, expected %s", response.ID.String(), college.ID.String())
	}
}

func TestCreateCollegeMissingName(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	requestBody := h.CreateCollegeRequest{
		State:        "Massachusetts",
		City:         "Boston",
		Website:      "https://www.northeastern.edu",
		DivisionRank: 1,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Unable to marshal request body: %s", err.Error())
	}

	resp := api.Post("/api/v1/college", "Authorization: Bearer mock-token", "Content-Type: application/json",
		bytes.NewReader(jsonBody))

	if resp.Code < 400 {
		t.Fatalf("Expected error response for missing name, got status %d: %s", resp.Code, resp.Body.String())
	}
}

func TestCreateCollegeMissingState(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	requestBody := h.CreateCollegeRequest{
		Name:         "Northeastern University",
		City:         "Boston",
		Website:      "https://www.northeastern.edu",
		DivisionRank: 1,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Unable to marshal request body: %s", err.Error())
	}

	resp := api.Post("/api/v1/college", "Authorization: Bearer mock-token", "Content-Type: application/json",
		bytes.NewReader(jsonBody))

	if resp.Code < 400 {
		t.Fatalf("Expected error response for missing state, got status %d: %s", resp.Code, resp.Body.String())
	}
}

func TestCreateCollegeMissingCity(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	requestBody := h.CreateCollegeRequest{
		Name:         "Northeastern University",
		State:        "Massachusetts",
		Website:      "https://www.northeastern.edu",
		DivisionRank: 1,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Unable to marshal request body: %s", err.Error())
	}

	resp := api.Post("/api/v1/college", "Authorization: Bearer mock-token", "Content-Type: application/json",
		bytes.NewReader(jsonBody))

	if resp.Code < 400 {
		t.Fatalf("Expected error response for missing city, got status %d: %s", resp.Code, resp.Body.String())
	}
}

func TestCreateCollegeMissingDivisionRank(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	requestBody := h.CreateCollegeRequest{
		Name:    "Northeastern University",
		State:   "Massachusetts",
		City:    "Boston",
		Website: "https://www.northeastern.edu",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Unable to marshal request body: %s", err.Error())
	}

	resp := api.Post("/api/v1/college", "Authorization: Bearer mock-token", "Content-Type: application/json",
		bytes.NewReader(jsonBody))

	if resp.Code < 400 {
		t.Fatalf("Expected error response for missing division rank, got status %d: %s", resp.Code, resp.Body.String())
	}
}

func TestCreateCollegeMissingWebsite(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	requestBody := h.CreateCollegeRequest{
		Name:         "Northeastern University",
		State:        "Massachusetts",
		City:         "Boston",
		DivisionRank: 1,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Unable to marshal request body: %s", err.Error())
	}

	resp := api.Post("/api/v1/college", "Authorization: Bearer mock-token", "Content-Type: application/json",
		bytes.NewReader(jsonBody))

	if resp.Code < 400 {
		t.Fatalf("Expected error response for missing website, got status %d: %s", resp.Code, resp.Body.String())
	}
}