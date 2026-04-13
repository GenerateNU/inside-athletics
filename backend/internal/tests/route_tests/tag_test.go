package routeTests

import (
	"inside-athletics/internal/handlers/tag"
	tagPackage "inside-athletics/internal/handlers/tag"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestGetTagByName(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	tag := models.Tag{
		ID:   uuid.New(),
		Name: "Hockey",
	}
	tagResp := testDB.DB.Create(&tag)
	_, err := utils.HandleDBError(&tag, tagResp.Error)

	if err != nil {
		t.Fatalf("Unable to add tag to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/tag/name/Hockey", "Authorization: Bearer "+mockUUID)

	var response tagPackage.GetTagResponse

	DecodeTo(&response, resp)

	if response.Name != "Hockey" {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestGetTagByID(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	newID := uuid.New()
	tag := models.Tag{
		ID:   newID,
		Name: "Women's Basketball",
	}
	tagResp := testDB.DB.Create(&tag)
	_, err := utils.HandleDBError(&tag, tagResp.Error)

	if err != nil {
		t.Fatalf("Unable to add tag to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/tag/"+newID.String(), "Authorization: Bearer "+mockUUID)

	var response tagPackage.GetTagResponse

	DecodeTo(&response, resp)

	if response.Name != "Women's Basketball" {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestCreateTag(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	payload := tagPackage.CreateTagBody{
		Name: "Basketball",
	}

	resp := api.Post("/api/v1/tag/", "Authorization: Bearer "+mockUUID, payload)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response tagPackage.CreateTagBody
	DecodeTo(&response, resp)

	if response.Name != "Basketball" {
		t.Fatalf("Unexpected response: %+v", response)
	}
}

func TestUpdateTag(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	tag := models.Tag{
		ID:   uuid.New(),
		Name: "Hockey",
	}
	tagResp := testDB.DB.Create(&tag)
	_, err := utils.HandleDBError(&tag, tagResp.Error)
	if err != nil {
		t.Fatalf("Unable to add tag to table: %s", err.Error())
	}

	update := tagPackage.UpdateTagBody{
		Name: "Updated",
	}

	resp := api.Patch("/api/v1/tag/"+tag.ID.String(), "Authorization: Bearer "+mockUUID, update)

	var response tagPackage.UpdateTagResponse
	DecodeTo(&response, resp)
	if response.Name != "Updated" {
		t.Fatalf("Unexpected response: %+v", response)
	}
}

func TestDeleteTag(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	tag := models.Tag{
		ID:   uuid.New(),
		Name: "Suli",
	}
	tagResp := testDB.DB.Create(&tag)
	_, err := utils.HandleDBError(&tag, tagResp.Error)
	if err != nil {
		t.Fatalf("Unable to add tag to table: %s", err.Error())
	}

	resp := api.Delete("/api/v1/tag/"+tag.ID.String(), "Authorization: Bearer "+mockUUID)

	var response tagPackage.DeleteTagResponse
	DecodeTo(&response, resp)

	if response.ID != tag.ID {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestTagSearch(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	tag1 := models.Tag{
		ID:   uuid.New(),
		Name: "Suli",
	}
	tagResp := testDB.DB.Create(&tag1)
	_, err := utils.HandleDBError(&tag1, tagResp.Error)
	if err != nil {
		t.Fatalf("Unable to add tag to table: %s", err.Error())
	}

	tag2 := models.Tag{
		ID:   uuid.New(),
		Name: "Erm",
	}
	tagResp2 := testDB.DB.Create(&tag2)
	_, err2 := utils.HandleDBError(&tag2, tagResp2.Error)
	if err2 != nil {
		t.Fatalf("Unable to add tag to table: %s", err2.Error())
	}

	resp := api.Get("/api/v1/tags/search?search_str=Erm", "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected code 200 but got %d", resp.Code)
	}

	var searchResults utils.SearchResults[*tag.GetTagResponse]
	DecodeTo(&searchResults, resp)

	n := len(searchResults.Results)
	if n != 1 {
		t.Fatalf("Expected only 1 search result got %d", n)
	}

	if searchResults.Results[0].Name != tag2.Name {
		t.Fatalf("Expected to get the erm tag but got %s", searchResults.Results[0].Name)
	}
}
