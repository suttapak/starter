package logger

import (
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

}

// Error implements AppLogger.
func (l *loggerMock) Error(message interface{}, fields ...zap.Field) {

}

// Info implements AppLogger.
func (l *loggerMock) Info(message string, fields ...zap.Field) {

}

func NewLoggerMock() AppLogger {
	return &loggerMock{}
}
