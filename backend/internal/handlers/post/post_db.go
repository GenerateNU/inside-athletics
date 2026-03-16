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

// CountViewedPostsByUser returns the number of distinct posts the user has viewed (for free-tier limit).
func (s *PostDB) CountViewedPostsByUser(userID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.Model(&models.ViewedPost{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// HasUserViewedPost returns true if the user has already viewed this post.
func (s *PostDB) HasUserViewedPost(userID, postID uuid.UUID) (bool, error) {
	var count int64
	err := s.db.Model(&models.ViewedPost{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

// RecordPostView records that the user viewed the post (idempotent: safe to call if already viewed).
func (s *PostDB) RecordPostView(userID, postID uuid.UUID) error {
	v := &models.ViewedPost{UserID: userID, PostID: postID}
	return s.db.Where("user_id = ? AND post_id = ?", userID, postID).FirstOrCreate(v).Error
}

// CountPostsByAuthor returns how many posts the user has created (for free-tier create limit).
func (s *PostDB) CountPostsByAuthor(authorID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.Model(&models.Post{}).Where("author_id = ?", authorID).Count(&count).Error
	return count, err
}

// GetPostByID retrieves a post by its ID
func (s *PostDB) GetPostByID(id uuid.UUID, userID uuid.UUID) (*models.Post, error) {
	var post models.Post
	dbResponse := s.db.
		Model(&models.Post{}).
		Select(`posts.*,
            (SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) AS like_count,
            (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
            (SELECT COUNT(*) > 0 FROM post_likes WHERE post_likes.post_id = posts.id AND post_likes.user_id = ?) AS is_liked`,
			userID).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tag_posts AS tp").Joins("JOIN tags t ON t.id = tp.tag_id")
		}).
		First(&post, "posts.id = ?", id)

	return utils.HandleDBError(&post, dbResponse.Error)
}

// GetPostsBySportID retrieves all posts with the given sport ID
func (s *PostDB) GetPostsBySportID(limit, offset int, sportID uuid.UUID, userID uuid.UUID) ([]models.Post, int64, error) {
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
		Table("posts").
		Select(`posts.*,
            (SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) AS like_count,
            (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
            (SELECT COUNT(*) > 0 FROM post_likes WHERE post_likes.post_id = posts.id AND post_likes.user_id = ?) AS is_liked`,
			userID).
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
func (s *PostDB) GetPostsByAuthorID(limit, offset int, authorID uuid.UUID, userID uuid.UUID) ([]models.Post, int64, error) {
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
		Table("posts").
		Select(`posts.*,
            (SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) AS like_count,
            (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
            (SELECT COUNT(*) > 0 FROM post_likes WHERE post_likes.post_id = posts.id AND post_likes.user_id = ?) AS is_liked`,
			userID).
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
func (p *PostDB) GetAllPosts(limit int, offset int, userID uuid.UUID) ([]models.Post, int, error) {
	var posts []models.Post
	var total int64

	// Get total count
	if err := p.db.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated posts
	dbResponse := p.db.
		Table("posts").
		Select(`posts.*,
            (SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) AS like_count,
            (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
            (SELECT COUNT(*) > 0 FROM post_likes WHERE post_likes.post_id = posts.id AND post_likes.user_id = ?) AS is_liked`,
			userID).
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
func (p *PostDB) UpdatePost(id uuid.UUID, updates UpdatePostRequest, userID uuid.UUID) (*models.Post, error) {
	var updatedPost models.Post
	dbResponse := p.db.Model(&models.Post{}).
		Clauses(clause.Returning{}).
		Select(`posts.*,
            (SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) AS like_count,
            (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
            (SELECT COUNT(*) > 0 FROM post_likes WHERE post_likes.post_id = posts.id AND post_likes.user_id = ?) AS is_liked`,
			userID).
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
