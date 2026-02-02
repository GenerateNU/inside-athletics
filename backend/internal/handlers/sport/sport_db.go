package sport

import (
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SportDB struct {
	db *gorm.DB
}

// NewSportDB creates a new SportDB instance
func NewSportDB(db *gorm.DB) *SportDB {
	return &SportDB{db: db}
}

// CreateSport creates a new sport in the database
func (s *SportDB) CreateSport(name string, popularity *int32) (*models.Sport, error) {
	sport := models.Sport{
		Name:       name,
		Popularity: popularity,
	}
	dbResponse := s.db.Create(&sport)
	return utils.HandleDBError(&sport, dbResponse.Error)
}

// GetSportByID retrieves a sport by its ID
func (s *SportDB) GetSportByID(id uuid.UUID) (*models.Sport, error) {
	var sport models.Sport
	result := s.db.First(&sport, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sport, nil
}

// GetSportByName retrieves a sport by its name
func (s *SportDB) GetSportByName(name string) (*models.Sport, error) {
	var sport models.Sport
	dbResponse := s.db.Where("name = ?", name).First(&sport)
	return utils.HandleDBError(&sport, dbResponse.Error)
}

// GetAllSports retrieves all sports with optional pagination
func (s *SportDB) GetAllSports(limit, offset int) ([]models.Sport, int64, error) {
	var sports []models.Sport
	var total int64

	// Get total count
	if err := s.db.Model(&models.Sport{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := s.db.Limit(limit).Offset(offset).Find(&sports).Error; err != nil {
		return nil, 0, err
	}

	return sports, total, nil
}

// UpdateSport updates an existing sport
func (s *SportDB) UpdateSport(sport *models.Sport) (*models.Sport, error) {
	dbResponse := s.db.Save(sport)
	return utils.HandleDBError(sport, dbResponse.Error)
}

// DeleteSport soft deletes a sport by ID
func (s *SportDB) DeleteSport(id uuid.UUID) error {
	result := s.db.Delete(&models.Sport{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
