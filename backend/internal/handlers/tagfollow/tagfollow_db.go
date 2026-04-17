package tagfollow

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagFollowDB struct {
	db *gorm.DB
}

func (u *TagFollowDB) GetTagFollowsByUser(user_id uuid.UUID) (*[]uuid.UUID, error) {
	var tagIDs []uuid.UUID
	dbResponse := u.db.Model(&models.TagFollow{}).
		Where("user_id = ?", user_id).
		Pluck("tag_id", &tagIDs)
	return utils.HandleDBError(&tagIDs, dbResponse.Error)
}

func (u *TagFollowDB) GetFollowingUsersByTag(tag_id uuid.UUID) (*[]uuid.UUID, error) {
	var userIDs []uuid.UUID
	dbResponse := u.db.Model(&models.TagFollow{}).
		Where("tag_id = ?", tag_id).
		Pluck("user_id", &userIDs)
	return utils.HandleDBError(&userIDs, dbResponse.Error)
}

func (u *TagFollowDB) CreateTagFollow(tagfollow *models.TagFollow) (*models.TagFollow, error) {
	// Checking if this user already followed the tag
	result := u.db.Where("user_id = ? AND tag_id = ?", tagfollow.UserID, tagfollow.TagID).First(tagfollow)
	if result.Error == nil {
		return nil, huma.Error409Conflict("User has already followed this tag")
	}
	dbResponse := u.db.Create(tagfollow)
	return utils.HandleDBError(tagfollow, dbResponse.Error)
}

func (u *TagFollowDB) DeleteTagFollow(id uuid.UUID) error {
	dbResponse := u.db.Delete(&models.TagFollow{}, "id = ?", id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.TagFollow{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}

func (u *TagFollowDB) DeleteTagFollowByUserAndTag(userID uuid.UUID, tagID uuid.UUID) error {
	dbResponse := u.db.Delete(&models.TagFollow{}, "user_id = ? AND tag_id = ?", userID, tagID)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.TagFollow{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
