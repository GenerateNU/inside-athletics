package athlete

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

func (a *AthleteDB) GetAthlete(name string, collegeID *uuid.UUID, sportID *uuid.UUID) (*models.Athlete, bool, error) {
	var athlete models.Athlete
	result := a.db.Where("name = ? AND sport_id = ? AND college_id = ?", name, sportID, collegeID).Limit(1).Find(&athlete)
	if err := result.Error; err != nil {
		return nil, false, err
	}
	return &athlete, result.RowsAffected > 0, nil
}
