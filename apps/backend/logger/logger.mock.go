package logger

import (
	"log"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type (
	loggerMock struct {
		*mock.Mock
	}
)

// Debug implements AppLogger.
func (l *loggerMock) Debug(message string, fields ...zap.Field) {
	log.Println(message)
}

// Error implements AppLogger.
func (l *loggerMock) Error(message interface{}, fields ...zap.Field) {
	log.Println(message)
}

// Info implements AppLogger.
func (l *loggerMock) Info(message string, fields ...zap.Field) {
	log.Println(message)
}

func NewLoggerMock() AppLogger {
	return &loggerMock{}
}
