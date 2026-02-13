package role

import (
	"inside-athletics/internal/handlers/permission"
	"inside-athletics/internal/models"

	"github.com/google/uuid"
)

type CreateRoleRequest struct {
	Name        string                `json:"name" example:"moderator" doc:"Name of the role"`
	Permissions []PermissionSpecInput `json:"permissions,omitempty" doc:"Permissions to attach to the role"`
}

type PermissionSpecInput struct {
	Action   models.PermissionAction `json:"action" example:"create" doc:"Permission action"`
	Resource string                  `json:"resource" example:"sport" doc:"Permission resource"`
}

type UpdateRoleRequest struct {
	Name        *string                `json:"name,omitempty" example:"moderator" doc:"Name of the role"`
	Permissions *[]PermissionSpecInput `json:"permissions,omitempty" doc:"Permissions to attach to the role"`
}

type RoleResponse struct {
	ID          uuid.UUID                       `json:"id" example:"1" doc:"ID of the role"`
	Name        models.RoleName                 `json:"name" example:"admin" doc:"Name of the role"`
	Permissions []permission.PermissionResponse `json:"permissions,omitempty" doc:"Permissions attached to the role"`
}

type GetRoleByIDParams struct {
	ID uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID of the role"`
}

type DeleteRoleRequest struct {
	ID uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID of the role"`
}

type GetAllRolesParams struct {
	Limit  int `query:"limit" example:"20" doc:"Limit number of roles"`
	Offset int `query:"offset" example:"0" doc:"Offset for pagination"`
}

type GetAllRolesResponse struct {
	Roles []RoleResponse `json:"roles" doc:"List of roles"`
	Total int            `json:"total" doc:"Total number of roles"`
}

// HasPermission checks whether the given role grants the action/resource.
// It expects Permissions to be preloaded on the role.
func HasPermission(role *models.Role, action models.PermissionAction, resource string) bool {
	if role == nil {
		return false
	}
	for _, perm := range role.Permissions {
		if perm.Action == action && perm.Resource == resource {
			return true
		}
	}
	return false
}
