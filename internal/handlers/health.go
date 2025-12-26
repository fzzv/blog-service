package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	PingDB func() error
}

func (h HealthHandler) Healthz(c *gin.Context) {
	if h.PingDB != nil {
		if err := h.PingDB(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "degraded",
				"db":     "down",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"db":     "up",
	})
}
