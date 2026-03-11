package routeTests

import (
	"inside-athletics/internal/handlers/permission"
	"inside-athletics/internal/models"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestPermissionCRUD(t *testing.T) {
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
	}
	if err := testDB.DB.Create(&adminUser).Error; err != nil {
		t.Fatalf("failed to create admin user: %v", err)
	}
	assignRoleToUser(t, testDB.DB, adminUserID, adminRoleID)

	createBody := permission.CreatePermissionRequest{
		Action:   "create",
		Resource: "post",
	}
	createResp := api.Post("/api/v1/permission/", createBody, "Authorization: Bearer "+adminUserID.String())
	var created permission.PermissionResponse
	switch createResp.Code {
	case http.StatusOK:
		DecodeTo(&created, createResp)
		if created.Action != models.PermissionAction("create") || created.Resource != "post" {
			t.Fatalf("unexpected permission response: %+v", created)
		}
	case http.StatusConflict:
		var existing models.Permission
		if err := testDB.DB.Where("action = ? AND resource = ?", "create", "post").First(&existing).Error; err != nil {
			t.Fatalf("expected existing permission after conflict: %v", err)
		}
		created = permission.PermissionResponse{
			ID:       existing.ID,
			Action:   existing.Action,
			Resource: existing.Resource,
		}
	default:
		t.Fatalf("expected status 200 or 409, got %d: %s", createResp.Code, createResp.Body.String())
	}

	getResp := api.Get("/api/v1/permission/"+created.ID.String(), "Authorization: Bearer "+adminUserID.String())
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", getResp.Code, getResp.Body.String())
	}

	var fetched permission.PermissionResponse
	DecodeTo(&fetched, getResp)
	if fetched.ID != created.ID || fetched.Resource != "post" {
		t.Fatalf("unexpected permission response: %+v", fetched)
	}

	updateAction := models.PermissionUpdate
	updateResource := "post"
	updateBody := permission.UpdatePermissionRequest{
		Action:   &updateAction,
		Resource: &updateResource,
	}
	updateResp := api.Patch("/api/v1/permission/"+created.ID.String(), updateBody, "Authorization: Bearer "+adminUserID.String())
	var updated permission.PermissionResponse
	switch updateResp.Code {
	case http.StatusOK:
		DecodeTo(&updated, updateResp)
		if updated.Action != models.PermissionAction("update") {
			t.Fatalf("unexpected updated permission: %+v", updated)
		}
	case http.StatusConflict:
		var existing models.Permission
		if err := testDB.DB.Where("action = ? AND resource = ?", "update", "post").First(&existing).Error; err != nil {
			t.Fatalf("expected existing permission after conflict: %v", err)
		}
		updated = permission.PermissionResponse{
			ID:       existing.ID,
			Action:   existing.Action,
			Resource: existing.Resource,
		}
	default:
		t.Fatalf("expected status 200 or 409, got %d: %s", updateResp.Code, updateResp.Body.String())
	}

	listResp := api.Get("/api/v1/permissions/", "Authorization: Bearer "+adminUserID.String())
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", listResp.Code, listResp.Body.String())
	}

	var list permission.GetAllPermissionsResponse
	DecodeTo(&list, listResp)
	if list.Total < 1 || len(list.Permissions) == 0 {
		t.Fatalf("expected at least one permission, got %+v", list)
	}

	deleteResp := api.Delete("/api/v1/permission/"+created.ID.String(), "Authorization: Bearer "+adminUserID.String())
	if deleteResp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", deleteResp.Code, deleteResp.Body.String())
	}
}
