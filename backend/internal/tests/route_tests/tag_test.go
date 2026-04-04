package routeTests

import (
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
		Type: models.TagTypeSports,
	}
	tagResp := testDB.DB.Create(&tag)
	_, err := utils.HandleDBError(&tag, tagResp.Error)
	if err != nil {
		t.Fatalf("unable to add tag to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/tag/name/Hockey", "Authorization: Bearer "+mockUUID)

	var response tagPackage.GetTagResponse
	DecodeTo(&response, resp)

	if response.Name != "Hockey" {
		t.Fatalf("unexpected response: %s", resp.Body.String())
	}
	if response.Type != models.TagTypeSports {
		t.Fatalf("unexpected type: %s", response.Type)
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
		Type: models.TagTypeSports,
	}
	tagResp := testDB.DB.Create(&tag)
	_, err := utils.HandleDBError(&tag, tagResp.Error)
	if err != nil {
		t.Fatalf("unable to add tag to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/tag/"+newID.String(), "Authorization: Bearer "+mockUUID)

	var response tagPackage.GetTagResponse
	DecodeTo(&response, resp)

	if response.Name != "Women's Basketball" {
		t.Fatalf("unexpected response: %s", resp.Body.String())
	}
	if response.Type != models.TagTypeSports {
		t.Fatalf("unexpected type: %s", response.Type)
	}
}

func TestGetTagByType(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	newID1 := uuid.New()
	tag1 := models.Tag{
		ID:   newID1,
		Name: "Women's Basketball",
		Type: models.TagTypeSports,
	}
	tagResp1 := testDB.DB.Create(&tag1)
	_, err := utils.HandleDBError(&tag1, tagResp1.Error)
	if err != nil {
		t.Fatalf("unable to add tag to table: %s", err.Error())
	}

	newID2 := uuid.New()
	tag2 := models.Tag{
		ID:   newID2,
		Name: "D1",
		Type: models.TagTypeDivisions,
	}
	tagResp2 := testDB.DB.Create(&tag2)
	_, err = utils.HandleDBError(&tag2, tagResp2.Error)
	if err != nil {
		t.Fatalf("unable to add tag to table: %s", err.Error())
	}

	newID3 := uuid.New()
	tag3 := models.Tag{
		ID:   newID3,
		Name: "D2",
		Type: models.TagTypeDivisions,
	}
	tagResp3 := testDB.DB.Create(&tag3)
	_, err = utils.HandleDBError(&tag3, tagResp3.Error)
	if err != nil {
		t.Fatalf("unable to add tag to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/tag/type/sports", "Authorization: Bearer "+mockUUID)

	var response []tagPackage.GetTagResponse
	DecodeTo(&response, resp)

	if len(response) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(response))
	}
	if response[0].Name != "Women's Basketball" {
		t.Fatalf("unexpected response: %s", resp.Body.String())
	}
	if response[0].Type != models.TagTypeSports {
		t.Fatalf("unexpected type: %s", response[0].Type)
	}
}

func TestCreateTag(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	payload := tagPackage.CreateTagBody{
		Name: "Basketball",
		Type: models.TagTypeSports,
	}

	resp := api.Post("/api/v1/tag/", "Authorization: Bearer "+mockUUID, payload)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response tagPackage.CreateTagBody
	DecodeTo(&response, resp)

	if response.Name != "Basketball" {
		t.Fatalf("unexpected response: %+v", response)
	}
	if response.Type != models.TagTypeSports {
		t.Fatalf("unexpected type: %s", response.Type)
	}
}

func TestUpdateTag(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	tag := models.Tag{
		ID:   uuid.New(),
		Name: "Hockey",
		Type: models.TagTypeSports,
	}
	tagResp := testDB.DB.Create(&tag)
	_, err := utils.HandleDBError(&tag, tagResp.Error)
	if err != nil {
		t.Fatalf("unable to add tag to table: %s", err.Error())
	}

	update := tagPackage.UpdateTagBody{
		Name: "Updated",
		Type: models.TagTypeDivisions,
	}

	resp := api.Patch("/api/v1/tag/"+tag.ID.String(), "Authorization: Bearer "+mockUUID, update)

	var response tagPackage.UpdateTagResponse
	DecodeTo(&response, resp)

	if response.Name != "Updated" {
		t.Fatalf("unexpected response: %+v", response)
	}
	if response.Type != models.TagTypeDivisions {
		t.Fatalf("unexpected type: %s", response.Type)
	}
}

func TestDeleteTag(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	tag := models.Tag{
		ID:   uuid.New(),
		Name: "Suli",
		Type: models.TagTypeSports,
	}
	tagResp := testDB.DB.Create(&tag)
	_, err := utils.HandleDBError(&tag, tagResp.Error)
	if err != nil {
		t.Fatalf("unable to add tag to table: %s", err.Error())
	}

	resp := api.Delete("/api/v1/tag/"+tag.ID.String(), "Authorization: Bearer "+mockUUID)

	var response tagPackage.DeleteTagResponse
	DecodeTo(&response, resp)

	if response.ID != tag.ID {
		t.Fatalf("unexpected response: %s", resp.Body.String())
	}
}