package post

import (
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostDB struct {
	db *gorm.DB
}

// NewPostDB creates a new postDB instance
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

// GetPostByAuthorID retrieves a post by its author ID
func (s *PostDB) GetPostByAuthorID(author_id uuid.UUID) (*models.Post, error) {
	var post models.Post
	dbResponse := s.db.First(&post, "author_id = ?", author_id)
	return utils.HandleDBError(&post, dbResponse.Error)
}
