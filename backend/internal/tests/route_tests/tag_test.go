package routeTests

import (
	"fmt"
	tagPackage "inside-athletics/internal/handlers/tag"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestGetTagByName(t *testing.T) {
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

	resp := api.Get("/api/v1/tag/name/Hockey", "Authorization: Bearer mock-token")
	fmt.Println("Raw Response:", resp.Body.String())

	var response tagPackage.GetTagResponse

	DecodeTo(&response, resp)

	if response.Name != "Hockey" {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestGetTagByID(t *testing.T) {
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

	resp := api.Get("/api/v1/tag/"+newID.String(), "Authorization: Bearer mock-token")
	fmt.Println("Raw Response:", resp.Body.String())

	var response tagPackage.GetTagResponse

	DecodeTo(&response, resp)

	if response.Name != "Women's Basketball" {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestCreateTag(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	payload := tagPackage.CreateTagBody{
		Name: "Basketball",
	}

	resp := api.Post("/api/v1/tag/", "Authorization: Bearer mock-token", payload)
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

	resp := api.Patch("/api/v1/tag/"+tag.ID.String(), "Authorization: Bearer mock-token", update)

	var response tagPackage.UpdateTagResponse
	DecodeTo(&response, resp)
	if response.Name != "Updated" {
		t.Fatalf("Unexpected response: %+v", response)
	}
}

func TestDeleteTag(t *testing.T) {
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

	resp := api.Delete("/api/v1/tag/"+tag.ID.String(), "Authorization: Bearer mock-token")

	var response tagPackage.DeleteTagResponse
	DecodeTo(&response, resp)

	if response.ID != tag.ID {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}
