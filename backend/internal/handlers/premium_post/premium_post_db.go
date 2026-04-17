package premiumpost

import (
	"fmt"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"strings"

	"github.com/danielgtaylor/huma/v2"
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
	result, err := utils.HandleDBError(premiumPost, dbResponse.Error)
	if err != nil {
		return nil, err
	}
	if err := s.db.Preload("Media").First(result, result.ID).Error; err != nil {
		return nil, err
	}
	return result, nil
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
		Preload("Media").
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
		Preload("Media").
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
		Preload("Media").
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
		Preload("Media").
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
		Preload("Media").
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

// FuzzySearchForPremiumPost returns premium posts whose title fuzzy-matches the given search string
func (s *PremiumPostDB) FuzzySearchForPremiumPost(searchStr string, limit, offset int) ([]models.PremiumPost, int64, error) {
	var posts []models.PremiumPost
	var total int64

	selectQuery, whereQuery, orderQuery := utils.FuzzySearchByQueries("title", searchStr)

	if err := s.db.Model(&models.PremiumPost{}).
		Select("premium_posts.*, "+selectQuery).
		Where(whereQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []models.PremiumPost{}, 0, nil
	}

	if err := s.db.
		Model(&models.PremiumPost{}).
		Select("premium_posts.*, "+selectQuery).
		Where(whereQuery).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Media").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id AND tp.postable_type = 'premium_post'")
		}).
		Order(orderQuery).
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// FilterPremiumPosts returns premium posts filtered by the given college, sport, and tag IDs
func (s *PremiumPostDB) FilterPremiumPosts(colleges []uuid.UUID, sports []uuid.UUID, tags []uuid.UUID, limit, offset int) ([]models.PremiumPost, int64, error) {
	var posts []models.PremiumPost
	var total int64

	filters := make([]string, 0)
	if len(colleges) > 0 {
		mappedColleges := utils.MapList(colleges, func(c uuid.UUID) string {
			return fmt.Sprintf("premium_posts.college_id = '%s'", c.String())
		})
		filters = append(filters, strings.Join(mappedColleges, " OR "))
	}
	if len(sports) > 0 {
		mappedSports := utils.MapList(sports, func(c uuid.UUID) string {
			return fmt.Sprintf("premium_posts.sport_id = '%s'", c.String())
		})
		filters = append(filters, strings.Join(mappedSports, " OR "))
	}
	if len(tags) > 0 {
		mappedTags := utils.MapList(tags, func(c uuid.UUID) string {
			return fmt.Sprintf("tag_posts.tag_id = '%s'", c.String())
		})
		filters = append(filters, strings.Join(mappedTags, " OR "))
	}

	whereQuery := strings.Join(filters, " OR ")

	if err := s.db.
		Model(&models.PremiumPost{}).
		Joins("JOIN tag_posts ON premium_posts.id = tag_posts.postable_id AND tag_posts.postable_type = 'premium_post'").
		Where(whereQuery).
		Group("premium_posts.id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.
		Model(&models.PremiumPost{}).
		Joins("JOIN tag_posts ON premium_posts.id = tag_posts.postable_id AND tag_posts.postable_type = 'premium_post'").
		Where(whereQuery).
		Group("premium_posts.id").
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Media").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id AND tp.postable_type = 'premium_post'")
		}).
		Limit(limit).
		Offset(offset).
		Order("premium_posts.created_at DESC").
		Find(&posts).Error; err != nil {
		return posts, 0, err
	}

	return posts, total, nil
}

// UpdatePremiumPost updates an existing premium post
func (s *PremiumPostDB) UpdatePremiumPost(id uuid.UUID, updates UpdatePremiumPostRequest, userID uuid.UUID) (*models.PremiumPost, error) {
	dbResponse := s.db.Model(&models.PremiumPost{}).
		Where("id = ?", id).
		Updates(updates)

	if dbResponse.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Resource not found")
	}
	if dbResponse.Error != nil {
		return utils.HandleDBError(&models.PremiumPost{}, dbResponse.Error)
	}

	// reload with associations
	var updatedPost models.PremiumPost
	if err := s.db.
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id AND tp.postable_type = 'premium_post'")
		}).
		First(&updatedPost, "id = ?", id).Error; err != nil {
		return utils.HandleDBError(&updatedPost, err)
	}

	return &updatedPost, nil
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
