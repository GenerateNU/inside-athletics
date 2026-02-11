package routeTests

import (
	"testing"

	"inside-athletics/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type permissionSpec struct {
	Action   models.PermissionAction
	Resource string
}

func seedUserWithRoleAndPermissions(t *testing.T, db *gorm.DB, roleName models.RoleName, perms []permissionSpec) (uuid.UUID, string) {
	t.Helper()

	roleID := getRoleID(t, db, roleName)
	userID := uuid.New()
	user := models.User{
		ID:                      userID,
		FirstName:               "Test",
		LastName:                "User",
		Email:                   "testuser@example.com",
		Username:                "testuser",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	assignRoleToUser(t, db, userID, roleID)

	for _, perm := range perms {
		ensurePermissionForRole(t, db, roleID, perm.Action, perm.Resource)
	}

	return userID, "Authorization: Bearer " + userID.String()
}

func ensurePermissionForRole(t *testing.T, db *gorm.DB, roleID uuid.UUID, action models.PermissionAction, resource string) {
	t.Helper()

	permission := models.Permission{
		Action:   action,
		Resource: resource,
	}
	if err := db.Where("action = ? AND resource = ?", action, resource).FirstOrCreate(&permission).Error; err != nil {
		t.Fatalf("failed to ensure permission %s %s: %v", action, resource, err)
	}

	rolePermission := models.RolePermission{
		RoleID:       roleID,
		PermissionID: permission.ID,
	}
	if err := db.Where("role_id = ? AND permission_id = ?", roleID, permission.ID).FirstOrCreate(&rolePermission).Error; err != nil {
		t.Fatalf("failed to ensure role permission: %v", err)
	}
}

func assignRoleToUser(t *testing.T, db *gorm.DB, userID, roleID uuid.UUID) {
	t.Helper()

	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	if err := db.Create(&userRole).Error; err != nil {
		t.Fatalf("failed to assign role to user: %v", err)
	}
}
