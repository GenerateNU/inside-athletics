package role

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleDB struct {
	db *gorm.DB
}

func NewRoleDB(db *gorm.DB) *RoleDB {
	return &RoleDB{db: db}
}

func (r *RoleDB) CreateRoleWithPermissionsStrict(spec models.RoleSpec) (*models.Role, error) {
	var created *models.Role

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Role
		if err := tx.Where("name = ?", spec.Name).First(&existing).Error; err == nil {
			return huma.Error400BadRequest("Role already exists")
		} else if err != gorm.ErrRecordNotFound {
			return huma.Error500InternalServerError("Database error", err)
		}

		role := &models.Role{Name: spec.Name}
		if err := tx.Create(role).Error; err != nil {
			return huma.Error500InternalServerError("Failed to create role", err)
		}

		for _, p := range spec.Permissions {
			var perm models.Permission
			if err := tx.
				Where("action = ? AND resource = ?", p.Action, p.Resource).
				First(&perm).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return huma.Error400BadRequest("Permission does not exist")
				}
				return huma.Error500InternalServerError("Database error", err)
			}

			link := models.RolePermission{
				RoleID:       role.ID,
				PermissionID: perm.ID,
			}
			if err := tx.Create(&link).Error; err != nil {
				return huma.Error500InternalServerError("Failed to link role permission", err)
			}
		}

		created = role
		return nil
	})

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *RoleDB) GetRoleByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	dbResponse := r.db.Where("id = ?", id).First(&role)
	return utils.HandleDBError(&role, dbResponse.Error)
}

func (r *RoleDB) GetAllRoles(limit, offset int) ([]models.Role, int64, error) {
	if limit == 0 {
		limit = 20
	}

	var roles []models.Role
	var total int64
	query := r.db.Model(&models.Role{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, huma.Error500InternalServerError("Database error", err)
	}

	if err := query.Limit(limit).Offset(offset).Find(&roles).Error; err != nil {
		return nil, 0, huma.Error500InternalServerError("Database error", err)
	}

	return roles, total, nil
}

func (r *RoleDB) UpdateRole(role *models.Role) (*models.Role, error) {
	dbResponse := r.db.Save(role)
	return utils.HandleDBError(role, dbResponse.Error)
}

func (r *RoleDB) UpdateRoleWithPermissionsStrict(id uuid.UUID, updates UpdateRoleRequest) (*models.Role, error) {
	var updated *models.Role

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var role models.Role
		if err := tx.First(&role, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return huma.Error404NotFound("Resource not found")
			}
			return huma.Error500InternalServerError("Database error", err)
		}

		if updates.Name != nil {
			if err := tx.Model(&role).Update("name", *updates.Name).Error; err != nil {
				_, err := utils.HandleDBError(&models.Role{}, err)
				return err
			}
			role.Name = models.RoleName(*updates.Name)
		}

		if updates.Permissions != nil {
			permIDs := make([]uuid.UUID, 0, len(*updates.Permissions))
			for _, p := range *updates.Permissions {
				var perm models.Permission
				if err := tx.Select("id").
					Where("action = ? AND resource = ?", p.Action, p.Resource).
					First(&perm).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						return huma.Error400BadRequest("Permission does not exist")
					}
					return huma.Error500InternalServerError("Database error", err)
				}
				permIDs = append(permIDs, perm.ID)
			}

			if err := tx.Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error; err != nil {
				return huma.Error500InternalServerError("Failed to update role permissions", err)
			}

			for _, permID := range permIDs {
				link := models.RolePermission{
					RoleID:       role.ID,
					PermissionID: permID,
				}
				if err := tx.Create(&link).Error; err != nil {
					return huma.Error500InternalServerError("Failed to update role permissions", err)
				}
			}
		}

		updated = &role
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *RoleDB) DeleteRole(id uuid.UUID) error {
	dbResponse := r.db.Delete(&models.Role{}, "id = ?", id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.Role{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
