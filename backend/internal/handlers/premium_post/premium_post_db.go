package premiumpost

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PremiumPostDB struct {
	db *gorm.DB
}

// Create a new PremiumPostDB instance
func NewPremiumPostDB(db *gorm.DB) *PremiumPostDB {
	return &PremiumPostDB{db: db}
}

// CreatePremiumPost creates a new premium post in the database
func (s *PremiumPostDB) CreatePremiumPost(premiumPost *models.PremiumPost) (*models.PremiumPost, error) {
	dbResponse := s.db.Create(premiumPost)
	return utils.HandleDBError(premiumPost, dbResponse.Error)
}

// GetAllPremiumPosts returns all premium posts in the database
func (s *PremiumPostDB) GetAllPremiumPosts(limit, offset int) ([]models.PremiumPost, int64, error) {
	var posts []models.PremiumPost
	var total int64

	// Get total count
	if err := s.db.Model(&models.PremiumPost{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated posts
	if err := s.db.
		Model(&models.PremiumPost{}).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetPremiumPostsBySportID returns all premium posts related to a given sport
func (s *PremiumPostDB) GetPremiumPostsBySportID(limit, offset int, sportID uuid.UUID) ([]models.PremiumPost, int64, error) {
	var posts []models.PremiumPost
	var total int64

	if err := s.db.Model(&models.PremiumPost{}).
		Where("sport_id = ?", sportID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.
		Model(&models.PremiumPost{}).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags").
		Where("sport_id = ?", sportID).
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetPremiumPostsByCollegeID returns all premium posts related to a given college
func (s *PremiumPostDB) GetPremiumPostsByCollegeID(limit, offset int, collegeID uuid.UUID) ([]models.PremiumPost, int64, error) {
	var posts []models.PremiumPost
	var total int64

	if err := s.db.Model(&models.PremiumPost{}).
		Where("college_id = ?", collegeID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.
		Model(&models.PremiumPost{}).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags").
		Where("college_id = ?", collegeID).
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetPremiumPostsByTagID returns all premium posts related to a given tag
func (s *PremiumPostDB) GetPremiumPostsByTagID(limit, offset int, tagID uuid.UUID) ([]models.PremiumPost, int64, error) {
	var posts []models.PremiumPost
	var total int64

	base := s.db.Model(&models.PremiumPost{}).
		Joins("JOIN tag_posts tp ON tp.premium_post_id = premium_posts.id").
		Where("tp.tag_id = ?", tagID)

	// if this jointable count is 0, there are no premium posts with the given tag
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.
		Model(&models.PremiumPost{}).
		Joins("JOIN tag_posts tp ON tp.premium_post_id = premium_posts.id").
		Where("tp.tag_id = ?", tagID).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
