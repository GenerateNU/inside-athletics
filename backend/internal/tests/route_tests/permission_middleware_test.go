package routeTests

import (
	"inside-athletics/internal/models"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestPermissionMiddleware_DeniesWithoutRolePermission(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	roleID := getRoleID(t, testDB.DB, models.RoleUser)
	userID := uuid.New()
	user := models.User{
		ID:                      userID,
		FirstName:               "Suli",
		LastName:                "Test",
		Email:                   "suli@example.com",
		Username:                "suli",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	assignRoleToUser(t, testDB.DB, userID, roleID)

	permission := models.Permission{
		Action:   models.PermissionCreate,
		Resource: "sport",
	}
	if err := testDB.DB.Create(&permission).Error; err != nil {
		t.Fatalf("failed to create permission: %v", err)
	}

	body := map[string]any{
		"name":       "Women's Basketball",
		"popularity": int32(100000),
	}

	resp := api.Post("/api/v1/sport/", body, "Authorization: Bearer "+userID.String())
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", resp.Code, resp.Body.String())
	}
}

func TestPermissionMiddleware_AllowsWithRolePermission(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	roleID := getRoleID(t, testDB.DB, models.RoleUser)
	userID := uuid.New()
	user := models.User{
		ID:                      userID,
		FirstName:               "Suli",
		LastName:                "Test",
		Email:                   "suli@example.com",
		Username:                "suli",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	assignRoleToUser(t, testDB.DB, userID, roleID)

	permission := models.Permission{
		Action:   models.PermissionCreate,
		Resource: "sport",
	}
	if err := testDB.DB.Create(&permission).Error; err != nil {
		t.Fatalf("failed to create permission: %v", err)
	}

	if err := testDB.DB.Exec(
		`INSERT INTO "public"."role_permissions" ("role_id", "permission_id") VALUES (?, ?)`,
		roleID, permission.ID,
	).Error; err != nil {
		t.Fatalf("failed to create role permission: %v", err)
	}

	body := map[string]any{
		"name":       "Women's Basketball",
		"popularity": int32(100000),
	}

	resp := api.Post("/api/v1/sport/", body, "Authorization: Bearer "+userID.String())
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", resp.Code, resp.Body.String())
	}
}
