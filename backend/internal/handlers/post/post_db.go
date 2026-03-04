package post

import (
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostDB struct {
	db *gorm.DB
}

// NewPostDB creates a new PostDB instance
func NewPostDB(db *gorm.DB) *PostDB {
	return &PostDB{db: db}
}

// CreatePost creates a new sport in the database
func (s *PostDB) CreatePost(post *models.Post, tags []TagRequest) (*models.Post, error) {
	dbError := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&post).Error; err != nil {
			return err
		}
		var tagModels []models.Tag
		for _, t := range tags {
			tagModels = append(tagModels, models.Tag{ID: t.ID})
		}
		if err := tx.Model(&post).Association("Tags").Append(&tagModels); err != nil {
			return err
		}

		return nil
	})
	return utils.HandleDBError(post, dbError)
}

// GetPostByID retrieves a post by its ID
func (s *PostDB) GetPostByID(id uuid.UUID) (*models.Post, error) {
	var post models.Post
	dbResponse := s.db.
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tag_posts AS tp").Joins("JOIN tags t ON t.id = tp.tag_id")
		}).
		First(&post, "id = ?", id)
	return utils.HandleDBError(&post, dbResponse.Error)
}

// GetPostsBySportID retrieves all posts with the given sport ID
func (s *PostDB) GetPostsBySportID(limit, offset int, sportID uuid.UUID) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// Count total matching posts
	if err := s.db.Model(&models.Post{}).
		Where("sport_id = ?", sportID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := s.db.
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tag_posts AS tp").Joins("JOIN tags t ON t.id = tp.tag_id")
		}).
		Where("sport_id = ?", sportID).
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetPostByAuthorID retrieves a post by its author ID
func (s *PostDB) GetPostsByAuthorID(limit, offset int, authorID uuid.UUID) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// Count total matching posts
	if err := s.db.Model(&models.Post{}).
		Where("author_id = ?", authorID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := s.db.
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tag_posts AS tp").Joins("JOIN tags t ON t.id = tp.tag_id")
		}).
		Where("author_id = ?", authorID).
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// DeletePost soft deletes a post by ID
func (p *PostDB) DeletePost(id uuid.UUID) error {
	dbResponse := p.db.Delete(&models.Post{}, "id = ?", id)
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.User{}, dbResponse.Error)
		return err
	}
	return nil
}

// GetAllPosts retrieves all posts with pagination
func (p *PostDB) GetAllPosts(limit int, offset int) ([]models.Post, int, error) {
	var posts []models.Post
	var total int64

	// Get total count
	if err := p.db.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated posts
	dbResponse := p.db.
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tag_posts AS tp").Joins("JOIN tags t ON t.id = tp.tag_id")
		}).
		Limit(limit).
		Offset(offset).
		Find(&posts)
	if dbResponse.Error != nil {
		return nil, 0, dbResponse.Error
	}

	return posts, int(total), nil
}

// UpdatePost updates an existing post
func (p *PostDB) UpdatePost(id uuid.UUID, updates UpdatePostRequest) (*models.Post, error) {
	var updatedPost models.Post
	dbResponse := p.db.Model(&models.Post{}).
		Clauses(clause.Returning{}).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tag_posts AS tp").Joins("JOIN tags t ON t.id = tp.tag_id")
		}).
		Where("id = ?", id).
		Updates(updates).
		Scan(&updatedPost)
	if dbResponse.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Resource not found")
	}
	return utils.HandleDBError(&updatedPost, dbResponse.Error)
}
