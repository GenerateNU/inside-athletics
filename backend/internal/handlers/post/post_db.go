package post

import (
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostDB struct {
	db *gorm.DB
}

// NewPostDB creates a new PostDB instance
func NewPostDB(db *gorm.DB) *PostDB {
	return &PostDB{db: db}
}

// CreatePost creates a new sport in the database
func (s *PostDB) CreatePost(author_id uuid.UUID, sport_id uuid.UUID, title string, content string, is_anonymous bool) (*models.Post, error) {
	post := models.Post{
		AuthorId:    author_id,
		SportId:     sport_id,
		Title:       title,
		Content:     content,
		IsAnonymous: is_anonymous,
	}
	dbResponse := s.db.Create(&post)
	return utils.HandleDBError(&post, dbResponse.Error)
}

// GetPostByID retrieves a post by its ID
func (s *PostDB) GetPostByID(id uuid.UUID) (*models.Post, error) {
	var post models.Post
	dbResponse := s.db.First(&post, "id = ?", id)
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
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.Post{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
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
	dbResponse := p.db.Limit(limit).Offset(offset).Find(&posts)
	if dbResponse.Error != nil {
		return nil, 0, dbResponse.Error
	}

	return posts, int(total), nil
}

// UpdatePost updates an existing post
func (p *PostDB) UpdatePost(post *models.Post) (*models.Post, error) {
	dbResponse := p.db.Save(post)
	return utils.HandleDBError(post, dbResponse.Error)
}
