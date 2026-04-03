package media

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaDB struct {
	db *gorm.DB
}

// NewMediaDB creates a new MediaDB instance
func NewMediaDB(db *gorm.DB) *MediaDB {
	return &MediaDB{db: db}
}

func (v *MediaDB) GetMedia(id *uuid.UUID) (*models.Media, error) {
	var media models.Media
	dbResponse := v.db.Where("id = ?", id).First(&media)
	if dbResponse.Error != nil {
		return utils.HandleDBError(&media, dbResponse.Error)
	}
	return &media, nil
}

func (v *MediaDB) CreateMedia(media *models.Media) (*models.Media, error) {
	dbResponse := v.db.Create(media)
	return utils.HandleDBError(media, dbResponse.Error)
}

func (v *MediaDB) DeleteMedia(id *uuid.UUID) error {
	dbResponse := v.db.Delete(&models.Media{}, id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.Media{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
