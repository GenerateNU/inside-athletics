package routeTests

import (
	"inside-athletics/internal/models"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getRoleID(t *testing.T, db *gorm.DB, name models.RoleName) uuid.UUID {
	t.Helper()

	var role models.Role
	if err := db.Select("id").Where("name = ?", name).First(&role).Error; err != nil {
		t.Fatalf("failed to get role %s: %v", name, err)
	}

	return role.ID
}

func authHeader() string {
	return "Authorization: Bearer " + uuid.NewString()
}
