package models

import (
	"time"

	"gorm.io/gorm"
)

// Course represents an educational course
type Course struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title" gorm:"not null"`
	Description  string         `json:"description" gorm:"type:text"`
	InstructorID uint           `json:"instructor_id" gorm:"not null"`
	Instructor   User           `json:"instructor" gorm:"foreignKey:InstructorID"`
	Price        int64          `json:"price" gorm:"not null;default:0"` // Price in cents
	Currency     string         `json:"currency" gorm:"default:usd"`
	Category     string         `json:"category"`
	Level        string         `json:"level"`    // beginner, intermediate, advanced
	Duration     int            `json:"duration"` // Duration in hours
	IsPublished  bool           `json:"is_published" gorm:"default:false"`
	ThumbnailURL string         `json:"thumbnail_url"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Contents    []CourseContent `json:"contents,omitempty" gorm:"foreignKey:CourseID"`
	Enrollments []Enrollment    `json:"enrollments,omitempty" gorm:"foreignKey:CourseID"`
}

// CourseContent represents individual lessons/modules within a course
type CourseContent struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CourseID    uint           `json:"course_id" gorm:"not null"`
	Course      Course         `json:"course,omitempty" gorm:"foreignKey:CourseID"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Content     string         `json:"content" gorm:"type:text"`         // Could be markdown, HTML, or text
	ContentType string         `json:"content_type" gorm:"default:text"` // text, video, quiz, assignment
	Order       int            `json:"order" gorm:"not null"`            // Order within the course
	Duration    int            `json:"duration"`                         // Duration in minutes
	IsPublished bool           `json:"is_published" gorm:"default:false"`
	VideoURL    string         `json:"video_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Request/Response DTOs
type CreateCourseRequest struct {
	Title        string `json:"title" validate:"required"`
	Description  string `json:"description"`
	Price        int64  `json:"price" validate:"min=0"`
	Currency     string `json:"currency" validate:"len=3"`
	Category     string `json:"category"`
	Level        string `json:"level"` // beginner, intermediate, advanced
	Duration     int    `json:"duration" validate:"min=0"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type UpdateCourseRequest struct {
	Title        *string `json:"title,omitempty"`
	Description  *string `json:"description,omitempty"`
	Price        *int64  `json:"price,omitempty"`
	Currency     *string `json:"currency,omitempty"`
	Category     *string `json:"category,omitempty"`
	Level        *string `json:"level,omitempty"`
	Duration     *int    `json:"duration,omitempty"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
	IsPublished  *bool   `json:"is_published,omitempty"`
}

type CourseResponse struct {
	ID           uint         `json:"id"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	InstructorID uint         `json:"instructor_id"`
	Instructor   UserResponse `json:"instructor"`
	Price        int64        `json:"price"`
	Currency     string       `json:"currency"`
	Category     string       `json:"category"`
	Level        string       `json:"level"`
	Duration     int          `json:"duration"`
	IsPublished  bool         `json:"is_published"`
	ThumbnailURL string       `json:"thumbnail_url"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type CreateCourseContentRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"` // text, video, quiz, assignment
	Order       int    `json:"order" validate:"min=0"`
	Duration    int    `json:"duration" validate:"min=0"`
	VideoURL    string `json:"video_url"`
}

type UpdateCourseContentRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Content     *string `json:"content,omitempty"`
	ContentType *string `json:"content_type,omitempty"`
	Order       *int    `json:"order,omitempty"`
	Duration    *int    `json:"duration,omitempty"`
	VideoURL    *string `json:"video_url,omitempty"`
	IsPublished *bool   `json:"is_published,omitempty"`
}

type CourseContentResponse struct {
	ID          uint      `json:"id"`
	CourseID    uint      `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	ContentType string    `json:"content_type"`
	Order       int       `json:"order"`
	Duration    int       `json:"duration"`
	IsPublished bool      `json:"is_published"`
	VideoURL    string    `json:"video_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListCoursesRequest struct {
	Category     string `json:"category"`
	Level        string `json:"level"`
	InstructorID uint   `json:"instructor_id"`
	IsPublished  *bool  `json:"is_published"`
	Page         int    `json:"page" validate:"min=1"`
	Limit        int    `json:"limit" validate:"min=1,max=100"`
}


