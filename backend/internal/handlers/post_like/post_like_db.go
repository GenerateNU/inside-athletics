package like

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostLikeDB struct {
	db *gorm.DB
}

// GetPostLike retrieves a post like given an ID
func (u *PostLikeDB) GetPostLike(id uuid.UUID) (*models.PostLike, error) {
	var like models.PostLike
	dbResponse := u.db.Where("id = ?", id).First(&like)
	return utils.HandleDBError(&like, dbResponse.Error)
}

// CreatePostLike creates a new like on a post in the database
func (u *PostLikeDB) CreatePostLike(postLike *models.PostLike) (*models.PostLike, error) {
	dbResponse := u.db.Create(postLike)
	return utils.HandleDBError(postLike, dbResponse.Error)
}

// DeletePostLike soft deletes a sport by ID
func (u *PostLikeDB) DeletePostLike(id uuid.UUID) error {
	result := u.db.Delete(&models.PostLike{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
