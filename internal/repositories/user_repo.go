package repositories

import (
	"errors"
	"time"

	"blog-service/internal/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

// 创建用户仓库
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

// 创建用户
func (r *UserRepo) Create(u *models.User) error {
	return r.DB.Create(u).Error
}

// 根据ID查找用户
func (r *UserRepo) FindByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// 根据邮箱或用户名查找用户
func (r *UserRepo) FindByEmailOrUsername(s string) (*models.User, error) {
	var u models.User
	err := r.DB.Where("email = ? OR username = ?", s, s).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// 用于注册时检查唯一性
func (r *UserRepo) ExistsEmail(email string) (bool, error) {
	var cnt int64
	if err := r.DB.Model(&models.User{}).Where("email = ?", email).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// 检查用户名是否存在
func (r *UserRepo) ExistsUsername(username string) (bool, error) {
	var cnt int64
	if err := r.DB.Model(&models.User{}).Where("username = ?", username).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// 检查用户是否存在
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// 更新最后登录时间
func (r *UserRepo) UpdateLastLoginAt(userID uint, t time.Time) error {
	return r.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", t).Error
}
