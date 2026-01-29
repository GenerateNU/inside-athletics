package tag

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagDB struct {
	db *gorm.DB
}

func (u *TagDB) GetTagByName(name string) (*models.Tag, error) {
	var tag models.Tag
	dbResponse := u.db.Where("name = ?", name).First(&tag)
	return utils.HandleDBError(&tag, dbResponse.Error) // helper function that maps GORM errors to Huma errors
}

func (u *TagDB) GetTagByID(id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	dbResponse := u.db.Where("id = ?", id).First(&tag)
	return utils.HandleDBError(&tag, dbResponse.Error) // helper function that maps GORM errors to Huma errors
}

func (u *TagDB) CreateTag(tag *models.Tag) (*models.Tag, error) {
	dbResponse := u.db.Create(tag)
	return utils.HandleDBError(tag, dbResponse.Error)
}

func (u *TagDB) UpdateTag(tag *models.Tag) (*models.Tag, error) {
	dbResponse := u.db.Save(tag)
	return utils.HandleDBError(tag, dbResponse.Error)
}

func (u *TagDB) DeleteTag(id uuid.UUID) error {
	dbResponse := u.db.Delete(&models.Tag{}, "id = ?", id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.Tag{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
