package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Post represents a post entity in the system
type Post struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AuthorID    uuid.UUID      `json:"author_id" type:"uuid"`
	Author      User           `json:"-" gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:CASCADE"`
	SportID     *uuid.UUID     `json:"sport_id" gorm:"type:uuid;default:null"`
	Sport       *Sport         `json:"-" gorm:"foreignKey:SportID;references:ID;constraint:OnDelete:SET NULL;"`
	CollegeID   *uuid.UUID     `json:"college_id" gorm:"type:uuid;default:null"`
	College     *College       `json:"-" gorm:"foreignKey:CollegeID;references:ID;constraint:OnDelete:SET NULL;"`
	Tags        []Tag          `json:"tags" gorm:"many2many:tag_posts;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Title       string         `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content     string         `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`
	IsAnonymous bool           `json:"isAnonymous" gorm:"default:false"`

	// only used for db queries -> ignored for migrations
	LikeCount       int64   `json:"like_count" gorm:"column:like_count;->;-:migration"`
	CommentCount    int64   `json:"comment_count" gorm:"column:comment_count;->;-:migration"`
	IsLiked         bool    `json:"is_liked" gorm:"column:is_liked;->;-:migration"`
	PopularityScore float64 `json:"popularity_score" gorm:"column:popularity_score;->;-:migration"`
}
