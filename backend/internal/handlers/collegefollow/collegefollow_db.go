package collegefollow

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollegeFollowDB struct {
	db *gorm.DB
}

func (u *CollegeFollowDB) GetCollegeFollowsByUser(user_id uuid.UUID) (*[]uuid.UUID, error) {
	var collegeIDs []uuid.UUID
	dbResponse := u.db.Model(&models.CollegeFollow{}).
		Where("user_id = ?", user_id).
		Pluck("college_id", &collegeIDs)
	return utils.HandleDBError(&collegeIDs, dbResponse.Error)
}

func (u *CollegeFollowDB) GetFollowingUsersByCollege(college_id uuid.UUID) (*[]uuid.UUID, error) {
	var userIDs []uuid.UUID
	dbResponse := u.db.Model(&models.CollegeFollow{}).
		Where("college_id = ?", college_id).
		Pluck("user_id", &userIDs)
	return utils.HandleDBError(&userIDs, dbResponse.Error)
}

func (u *CollegeFollowDB) CreateCollegeFollow(collegefollow *models.CollegeFollow) (*models.CollegeFollow, error) {
	// checking if this user already followed the college
	result := u.db.Where("user_id = ? AND college_id = ?", collegefollow.UserID, collegefollow.CollegeID).First(collegefollow)
	if result.Error == nil {
		return nil, huma.Error409Conflict("User has already followed this college")
	}
	dbResponse := u.db.Create(collegefollow)
	return utils.HandleDBError(collegefollow, dbResponse.Error)
}

func (u *CollegeFollowDB) DeleteCollegeFollow(userID uuid.UUID, collegeID uuid.UUID) error {
	dbResponse := u.db.Delete(&models.CollegeFollow{}, "user_id = ? AND college_id = ?", userID, collegeID)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.CollegeFollow{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
