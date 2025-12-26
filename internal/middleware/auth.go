package middleware

import (
	"net/http"
	"strings"

	jwtutil "blog-service/internal/utils/jwt"

	"github.com/gin-gonic/gin"
)

const (
	ctxUserIDKey = "auth.user_id"
	ctxRoleKey   = "auth.role"
)

type AuthMiddleware struct {
	JWT jwtutil.Manager
}

func NewAuthMiddleware(m jwtutil.Manager) AuthMiddleware {
	return AuthMiddleware{JWT: m}
}

func (a AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		token := strings.TrimSpace(auth[len("Bearer "):])

		claims, err := a.JWT.Parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set(ctxUserIDKey, claims.UserID)
		c.Set(ctxRoleKey, claims.Role)
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ctxRoleKey)
		if roleStr, ok := role.(string); !ok || roleStr != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

func GetAuthUserID(c *gin.Context) (uint, bool) {
	v, ok := c.Get(ctxUserIDKey)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}

func GetAuthRole(c *gin.Context) (string, bool) {
	v, ok := c.Get(ctxRoleKey)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
