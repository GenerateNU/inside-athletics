package sport

import (
	models "inside-athletics/internal/models"

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
	sport := &models.Sport{
		Name:       name,
		Popularity: popularity,
	}

	if err := s.db.Create(sport).Error; err != nil {
		return nil, err
	}

	return sport, nil
}

// GetSportByID retrieves a sport by its ID
func (s *SportDB) GetSportByID(id uuid.UUID) (*models.Sport, error) {
	var sport models.Sport
	if err := s.db.First(&sport, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &sport, nil
}

// GetSportByName retrieves a sport by its name
func (s *SportDB) GetSportByName(name string) (*models.Sport, error) {
	var sport models.Sport
	if err := s.db.Where("name = ?", name).First(&sport).Error; err != nil {
		return nil, err
	}
	return &sport, nil
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
	if err := s.db.Save(sport).Error; err != nil {
		return nil, err
	}
	return sport, nil
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
