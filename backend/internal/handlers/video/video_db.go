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

func (v *VideoDB) GetVideo(id uuid.UUID) (*models.Video, error) {
	var video models.Video
	dpResponse := v.db.Where("id = ?", id).First(&video)
	if dpResponse.Error != nil {
		return nil, dpResponse.Error
	}
	return &video, nil
}

func (v *VideoDB) CreateVideo(video *models.Video) (*models.Video, error) {
	dbResponse := v.db.Create(video)
	return utils.HandleDBError(video, dbResponse.Error)
}

func (v *VideoDB) DeleteVideo(id uuid.UUID) error {
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
