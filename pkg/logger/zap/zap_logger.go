package zap

import (
	"os"
	"time"

	"github.com/tuanp/go-gin-boilerplate/pkg/config"
	"github.com/tuanp/go-gin-boilerplate/pkg/constants"
	"github.com/tuanp/go-gin-boilerplate/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	level       string
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
}

type ZapLogger interface {
	logger.Logger
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Sync() error
}

// For mapping config logger
var loggerLevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}

// NewZapLogger create new zap logger
func NewZapLogger(lc *config.LoggerConfig, sc *config.ServerConfig) ZapLogger {
	zapLogger := &zapLogger{level: lc.LogLevel}
	zapLogger.initLogger(sc)

	return zapLogger
}

func (l *zapLogger) getLoggerLevel() zapcore.Level {
	level, exist := loggerLevelMap[l.level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// InitLogger Init logger
func (l *zapLogger) initLogger(cfg *config.ServerConfig) {
	logLevel := l.getLoggerLevel()

	logWriter := zapcore.AddSync(os.Stdout)

	var encoderCfg zapcore.EncoderConfig
	var encoder zapcore.Encoder

	if cfg.Mode == constants.EnvProd {
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.NameKey = "[SERVICE]"
		encoderCfg.TimeKey = "[TIME]"
		encoderCfg.LevelKey = "[LEVEL]"
		encoderCfg.FunctionKey = "[CALLER]"
		encoderCfg.CallerKey = "[LINE]"
		encoderCfg.MessageKey = "[MESSAGE]"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
		encoderCfg.EncodeName = zapcore.FullNameEncoder
		encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		encoderCfg.NameKey = "[SERVICE]"
		encoderCfg.TimeKey = "[TIME]"
		encoderCfg.LevelKey = "[LEVEL]"
		encoderCfg.FunctionKey = "[CALLER]"
		encoderCfg.CallerKey = "[LINE]"
		encoderCfg.MessageKey = "[MESSAGE]"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeName = zapcore.FullNameEncoder
		encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeCaller = zapcore.FullCallerEncoder
		encoderCfg.ConsoleSeparator = " | "
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.logger = zapLogger
	l.sugarLogger = zapLogger.Sugar()
}

func (l *zapLogger) Configure(cfg func(internalLog interface{})) {
	cfg(l.logger)
}

func (l *zapLogger) LogType() config.LogType {
	return config.Zap
}

// WithName add logger microservice name
func (l *zapLogger) WithName(name string) {
	l.logger = l.logger.Named(name)
	l.sugarLogger = l.sugarLogger.Named(name)
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message
func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *zapLogger) Debugw(msg string, fields logger.Fields) {
	zapFields := mapToFields(fields)
	l.logger.Debug(msg, zapFields...)
}

// Info uses fmt.Sprint to construct and log a message
func (l *zapLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Infow logs a message with some additional context.
func (l *zapLogger) Infow(msg string, fields logger.Fields) {
	zapFields := mapToFields(fields)
	l.logger.Info(msg, zapFields...)
}

// Printf uses fmt.Sprintf to log a templated message
func (l *zapLogger) Printf(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

// WarnMsg log error message with warn level.
func (l *zapLogger) WarnMsg(msg string, err error) {
	l.logger.Warn(msg, zap.String("error", err.Error()))
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *zapLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

// Errorw logs a message with some additional context.
func (l *zapLogger) Errorw(msg string, fields logger.Fields) {
	zapFields := mapToFields(fields)
	l.logger.Error(msg, zapFields...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

// Err uses error to log a message.
func (l *zapLogger) Err(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *zapLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *zapLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *zapLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics
func (l *zapLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	go func() {
		err := l.logger.Sync()
		if err != nil {
			l.logger.Error("error while syncing", zap.Error(err))
		}
	}() // nolint: errcheck
	return l.sugarLogger.Sync()
}

func (l *zapLogger) GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error) {
	l.Info(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
		zap.Error(err),
	)
}

func (l *zapLogger) GrpcClientInterceptorLogger(method string, req, reply interface{}, time time.Duration, metaData map[string][]string, err error) {
	l.Info(
		constants.GRPC,
		zap.String(constants.METHOD, method),
		zap.Any(constants.REQUEST, req),
		zap.Any(constants.REPLY, reply),
		zap.Duration(constants.TIME, time),
		zap.Any(constants.METADATA, metaData),
		zap.Error(err),
	)
}

func mapToFields(fields map[string]interface{}) []zap.Field {
	var zapFields []zap.Field
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return zapFields
}
