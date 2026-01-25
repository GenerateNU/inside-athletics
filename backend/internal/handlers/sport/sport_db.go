package sport

import (
	"fmt"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SportDB struct {
	db *gorm.DB
}

// Create sport
func (u *SportDB) CreateSport(name string, popularity int32) (*models.Sport, error) {
	if name == "" {
		return nil, fmt.Errorf("Name cannot be an empty string!")
	}
	if popularity < 0 {
		return nil, fmt.Errorf("Popularity cannot be a negative number!")
	}
	sport := models.Sport{Name: name, Popularity: popularity}
	dbResponse := u.db.Create(&sport)
	return utils.HandleDBError(&sport, dbResponse.Error)
}

// Read sport (by id)
func (u *SportDB) GetSportById(id uuid.UUID) (*models.Sport, error) {
	var sport models.Sport
	dbResponse := u.db.Where("id = ?", id).First(&sport)
	return utils.HandleDBError(&sport, dbResponse.Error)
}

// Read sport (by name)
func (u *SportDB) GetSportByName(name string) (*models.Sport, error) {
	var sport models.Sport
	dbResponse := u.db.Where("name = ?", name).First(&sport)
	return utils.HandleDBError(&sport, dbResponse.Error)
}

// Get all sports
func (u *SportDB) GetAllSports() ([]models.Sport, error) {
	var sports []models.Sport
	dbResponse := u.db.Find(&sports)
	return utils.HandleDBError(&sports, dbResponse.Error)
}

// Update sport
func (u *SportDB) UpdateSport(id uuid.UUID, name string, popularity int32) (*models.Sport, error) {
	if name == "" {
		return nil, fmt.Errorf("Name cannot be an empty string!")
	}
	if popularity < 0 {
		return nil, fmt.Errorf("Popularity cannot be a negative number!")
	}
	var sport models.Sport
	dbResponse := u.db.Where("id = ?", id).First(&sport)
	if dbResponse.Error != nil {
		return nil, fmt.Errorf("Sport does not exist in database!")
	}
	sport.Name = name
	sport.Popularity = popularity
	if err := u.db.Save(&sport).Error; err != nil {
		return nil, err
	}
	return utils.HandleDBError(&sport, dbResponse.Error)
}

// delete sport
func (u *SportDB) DeleteSport(id uuid.UUID) (*models.Sport, error) {
	var sport models.Sport
	dbResponse := u.db.Where("id = ?", id).Delete(&sport)
	return utils.HandleDBError(&sport, dbResponse.Error)
}
