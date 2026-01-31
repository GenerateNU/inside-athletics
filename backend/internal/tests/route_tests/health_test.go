package routeTests

import (
	"inside-athletics/internal/handlers/health"
	"strings"
	"testing"
)

func TestGetGreeting(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	resp := api.Get("/api/v1/health/", "Authorization: Bearer mock-token")

	var health health.HealthResponse

	DecodeTo(&health, resp)

	if !strings.Contains(health.Message, "Welcome to Inside Athletics API Version 1.0.0") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}
