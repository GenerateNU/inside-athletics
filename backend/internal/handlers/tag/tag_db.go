package tag

import (
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagDB struct {
	db *gorm.DB
}

func (u *TagDB) GetPostsByTag(tag_id uuid.UUID, limit int, offset int, userID uuid.UUID) (*[]models.Post, error) {
	var posts []models.Post
	dbResponse := u.db.
		Table("posts").
		Select(`posts.*,
            (SELECT COUNT(*) FROM post_likes WHERE post_likes.post_id = posts.id) AS like_count,
            (SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
            (SELECT COUNT(*) > 0 FROM post_likes WHERE post_likes.post_id = posts.id AND post_likes.user_id = ?) AS is_liked`,
			userID).
		Joins("JOIN tag_posts tp ON tp.post_id = posts.id").
		Where("tp.tag_id = ?", tag_id).
		Preload("Author").
		Preload("Sport", "id IS NOT NULL").
		Preload("College", "id IS NOT NULL").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Table("tags t").Joins("JOIN tag_posts tp ON tp.tag_id = t.id")
		}).
		Limit(limit).
		Offset(offset).
		Find(&posts)
	return utils.HandleDBError(&posts, dbResponse.Error)
}

func (u *TagDB) GetTagByName(name string) (*models.Tag, error) {
	var tag models.Tag
	dbResponse := u.db.Where("name = ?", name).First(&tag)
	return utils.HandleDBError(&tag, dbResponse.Error) // helper function that maps GORM errors to Huma errors
}

func (u *TagDB) GetTagByID(id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	dbResponse := u.db.Where("id = ?", id).First(&tag)
	return utils.HandleDBError(&tag, dbResponse.Error) // helper function that maps GORM errors to Huma errors
}

func (u *TagDB) GetTagByType(tagType models.TagType) (*models.Tag, error) {
	var tag models.Tag
	dbResponse := u.db.Where("type = ?", tagType).First(&tag)
	return utils.HandleDBError(&tag, dbResponse.Error) // helper function that maps GORM errors to Huma errors
}

func (u *TagDB) CreateTag(tag *models.Tag) (*models.Tag, error) {
	dbResponse := u.db.Create(tag)
	return utils.HandleDBError(tag, dbResponse.Error)
}

func (u *TagDB) UpdateTag(id uuid.UUID, updates *UpdateTagBody) (*models.Tag, error) {
	dbResponse := u.db.Model(&models.Tag{}).
		Where("id = ?", id).
		Updates(updates)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.Tag{}, dbResponse.Error)
		return nil, err
	}
	if dbResponse.RowsAffected == 0 {
		return nil, huma.Error404NotFound("Resource not found")
	}
	return u.GetTagByID(id)
}

func (u *TagDB) DeleteTag(id uuid.UUID) error {
	dbResponse := u.db.Delete(&models.Tag{}, "id = ?", id)
	if dbResponse.Error != nil {
		_, err := utils.HandleDBError(&models.Tag{}, dbResponse.Error)
		return err
	}
	if dbResponse.RowsAffected == 0 {
		return huma.Error404NotFound("Resource not found")
	}
	return nil
}
