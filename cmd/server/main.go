package main

import (
	"log"
	"os"
	"path/filepath"

	"blog-service/internal/config"
	"blog-service/internal/db"
	"blog-service/internal/router"
)

func main() {
	cfg := config.Load()

	if cfg.UploadDir != "" {
		if err := os.MkdirAll(filepath.Clean(cfg.UploadDir), 0o755); err != nil {
			log.Fatalf("create upload dir failed: %v", err)
		}
	}

	var pingDB func() error

	if cfg.MySQLDSN != "" {
		d, err := db.Open(cfg.MySQLDSN)
		if err != nil {
			log.Fatalf("mysql connect failed: %v", err)
		}
		pingDB = d.SQL.Ping
	} else {
		log.Println("MYSQL_DSN empty: running without database")
	}

	r := router.New(pingDB)

	log.Printf("server listening on %s", cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}
