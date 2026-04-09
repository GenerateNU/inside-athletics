package routeTests

import (
	surveyPackage "inside-athletics/internal/handlers/survey"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

// seedUser inserts a User row so surveys can satisfy fk_surveys_user.
func seedUser(t *testing.T, testDB *TestDatabase) *models.User {
	t.Helper()
	user := models.User{
		ID:                      uuid.New(),
		FirstName:               "Test",
		LastName:                "User",
		Email:                   uuid.NewString() + "@example.com",
		Username:                "testuser-" + uuid.NewString(),
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	return &user
}

// seedCollege inserts a College row so surveys can satisfy fk_surveys_college.
func seedCollege(t *testing.T, testDB *TestDatabase) *models.College {
	t.Helper()
	college := models.College{
		ID:           uuid.New(),
		Name:         "Test College",
		State:        "Massachusetts",
		City:         "Boston",
		Website:      "https://www.testcollege.edu",
		DivisionRank: 1,
	}
	if err := testDB.DB.Create(&college).Error; err != nil {
		t.Fatalf("Unable to seed college: %s", err.Error())
	}
	return &college
}

// seedSport inserts a Sport row so surveys can satisfy fk_surveys_sport.
func seedSport(t *testing.T, testDB *TestDatabase) *models.Sport {
	t.Helper()
	sport := models.Sport{
		ID:   uuid.New(),
		Name: "Test Sport",
	}
	if err := testDB.DB.Create(&sport).Error; err != nil {
		t.Fatalf("Unable to seed sport: %s", err.Error())
	}
	return &sport
}

// seedSurvey inserts a Survey row. userID, collegeID and sportID must already exist in the DB.
func seedSurvey(t *testing.T, testDB *TestDatabase, userID, collegeID, sportID uuid.UUID) *models.Survey {
	t.Helper()
	survey := models.Survey{
		ID:                         uuid.New(),
		UserID:                     userID,
		CollegeID:                  collegeID,
		SportID:                    sportID,
		PlayerDev:                  4,
		AcademicsAthleticsPriority: 3,
		AcademicCareerResources:    4,
		MentalHealthPriority:       3,
		Environment:                5,
		Culture:                    4,
		Transparency:               3,
	}
	dbResp := testDB.DB.Create(&survey)
	_, err := utils.HandleDBError(&survey, dbResp.Error)
	if err != nil {
		t.Fatalf("Unable to seed survey: %s", err.Error())
	}
	return &survey
}

func TestCreateSurvey(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	// mockUUID user must exist for the FK
	user := models.User{
		ID:                      uuid.MustParse(mockUUID),
		FirstName:               "Test",
		LastName:                "User",
		Email:                   "test@example.com",
		Username:                "testuser-mock",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	college := seedCollege(t, testDB)
	sport := seedSport(t, testDB)

	payload := surveyPackage.CreateSurveyRequest{
		UserID:                     uuid.MustParse(mockUUID),
		CollegeID:                  college.ID,
		SportID:                    sport.ID,
		PlayerDev:                  4,
		AcademicsAthleticsPriority: 3,
		AcademicCareerResources:    4,
		MentalHealthPriority:       3,
		Environment:                5,
		Culture:                    4,
		Transparency:               3,
	}

	resp := api.Post("/api/v1/survey/", "Authorization: Bearer "+mockUUID, payload)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response surveyPackage.SurveyResponse
	DecodeTo(&response, resp)

	if response.PlayerDev != payload.PlayerDev {
		t.Fatalf("unexpected player_dev: got %d, want %d", response.PlayerDev, payload.PlayerDev)
	}
	if response.UserID != payload.UserID {
		t.Fatalf("unexpected user_id: got %s, want %s", response.UserID, payload.UserID)
	}
}

func TestCreateSurveyInvalidRating(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	college := seedCollege(t, testDB)
	sport := seedSport(t, testDB)

	payload := surveyPackage.CreateSurveyRequest{
		UserID:                     uuid.MustParse(mockUUID),
		CollegeID:                  college.ID,
		SportID:                    sport.ID,
		PlayerDev:                  6, // out of range
		AcademicsAthleticsPriority: 3,
		AcademicCareerResources:    4,
		MentalHealthPriority:       3,
		Environment:                5,
		Culture:                    4,
		Transparency:               3,
	}

	resp := api.Post("/api/v1/survey/", "Authorization: Bearer "+mockUUID, payload)
	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422 for invalid rating, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestDeleteSurvey(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user := seedUser(t, testDB)
	college := seedCollege(t, testDB)
	sport := seedSport(t, testDB)
	survey := seedSurvey(t, testDB, user.ID, college.ID, sport.ID)

	resp := api.Delete("/api/v1/survey/"+survey.ID.String(), "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response surveyPackage.SurveyResponse
	DecodeTo(&response, resp)

	if response.ID != survey.ID {
		t.Fatalf("unexpected id: got %s, want %s", response.ID, survey.ID)
	}
}

func TestDeleteSurveyNotFound(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	resp := api.Delete("/api/v1/survey/"+uuid.New().String(), "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for missing survey, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestGetSurveysByUser(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user := seedUser(t, testDB)

	college1 := seedCollege(t, testDB)
	sport1 := seedSport(t, testDB)
	college2 := seedCollege(t, testDB)
	sport2 := seedSport(t, testDB)

	seedSurvey(t, testDB, user.ID, college1.ID, sport1.ID)
	seedSurvey(t, testDB, user.ID, college2.ID, sport2.ID)

	resp := api.Get("/api/v1/survey/user/"+user.ID.String(), "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response surveyPackage.GetSurveysByUserResponse
	DecodeTo(&response, resp)

	if response.Total != 2 {
		t.Fatalf("expected 2 surveys, got %d", response.Total)
	}
	if len(response.Surveys) != 2 {
		t.Fatalf("expected 2 surveys in body, got %d", len(response.Surveys))
	}
}

func TestGetSurveysByUserEmpty(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	resp := api.Get("/api/v1/survey/user/"+uuid.New().String(), "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200 for empty result, got %d: %s", resp.Code, resp.Body.String())
	}

	var response surveyPackage.GetSurveysByUserResponse
	DecodeTo(&response, resp)

	if response.Total != 0 {
		t.Fatalf("expected 0 surveys, got %d", response.Total)
	}
}

func TestGetAverageRatingsNoFilter(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	college := seedCollege(t, testDB)
	sport := seedSport(t, testDB)
	user1 := seedUser(t, testDB)
	user2 := seedUser(t, testDB)

	seedSurvey(t, testDB, user1.ID, college.ID, sport.ID)
	seedSurvey(t, testDB, user2.ID, college.ID, sport.ID)

	resp := api.Get("/api/v1/survey/averages", "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response surveyPackage.AverageRatingsResponse
	DecodeTo(&response, resp)

	if len(response.Averages) == 0 {
		t.Fatalf("expected at least one averages row, got 0")
	}
}

func TestGetAverageRatingsFilterBySport(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	college := seedCollege(t, testDB)
	sport := seedSport(t, testDB)
	otherSport := seedSport(t, testDB)
	user1 := seedUser(t, testDB)
	user2 := seedUser(t, testDB)

	seedSurvey(t, testDB, user1.ID, college.ID, sport.ID)
	seedSurvey(t, testDB, user2.ID, college.ID, otherSport.ID) // should not appear

	resp := api.Get("/api/v1/survey/averages?sport_id="+sport.ID.String(), "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response surveyPackage.AverageRatingsResponse
	DecodeTo(&response, resp)

	if len(response.Averages) != 1 {
		t.Fatalf("expected 1 averages row for sport filter, got %d", len(response.Averages))
	}
	if response.Averages[0].SportID != sport.ID {
		t.Fatalf("unexpected sport_id in averages: got %s, want %s", response.Averages[0].SportID, sport.ID)
	}
}

func TestGetAverageRatingsFilterByCollege(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	college := seedCollege(t, testDB)
	otherCollege := seedCollege(t, testDB)
	sport := seedSport(t, testDB)
	user1 := seedUser(t, testDB)
	user2 := seedUser(t, testDB)

	seedSurvey(t, testDB, user1.ID, college.ID, sport.ID)
	seedSurvey(t, testDB, user2.ID, otherCollege.ID, sport.ID) // should not appear

	resp := api.Get("/api/v1/survey/averages?college_id="+college.ID.String(), "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response surveyPackage.AverageRatingsResponse
	DecodeTo(&response, resp)

	if len(response.Averages) != 1 {
		t.Fatalf("expected 1 averages row for college filter, got %d", len(response.Averages))
	}
	if response.Averages[0].CollegeID != college.ID {
		t.Fatalf("unexpected college_id in averages: got %s, want %s", response.Averages[0].CollegeID, college.ID)
	}
}

func TestGetAverageRatingsFilterBySportAndCollege(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	college := seedCollege(t, testDB)
	sport := seedSport(t, testDB)
	user1 := seedUser(t, testDB)
	user2 := seedUser(t, testDB)

	seedSurvey(t, testDB, user1.ID, college.ID, sport.ID)

	noiseCollege := seedCollege(t, testDB)
	noiseSport := seedSport(t, testDB)
	seedSurvey(t, testDB, user2.ID, noiseCollege.ID, noiseSport.ID)

	url := "/api/v1/survey/averages?sport_id=" + sport.ID.String() + "&college_id=" + college.ID.String()
	resp := api.Get(url, "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response surveyPackage.AverageRatingsResponse
	DecodeTo(&response, resp)

	if len(response.Averages) != 1 {
		t.Fatalf("expected 1 averages row for sport+college filter, got %d", len(response.Averages))
	}
	row := response.Averages[0]
	if row.SportID != sport.ID || row.CollegeID != college.ID {
		t.Fatalf("unexpected row: got sport=%s college=%s", row.SportID, row.CollegeID)
	}
	if row.ResponseCount != 1 {
		t.Fatalf("expected response_count 1, got %d", row.ResponseCount)
	}
}