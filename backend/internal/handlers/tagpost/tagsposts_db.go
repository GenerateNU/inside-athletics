package tagpost

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagPostDB struct {
	db *gorm.DB
}

func (u *TagPostDB) GetPostsByTag(tag_id uuid.UUID) (*[]uuid.UUID, error) {
	var postIDs []uuid.UUID
	dbResponse := u.db.Model(&models.TagPost{}).
		Where("tag_id = ?", tag_id).
		Pluck("post_id", &postIDs)
	return utils.HandleDBError(&postIDs, dbResponse.Error)
}

func (u *TagPostDB) GetTagsByPost(post_id uuid.UUID) (*[]uuid.UUID, error) {
	var tagIDs []uuid.UUID
	dbResponse := u.db.Model(&models.TagPost{}).
		Where("post_id = ?", post_id).
		Pluck("tag_id", &tagIDs)
	return utils.HandleDBError(&tagIDs, dbResponse.Error)
}

func (u *TagPostDB) GetTagPostById(id uuid.UUID) (*models.TagPost, error) {
	var tagpost models.TagPost
	dbResponse := u.db.Where("id = ?", id).First(&tagpost)
	return utils.HandleDBError(&tagpost, dbResponse.Error) // helper function that maps GORM errors to Huma errors
}

func (u *TagPostDB) CreateTagPost(tagpost *models.TagPost) (*models.TagPost, error) {
	dbResponse := u.db.Create(tagpost)
	return utils.HandleDBError(tagpost, dbResponse.Error)
}

func (u *TagPostDB) UpdateTagPost(id uuid.UUID, updates *UpdateTagPostBody) (*models.TagPost, error) {
	dbResponse := u.db.Model(&models.TagPost{}).
		Where("id = ?", id).
		Updates(updates)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.TagPost{}, dbResponse.Error)
		return nil, err
	}
	if dbResponse.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Resource not found")
	}
	return u.GetTagPostById(id)
}

func (u *TagPostDB) DeleteTagPost(id uuid.UUID) error {
	dbResponse := u.db.Delete(&models.TagPost{}, "id = ?", id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.TagPost{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
