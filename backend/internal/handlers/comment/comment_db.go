package comment

import (
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentDB struct {
	db *gorm.DB
}

// Creates a new CommentDB instance
func NewCommentDB(db *gorm.DB) *CommentDB {
	return &CommentDB{db: db}
}

// Retrieves a comment by its ID
func (c *CommentDB) GetCommentByID(id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	dbResponse := c.db.Where("id = ?", id).First(&comment)
	return utils.HandleDBError(&comment, dbResponse.Error)
}

// Creates a new comment in the database
func (c *CommentDB) CreateComment(comment *models.Comment) (*models.Comment, error) {
	dbResponse := c.db.Create(comment)
	return utils.HandleDBError(comment, dbResponse.Error)
}

// Retrieves top-level comments for a post
func (c *CommentDB) GetCommentsByPost(postID uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	res := c.db.Where("post_id = ? AND parent_comment_id IS NULL", postID).Order("created_at ASC").Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	return comments, nil
}

// Retrieves replies to a comment
func (c *CommentDB) GetReplies(commentID uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	res := c.db.Where("parent_comment_id = ?", commentID).Order("created_at ASC").Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	return comments, nil
}

// Updates an existing comment
func (c *CommentDB) UpdateComment(comment *models.Comment) (*models.Comment, error) {
	dbResponse := c.db.Save(comment)
	return utils.HandleDBError(comment, dbResponse.Error)
}

// Soft deletes a comment by ID
func (c *CommentDB) DeleteComment(id uuid.UUID) error {
	dbResponse := c.db.Delete(&models.Comment{}, "id = ?", id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError((*models.Comment)(nil), dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
