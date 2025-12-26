package models

import "time"

type PostStatus string

const (
	PostDraft     PostStatus = "draft"
	PostPublished PostStatus = "published"
)

type Post struct {
	ID uint `gorm:"primaryKey"`

	Title       string     `gorm:"size:200;not null;index"`
	Slug        string     `gorm:"size:220;not null;uniqueIndex"`
	ContentMD   string     `gorm:"type:longtext;not null"`
	ContentHTML string     `gorm:"type:longtext;not null"`
	Status      PostStatus `gorm:"type:enum('draft','published');not null;default:'draft';index"`

	PublishedAt *time.Time `gorm:"index"`

	AuthorID uint `gorm:"not null;index"`
	Author   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	ViewCount    uint64 `gorm:"not null;default:0"`
	LikeCount    uint64 `gorm:"not null;default:0"`
	CommentCount uint64 `gorm:"not null;default:0"`

	Tags []Tag `gorm:"many2many:post_tags;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
