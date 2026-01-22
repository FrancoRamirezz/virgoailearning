package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Content     string         `json:"content" gorm:"type:text"`
	AuthorID    uint           `json:"author_id" gorm:"not null"`
	Author      User           `json:"author" gorm:"foreignKey:AuthorID"`
	IsPublished bool           `json:"is_published" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreatePostRequest struct {
	Title       string `json:"title" validate:"required"`
	Content     string `json:"content" validate:"required"`
	IsPublished bool   `json:"is_published"`
}

type UpdatePostRequest struct {
	Title       *string `json:"title,omitempty"`
	Content     *string `json:"content,omitempty"`
	IsPublished *bool   `json:"is_published,omitempty"`
}
