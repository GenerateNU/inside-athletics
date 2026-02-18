package permission

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PermissionDB struct {
	db *gorm.DB
}

func NewPermissionDB(db *gorm.DB) *PermissionDB {
	return &PermissionDB{db: db}
}

func (p *PermissionDB) CreatePermission(action, resource string) (*models.Permission, error) {
	perm := &models.Permission{
		Action:   models.PermissionAction(action),
		Resource: resource,
	}
	dbResponse := p.db.Create(perm)
	return utils.HandleDBError(perm, dbResponse.Error)
}

func (p *PermissionDB) GetPermissionByID(id uuid.UUID) (*models.Permission, error) {
	var perm models.Permission
	dbResponse := p.db.Where("id = ?", id).First(&perm)
	return utils.HandleDBError(&perm, dbResponse.Error)
}

func (p *PermissionDB) GetAllPermissions(limit, offset int) ([]models.Permission, int64, error) {
	if limit == 0 {
		limit = 20
	}

	var perms []models.Permission
	var total int64
	query := p.db.Model(&models.Permission{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, huma.Error500InternalServerError("Database error", err)
	}

	if err := query.Limit(limit).Offset(offset).Find(&perms).Error; err != nil {
		return nil, 0, huma.Error500InternalServerError("Database error", err)
	}

	return perms, total, nil
}

func (p *PermissionDB) UpdatePermission(perm *models.Permission) (*models.Permission, error) {
	dbResponse := p.db.Save(perm)
	return utils.HandleDBError(perm, dbResponse.Error)
}

func (p *PermissionDB) UpdatePermissionByID(id uuid.UUID, updates UpdatePermissionRequest) (*models.Permission, error) {
	var updatedPerm models.Permission
	dbResponse := p.db.Model(&models.Permission{}).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updates).
		Scan(&updatedPerm)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.Permission{}, dbResponse.Error)
		return nil, err
	}
	if dbResponse.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Resource not found")
	}
	return &updatedPerm, nil
}

func (p *PermissionDB) DeletePermission(id uuid.UUID) error {
	dbResponse := p.db.Delete(&models.Permission{}, "id = ?", id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.Permission{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
