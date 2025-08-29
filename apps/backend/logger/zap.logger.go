package logger

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type AppLogger interface {
	Info(message string, fields ...zap.Field)
	Debug(message string, fields ...zap.Field)
	Error(message interface{}, fields ...zap.Field)
}

type appLogger struct {
	logger *zap.Logger
}

// ensureDir checks if a directory exists and creates it if not.
func ensureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755) // Creates parent directories if needed
		if err != nil {
			return err
		}
	}
	return nil
}

func newAppLogger() (AppLogger, error) {
	// Ensure logs directory exists
	if err := ensureDir("./logs"); err != nil {
		log.Fatalf("Failed to create log directory: %s, error: %v", "logs", err)
	}

	// Lumberjack for log rotation
	logWriter := &lumberjack.Logger{
		Filename: "./logs/app.log",
		MaxAge:   90, // keep logs for 90 days
		// MaxSize, MaxBackups not set (ignored)
	}

	// Zap encoder config
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // structured JSON logs
		zapcore.AddSync(logWriter),
		zap.DebugLevel, // log level
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &appLogger{logger: logger}, nil
}
func (a *appLogger) Info(message string, fields ...zap.Field) {
	a.logger.Info(message, fields...)
}

func (a *appLogger) Debug(message string, fields ...zap.Field) {
	a.logger.Debug(message, fields...)
}

func (a *appLogger) Error(message interface{}, fields ...zap.Field) {
	switch v := message.(type) {
	case error:
		a.logger.Error(v.Error(), fields...)
	case string:
		a.logger.Error(v, fields...)
	}
}

// GinLogger returns a gin.HandlerFunc that logs requests via zap
func GinLogger(appLogger AppLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log details
		latency := time.Since(start)
		status := c.Writer.Status()

		appLogger.Info("HTTP Request",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		)
	}
}
