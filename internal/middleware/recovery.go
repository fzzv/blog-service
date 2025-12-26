package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoveryJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 统一 JSON 错误返回（避免泄露内部细节）
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal_server_error",
				})
			}
		}()
		c.Next()
	}
}
