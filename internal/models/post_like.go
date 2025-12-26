package models

import "time"

// 点赞表：用唯一约束保证一个用户对同一篇文章只能点一次
type PostLike struct {
	ID uint `gorm:"primaryKey"`

	PostID uint `gorm:"not null;index"`
	UserID uint `gorm:"not null;index"`

	CreatedAt time.Time

	// 联合唯一索引在迁移里加（更直观）
}
