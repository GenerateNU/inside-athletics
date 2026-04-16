package post

import (
	"errors"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"math"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostDB struct {
	db *gorm.DB
}

var (
	ErrFreePostCreationLimitReached = errors.New("free-tier post creation limit reached")
	ErrFreePostViewLimitReached     = errors.New("free-tier post view limit reached")
)

// NewPostDB creates a new PostDB instance
func NewPostDB(db *gorm.DB) *PostDB {
	return &PostDB{db: db}
}

// CreatePost creates a new sport in the database
func (s *PostDB) CreatePost(post *models.Post, tags []TagRequest) (*models.Post, error) {
	dbError := s.db.Transaction(func(tx *gorm.DB) error {
		return s.createPostTx(tx, post, tags)
	})
	return utils.HandleDBError(post, dbError)
}

// CreatePostWithAuthorLimit creates a post while enforcing the author's post cap atomically.
func (s *PostDB) CreatePostWithAuthorLimit(post *models.Post, tags []TagRequest, maxPosts int64) (*models.Post, error) {
	dbError := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.lockUserForUpdate(tx, post.AuthorID); err != nil {
			return err
		}

		var count int64
		if err := tx.Model(&models.Post{}).Where("author_id = ?", post.AuthorID).Count(&count).Error; err != nil {
			return err
		}
		if count >= maxPosts {
			return ErrFreePostCreationLimitReached
		}

		return s.createPostTx(tx, post, tags)
	})
	if dbError != nil {
		if errors.Is(dbError, ErrFreePostCreationLimitReached) {
			return nil, dbError
		}
		return utils.HandleDBError(post, dbError)
	}
	return post, nil
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

// RecordPostViewIfAllowed records a new post view while enforcing the user's distinct-view cap atomically.
func (s *PostDB) RecordPostViewIfAllowed(userID, postID uuid.UUID, maxViews int64) error {
	dbError := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.lockUserForUpdate(tx, userID); err != nil {
			return err
		}

		var existing models.ViewedPost
		err := tx.Where("user_id = ? AND post_id = ?", userID, postID).Take(&existing).Error
		switch {
		case err == nil:
			return nil
		case !errors.Is(err, gorm.ErrRecordNotFound):
			return err
		}

		var count int64
		if err := tx.Model(&models.ViewedPost{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
			return err
		}
		if count >= maxViews {
			return ErrFreePostViewLimitReached
		}

		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "post_id"}},
			DoNothing: true,
		}).Create(&models.ViewedPost{UserID: userID, PostID: postID}).Error
	})
	if dbError != nil {
		if errors.Is(dbError, ErrFreePostViewLimitReached) {
			return dbError
		}
		_, err := utils.HandleDBError((*models.ViewedPost)(nil), dbError)
		return err
	}
	return nil
}

// CountPostsByAuthor returns how many posts the user has created (for free-tier create limit).
func (s *PostDB) CountPostsByAuthor(authorID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.Model(&models.Post{}).Where("author_id = ?", authorID).Count(&count).Error
	return count, err
}

func (s *PostDB) createPostTx(tx *gorm.DB, post *models.Post, tags []TagRequest) error {
	if err := tx.Create(post).Error; err != nil {
		return err
	}
	var tagModels []models.Tag
	for _, t := range tags {
		tagModels = append(tagModels, models.Tag{ID: t.ID})
	}
	if err := tx.Model(post).Association("Tags").Append(&tagModels); err != nil {
		return err
	}

	return nil
}

func (s *PostDB) lockUserForUpdate(tx *gorm.DB, userID uuid.UUID) error {
	return tx.
		Model(&models.User{}).
		Select("id").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", userID).
		Take(&models.User{}).Error
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
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id")
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
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id")
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
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id")
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
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id")
		}).
		Limit(limit).
		Offset(offset).
		Find(&posts)
	if dbResponse.Error != nil {
		return nil, 0, dbResponse.Error
	}

	return posts, int(total), nil
}

func (p *PostDB) GetPopularPosts(limit int, offset int, windowHours int, userID uuid.UUID) ([]models.Post, int, error) {
	var posts []models.Post
	var total int64

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	if windowHours <= 0 {
		windowHours = 72
	}
	windowHours = min(windowHours, 24*30)
	recencyWindow := float64(windowHours)

	if err := p.db.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	dbResponse := p.db.
		Table("posts").
		Select(`
			posts.*,
			COALESCE(all_likes.total_likes, 0) AS like_count,
			COALESCE(all_comments.total_comments, 0) AS comment_count,
			COALESCE(user_likes.is_liked, false) AS is_liked,
			(
				COALESCE(recent_comments.recent_comments, 0) * 8.0 +
				COALESCE(all_comments.total_comments, 0) * 2.0 +
				COALESCE(recent_likes.recent_likes, 0) * 3.0 +
				COALESCE(all_likes.total_likes, 0) * 1.0 +
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM tag_follows tf
						JOIN tag_posts tp_sub ON tp_sub.tag_id = tf.tag_id
						WHERE tf.user_id = ? AND tp_sub.post_id = posts.id
					) THEN 12.0
					ELSE 0.0
				END +
				CASE
					WHEN posts.sport_id IS NOT NULL AND EXISTS (
						SELECT 1
						FROM sport_follows sf
						WHERE sf.user_id = ? AND sf.sport_id = posts.sport_id
					) THEN 4.0
					ELSE 0.0
				END +
				CASE
					WHEN posts.college_id IS NOT NULL AND EXISTS (
						SELECT 1
						FROM college_follows cf
						WHERE cf.user_id = ? AND cf.college_id = posts.college_id
					) THEN 2.0
					ELSE 0.0
				END +
				GREATEST(0.0, ? - (EXTRACT(EPOCH FROM (NOW() - posts.created_at)) / 3600.0)) * 0.15
			) AS popularity_score`,
			userID, userID, userID, recencyWindow).
		Joins(`
			LEFT JOIN (
				SELECT post_id, COUNT(*) AS total_likes
				FROM post_likes
				GROUP BY post_id
			) AS all_likes ON all_likes.post_id = posts.id`).
		Joins(`
			LEFT JOIN (
				SELECT post_id, COUNT(*) AS total_comments
				FROM comments
				GROUP BY post_id
			) AS all_comments ON all_comments.post_id = posts.id`).
		Joins(`
			LEFT JOIN (
				SELECT post_id, COUNT(*) AS recent_likes
				FROM post_likes
				WHERE created_at >= NOW() - (? * INTERVAL '1 hour')
				GROUP BY post_id
			) AS recent_likes ON recent_likes.post_id = posts.id`, windowHours).
		Joins(`
			LEFT JOIN (
				SELECT post_id, COUNT(*) AS recent_comments
				FROM comments
				WHERE created_at >= NOW() - (? * INTERVAL '1 hour')
				GROUP BY post_id
			) AS recent_comments ON recent_comments.post_id = posts.id`, windowHours).
		Joins(`
			LEFT JOIN (
				SELECT post_id, true AS is_liked
				FROM post_likes
				WHERE user_id = ?
				GROUP BY post_id
			) AS user_likes ON user_likes.post_id = posts.id`, userID).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id")
		}).
		Order("popularity_score DESC").
		Order("comment_count DESC").
		Order("like_count DESC").
		Order("posts.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts)
	if dbResponse.Error != nil {
		return nil, 0, dbResponse.Error
	}

	for i := range posts {
		posts[i].PopularityScore = math.Round(posts[i].PopularityScore*100) / 100
	}

	return posts, int(total), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id")
		}).
		Where("id = ?", id).
		Updates(updates).
		Scan(&updatedPost)
	if dbResponse.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Resource not found")
	}
	return utils.HandleDBError(&updatedPost, dbResponse.Error)
}
