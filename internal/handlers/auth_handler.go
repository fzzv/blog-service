package handlers

import (
	"net/http"

	"blog-service/internal/middleware"
	"blog-service/internal/repositories"
	"blog-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	Auth  *services.AuthService
	Users *repositories.UserRepo
	V     *validator.Validate
}

type registerReq struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (h AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	if err := h.V.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error"})
		return
	}

	u, err := h.Auth.Register(services.RegisterInput(req))
	if err != nil {
		switch err {
		case services.ErrEmailTaken:
			c.JSON(http.StatusConflict, gin.H{"error": "email_taken"})
		case services.ErrUsernameTaken:
			c.JSON(http.StatusConflict, gin.H{"error": "username_taken"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":       u.ID,
			"email":    u.Email,
			"username": u.Username,
			"role":     u.Role,
		},
	})
}

type loginReq struct {
	EmailOrUsername string `json:"email_or_username" validate:"required,max=255"`
	Password        string `json:"password" validate:"required,max=72"`
}

func (h AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	if err := h.V.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error"})
		return
	}

	token, u, err := h.Auth.Login(req.EmailOrUsername, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"user": gin.H{
			"id":       u.ID,
			"email":    u.Email,
			"username": u.Username,
			"role":     u.Role,
		},
	})
}

func (h AuthHandler) Me(c *gin.Context) {
	uid, ok := middleware.GetAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	u, err := h.Users.FindByID(uid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       u.ID,
			"email":    u.Email,
			"username": u.Username,
			"role":     u.Role,
		},
	})
}
