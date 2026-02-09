package role

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleDB struct {
	db *gorm.DB
}

func NewRoleDB(db *gorm.DB) *RoleDB {
	return &RoleDB{db: db}
}

func (r *RoleDB) CreateRole(name string) (*models.Role, error) {
	role := &models.Role{Name: models.RoleName(name)}
	dbResponse := r.db.Create(role)
	return utils.HandleDBError(role, dbResponse.Error)
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

func (r *RoleDB) UpdateRoleByID(id uuid.UUID, updates UpdateRoleRequest) (*models.Role, error) {
	var updatedRole models.Role
	dbResponse := r.db.Model(&models.Role{}).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updates).
		Scan(&updatedRole)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.Role{}, dbResponse.Error)
		return nil, err
	}
	if dbResponse.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Resource not found")
	}
	return &updatedRole, nil
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
