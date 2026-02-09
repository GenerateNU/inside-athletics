package role

import (
	"inside-athletics/internal/models"

	"github.com/google/uuid"
)

type CreateRoleRequest struct {
	Name string `json:"name" example:"moderator" doc:"Name of the role"`
}

type UpdateRoleRequest struct {
	Name *string `json:"name,omitempty" example:"moderator" doc:"Name of the role"`
}

type RoleResponse struct {
	ID   uuid.UUID       `json:"id" example:"1" doc:"ID of the role"`
	Name models.RoleName `json:"name" example:"admin" doc:"Name of the role"`
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
