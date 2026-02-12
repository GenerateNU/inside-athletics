package role

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleService struct {
	roleDB *RoleDB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{
		roleDB: NewRoleDB(db),
	}
}

func (r *RoleService) CreateRole(ctx context.Context, input *struct{ Body CreateRoleRequest }) (*utils.ResponseBody[RoleResponse], error) {
	if input.Body.Name == "" {
		return nil, huma.Error422UnprocessableEntity("name cannot be empty")
	}

	builder := models.NewRoleBuilder(models.RoleName(input.Body.Name))
	for _, perm := range input.Body.Permissions {
		if perm.Action == "" || perm.Resource == "" {
			return nil, huma.Error422UnprocessableEntity("permissions must include action and resource")
		}
		builder.WithPermission(perm.Action, perm.Resource)
	}

	spec := builder.Build()

	role, err := r.roleDB.CreateRoleWithPermissionsStrict(spec)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[RoleResponse]{
		Body: toRoleResponse(role),
	}, nil
}

func (r *RoleService) CreateRoleNameOnly(ctx context.Context, input *struct{ Body CreateRoleRequest }) (*utils.ResponseBody[RoleResponse], error) {
	if input.Body.Name == "" {
		return nil, huma.Error422UnprocessableEntity("name cannot be empty")
	}
	if len(input.Body.Permissions) > 0 {
		return nil, huma.Error422UnprocessableEntity("permissions must be empty for this endpoint")
	}

	spec := models.NewRoleBuilder(models.RoleName(input.Body.Name)).Build()

	role, err := r.roleDB.CreateRoleWithPermissionsStrict(spec)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[RoleResponse]{
		Body: toRoleResponse(role),
	}, nil
}

func (r *RoleService) GetRoleByID(ctx context.Context, input *GetRoleByIDParams) (*utils.ResponseBody[RoleResponse], error) {
	role, err := r.roleDB.GetRoleByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[RoleResponse]{
		Body: toRoleResponse(role),
	}, nil
}

func (r *RoleService) GetAllRoles(ctx context.Context, input *GetAllRolesParams) (*utils.ResponseBody[GetAllRolesResponse], error) {
	roles, total, err := r.roleDB.GetAllRoles(input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	responses := make([]RoleResponse, 0, len(roles))
	for i := range roles {
		responses = append(responses, *toRoleResponse(&roles[i]))
	}

	return &utils.ResponseBody[GetAllRolesResponse]{
		Body: &GetAllRolesResponse{
			Roles: responses,
			Total: int(total),
		},
	}, nil
}

func (r *RoleService) UpdateRole(ctx context.Context, input *struct {
	ID   uuid.UUID `path:"id"`
	Body UpdateRoleRequest
}) (*utils.ResponseBody[RoleResponse], error) {
	if input.Body.Name != nil {
		if *input.Body.Name == "" {
			return nil, huma.Error422UnprocessableEntity("name cannot be empty")
		}
	}
	if input.Body.Permissions != nil {
		for _, perm := range *input.Body.Permissions {
			if perm.Action == "" || perm.Resource == "" {
				return nil, huma.Error422UnprocessableEntity("permissions must include action and resource")
			}
		}
	}

	updated, err := r.roleDB.UpdateRoleWithPermissionsStrict(input.ID, input.Body)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[RoleResponse]{
		Body: toRoleResponse(updated),
	}, nil
}

func (r *RoleService) DeleteRole(ctx context.Context, input *DeleteRoleRequest) (*utils.ResponseBody[RoleResponse], error) {
	role, err := r.roleDB.GetRoleByID(input.ID)
	if err != nil {
		return nil, err
	}

	if err := r.roleDB.DeleteRole(input.ID); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[RoleResponse]{
		Body: toRoleResponse(role),
	}, nil
}

func toRoleResponse(role *models.Role) *RoleResponse {
	return &RoleResponse{
		ID:   role.ID,
		Name: role.Name,
	}
}
