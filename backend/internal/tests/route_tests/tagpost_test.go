package routeTests

import (
	postPackage "inside-athletics/internal/handlers/post"
	tagPackage "inside-athletics/internal/handlers/tag"
	tagpostPackage "inside-athletics/internal/handlers/tagpost"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestGetPostsByTag(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	postId := uuid.New()
	tagId := uuid.New()
	tagpost := models.TagPost{
		ID:     uuid.New(),
		PostID: postId,
		TagID:  tagId,
	}
	tagpostResp := testDB.DB.Create(&tagpost)
	_, err := utils.HandleDBError(&tagpost, tagpostResp.Error)

	if err != nil {
		t.Fatalf("Unable to add tag to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/tag/"+tagId.String()+"/posts", "Authorization: Bearer mock-token")

	var response tagPackage.GetPostsByTagResponse

	DecodeTo(&response, resp)

	postIds := []uuid.UUID{postId}
	if !reflect.DeepEqual(response.PostIDs, postIds) {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestGetTagsByPost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	postId := uuid.New()
	tagId := uuid.New()
	tagpost := models.TagPost{
		ID:     uuid.New(),
		PostID: postId,
		TagID:  tagId,
	}
	tagpostResp := testDB.DB.Create(&tagpost)
	_, err := utils.HandleDBError(&tagpost, tagpostResp.Error)

	if err != nil {
		t.Fatalf("Unable to add tag to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/post/"+postId.String()+"/tags", "Authorization: Bearer mock-token")

	var response postPackage.GetTagsByPostResponse

	DecodeTo(&response, resp)

	tagIds := []uuid.UUID{tagId}
	if !reflect.DeepEqual(response.TagIDs, tagIds) {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestGetTagpostByID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	newID := uuid.New()
	tagpost := models.TagPost{
		ID:     newID,
		PostID: uuid.New(),
		TagID:  uuid.New(),
	}
	tagpostResp := testDB.DB.Create(&tagpost)
	_, err := utils.HandleDBError(&tagpost, tagpostResp.Error)

	if err != nil {
		t.Fatalf("Unable to add tag to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/post/tag/"+newID.String(), "Authorization: Bearer mock-token")

	var response tagpostPackage.GetTagPostByIDResponse

	DecodeTo(&response, resp)

	if response.ID.String() != newID.String() {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestCreateTagPost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	tagId := uuid.New()
	postId := uuid.New()
	payload := tagpostPackage.CreateTagPostBody{
		TagID:  tagId,
		PostID: postId,
	}

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionCreate, Resource: "post"},
	})
	resp := api.Post("/api/v1/post/tag/", authHeader, payload)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response tagpostPackage.CreateTagPostsResponse
	DecodeTo(&response, resp)

	if response.PostID.String() != postId.String() && response.TagID.String() != tagId.String() {
		t.Fatalf("Unexpected response: %+v", response)
	}
}

func TestUpdateTagPost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	tagpostId := uuid.New()
	tagpost := models.TagPost{
		ID:     tagpostId,
		PostID: uuid.New(),
		TagID:  uuid.New(),
	}
	tagpostResp := testDB.DB.Create(&tagpost)
	_, err := utils.HandleDBError(&tagpost, tagpostResp.Error)
	if err != nil {
		t.Fatalf("Unable to add tagpost to table: %s", err.Error())
	}

	updatedId := uuid.New()
	update := tagpostPackage.UpdateTagPostBody{
		PostID: updatedId,
	}

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionUpdate, Resource: "post"},
	})
	resp := api.Patch("/api/v1/post/tag/"+tagpost.ID.String(), authHeader, update)

	var response tagpostPackage.UpdateTagPostResponse
	DecodeTo(&response, resp)
	if response.PostID.String() != updatedId.String() {
		t.Fatalf("Unexpected response: %+v", response)
	}
}

func TestDeleteTagPost(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	tagpost := models.TagPost{
		ID:     uuid.New(),
		PostID: uuid.New(),
		TagID:  uuid.New(),
	}
	tagResp := testDB.DB.Create(&tagpost)
	_, err := utils.HandleDBError(&tagpost, tagResp.Error)
	if err != nil {
		t.Fatalf("Unable to add tag to table: %s", err.Error())
	}

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionDelete, Resource: "post"},
	})
	resp := api.Delete("/api/v1/post/tag/"+tagpost.ID.String(), authHeader)

	var response tagpostPackage.DeleteTagPostResponse
	DecodeTo(&response, resp)

	if response.ID != tagpost.ID {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}
