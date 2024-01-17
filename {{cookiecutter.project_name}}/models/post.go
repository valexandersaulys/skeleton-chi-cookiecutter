package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Uuid     string `gorm:"type:CHAR(35);column:uuid;index;not null"`
	Title    string
	Content  string `gorm:"type:text"`
	IsPublic bool   `gorm:"default:false"`
	AuthorID uint   `gorm:"not null"`
	Author   User
}

// Instantiate uuid column in Post
func (post *Post) BeforeCreate(tx *gorm.DB) (err error) {
	post.Uuid = uuid.New().String()
	return
}
