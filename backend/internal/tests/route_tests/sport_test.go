package routeTests

import (
	"inside-athletics/internal/handlers/sport"
	"inside-athletics/internal/models"
	"net/http"
	"testing"
)

func TestCreateSport(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Sport{}); err != nil {
		t.Fatalf("failed to migrate sports table: %v", err)
	}

	sport.Route(testDB.API, testDB.DB)
	api := testDB.API

	popularity := int32(100000)
	body := map[string]any{
		"name":       "Women's Basketball",
		"popularity": popularity,
	}

	resp := api.Post("/api/v1/sport/", body, "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result sport.SportResponse
	DecodeTo(&result, resp)

	if result.Name != "Women's Basketball" {
		t.Errorf("expected name Women's Basketball, got %s", result.Name)
	}

	if result.Popularity == nil || *result.Popularity != 100000 {
		t.Errorf("expected popularity 100000, got %v", result.Popularity)
	}
}

func TestGetSportById(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Sport{}); err != nil {
		t.Fatalf("failed to migrate sports table: %v", err)
	}

	sport.Route(testDB.API, testDB.DB)
	api := testDB.API
	sportDB := sport.NewSportDB(testDB.DB)

	popularity := int32(100000)
	createdSport, err := sportDB.CreateSport("Women's Basketball", &popularity)
	if err != nil {
		t.Fatalf("failed to create sport: %v", err)
	}

	resp := api.Get("/api/v1/sport/"+createdSport.ID.String(), "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result sport.SportResponse
	DecodeTo(&result, resp)

	if result.Name != "Women's Basketball" {
		t.Errorf("expected name Women's Basketball, got %s", result.Name)
	}

	if result.Popularity == nil || *result.Popularity != 100000 {
		t.Errorf("expected popularity 100000, got %v", result.Popularity)
	}
}

func TestGetSportByName(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Sport{}); err != nil {
		t.Fatalf("failed to migrate sports table: %v", err)
	}

	sport.Route(testDB.API, testDB.DB)
	api := testDB.API
	sportDB := sport.NewSportDB(testDB.DB)

	popularity := int32(100000)
	_, err := sportDB.CreateSport("Women's Basketball", &popularity)
	if err != nil {
		t.Fatalf("failed to create sport: %v", err)
	}

	resp := api.Get("/api/v1/sport/by-name/Women's Basketball", "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result sport.SportResponse
	DecodeTo(&result, resp)

	if result.Name != "Women's Basketball" {
		t.Errorf("expected name Women's Basketball, got %s", result.Name)
	}

	if result.Popularity == nil || *result.Popularity != 100000 {
		t.Errorf("expected popularity 100000, got %v", result.Popularity)
	}
}

func TestGetAllSports(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Sport{}); err != nil {
		t.Fatalf("failed to migrate sports table: %v", err)
	}

	sport.Route(testDB.API, testDB.DB)
	api := testDB.API
	sportDB := sport.NewSportDB(testDB.DB)

	pop1 := int32(100000)
	pop2 := int32(200000)
	_, err1 := sportDB.CreateSport("Women's Basketball", &pop1)
	if err1 != nil {
		t.Fatalf("failed to create sport 1: %v", err1)
	}

	_, err2 := sportDB.CreateSport("Women's Hockey", &pop2)
	if err2 != nil {
		t.Fatalf("failed to create sport 2: %v", err2)
	}

	resp := api.Get("/api/v1/sports/", "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result sport.GetAllSportsResponse
	DecodeTo(&result, resp)

	if result.Total < 2 {
		t.Errorf("expected at least 2 sports, got %d", result.Total)
	}

	if len(result.Sports) < 2 {
		t.Errorf("expected at least 2 sports in response, got %d", len(result.Sports))
	}
}

func TestUpdateSport(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Sport{}); err != nil {
		t.Fatalf("failed to migrate sports table: %v", err)
	}

	sport.Route(testDB.API, testDB.DB)
	api := testDB.API
	sportDB := sport.NewSportDB(testDB.DB)

	popularity := int32(100000)
	createdSport, err := sportDB.CreateSport("Women's Basketball", &popularity)
	if err != nil {
		t.Fatalf("failed to create sport: %v", err)
	}

	updateBody := map[string]any{
		"name":       "Men's Basketball",
		"popularity": int32(200000),
	}

	resp := api.Patch("/api/v1/sport/"+createdSport.ID.String(), updateBody, "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var result sport.SportResponse
	DecodeTo(&result, resp)

	if result.Name != "Men's Basketball" {
		t.Errorf("expected name Men's Basketball, got %s", result.Name)
	}

	if result.Popularity == nil || *result.Popularity != 200000 {
		t.Errorf("expected popularity 200000, got %v", result.Popularity)
	}
}

func TestDeleteSport(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	if err := testDB.DB.AutoMigrate(&models.Sport{}); err != nil {
		t.Fatalf("failed to migrate sports table: %v", err)
	}

	sport.Route(testDB.API, testDB.DB)
	api := testDB.API
	sportDB := sport.NewSportDB(testDB.DB)

	popularity := int32(100000)
	createdSport, err := sportDB.CreateSport("Women's Basketball", &popularity)
	if err != nil {
		t.Fatalf("failed to create sport: %v", err)
	}

	resp := api.Delete("/api/v1/sport/"+createdSport.ID.String(), "Authorization: Bearer mock-token")

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	// Verify the sport is deleted
	getResp := api.Get("/api/v1/sport/"+createdSport.ID.String(), "Authorization: Bearer mock-token")
	if getResp.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", getResp.Code)
	}
}
