package models

type PermissionSpec struct {
	Action   PermissionAction
	Resource string
}

type RoleSpec struct {
	Name        RoleName
	Permissions []PermissionSpec
}

type RoleBuilder struct {
	spec RoleSpec
}

func NewRoleBuilder(name RoleName) *RoleBuilder {
	return &RoleBuilder{
		spec: RoleSpec{Name: name},
	}
}

func (b *RoleBuilder) WithPermission(action PermissionAction, resource string) *RoleBuilder {
	b.spec.Permissions = append(b.spec.Permissions, PermissionSpec{
		Action:   action,
		Resource: resource,
	})
	return b
}

func (b *RoleBuilder) Build() RoleSpec {
	return b.spec
}
