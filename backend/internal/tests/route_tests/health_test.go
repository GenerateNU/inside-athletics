package routeTests

import (
	"inside-athletics/internal/handlers/health/types"
	"strings"
	"testing"
)

func TestGetGreeting(t *testing.T) {
	api := SetupTestAPI(t)

	resp := api.Get("/api/v1/health/")

	var health types.HealthResponse

	DecodeTo(&health, resp)

	if !strings.Contains(health.Message, "Welcome to Inside Athletics API Version 1.0.0") {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}
