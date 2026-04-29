package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
	SubscriptionStatusPastDue  SubscriptionStatus = "past_due"
	SubscriptionStatusTrialing SubscriptionStatus = "trialing"
)

type UserSubscription struct {
	ID                   uuid.UUID          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
	DeletedAt            gorm.DeletedAt     `json:"deleted_at,omitempty" gorm:"index"`
	UserID               uuid.UUID          `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	User                 User               `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	StripeSubscriptionID string             `json:"stripe_subscription_id" gorm:"type:varchar(255);not null"`
	StripePriceID        string             `json:"stripe_price_id" gorm:"type:varchar(255);not null"`
	Status               SubscriptionStatus `json:"status" gorm:"type:varchar(50);not null"`
	CurrentPeriodStart   time.Time          `json:"current_period_start"`
	CurrentPeriodEnd     time.Time          `json:"current_period_end"`
	CanceledAt           *time.Time         `json:"canceled_at,omitempty"`
}
