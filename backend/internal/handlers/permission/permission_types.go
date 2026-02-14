package permission

import (
	"inside-athletics/internal/models"

	"github.com/google/uuid"
)

type CreatePermissionRequest struct {
	Action   string `json:"action" example:"create" doc:"Action for the permission"`
	Resource string `json:"resource" example:"user" doc:"Resource for the permission"`
}

type UpdatePermissionRequest struct {
	Action   *string `json:"action,omitempty" example:"update" doc:"Action for the permission"`
	Resource *string `json:"resource,omitempty" example:"user" doc:"Resource for the permission"`
}

type PermissionResponse struct {
	ID       uuid.UUID               `json:"id" example:"1" doc:"ID of the permission"`
	Action   models.PermissionAction `json:"action" example:"create" doc:"Action for the permission"`
	Resource string                  `json:"resource" example:"user" doc:"Resource for the permission"`
}

func ToPermissionResponse(perm *models.Permission) *PermissionResponse {
	return &PermissionResponse{
		ID:       perm.ID,
		Action:   perm.Action,
		Resource: perm.Resource,
	}
}

func ToPermissionResponses(perms []models.Permission) []PermissionResponse {
	if len(perms) == 0 {
		return nil
	}
	responses := make([]PermissionResponse, 0, len(perms))
	for i := range perms {
		responses = append(responses, *ToPermissionResponse(&perms[i]))
	}
	return responses
}

type GetPermissionByIDParams struct {
	ID uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID of the permission"`
}

type DeletePermissionRequest struct {
	ID uuid.UUID `path:"id" maxLength:"36" example:"1" doc:"ID of the permission"`
}

type GetAllPermissionsParams struct {
	Limit  int `query:"limit" example:"20" doc:"Limit number of permissions"`
	Offset int `query:"offset" example:"0" doc:"Offset for pagination"`
}

type GetAllPermissionsResponse struct {
	Permissions []PermissionResponse `json:"permissions" doc:"List of permissions"`
	Total       int                  `json:"total" doc:"Total number of permissions"`
}
