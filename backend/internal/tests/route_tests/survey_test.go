package routeTests

import (
	surveyPackage "inside-athletics/internal/handlers/survey"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

// helpers to seed a survey into the DB
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

	payload := surveyPackage.CreateSurveyRequest{
		UserID:                     uuid.MustParse(mockUUID),
		CollegeID:                  uuid.New(),
		SportID:                    uuid.New(),
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

	payload := surveyPackage.CreateSurveyRequest{
		UserID:                     uuid.MustParse(mockUUID),
		CollegeID:                  uuid.New(),
		SportID:                    uuid.New(),
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

	survey := seedSurvey(t, testDB, uuid.MustParse(mockUUID), uuid.New(), uuid.New())

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

	userID := uuid.MustParse(mockUUID)
	seedSurvey(t, testDB, userID, uuid.New(), uuid.New())
	seedSurvey(t, testDB, userID, uuid.New(), uuid.New())

	resp := api.Get("/api/v1/survey/user/"+userID.String(), "Authorization: Bearer "+mockUUID)
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

	collegeID := uuid.New()
	sportID := uuid.New()
	seedSurvey(t, testDB, uuid.New(), collegeID, sportID)
	seedSurvey(t, testDB, uuid.New(), collegeID, sportID)

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

	sportID := uuid.New()
	otherSportID := uuid.New()
	collegeID := uuid.New()

	seedSurvey(t, testDB, uuid.New(), collegeID, sportID)
	seedSurvey(t, testDB, uuid.New(), collegeID, otherSportID) // should not appear

	resp := api.Get("/api/v1/survey/averages?sport_id="+sportID.String(), "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response surveyPackage.AverageRatingsResponse
	DecodeTo(&response, resp)

	if len(response.Averages) != 1 {
		t.Fatalf("expected 1 averages row for sport filter, got %d", len(response.Averages))
	}
	if response.Averages[0].SportID != sportID {
		t.Fatalf("unexpected sport_id in averages: got %s, want %s", response.Averages[0].SportID, sportID)
	}
}

func TestGetAverageRatingsFilterByCollege(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	collegeID := uuid.New()
	otherCollegeID := uuid.New()
	sportID := uuid.New()

	seedSurvey(t, testDB, uuid.New(), collegeID, sportID)
	seedSurvey(t, testDB, uuid.New(), otherCollegeID, sportID) // should not appear

	resp := api.Get("/api/v1/survey/averages?college_id="+collegeID.String(), "Authorization: Bearer "+mockUUID)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}

	var response surveyPackage.AverageRatingsResponse
	DecodeTo(&response, resp)

	if len(response.Averages) != 1 {
		t.Fatalf("expected 1 averages row for college filter, got %d", len(response.Averages))
	}
	if response.Averages[0].CollegeID != collegeID {
		t.Fatalf("unexpected college_id in averages: got %s, want %s", response.Averages[0].CollegeID, collegeID)
	}
}

func TestGetAverageRatingsFilterBySportAndCollege(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	collegeID := uuid.New()
	sportID := uuid.New()

	seedSurvey(t, testDB, uuid.New(), collegeID, sportID)
	seedSurvey(t, testDB, uuid.New(), uuid.New(), uuid.New()) // noise — different college+sport

	url := "/api/v1/survey/averages?sport_id=" + sportID.String() + "&college_id=" + collegeID.String()
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
	if row.SportID != sportID || row.CollegeID != collegeID {
		t.Fatalf("unexpected row: got sport=%s college=%s", row.SportID, row.CollegeID)
	}
	if row.ResponseCount != 1 {
		t.Fatalf("expected response_count 1, got %d", row.ResponseCount)
	}
}