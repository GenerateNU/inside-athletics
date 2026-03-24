package premiumpost

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"gorm.io/gorm"
)

type PremiumPostDB struct {
	db *gorm.DB
}

// CreatePremiumPost creates a new premium post in the database
func (s *PremiumPostDB) CreatePremiumPost(premiumPost *models.PremiumPost) (*models.PremiumPost, error) {
	dbResponse := s.db.Create(premiumPost)
	return utils.HandleDBError(premiumPost, dbResponse.Error)
}
