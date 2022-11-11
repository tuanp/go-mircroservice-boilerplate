package db

import (
	"time"

	"github.com/tuanp/go-gin-boilerplate/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	zapgorm "moul.io/zapgorm2"
)

const (
	maxDBIdleConns  = 10
	maxDBOpenConns  = 100
	maxConnLifeTime = 30 * time.Minute
)

func ConnectMySQL(cfg *config.MysqlConfig, zapLogger *zap.Logger) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: zapgorm.New(zapLogger).LogMode(logger.Silent),
	})
	if err != nil {
		zapLogger.Sugar().Fatalf("Error open mysql: %v", err)
	}

	err = db.Raw("SELECT 1").Error
	if err != nil {
		zapLogger.Sugar().Fatal("Error querying SELECT 1", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		zapLogger.Sugar().Fatal("Error get sql DB", err)
	}
	sqlDB.SetMaxIdleConns(maxDBIdleConns)
	sqlDB.SetMaxOpenConns(maxDBOpenConns)
	sqlDB.SetConnMaxLifetime(maxConnLifeTime)
	return db
}
