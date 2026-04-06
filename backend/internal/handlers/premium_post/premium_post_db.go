package premiumpost

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id AND tp.postable_type = 'premium_post'")
		}).
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetPremiumPostsByAuthorID returns all premium posts related to a given author
func (s *PremiumPostDB) GetPremiumPostsByAuthorID(limit, offset int, authorID uuid.UUID) ([]models.PremiumPost, int64, error) {
	var posts []models.PremiumPost
	var total int64

	// check if there are actually premium posts where the given author is the author
	if err := s.db.Model(&models.PremiumPost{}).
		Where("author_id = ?", authorID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.
		Model(&models.PremiumPost{}).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id AND tp.postable_type = 'premium_post'")
		}).
		Where("author_id = ?", authorID).
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
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id AND tp.postable_type = 'premium_post'")
		}).
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
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id AND tp.postable_type = 'premium_post'")
		}).
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
		Joins("JOIN tag_posts tp ON tp.postable_id = premium_posts.id AND tp.postable_type = 'premium_post'").
		Where("tp.tag_id = ?", tagID)

	// if this jointable count is 0, there are no premium posts with the given tag
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.
		Model(&models.PremiumPost{}).
		Joins("JOIN tag_posts tp ON tp.postable_id = premium_posts.id AND tp.postable_type = 'premium_post'").
		Where("tp.tag_id = ?", tagID).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id AND tp.postable_type = 'premium_post'")
		}).
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// UpdatePremiumPost updates an existing premium post
func (s *PremiumPostDB) UpdatePremiumPost(id uuid.UUID, updates UpdatePremiumPostRequest, userID uuid.UUID) (*models.PremiumPost, error) {
	var updatedPost models.PremiumPost
	dbResponse := s.db.Model(&models.PremiumPost{}).
		Clauses(clause.Returning{}).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id AND tp.postable_type = 'premium_post'")
		}).
		Where("id = ?", id).
		Updates(updates).
		Scan(&updatedPost)

	if dbResponse.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Resource not found")
	}
	return utils.HandleDBError(&updatedPost, dbResponse.Error)
}

// DeletePremiumPost soft deletes a premium post by ID
func (s *PremiumPostDB) DeletePremiumPost(id uuid.UUID) error {
	dbResponse := s.db.Delete(&models.PremiumPost{}, "id = ?", id)
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.PremiumPost{}, dbResponse.Error)
		return err
	}
	return nil
}
