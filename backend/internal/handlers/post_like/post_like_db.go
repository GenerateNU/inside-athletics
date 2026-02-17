package post_like

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostLikeDB struct {
	db *gorm.DB
}

// Retrieves a post like given an ID
func (u *PostLikeDB) GetPostLike(id uuid.UUID) (*models.PostLike, error) {
	var like models.PostLike
	dbResponse := u.db.Where("id = ?", id).First(&like)
	return utils.HandleDBError(&like, dbResponse.Error)
}

// Creates a new like on a post in the database
func (u *PostLikeDB) CreatePostLike(postLike *models.PostLike) (*models.PostLike, bool, error) {
	dbResponse := u.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "post_id"}},
			DoNothing: true,
		},
		clause.Returning{},
	).Create(postLike)
	if dbResponse.Error != nil {
		return nil, false, dbResponse.Error
	}
	if dbResponse.RowsAffected == 0 {
		return nil, false, nil
	}
	return postLike, true, nil
}

/// Permanently deletes a like by ID
func (u *PostLikeDB) DeletePostLike(id uuid.UUID) (postID uuid.UUID, err error) {
	var like models.PostLike
	// checking if like that needs to be deleted exists
	if err := u.db.Where("id = ?", id).First(&like).Error; err != nil {
		_, handleErr := utils.HandleDBError(&like, err)
		return uuid.Nil, handleErr
	}
	result := u.db.Delete(&like)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return like.PostID, nil
}

// Returns like count for the post and whether the given user has liked it. If userID is zero, liked is false.
func (u *PostLikeDB) GetPostLikeInfo(postID, userID uuid.UUID) (count int64, liked bool, err error) {
	err = u.db.Model(&models.PostLike{}).Where("post_id = ?", postID).Count(&count).Error
	if err != nil {
		return 0, false, err
	}
	if userID != uuid.Nil {
		var userCount int64
		err = u.db.Model(&models.PostLike{}).
			Where("user_id = ? AND post_id = ?", userID, postID).
			Count(&userCount).Error
		if err != nil {
			return 0, false, err
		}
		liked = userCount > 0
	}
	return count, liked, nil
}
