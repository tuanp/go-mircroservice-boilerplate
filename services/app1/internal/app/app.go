package app

import (
	"github.com/tuanp/go-gin-boilerplate/pkg/db"
	"log"

	"github.com/tuanp/go-gin-boilerplate/config"
	"github.com/tuanp/go-gin-boilerplate/pkg/logger"
)

// Run initializes whole application.
func Run(configPath string) {
	log.Println("Starting api server")
	cfg, err := config.Init(configPath)

	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	psqlDB, err := db.ConnectMySQL()
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	appLogger.Info("Redis connected")

}
