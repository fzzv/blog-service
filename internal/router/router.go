package router

import (
	"net/http"

	"blog-service/internal/handlers"
	"blog-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New(pingDB func() error) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.RecoveryJSON())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
	})

	hh := handlers.HealthHandler{PingDB: pingDB}
	r.GET("/healthz", hh.Healthz)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}

	return r
}
