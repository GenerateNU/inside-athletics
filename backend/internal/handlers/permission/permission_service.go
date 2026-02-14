package permission

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionService struct {
	permissionDB *PermissionDB
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{
		permissionDB: NewPermissionDB(db),
	}
}

func (p *PermissionService) CreatePermission(ctx context.Context, input *struct{ Body CreatePermissionRequest }) (*utils.ResponseBody[PermissionResponse], error) {
	if input.Body.Action == "" || input.Body.Resource == "" {
		return nil, huma.Error422UnprocessableEntity("action and resource are required")
	}
	if !models.IsValidPermissionAction(models.PermissionAction(input.Body.Action)) {
		return nil, huma.Error422UnprocessableEntity("invalid permission action")
	}

	perm, err := p.permissionDB.CreatePermission(input.Body.Action, input.Body.Resource)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PermissionResponse]{
		Body: ToPermissionResponse(perm),
	}, nil
}

func (p *PermissionService) GetPermissionByID(ctx context.Context, input *GetPermissionByIDParams) (*utils.ResponseBody[PermissionResponse], error) {
	perm, err := p.permissionDB.GetPermissionByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PermissionResponse]{
		Body: ToPermissionResponse(perm),
	}, nil
}

func (p *PermissionService) GetAllPermissions(ctx context.Context, input *GetAllPermissionsParams) (*utils.ResponseBody[GetAllPermissionsResponse], error) {
	perms, total, err := p.permissionDB.GetAllPermissions(input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	responses := make([]PermissionResponse, 0, len(perms))
	for i := range perms {
		responses = append(responses, *ToPermissionResponse(&perms[i]))
	}

	return &utils.ResponseBody[GetAllPermissionsResponse]{
		Body: &GetAllPermissionsResponse{
			Permissions: responses,
			Total:       int(total),
		},
	}, nil
}

func (p *PermissionService) UpdatePermission(ctx context.Context, input *struct {
	ID   uuid.UUID `path:"id"`
	Body UpdatePermissionRequest
}) (*utils.ResponseBody[PermissionResponse], error) {
	if input.Body.Action != nil {
		if *input.Body.Action == "" {
			return nil, huma.Error422UnprocessableEntity("action cannot be empty")
		}
		if !models.IsValidPermissionAction(models.PermissionAction(*input.Body.Action)) {
			return nil, huma.Error422UnprocessableEntity("invalid permission action")
		}
	}
	if input.Body.Resource != nil {
		if *input.Body.Resource == "" {
			return nil, huma.Error422UnprocessableEntity("resource cannot be empty")
		}
	}

	updated, err := p.permissionDB.UpdatePermissionByID(input.ID, input.Body)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PermissionResponse]{
		Body: ToPermissionResponse(updated),
	}, nil
}

func (p *PermissionService) DeletePermission(ctx context.Context, input *DeletePermissionRequest) (*utils.ResponseBody[PermissionResponse], error) {
	perm, err := p.permissionDB.GetPermissionByID(input.ID)
	if err != nil {
		return nil, err
	}

	if err := p.permissionDB.DeletePermission(input.ID); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[PermissionResponse]{
		Body: ToPermissionResponse(perm),
	}, nil
}
