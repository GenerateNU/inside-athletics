package video

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoDB struct {
	db *gorm.DB
}

// NewVideoDB creates a new VideoDB instance
func NewVideoDB(db *gorm.DB) *VideoDB {
	return &VideoDB{db: db}
}

func (v *VideoDB) GetVideo(id *uuid.UUID) (*models.Video, error) {
	var video models.Video
	dbResponse := v.db.Where("id = ?", id).First(&video)
	if dbResponse.Error != nil {
		return utils.HandleDBError(&video, dbResponse.Error)
	}
	return &video, nil
}

func (v *VideoDB) CreateVideo(video *models.Video) (*models.Video, error) {
	dbResponse := v.db.Create(video)
	return utils.HandleDBError(video, dbResponse.Error)
}

func (v *VideoDB) DeleteVideo(id *uuid.UUID) error {
	dbResponse := v.db.Delete(&models.Video{}, id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.Video{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
