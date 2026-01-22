package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	UserID            uint           `json:"user_id" gorm:"not null"`
	User              User           `json:"user" gorm:"foreignKey:UserID"`
	StripePaymentID   string         `json:"stripe_payment_id" gorm:"uniqueIndex;not null"`
	Amount            int64          `json:"amount" gorm:"not null"` // Amount in cents
	Currency          string         `json:"currency" gorm:"not null;default:usd"`
	Status            string         `json:"status" gorm:"not null"` // pending, succeeded, failed, canceled
	Description       string         `json:"description"`
	Metadata          string         `json:"metadata" gorm:"type:jsonb"` // Additional data as JSON
	PaymentMethodID   string         `json:"payment_method_id"`
	ClientSecret      string         `json:"client_secret"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

type PaymentIntent struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	UserID            uint           `json:"user_id" gorm:"not null"`
	User              User           `json:"user" gorm:"foreignKey:UserID"`
	StripeIntentID    string         `json:"stripe_intent_id" gorm:"uniqueIndex;not null"`
	Amount            int64          `json:"amount" gorm:"not null"` // Amount in cents
	Currency          string         `json:"currency" gorm:"not null;default:usd"`
	Status            string         `json:"status" gorm:"not null"` // requires_payment_method, requires_confirmation, requires_action, processing, requires_capture, canceled, succeeded
	Description       string         `json:"description"`
	ClientSecret      string         `json:"client_secret"`
	PaymentMethodID   string         `json:"payment_method_id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

type Customer struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	User            User           `json:"user" gorm:"foreignKey:UserID"`
	StripeCustomerID string        `json:"stripe_customer_id" gorm:"uniqueIndex;not null"`
	Email           string         `json:"email" gorm:"not null"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// Request/Response DTOs
type CreatePaymentIntentRequest struct {
	Amount      int64  `json:"amount" validate:"required,min=50"` // Minimum $0.50
	Currency    string `json:"currency" validate:"required,len=3"`
	Description string `json:"description"`
	Metadata    string `json:"metadata,omitempty"`
}

type CreatePaymentIntentResponse struct {
	ClientSecret string `json:"client_secret"`
	IntentID     string `json:"intent_id"`
}

type ConfirmPaymentRequest struct {
	PaymentIntentID string `json:"payment_intent_id" validate:"required"`
	PaymentMethodID string `json:"payment_method_id" validate:"required"`
}

type PaymentResponse struct {
	ID              uint      `json:"id"`
	StripePaymentID string    `json:"stripe_payment_id"`
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}

type WebhookEvent struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	EventID   string         `json:"event_id" gorm:"uniqueIndex;not null"`
	EventType string         `json:"event_type" gorm:"not null"`
	Processed bool           `json:"processed" gorm:"default:false"`
	Data      string         `json:"data" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
