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
<<<<<<< HEAD
func (s *PostDB) CreatePost(author_id uuid.UUID, sport_id uuid.UUID, title string, content string, is_anonymous bool) (*models.Post, error) {
	post := models.Post{
		AuthorId: author_id,
		SportId: sport_id,
		Title: title,
		Content: content,
		UpVotes: 0, //default 0
		DownVotes: 0, //default 0
		IsAnonymous: is_anonymous,
=======
func (s *PostDB) CreatePost(uuid.UUID author_id, uuid.UUID sport_id, ) (*models.Post, error) {
	sport := models.Post{
		AuthorId   uuid.UUID      `json:"author_id" type:"uuid" default:"gen_random_uuid()"`
		SportId    uuid.UUID      `json:"sport_id" type:"uuid" default:"gen_random_uuid()"`
		Title      string         `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
		Content    string         `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
		UpVotes	   int64          `json:"numUpVotes,omitempty" example:"20000" gorm:"type:int"`
		DownVotes  int64          `json:"numDownVotes,omitempty" example:"20000" gorm:"type:int"`
		IsAnonymous bool          `json:"isAnonymous"`
>>>>>>> acd49d2 (IsAnanymous)
	}
	dbResponse := s.db.Create(&post)
	return utils.HandleDBError(&post, dbResponse.Error)
}

// GetSportByID retrieves a sport by its ID
func (s *SportDB) GetSportByID(id uuid.UUID) (*models.Sport, error) {
	var sport models.Sport
	dbResponse := s.db.First(&sport, "id = ?", id)
	return utils.HandleDBError(&sport, dbResponse.Error)
}

// GetPostByAuthorID retrieves a post by its author ID
func (s *PostDB) GetPostByAuthorID(author_id uuid.UUID) (*models.Post, error) {
	var post models.Post
	dbResponse := s.db.First(&post, "author_id = ?", id)
	return utils.HandleDBError(&post, dbResponse.Error)
}
