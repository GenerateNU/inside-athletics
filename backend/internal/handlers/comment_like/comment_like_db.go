package comment_like

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentLikeDB struct {
	db *gorm.DB
}

// Retrieves a comment like given an ID
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

// Soft deletes a like by ID.
func (u *CommentLikeDB) DeleteCommentLike(id uuid.UUID) error {
	result := u.db.Delete(&models.CommentLike{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}

// Retrieves like count for the comment and whether the given user has liked it. If userID is zero, liked is false.
func (u *CommentLikeDB) GetCommentLikeInfo(commentID, userID uuid.UUID) (count int64, liked bool, err error) {
	err = u.db.Model(&models.CommentLike{}).Where("comment_id = ?", commentID).Count(&count).Error
	if err != nil {
		return 0, false, err
	}
	if userID != uuid.Nil {
		var userCount int64
		err = u.db.Model(&models.CommentLike{}).
			Where("user_id = ? AND comment_id = ?", userID, commentID).
			Count(&userCount).Error
		if err != nil {
			return 0, false, err
		}
		liked = userCount > 0
	}
	return count, liked, nil
}
