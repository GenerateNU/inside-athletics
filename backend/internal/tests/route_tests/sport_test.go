package routeTests

import (
	"inside-athletics/internal/handlers/sport"
	"strings"
	"testing"
)

func TestCreateSport(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	body := map[string]any{
    "name":       "Women's Basketball",
    "popularity": 100000,
	}

	resp := api.Post("/api/v1/sport/", body, "Authorization: Bearer mock-token",)
	var sport sport.SportResponse
	DecodeTo(&sport, resp)

	if sport.Name != "Women's Basketball" {
    	t.Errorf("expected name Women's Basketball, got %s", sport.Name)
	}

	if sport.Popularity != 100000 {
		t.Errorf("expected popularity 10, got %d", sport.Popularity)
	}
}

func TestGetSportById(t *testing.T) {
    testDB := SetupTestDB(t)
    defer testDB.Teardown(t)
    api := testDB.API

    createdSport, err := api.CreateSport("Women's Basketball", 100000)
    if err != nil {
        t.Fatalf("failed to create sport", err)
    }

    resp := api.Get("/api/v1/sport/" + createdSport.ID.String(), "Authorization: Bearer mock-token")
    if resp.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", resp.Code)
    }

    var sport sport.SportResponse
    DecodeTo(&sport, resp)

    if sport.Name != "Women's Basketball" {
        t.Errorf("expected name Women's Basketball, got %s", sport.Name)
    }

    if sport.Popularity != 100000 {
        t.Errorf("expected popularity 100000, got %d", sport.Popularity)
    }
}

func TestGetSportByName(t *testing.T) {
    testDB := SetupTestDB(t)
    defer testDB.Teardown(t)
    api := testDB.API

    createdSport, err := api.CreateSport("Women's Basketball", 100000)
    if err != nil {
        t.Fatalf("failed to create sport", err)
    }

    resp := api.Get("/api/v1/sport/" + createdSport.Name.String(), "Authorization: Bearer mock-token")
    if resp.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", resp.Code)
    }

    var sport sport.SportResponse
    DecodeTo(&sport, resp)

    if sport.Name != "Women's Basketball" {
        t.Errorf("expected name Women's Basketball, got %s", sport.Name)
    }

    if sport.Popularity != 100000 {
        t.Errorf("expected popularity 100000, got %d", sport.Popularity)
    }
}

func TestGetAllSports(t *testing.T) {
    testDB := SetupTestDB(t)
    defer testDB.Teardown(t)
    api := testDB.API

    createdSport1, err1 := api.CreateSport("Basketball", 100000)
    if err1 != nil {
        t.Fatalf("failed to create sport 1", err)
    }

	createdSport2, err2 := api.CreateSport("Hockey", 100000)
    if err2 != nil {
        t.Fatalf("failed to create sport 1", err)
    }

    resp := api.Get("/api/v1/sport/" + createdSport.ID.String() "Authorization: Bearer mock-token")
    if resp.Code != http.StatusOK {
        t.Fatalf("expected status 200, got %d", resp.Code)
    }

    var sport sport.SportResponse
    DecodeTo(&sport, resp)

    if sport.Name != "Basketball" {
        t.Errorf("expected name Basketball, got %s", sport.Name)
    }

    if sport.Popularity != 100000 {
        t.Errorf("expected popularity 100000, got %d", sport.Popularity)
    }
}

func TestUpdateSport(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	resp := api.Get("/api/v1/health/", "Authorization: Bearer mock-token",)

	var health health.HealthResponse

	DecodeTo(&health, resp)

	if !strings.Contains(health.Message, "Welcome to Inside Athletics API Version 1.0.0") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestDeleteSport(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	resp := api.Get("/api/v1/health/", "Authorization: Bearer mock-token",)

	var health health.HealthResponse

	DecodeTo(&health, resp)

	if !strings.Contains(health.Message, "Welcome to Inside Athletics API Version 1.0.0") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}