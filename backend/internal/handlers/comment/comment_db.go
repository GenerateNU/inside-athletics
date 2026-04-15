package comment

import (
	"errors"
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
func (c *CommentDB) GetCommentByID(id uuid.UUID, userID uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	dbResponse := c.db.
		Model(&models.Comment{}).
		Select(`comments.*,
            (SELECT COUNT(*) FROM comment_likes WHERE comment_likes.comment_id = comments.id) AS like_count,
            (SELECT COUNT(*) > 0 FROM comment_likes WHERE comment_likes.comment_id = comments.id AND comment_likes.user_id = ?) AS is_liked,
            (SELECT COUNT(*) > 0 FROM comments AS replies WHERE replies.parent_comment_id = comments.id AND replies.deleted_at IS NULL) AS has_replies`,
			userID).
		Preload("User").
		Where("id = ?", id).
		First(&comment)
	return utils.HandleDBError(&comment, dbResponse.Error)
}

// Creates a new comment in the database
func (c *CommentDB) CreateComment(comment *models.Comment) (*models.Comment, error) {
	dbResponse := c.db.Create(comment)
	return utils.HandleDBError(comment, dbResponse.Error)
}

// Retrieves top-level comments for a post
func (c *CommentDB) GetCommentsByPost(postID uuid.UUID, userID uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	res := c.db.
		Model(&models.Comment{}).
		Select(`comments.*,
            (SELECT COUNT(*) FROM comment_likes WHERE comment_likes.comment_id = comments.id) AS like_count,
            (SELECT COUNT(*) > 0 FROM comment_likes WHERE comment_likes.comment_id = comments.id AND comment_likes.user_id = ?) AS is_liked,
            (SELECT COUNT(*) > 0 FROM comments AS replies WHERE replies.parent_comment_id = comments.id AND replies.deleted_at IS NULL) AS has_replies`,
			userID).
		Preload("User").
		Where("post_id = ? AND parent_comment_id IS NULL", postID).
		Order("created_at ASC").
		Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	return comments, nil
}

// Retrieves replies to a comment
func (c *CommentDB) GetReplies(commentID uuid.UUID, userID uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	res := c.db.
		Select(`comments.*,
            (SELECT COUNT(*) FROM comment_likes WHERE comment_likes.comment_id = comments.id) AS like_count,
            (SELECT COUNT(*) > 0 FROM comment_likes WHERE comment_likes.comment_id = comments.id AND comment_likes.user_id = ?) AS is_liked`,
			userID).
		Preload("User").
		Where("parent_comment_id = ?", commentID).
		Order("created_at ASC").
		Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	return comments, nil
}

// Updates an existing comment by ID
func (c *CommentDB) UpdateComment(id uuid.UUID, updates UpdateCommentBody, userID uuid.UUID) (*models.Comment, error) {
	dbResponse := c.db.Model(&models.Comment{}).
		Model(&models.Comment{}).
		Select(`comments.*,
            (SELECT COUNT(*) FROM comment_likes WHERE comment_likes.comment_id = comments.id) AS like_count,
            (SELECT COUNT(*) > 0 FROM comment_likes WHERE comment_likes.comment_id = comments.id AND comment_likes.user_id = ?) AS is_liked`,
			userID).
		Preload("User").
		Where("id = ?", id).
		Updates(updates)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError((*models.Comment)(nil), dbResponse.Error)
		return nil, err
	}
	if dbResponse.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Resource not found")
	}
	return c.GetCommentByID(id, userID)
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

// IsUserPremium returns true when the user has premium access.
// If the user record is missing, we do not apply free-tier restrictions.
func (c *CommentDB) IsUserPremium(userID uuid.UUID) (bool, error) {
	var user models.User
	err := c.db.Select("account_type").Where("id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return user.Account_Type, nil
}

// IsPostWithinFirstViewedPosts returns true when the post is in the user's first N viewed posts.
func (c *CommentDB) IsPostWithinFirstViewedPosts(userID, postID uuid.UUID, maxPosts int) (bool, error) {
	var count int64
	subQuery := c.db.Model(&models.ViewedPost{}).
		Select("post_id").
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Limit(maxPosts)

	err := c.db.Model(&models.ViewedPost{}).
		Where("user_id = ? AND post_id = ? AND post_id IN (?)", userID, postID, subQuery).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
