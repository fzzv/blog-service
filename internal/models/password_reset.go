package models

import "time"

type PasswordResetToken struct {
	ID uint `gorm:"primaryKey"`

	UserID    uint      `gorm:"not null;index"`
	TokenHash string    `gorm:"size:255;not null;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	UsedAt    *time.Time

	CreatedAt time.Time
}
