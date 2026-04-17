package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PremiumPost represents a post that is only available to paid users & moderators/admins in the system
type PremiumPost struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	AuthorID  uuid.UUID  `json:"author_id" type:"uuid"`
	Author    User       `json:"-" gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:CASCADE"`
	SportID   *uuid.UUID `json:"sport_id" gorm:"type:uuid;default:null"`
	Sport     *Sport     `json:"-" gorm:"foreignKey:SportID;references:ID;constraint:OnDelete:SET NULL;"`
	CollegeID *uuid.UUID `json:"college_id" gorm:"type:uuid;default:null"`
	College   *College   `json:"-" gorm:"foreignKey:CollegeID;references:ID;constraint:OnDelete:SET NULL;"`

	Tags []Tag `json:"tags" gorm:"many2many:tag_posts;foreignKey:ID;joinForeignKey:PostableID;joinReferences:TagID"`

	Title   string `json:"title" example:"Looking for thoughts on NEU Fencing!" gorm:"type:varchar(100);not null" validate:"required,min=1,max=100"`
	Content string `json:"content" example:"My name is Bob Joe and I am a rising senior who just got into NEU. What is the fencing program like? Are they competitive?" gorm:"type:varchar(5000);not null" validate:"required,min=1,max=5000"`

	MediaID *uuid.UUID `json:"media_id,omitempty" gorm:"type:uuid;default:null"`
	Media   *Media     `json:"media,omitempty" gorm:"foreignKey:MediaID;references:ID;constraint:OnDelete:SET NULL"`
}
