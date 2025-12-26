package middleware

import (
	"strings"

	jwtutil "blog-service/internal/utils/jwt"

	"github.com/gin-gonic/gin"
)

func NewOptionalAuth(m jwtutil.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			c.Next()
			return
		}
		token := strings.TrimSpace(auth[len("Bearer "):])
		claims, err := m.Parse(token)
		if err != nil {
			c.Next()
			return
		}
		c.Set(ctxUserIDKey, claims.UserID)
		c.Set(ctxRoleKey, claims.Role)
		c.Next()
	}
}
