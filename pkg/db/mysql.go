package db

import (
	"github.com/tuanp/go-gin-boilerplate/pkg/logger"
	"time"

	"github.com/tuanp/go-gin-boilerplate/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	maxDBIdleConns  = 10
	maxDBOpenConns  = 100
	maxConnLifeTime = 30 * time.Minute
)

func ConnectMySQL(cfg *config.MysqlConfig, logger logger.Logger) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Error open mysql: %v", err)
	}

	err = db.Raw("SELECT 1").Error
	if err != nil {
		logger.Fatalf("Error querying SELECT 1", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("Error get sql DB", err)
	}
	sqlDB.SetMaxIdleConns(maxDBIdleConns)
	sqlDB.SetMaxOpenConns(maxDBOpenConns)
	sqlDB.SetConnMaxLifetime(maxConnLifeTime)
	return db
}
