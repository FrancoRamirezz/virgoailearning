package models

import (
	"time"

	"gorm.io/gorm"
)

// Contact represents a contact entry (could be for support, inquiries, etc.)
type Contact struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    *uint          `json:"user_id"` // Optional: contact may be from logged-in user or anonymous
	User      *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"not null"`
	Phone     string         `json:"phone"`
	Subject   string         `json:"subject" gorm:"not null"`
	Message   string         `json:"message" gorm:"type:text;not null"`
	Status    string         `json:"status" gorm:"default:pending"` // pending, responded, resolved, archived
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Request/Response DTOs
type CreateContactRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone"`
	Subject string `json:"subject" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type UpdateContactRequest struct {
	Name    *string `json:"name,omitempty"`
	Email   *string `json:"email,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Subject *string `json:"subject,omitempty"`
	Message *string `json:"message,omitempty"`
	Status  *string `json:"status,omitempty"` // pending, responded, resolved, archived
}

type ContactResponse struct {
	ID        uint      `json:"id"`
	UserID    *uint     `json:"user_id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SearchContactsRequest struct {
	Query  string `json:"query"`  // Search by name or email
	Status string `json:"status"` // Filter by status
	Page   int    `json:"page" validate:"min=1"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
}


