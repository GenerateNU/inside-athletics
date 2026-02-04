package like

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentLikeDB struct {
	db *gorm.DB
}

// GetCommentLike retrieves a comment like given an ID
func (u *CommentLikeDB) GetCommentLike(id uuid.UUID) (*models.CommentLike, error) {
	var like models.CommentLike
	dbResponse := u.db.Where("id = ?", id).First(&like)
	return utils.HandleDBError(&like, dbResponse.Error)
}

// CreateCommentLike creates a new like on a comment in the database
func (u *CommentLikeDB) CreateCommentLike(commentLike *models.CommentLike) (*models.CommentLike, error) {
	dbResponse := u.db.Create(commentLike)
	return utils.HandleDBError(commentLike, dbResponse.Error)
}

// DeleteCommentLike soft deletes a sport by ID
func (u *CommentLikeDB) DeleteCommentLike(id uuid.UUID) error {
	result := u.db.Delete(&models.CommentLike{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
