package models

import (
	"time"

	"gorm.io/gorm"
)

// Enrollment represents a user's enrollment in a course
type Enrollment struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	CourseID    uint           `json:"course_id" gorm:"not null"`
	Course      Course         `json:"course" gorm:"foreignKey:CourseID"`
	Status      string         `json:"status" gorm:"default:active"` // active, completed, cancelled
	Progress    float64        `json:"progress" gorm:"default:0"`    // Progress percentage (0-100)
	EnrolledAt  time.Time      `json:"enrolled_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Request/Response DTOs
type CreateEnrollmentRequest struct {
	CourseID uint `json:"course_id" validate:"required"`
}

type UpdateEnrollmentRequest struct {
	Status   *string  `json:"status,omitempty"` // active, completed, cancelled
	Progress *float64 `json:"progress,omitempty" validate:"min=0,max=100"`
}

type EnrollmentResponse struct {
	ID          uint           `json:"id"`
	UserID      uint           `json:"user_id"`
	User        UserResponse   `json:"user"`
	CourseID    uint           `json:"course_id"`
	Course      CourseResponse `json:"course"`
	Status      string         `json:"status"`
	Progress    float64        `json:"progress"`
	EnrolledAt  time.Time      `json:"enrolled_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}


