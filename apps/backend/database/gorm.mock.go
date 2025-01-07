package database

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type gormMock struct {
	mock.Mock
	gorm.DB
}

func NewGormMock() *gormMock {
	return &gormMock{}
}

// Mock Rollback
func (m *gormMock) Rollback() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// Mock Commit
func (m *gormMock) Commit() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

// Mock Other Methods
func (m *gormMock) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}

func (m *gormMock) First(dest interface{}, conds ...interface{}) *gorm.DB {
	callArgs := m.Called(dest, conds)
	return callArgs.Get(0).(*gorm.DB)
}
