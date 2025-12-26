package db

import (
	"fmt"

	"blog-service/internal/models"

	"gorm.io/gorm"
)

func Migrate(gdb *gorm.DB) error {
	// 1) 基础表
	if err := gdb.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Tag{},
		&models.Comment{},
		&models.PostLike{},
		&models.PasswordResetToken{},
	); err != nil {
		return err
	}

	// 2) post_likes 联合唯一
	if err := gdb.Exec(`
		ALTER TABLE post_likes
		ADD UNIQUE KEY uk_post_likes_post_user (post_id, user_id)
	`).Error; err != nil {
		// 已存在时会报错：忽略更稳妥（MySQL 没 IF NOT EXISTS for constraint in all versions）
		// 这里简单做：如果是重复键错误可以忽略；为了保持简洁先不做复杂错误判断
	}

	// 3) posts FULLTEXT（MySQL 8 InnoDB 支持）
	// 注意：FULLTEXT 不支持 IF NOT EXISTS，重复会报错，和上面同理可以容忍失败
	if err := gdb.Exec(`
		ALTER TABLE posts
		ADD FULLTEXT INDEX ft_posts_title_md (title, content_md)
	`).Error; err != nil {
	}

	return nil
}

func EnsureSchema(gdb *gorm.DB) {
	if err := Migrate(gdb); err != nil {
		panic(fmt.Errorf("migrate failed: %w", err))
	}
}
