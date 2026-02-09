package permission

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	permissionService := NewPermissionService(db)

	{
		grp := huma.NewGroup(api, "/api/v1/permission")
		huma.Post(grp, "/", permissionService.CreatePermission)
		huma.Get(grp, "/{id}", permissionService.GetPermissionByID)
		huma.Patch(grp, "/{id}", permissionService.UpdatePermission)
		huma.Delete(grp, "/{id}", permissionService.DeletePermission)
	}
	{
		grp := huma.NewGroup(api, "/api/v1/permissions")
		huma.Get(grp, "/", permissionService.GetAllPermissions)
	}
}
