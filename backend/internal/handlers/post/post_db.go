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

// CreatePost creates a new post in the database
func (p *PostDB) CreatePost(authorId, sportId uuid.UUID, title, content string, isAnonymous bool) (*models.Post, error) {
	post := models.Post{
		AuthorId:    authorId,
		SportId:     sportId,
		Title:       title,
		Content:     content,
		IsAnonymous: isAnonymous,
	}
	dbResponse := p.db.Create(&post)
	return utils.HandleDBError(&post, dbResponse.Error)
}

// GetPostByID retrieves a post by its ID
func (p *PostDB) GetPostByID(id uuid.UUID) (*models.Post, error) {
	var post models.Post
	dbResponse := p.db.First(&post, "id = ?", id)
	return utils.HandleDBError(&post, dbResponse.Error)
}

// GetPostBySportId retrieves all posts for a specific sport with pagination
func (p *PostDB) GetPostBySportId(sportId uuid.UUID, limit, offset int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	if err := p.db.Model(&models.Post{}).Where("sport_id = ?", sportId).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := p.db.Where("sport_id = ?", sportId).Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetAllPosts retrieves all posts with pagination
func (p *PostDB) GetAllPosts(limit, offset int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	// Get total count
	if err := p.db.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := p.db.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// UpdatePost updates an existing post
func (p *PostDB) UpdatePost(post *models.Post) (*models.Post, error) {
	dbResponse := p.db.Save(post)
	return utils.HandleDBError(post, dbResponse.Error)
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
