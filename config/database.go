package config

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

func ConnectGormDB() *gorm.DB {
	once.Do(func() {
		LoadEnv()

		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			GetEnv("DB_USER", ""),
			GetEnv("DB_PASSWORD", ""),
			GetEnv("DB_HOST", ""),
			GetEnv("DB_PORT", ""),
			GetEnv("DB_NAME", ""),
		)

		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err != nil {
			panic(fmt.Sprintf("Gagal koneksi database: %v", err))
		}

		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(50)
		sqlDB.SetMaxOpenConns(200)
		sqlDB.SetConnMaxLifetime(1 * time.Hour)

		dbInstance = db
	})

	return dbInstance
}
