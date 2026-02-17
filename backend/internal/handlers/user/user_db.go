package user

import (
	"inside-athletics/internal/handlers/permission"
	"inside-athletics/internal/handlers/role"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDB struct {
	db *gorm.DB
}

/*
*
Here we are using GORM to interact with the database. This is an ORM (Object Relational Mapping)
which allows us to interact with the database without having to write raw SQL queries
*/
func (u *UserDB) GetUser(id uuid.UUID) (*models.User, error) {
	var user models.User
	dbResponse := u.db.Where("id = ?", id).First(&user)
	return utils.HandleDBError(&user, dbResponse.Error) // helper function that maps GORM errors to Huma errors
}

func (u *UserDB) CreateUser(user *models.User) (*models.User, error) {
	dbResponse := u.db.Create(user)
	return utils.HandleDBError(user, dbResponse.Error)
}

// This function creates a link between a user and a role in the user_roles table
// We use FirstOrCreate to avoid duplicate entries if the user already has the role
func (u *UserDB) AddUserRole(userID, roleID uuid.UUID) error {
	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	if err := u.db.Where("user_id = ? AND role_id = ?", userID, roleID).Find(&userRole).Error; err != nil {
		return huma.Error500InternalServerError("Failed to assign role to user", err)
	}
	return nil
}

func (u *UserDB) GetAllRolesForUser(userID uuid.UUID) (*[]models.Role, error) {
	var userRoles []models.Role
	err := u.db.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&userRoles).Error
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get user roles", err)
	}
	return &userRoles, nil
}

type rolePermissionRow struct {
	RoleID             uuid.UUID                `gorm:"column:role_id"`
	RoleName           models.RoleName          `gorm:"column:role_name"`
	PermissionID       *uuid.UUID               `gorm:"column:permission_id"`
	PermissionAction   *models.PermissionAction `gorm:"column:permission_action"`
	PermissionResource *string                  `gorm:"column:permission_resource"`
}

// GetRolesWithPermissionsForUser returns the roles and permissions for a user using a single join query.
func (u *UserDB) GetRolesWithPermissionsForUser(userID uuid.UUID) (*[]role.RoleResponse, error) {
	var rows []rolePermissionRow
	err := u.db.Table("user_roles").
		Select("roles.id as role_id, roles.name as role_name, permissions.id as permission_id, permissions.action as permission_action, permissions.resource as permission_resource").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Joins("LEFT JOIN role_permissions rp ON rp.role_id = roles.id").
		Joins("LEFT JOIN permissions ON permissions.id = rp.permission_id").
		Where("user_roles.user_id = ?", userID).
		Order("roles.name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get user roles", err)
	}
	if len(rows) == 0 {
		return nil, nil
	}

	roleMap := make(map[uuid.UUID]*role.RoleResponse, len(rows))
	roleOrder := make([]uuid.UUID, 0, len(rows))
	permSeen := make(map[uuid.UUID]map[uuid.UUID]struct{}, len(rows))

	for _, row := range rows {
		r, ok := roleMap[row.RoleID]
		if !ok {
			r = &role.RoleResponse{
				ID:          row.RoleID,
				Name:        row.RoleName,
				Permissions: nil,
			}
			roleMap[row.RoleID] = r
			roleOrder = append(roleOrder, row.RoleID)
		}

		if row.PermissionID == nil || row.PermissionAction == nil || row.PermissionResource == nil {
			continue
		}

		if _, ok := permSeen[row.RoleID]; !ok {
			permSeen[row.RoleID] = make(map[uuid.UUID]struct{})
		}
		if _, ok := permSeen[row.RoleID][*row.PermissionID]; ok {
			continue
		}
		permSeen[row.RoleID][*row.PermissionID] = struct{}{}

		r.Permissions = append(r.Permissions, permission.PermissionResponse{
			ID:       *row.PermissionID,
			Action:   *row.PermissionAction,
			Resource: *row.PermissionResource,
		})
	}

	responses := make([]role.RoleResponse, 0, len(roleOrder))
	for _, id := range roleOrder {
		responses = append(responses, *roleMap[id])
	}

	return &responses, nil
}

func (u *UserDB) UpdateUser(id uuid.UUID, updates UpdateUserBody) (*models.User, error) {
	var updatedUser models.User
	dbResponse := u.db.Model(&models.User{}).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updates).
		Scan(&updatedUser)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.User{}, dbResponse.Error)
		return nil, err
	}
	if dbResponse.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Resource not found")
	}
	return &updatedUser, nil
}

func (u *UserDB) DeleteUser(id uuid.UUID) error {
	dbResponse := u.db.Delete(&models.User{}, "id = ?", id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.User{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
