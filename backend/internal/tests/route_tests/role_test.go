package routeTests

import (
	"inside-athletics/internal/handlers/role"
	"inside-athletics/internal/models"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestRoleCRUD(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	adminRoleID := getRoleID(t, testDB.DB, models.RoleAdmin)
	adminUserID := uuid.New()
	adminUser := models.User{
		ID:                      adminUserID,
		FirstName:               "Admin",
		LastName:                "User",
		Email:                   "admin@example.com",
		Username:                "admin",
		Account_Type:            true,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
		RoleID:                  adminRoleID,
	}
	if err := testDB.DB.Create(&adminUser).Error; err != nil {
		t.Fatalf("failed to create admin user: %v", err)
	}

	createBody := role.CreateRoleRequest{Name: "coach"}
	createResp := api.Post("/api/v1/role/", createBody, "Authorization: Bearer "+adminUserID.String())
	if createResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", createResp.Code, createResp.Body.String())
	}

	var created role.RoleResponse
	DecodeTo(&created, createResp)
	if created.Name != "coach" {
		t.Fatalf("unexpected role name: %+v", created)
	}

	getResp := api.Get("/api/v1/role/"+created.ID.String(), "Authorization: Bearer "+adminUserID.String())
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", getResp.Code, getResp.Body.String())
	}

	var fetched role.RoleResponse
	DecodeTo(&fetched, getResp)
	if fetched.ID != created.ID || fetched.Name != "coach" {
		t.Fatalf("unexpected role response: %+v", fetched)
	}

	updateName := "assistant_coach"
	updateBody := role.UpdateRoleRequest{Name: &updateName}
	updateResp := api.Patch("/api/v1/role/"+created.ID.String(), updateBody, "Authorization: Bearer "+adminUserID.String())
	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", updateResp.Code, updateResp.Body.String())
	}

	var updated role.RoleResponse
	DecodeTo(&updated, updateResp)
	if string(updated.Name) != updateName {
		t.Fatalf("unexpected updated role: %+v", updated)
	}

	listResp := api.Get("/api/v1/roles/", "Authorization: Bearer "+adminUserID.String())
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", listResp.Code, listResp.Body.String())
	}

	var list role.GetAllRolesResponse
	DecodeTo(&list, listResp)
	if list.Total < 1 || len(list.Roles) == 0 {
		t.Fatalf("expected at least one role, got %+v", list)
	}

	deleteResp := api.Delete("/api/v1/role/"+created.ID.String(), "Authorization: Bearer "+adminUserID.String())
	if deleteResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", deleteResp.Code, deleteResp.Body.String())
	}
}
