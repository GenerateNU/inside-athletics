package athleteverification

import (
	"inside-athletics/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AthleteDB struct {
	db *gorm.DB
}

func NewAthleteDB(db *gorm.DB) *AthleteDB {
	return &AthleteDB{
		db: db,
	}
}

func (a *AthleteDB) AthleteExists(name string, collegeID *uuid.UUID, sportID *uuid.UUID) (bool, error) {
	var athlete models.Athlete
	result := a.db.Where("name = ? AND sport_id = ? AND college_id = ?", name, sportID, collegeID).Limit(1).Find(&athlete)
	if err := result.Error; err != nil {
		return false, err
	}
	return result.RowsAffected > 0, nil
}

func (a *AthleteDB) ManualVerificationData() {

}
