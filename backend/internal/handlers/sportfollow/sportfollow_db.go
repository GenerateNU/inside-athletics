package sportfollow

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SportFollowDB struct {
	db *gorm.DB
}

func (u *SportFollowDB) GetSportFollowsByUser(user_id uuid.UUID) (*[]uuid.UUID, error) {
	var sportIDs []uuid.UUID
	dbResponse := u.db.Model(&models.SportFollow{}).
		Where("user_id = ?", user_id).
		Pluck("sport_id", &sportIDs)
	return utils.HandleDBError(&sportIDs, dbResponse.Error)
}

func (u *SportFollowDB) GetFollowingUsersBySport(sport_id uuid.UUID) (*[]uuid.UUID, error) {
	var userIDs []uuid.UUID
	dbResponse := u.db.Model(&models.SportFollow{}).
		Where("sport_id = ?", sport_id).
		Pluck("user_id", &userIDs)
	return utils.HandleDBError(&userIDs, dbResponse.Error)
}

func (u *SportFollowDB) CreateSportFollow(sportfollow *models.SportFollow) (*models.SportFollow, error) {
	// checking if user already followed the sport
	result := u.db.Where("user_id = ? AND sport_id = ?", sportfollow.UserID, sportfollow.SportID).First(sportfollow)
	if result.Error == nil {
		return nil, huma.Error409Conflict("User has already followed this sport")
	}
	dbResponse := u.db.Create(sportfollow)
	return utils.HandleDBError(sportfollow, dbResponse.Error)
}

func (u *SportFollowDB) DeleteSportFollow(id uuid.UUID) error {
	dbResponse := u.db.Delete(&models.SportFollow{}, "id = ?", id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.SportFollow{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
