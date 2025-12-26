package router

import (
	"net/http"
	"time"

	"blog-service/internal/handlers"
	"blog-service/internal/middleware"
	"blog-service/internal/repositories"
	"blog-service/internal/services"
	jwtutil "blog-service/internal/utils/jwt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Deps struct {
	DB     *gorm.DB
	PingDB func() error

	JWTSecret string
}

func New(d Deps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.RecoveryJSON())

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
	})

	// healthz
	hh := handlers.HealthHandler{PingDB: d.PingDB}
	r.GET("/healthz", hh.Healthz)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}

	// 只有 DB 存在时才注册需要 DB 的路由
	if d.DB != nil {
		v := validator.New()

		userRepo := repositories.NewUserRepo(d.DB)

		jm := jwtutil.Manager{
			Secret: []byte(d.JWTSecret),
			Issuer: "blog-service",
			TTL:    2 * time.Hour,
		}

		authSvc := &services.AuthService{
			Users: userRepo,
			JWT:   jm,
		}
		authMW := middleware.NewAuthMiddleware(jm)

		authHandler := handlers.AuthHandler{
			Auth:  authSvc,
			Users: userRepo,
			V:     v,
		}

		av1 := r.Group("/api/v1/auth")
		{
			av1.POST("/register", authHandler.Register)
			av1.POST("/login", authHandler.Login)
			av1.GET("/me", authMW.AuthRequired(), authHandler.Me)
		}

		// 示例：管理员保护路由（后续发文章就用这个）
		admin := r.Group("/api/v1/admin")
		admin.Use(authMW.AuthRequired(), middleware.RequireAdmin())
		{
			admin.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "admin pong"})
			})
		}
	}

	return r
}
