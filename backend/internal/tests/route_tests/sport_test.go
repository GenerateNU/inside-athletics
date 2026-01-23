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

	resp := api.Get("/api/v1/health/", "Authorization: Bearer mock-token",)

	var health health.HealthResponse

	DecodeTo(&health, resp)

	if !strings.Contains(health.Message, "Welcome to Inside Athletics API Version 1.0.0") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestGetSportById(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	resp := api.Get("/api/v1/sport/", "Authorization: Bearer mock-token",)

	var sport health.HealthResponse

	DecodeTo(&health, resp)

	if !strings.Contains(health.Message, "Welcome to Inside Athletics API Version 1.0.0") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestGetSportByName(t *testing.T) {
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

func TestGetAllSports(t *testing.T) {
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