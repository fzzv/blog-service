package services

import (
	"errors"
	"strings"
	"time"

	"blog-service/internal/models"
	"blog-service/internal/repositories"
	jwtutil "blog-service/internal/utils/jwt"
	"blog-service/internal/utils/password"
)

var (
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrEmailTaken         = errors.New("email_taken")
	ErrUsernameTaken      = errors.New("username_taken")
)

type AuthService struct {
	Users *repositories.UserRepo
	JWT   jwtutil.Manager
}

type RegisterInput struct {
	Email    string
	Username string
	Password string
}

/*
* 注册
* @param in RegisterInput 注册输入
* @return user 用户
* @return err 错误
 */
func (s *AuthService) Register(in RegisterInput) (*models.User, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	username := strings.TrimSpace(in.Username)

	// 检查邮箱是否已经注册过
	if ok, err := s.Users.ExistsEmail(email); err != nil {
		return nil, err
	} else if ok {
		return nil, ErrEmailTaken
	}

	// 检查用户名是否已经注册过
	if ok, err := s.Users.ExistsUsername(username); err != nil {
		return nil, err
	} else if ok {
		return nil, ErrUsernameTaken
	}

	// 生成密码哈希
	hash, err := password.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: hash,
		Role:         models.RoleUser,
	}
	// 创建用户
	if err := s.Users.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

/*
* 登录
* @param emailOrUsername 邮箱或用户名
* @param plainPassword 密码
* @return token JWT token
* @return user 用户
* @return err 错误
 */
func (s *AuthService) Login(emailOrUsername, plainPassword string) (token string, user *models.User, err error) {
	u, err := s.Users.FindByEmailOrUsername(strings.TrimSpace(emailOrUsername))
	if err != nil {
		if repositories.IsNotFound(err) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	// 验证密码
	if !password.Verify(u.PasswordHash, plainPassword) {
		return "", nil, ErrInvalidCredentials
	}
	// 更新最后登录时间
	_ = s.Users.UpdateLastLoginAt(u.ID, time.Now())

	// 生成JWT token
	tok, err := s.JWT.Sign(u.ID, string(u.Role))
	if err != nil {
		return "", nil, err
	}
	return tok, u, nil
}
