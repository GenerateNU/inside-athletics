package post

import (
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

const (
	POST_SELECT_QUERY string = `posts.*,
            (SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) AS like_count,
            (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
            (SELECT COUNT(*) > 0 FROM post_likes WHERE post_likes.post_id = posts.id AND post_likes.user_id = ?) AS is_liked`
)

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
func (s *PostDB) GetPostByID(id uuid.UUID, userID uuid.UUID) (*models.Post, error) {
	var post models.Post
	dbResponse := s.db.
		Model(&models.Post{}).
		Select(POST_SELECT_QUERY,
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
		Select(POST_SELECT_QUERY,
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
		Select(POST_SELECT_QUERY,
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
		Select(POST_SELECT_QUERY,
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
		Select(POST_SELECT_QUERY,
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

func (p *PostDB) FuzzySearchForPost(userID uuid.UUID, searchStr string) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	selectQuery, whereQuery, orderQuery := utils.FuzzySearchBy("title", searchStr)
	selectQuery = POST_SELECT_QUERY + "," + selectQuery
	if err := p.db.
		Model(&models.Post{}).
		Select(selectQuery,
			userID).
		Where(whereQuery).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags AS t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id")
		}).
		Order(orderQuery).
		Scan(&posts).Count(&total).Error; err != nil {
		return posts, 0, err
	}

	if total == 0 {
		return []models.Post{}, 0, nil
	}

	return posts, total, nil
}
