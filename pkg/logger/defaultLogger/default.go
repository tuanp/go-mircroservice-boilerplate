package defaultLogger

import (
	"os"

	"github.com/tuanp/go-gin-boilerplate/config"
	"github.com/tuanp/go-gin-boilerplate/pkg/logger"
	"github.com/tuanp/go-gin-boilerplate/pkg/logger/logrous"
	"github.com/tuanp/go-gin-boilerplate/pkg/logger/zap"
)

var (
	Logger logger.Logger
)

func init() {
	logType := os.Getenv("LogConfig_LogType")
	if logType == "" {
		logType = "Zap"
	}
	switch logType {
	case "Zap":
		Logger = zap.NewZapLogger(&config.LoggerConfig{
			LogLevel: "debug",
			LogType:  config.Zap,
		}, &config.ServerConfig{
			Mode: "Development",
		})
		break
	case "Logrus":
		Logger = logrous.NewLogrusLogger(&config.LoggerConfig{
			LogLevel: "debug",
			LogType:  config.Zap,
		}, &config.ServerConfig{
			Mode: "Development",
		})
		break
	}
}
