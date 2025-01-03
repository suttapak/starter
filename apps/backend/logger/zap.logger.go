package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AppLogger interface {
	Info(message string, fields ...zap.Field)
	Debug(message string, fields ...zap.Field)
	Error(message interface{}, fields ...zap.Field)
}

type appLogger struct {
	logger *zap.Logger
}

func newAppLogger() (AppLogger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &appLogger{logger: logger}, nil
}
func (a appLogger) Info(message string, fields ...zap.Field) {
	a.logger.Info(message, fields...)
}

func (a appLogger) Debug(message string, fields ...zap.Field) {
	a.logger.Debug(message, fields...)
}

func (a appLogger) Error(message interface{}, fields ...zap.Field) {
	switch v := message.(type) {
	case error:
		a.logger.Error(v.Error(), fields...)
	case string:
		a.logger.Error(v, fields...)
	}
}
