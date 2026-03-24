package models

import (
	"time"

	"github.com/google/uuid"
)

type UserTagSubscription struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:idx_user_tag_subscriptions_user_tag"`
	User      User      `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	TagID     uuid.UUID `json:"tag_id" gorm:"type:uuid;not null;uniqueIndex:idx_user_tag_subscriptions_user_tag"`
	Tag       Tag       `json:"-" gorm:"foreignKey:TagID;references:ID;constraint:OnDelete:CASCADE"`
}
