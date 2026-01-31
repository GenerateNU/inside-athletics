package post

import (
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postDB struct {
	db *gorm.DB
}

// NewPostDB creates a new postDB instance
func NewPostDB(db *gorm.DB) *PostDB {
	return &PostDB{db: db}
}

// CreatePost creates a new sport in the database
func (s *PostDB) CreatePost(uuid.UUID author_id, uuid.UUID sport_id, ) (*models.Post, error) {
	sport := models.Post{
		AuthorId   uuid.UUID      `json:"author_id" type:"uuid" default:"gen_random_uuid()"`
		SportId    uuid.UUID      `json:"sport_id" type:"uuid" default:"gen_random_uuid()"`
		Title      string         `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
		Content    string         `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
		UpVotes	   int64          `json:"numUpVotes,omitempty" example:"20000" gorm:"type:int"`
		DownVotes  int64          `json:"numDownVotes,omitempty" example:"20000" gorm:"type:int"`
		IsAnonymous bool          `json:"isAnanymous"`
	}
	dbResponse := s.db.Create(&sport)
	return utils.HandleDBError(&sport, dbResponse.Error) 
}

// GetSportByID retrieves a sport by its ID
func (s *SportDB) GetSportByID(id uuid.UUID) (*models.Sport, error) {
	var sport models.Sport
	dbResponse := s.db.First(&sport, "id = ?", id)
	return utils.HandleDBError(&sport, dbResponse.Error)
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