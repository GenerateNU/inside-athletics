package routeTests

import (
	tagPackage "inside-athletics/internal/handlers/tag"
	tagpostPackage "inside-athletics/internal/handlers/tagpost"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func CreateUserAndSportAndTag(testDB *TestDatabase, t *testing.T) models.Post {
	user := models.User{
		ID:                      JohnID,
		FirstName:               "Test",
		LastName:                "User",
		Email:                   "test-john@example.com",
		Username:                "testuser-john",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}

	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	popularity := int32(100)
	soccer := models.Sport{
		ID:         SoccerID,
		Name:       "Soccer",
		Popularity: &popularity,
	}

	if err := testDB.DB.Create(&soccer).Error; err != nil {
		t.Fatalf("failed to create sport: %v", err)
	}

	post := models.Post{
		ID:          Post1ID,
		AuthorID:    JohnID,
		SportID:     &SoccerID,
		Title:       "Looking for thoughts on NEU Fencing!",
		Content:     "My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?",
		IsAnonymous: false,
	}

	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	tag := models.Tag{
		ID:   HealthAndWellnessID,
		Name: "Health And Wellness",
	}

	if err := testDB.DB.Create(&tag).Error; err != nil {
		t.Fatalf("failed to create tag: %v", err)
	}

	return post
}

func TestGetPostsByTag(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	post := CreateUserAndSportAndTag(testDB, t)

	tagpost := models.TagPost{
		ID:     uuid.New(),
		PostID: Post1ID,
		TagID:  HealthAndWellnessID,
	}
	tagpostResp := testDB.DB.Create(&tagpost)
	_, err := utils.HandleDBError(&tagpost, tagpostResp.Error)

	if err != nil {
		t.Fatalf("Unable to add tag to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/tag/"+HealthAndWellnessID.String()+"/posts", "Authorization: Bearer "+mockUUID)

	var response tagPackage.GetPostsByTagResponse

	DecodeTo(&response, resp)

	if len(response.Posts) != 1 || response.Posts[0].ID != post.ID {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestGetTagpostByID(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	CreateUserAndSportAndTag(testDB, t)
	newId := uuid.New()
	tagpost := models.TagPost{
		ID:     newId,
		PostID: Post1ID,
		TagID:  HealthAndWellnessID,
	}
	tagpostResp := testDB.DB.Create(&tagpost)
	_, err := utils.HandleDBError(&tagpost, tagpostResp.Error)

	if err != nil {
		t.Fatalf("Unable to add tag to table: %s", err.Error())
	}

	resp := api.Get("/api/v1/post/tag/"+newId.String(), "Authorization: Bearer "+mockUUID)

	var response tagpostPackage.GetTagPostByIDResponse

	DecodeTo(&response, resp)

	if response.ID.String() != newId.String() {
		t.Fatalf("Unexpected response: %s", resp.Body.String())
	}
}

func TestCreateTagPost(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	CreateUserAndSportAndTag(testDB, t)
	payload := tagpostPackage.CreateTagPostBody{
		TagID:  HealthAndWellnessID,
		PostID: Post1ID,
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

	if response.PostID.String() != Post1ID.String() && response.TagID.String() != HealthAndWellnessID.String() {
		t.Fatalf("Unexpected response: %+v", response)
	}
}

func TestUpdateTagPost(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	CreateUserAndSportAndTag(testDB, t)
	tagpostId := uuid.New()
	tagpost := models.TagPost{
		ID:     tagpostId,
		PostID: Post1ID,
		TagID:  HealthAndWellnessID,
	}
	tagpostResp := testDB.DB.Create(&tagpost)
	_, err := utils.HandleDBError(&tagpost, tagpostResp.Error)
	if err != nil {
		t.Fatalf("Unable to add tagpost to table: %s", err.Error())
	}

	update := tagpostPackage.UpdateTagPostBody{
		PostID: Post1ID,
	}

	authHeader := authHeaderWithPermissions(t, testDB.DB, []permissionSpec{
		{Action: models.PermissionUpdate, Resource: "post"},
	})
	resp := api.Patch("/api/v1/post/tag/"+tagpost.ID.String(), authHeader, update)

	var response tagpostPackage.UpdateTagPostResponse
	DecodeTo(&response, resp)
	if response.PostID.String() != Post1ID.String() {
		t.Fatalf("Unexpected response: %+v", response)
	}
}

func TestDeleteTagPost(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	CreateUserAndSportAndTag(testDB, t)
	tagpost := models.TagPost{
		ID:     uuid.New(),
		PostID: Post1ID,
		TagID:  HealthAndWellnessID,
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
