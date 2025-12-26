package models

import "time"

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID           uint       `gorm:"primaryKey"`
	Email        string     `gorm:"size:255;not null;uniqueIndex"`
	Username     string     `gorm:"size:50;not null;uniqueIndex"`
	PasswordHash string     `gorm:"size:255;not null"`
	Role         UserRole   `gorm:"type:enum('admin','user');not null;default:'user';index"`
	LastLoginAt  *time.Time `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
