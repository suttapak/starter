package logger

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	loggerMock struct {
		*mock.Mock
		logger *zap.Logger
	}
)

// Debug implements AppLogger.
func (l *loggerMock) Debug(message string, fields ...zap.Field) {
	l.logger.Debug(message, fields...)
}

// Error implements AppLogger.
func (l *loggerMock) Error(message interface{}, fields ...zap.Field) {
	switch v := message.(type) {
	case error:
		l.logger.Error(v.Error(), fields...)
	case string:
		l.logger.Error(v, fields...)
	}
}

// Info implements AppLogger.
func (l *loggerMock) Info(message string, fields ...zap.Field) {
	l.logger.Info(message, fields...)
}

func NewLoggerMock() AppLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""

	var err error
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return &loggerMock{logger: logger}
}
