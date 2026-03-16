package unitTests

import (
	"testing"

	"inside-athletics/internal/handlers/role"
	"inside-athletics/internal/models"
)

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name     string
		perms    []models.Permission
		action   models.PermissionAction
		resource string
		want     bool
	}{
		{
			name:     "nil role",
			perms:    nil,
			action:   models.PermissionCreate,
			resource: "sport",
			want:     false,
		},
		{
			name:     "no permissions",
			perms:    []models.Permission{},
			action:   models.PermissionCreate,
			resource: "sport",
			want:     false,
		},
		{
			name: "matching permission",
			perms: []models.Permission{
				{
					Action:   models.PermissionCreate,
					Resource: "sport",
				},
			},
			action:   models.PermissionCreate,
			resource: "sport",
			want:     true,
		},
		{
			name: "non-matching resource",
			perms: []models.Permission{
				{
					Action:   models.PermissionCreate,
					Resource: "post",
				},
			},
			action:   models.PermissionCreate,
			resource: "sport",
			want:     false,
		},
		{
			name: "non-matching action",
			perms: []models.Permission{
				{
					Action:   models.PermissionUpdate,
					Resource: "sport",
				},
			},
			action:   models.PermissionCreate,
			resource: "sport",
			want:     false,
		},
		{
			name: "multiple permissions includes match",
			perms: []models.Permission{
				{
					Action:   models.PermissionCreate,
					Resource: "post",
				},
				{
					Action:   models.PermissionDelete,
					Resource: "sport",
				},
			},
			action:   models.PermissionDelete,
			resource: "sport",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := role.HasPermission(tt.perms, tt.action, tt.resource); got != tt.want {
				t.Fatalf("HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}
