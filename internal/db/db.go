package db

import (
	"database/sql"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	Gorm *gorm.DB
	SQL  *sql.DB
}

func Open(mysqlDSN string) (*DB, error) {
	gdb, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{
		// 后面可以按需开启：Logger、PrepareStmt 等
	})
	if err != nil {
		return nil, err
	}

	sqldb, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	// 连接池参数：生产可再调整
	sqldb.SetMaxOpenConns(50)
	sqldb.SetMaxIdleConns(25)
	sqldb.SetConnMaxLifetime(30 * time.Minute)

	// 启动时 ping 一次，尽早失败
	if err := sqldb.Ping(); err != nil {
		return nil, err
	}

	log.Println("mysql connected")
	return &DB{Gorm: gdb, SQL: sqldb}, nil
}
