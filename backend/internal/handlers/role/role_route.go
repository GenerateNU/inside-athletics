package role

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	roleService := NewRoleService(db)

	{
		grp := huma.NewGroup(api, "/api/v1/role")
		huma.Post(grp, "/", roleService.CreateRole)
		huma.Get(grp, "/{id}", roleService.GetRoleByID)
		huma.Patch(grp, "/{id}", roleService.UpdateRole)
		huma.Delete(grp, "/{id}", roleService.DeleteRole)
	}
	{
		grp := huma.NewGroup(api, "/api/v1/roles")
		huma.Get(grp, "/", roleService.GetAllRoles)
	}
}
