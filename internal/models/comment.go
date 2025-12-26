package models

import "time"

type CommentStatus string

const (
	CommentPending  CommentStatus = "pending"
	CommentApproved CommentStatus = "approved"
)

type Comment struct {
	ID uint `gorm:"primaryKey"`

	PostID uint `gorm:"not null;index"`
	Post   Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	AuthorID uint `gorm:"not null;index"`
	Author   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	ParentID *uint         `gorm:"index"`
	Content  string        `gorm:"type:text;not null"`
	Status   CommentStatus `gorm:"type:enum('pending','approved');not null;default:'approved';index"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
