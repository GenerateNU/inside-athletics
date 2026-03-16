package routeTests

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/server"
	"testing"

	"github.com/google/uuid"
)

func TestAuthorizationDBUserExists(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	authDB := server.NewAuthorizationDB(testDB.DB)

	user := models.User{
		ID:                      uuid.New(),
		FirstName:               "Test",
		LastName:                "User",
		Email:                   "test@example.com",
		Username:                "testuser",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if err := authDB.UserExists(user.ID); err != nil {
		t.Fatalf("expected user to exist, got error: %v", err)
	}

	if err := authDB.UserExists(uuid.New()); err == nil {
		t.Fatalf("expected error for missing user")
	}
}

func TestAuthorizationDBUserHasPermission(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)

	authDB := server.NewAuthorizationDB(testDB.DB)

	user := models.User{
		ID:                      uuid.New(),
		FirstName:               "Perm",
		LastName:                "User",
		Email:                   "perm@example.com",
		Username:                "permuser",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	role := models.Role{Name: models.RoleName("tester_auth")}
	if err := testDB.DB.Create(&role).Error; err != nil {
		t.Fatalf("failed to create role: %v", err)
	}

	perm := models.Permission{Action: models.PermissionCreate, Resource: "sport"}
	if err := testDB.DB.Create(&perm).Error; err != nil {
		t.Fatalf("failed to create permission: %v", err)
	}

	link := models.RolePermission{RoleID: role.ID, PermissionID: perm.ID}
	if err := testDB.DB.Create(&link).Error; err != nil {
		t.Fatalf("failed to create role permission link: %v", err)
	}

	userRole := models.UserRole{UserID: user.ID, RoleID: role.ID}
	if err := testDB.DB.Create(&userRole).Error; err != nil {
		t.Fatalf("failed to create user role link: %v", err)
	}

	hasPerm, err := authDB.UserHasPermission(user.ID, models.PermissionCreate, "sport")
	if err != nil {
		t.Fatalf("unexpected error checking permission: %v", err)
	}
	if !hasPerm {
		t.Fatalf("expected user to have permission")
	}

	hasPerm, err = authDB.UserHasPermission(user.ID, models.PermissionDelete, "sport")
	if err != nil {
		t.Fatalf("unexpected error checking permission: %v", err)
	}
	if hasPerm {
		t.Fatalf("expected user to not have permission")
	}
}
