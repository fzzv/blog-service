package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr string

	// 先预留：后续第二步接 MySQL 时会用到
	MySQLDSN string

	JWTSecret string
	UploadDir string

	CommentModeration bool
}

func Load() Config {
	// 本地开发优先加载 .env（没有也不报错）
	_ = godotenv.Load()

	return Config{
		Addr:              getEnv("APP_ADDR", ":8080"),
		MySQLDSN:          getEnv("MYSQL_DSN", ""),
		JWTSecret:         getEnv("JWT_SECRET", "dev-secret-change-me"),
		UploadDir:         getEnv("UPLOAD_DIR", "./uploads"),
		CommentModeration: getEnvBool("COMMENT_MODERATION", false),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
