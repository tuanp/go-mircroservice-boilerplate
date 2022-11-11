package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tuanp/go-gin-boilerplate/pkg/config"
	"github.com/tuanp/go-gin-boilerplate/pkg/db"
	"github.com/tuanp/go-gin-boilerplate/pkg/logger/zap"
	"github.com/tuanp/go-gin-boilerplate/services/app1/internal/handler"
	"github.com/tuanp/go-gin-boilerplate/services/app1/internal/repository"
	"github.com/tuanp/go-gin-boilerplate/services/app1/internal/server"
	"github.com/tuanp/go-gin-boilerplate/services/app1/internal/service"
)

// Run initializes whole application.
func Run(configPath string) {
	log.Println("Starting api server")
	cfg, err := config.Init(configPath)

	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	logger := zap.NewZapLogger(&cfg.Logger, &cfg.Server)
	logger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.LogLevel, cfg.Server.Mode)

	mysqlDB := db.ConnectMySQL(&cfg.Mysql, logger)
	if err != nil {
		logger.Fatalf("Postgresql init: %s", err)
	} else {
		logger.Infof("Mysql connected")
	}

	defer func() {
		dbInstance, _ := mysqlDB.DB()
		_ = dbInstance.Close()
		logger.Info("Mysql closed")
	}()

	redisClient := db.ConnectRedis(&cfg.Redis)
	defer func() {
		_ = redisClient.Close()
		logger.Info("Redis closed")
	}()
	logger.Info("Redis connected")

	// Services, Repos & API Handlers
	repos := repository.NewRepositories(mysqlDB)

	services := service.NewServices(service.Deps{
		Repos: repos,
		//Cache:                  memCache,
		CacheTTL:    int64(cfg.CacheTTL.Seconds()),
		Environment: cfg.Server.Mode,
		Domain:      cfg.HTTP.Host,
		Logger:      logger,
	})

	handlers := handler.NewHandler(services, logger)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")
	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 10 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
