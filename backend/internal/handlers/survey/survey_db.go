package survey

import (
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SurveyDB struct {
	db *gorm.DB
}

func NewSurveyDB(db *gorm.DB) *SurveyDB {
	return &SurveyDB{db: db}
}

// CreateSurvey creates a new survey response in the database
func (s *SurveyDB) CreateSurvey(req CreateSurveyRequest) (*models.Survey, error) {
	survey := models.Survey{
		UserID:                     req.UserID,
		CollegeID:                  req.CollegeID,
		SportID:                    req.SportID,
		PlayerDev:                  req.PlayerDev,
		AcademicsAthleticsPriority: req.AcademicsAthleticsPriority,
		AcademicCareerResources:    req.AcademicCareerResources,
		MentalHealthPriority:       req.MentalHealthPriority,
		Environment:                req.Environment,
		Culture:                    req.Culture,
		Transparency:               req.Transparency,
	}
	dbResponse := s.db.Create(&survey)
	return utils.HandleDBError(&survey, dbResponse.Error)
}

// GetSurveyByID retrieves a single survey by its ID
func (s *SurveyDB) GetSurveyByID(id uuid.UUID) (*models.Survey, error) {
	var survey models.Survey
	result := s.db.First(&survey, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &survey, nil
}

// GetSurveysByUserID retrieves all surveys submitted by a user with optional pagination
func (s *SurveyDB) GetSurveysByUserID(userID uuid.UUID, limit, offset int) ([]models.Survey, int64, error) {
	var surveys []models.Survey
	var total int64

	q := s.db.Model(&models.Survey{}).Where("user_id = ?", userID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Limit(limit).Offset(offset).Find(&surveys).Error; err != nil {
		return nil, 0, err
	}
	return surveys, total, nil
}

// DeleteSurvey soft deletes a survey by ID
func (s *SurveyDB) DeleteSurvey(id uuid.UUID) error {
	result := s.db.Delete(&models.Survey{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetAverageRatings returns average scores for each rating field,
// optionally filtered by sportID and/or collegeID, grouped by both.
func (s *SurveyDB) GetAverageRatings(sportID, collegeID uuid.UUID) ([]AverageRatingsRow, error) {
	q := s.db.Model(&models.Survey{}).
		Select(`
			sport_id,
			college_id,
			AVG(player_dev)                       AS player_dev,
			AVG(academics_athletics_priority)     AS academics_athletics_priority,
			AVG(academic_career_resources)        AS academic_career_resources,
			AVG(mental_health_priority)           AS mental_health_priority,
			AVG(environment)                      AS environment,
			AVG(culture)                          AS culture,
			AVG(transparency)                     AS transparency,
			COUNT(*)                              AS response_count
		`).
		Group("sport_id, college_id")

	if sportID != uuid.Nil {
		q = q.Where("sport_id = ?", sportID)
	}
	if collegeID != uuid.Nil {
		q = q.Where("college_id = ?", collegeID)
	}

	var rows []AverageRatingsRow
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
