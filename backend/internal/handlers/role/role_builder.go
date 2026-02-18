package role

import "inside-athletics/internal/models"

type PermissionSpec struct {
	Action   models.PermissionAction
	Resource string
}

type RoleSpec struct {
	Name        models.RoleName
	Permissions []PermissionSpec
}

type RoleBuilder struct {
	spec RoleSpec
}

func NewRoleBuilder(name models.RoleName) *RoleBuilder {
	return &RoleBuilder{
		spec: RoleSpec{Name: name},
	}
}

func (b *RoleBuilder) WithPermission(action models.PermissionAction, resource string) *RoleBuilder {
	b.spec.Permissions = append(b.spec.Permissions, PermissionSpec{
		Action:   action,
		Resource: resource,
	})
	return b
}

func (b *RoleBuilder) Build() RoleSpec {
	return b.spec
}
