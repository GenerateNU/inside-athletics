package comment_like

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (u *CommentLikeDB) CreateCommentLike(commentLike *models.CommentLike) (*models.CommentLike, bool, error) {
	dbResponse := u.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "comment_id"}},
			DoNothing: true,
		},
		clause.Returning{},
	).Create(commentLike)
	if dbResponse.Error != nil {
		return nil, false, dbResponse.Error
	}
	if dbResponse.RowsAffected == 0 {
		return nil, false, nil
	}
	return commentLike, true, nil
}

// Permanently deletes a like by ID
func (u *CommentLikeDB) DeleteCommentLike(id uuid.UUID) (commentID uuid.UUID, err error) {
	var like models.CommentLike

	// checks that like exists to be deleted
	if err := u.db.Where("id = ?", id).First(&like).Error; err != nil {
		_, handleErr := utils.HandleDBError(&like, err)
		return uuid.Nil, handleErr
	}
	result := u.db.Delete(&like)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return like.CommentID, nil
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
