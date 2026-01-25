package routeTests

import (
    "inside-athletics/internal/handlers/sport"
    "net/http"
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
    sportDB := &SportDB{db: testDB.DB}

    createdSport, err := sportDB.CreateSport("Women's Basketball", 100000)
    if err != nil {
        t.Fatalf("failed to create sport", err)
    }

    resp := api.Get("/api/v1/sport/" + createdSport.ID.String(), "Authorization: Bearer mock-token")

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
	sportDB := &SportDB{db: testDB.DB}

    createdSport, err := sportDB.CreateSport("Women's Basketball", 100000)
    if err != nil {
        t.Fatalf("failed to create sport", err)
    }

    resp := api.Get("/api/v1/sport/" + createdSport.Name.String(), "Authorization: Bearer mock-token")

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
	sportDB := &SportDB{db: testDB.DB}

    createdSport1, err1 := sportDB.CreateSport("Women's Basketball", 100000)
    if err1 != nil {
        t.Fatalf("failed to create sport 1", err1)
    }

	createdSport2, err2 := sportDB.CreateSport("Women's Hockey", 200000)
    if err2 != nil {
        t.Fatalf("failed to create sport 2", err2)
    }

    resp := api.Get("/api/v1/sports/", "Authorization: Bearer mock-token")

    var sports []sport.SportResponse
    DecodeTo(&sports, resp)

    if sports[0].Name != createdSport1.Name {
        t.Errorf("expected name Basketball, got %s", sports[0].Name)
    }

    if sports[0].Popularity != createdSport1.Popularity {
        t.Errorf("expected popularity 100000, got %d", sports[0].Popularity)
    }

	if sports[1].Name != createdSport2.Name {
        t.Errorf("expected name Basketball, got %s", sports[1].Name)
    }

    if sports[1].Popularity != createdSport2.Popularity {
        t.Errorf("expected popularity 100000, got %d", sports[1].Popularity)
    }
}

func TestUpdateSport(t *testing.T) {
    testDB := SetupTestDB(t)
    defer testDB.Teardown(t)
    api := testDB.API
	sportDB := &SportDB{db: testDB.DB}

    createdSport, err := sportDB.CreateSport("Women's Basketball", 100000)
    if err != nil {
        t.Fatalf("failed to create sport", err)
    }

    resp := api.Patch("/api/v1/sport/" + createdSport.ID.String(), "Authorization: Bearer mock-token")

    var sport sport.SportResponse
    DecodeTo(&sport, resp)

    if sport.Name != "Men's Basketball" {
        t.Errorf("expected name Men's Basketball, got %s", sport.Name)
    }

    if sport.Popularity != 200000 {
        t.Errorf("expected popularity 200000, got %d", sport.Popularity)
    }
}

func TestDeleteSport(t *testing.T) {
    testDB := SetupTestDB(t)
    defer testDB.Teardown(t)
    api := testDB.API
	sportDB := &SportDB{db: testDB.DB}

    createdSport, err := sportDB.CreateSport("Women's Basketball", 100000)
    if err != nil {
        t.Fatalf("failed to create sport", err)
    }

    resp := api.Delete("/api/v1/sport/" + createdSport.ID.String(), "Authorization: Bearer mock-token")

    var sport sport.SportResponse
    DecodeTo(&sport, resp)

	resp = api.Get("/api/v1/sport/"+createdSport.ID.String(), "Authorization: Bearer mock-token")
    if resp.Code != http.StatusNotFound {
        t.Errorf("expected 404 after delete, got %d", resp.Code)
    }
 
}