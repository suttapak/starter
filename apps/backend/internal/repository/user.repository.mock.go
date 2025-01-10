package repository

import (
	"context"

	"github.com/suttapak/starter/internal/model"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type (
	userMock struct {
		mock.Mock
	}
)

// CommitTx implements User.
func (u *userMock) CommitTx(tx *gorm.DB) {
	u.Called()
}

// RollbackTx implements User.
func (u *userMock) RollbackTx(tx *gorm.DB) {
	u.Called()
}

// GetUserByUserId implements User.
func (u *userMock) GetUserByUserId(ctx context.Context, tx *gorm.DB, uId uint) (user *model.User, err error) {
	args := u.Called()
	return args.Get(0).(*model.User), args.Error(1)
}

func (u *userMock) VerifyEmail(ctx context.Context, tx *gorm.DB, email string) (user *model.User, err error) {
	args := u.Called()
	return args.Get(0).(*model.User), args.Error(1)
}

// GetUserByEmailOrUsername implements User.
func (u *userMock) GetUserByEmailOrUsername(ctx context.Context, tx *gorm.DB, emailOrUsername string) (user *model.User, err error) {
	args := u.Called()
	return args.Get(0).(*model.User), args.Error(1)
}

// CheckEmail implements User.
func (u *userMock) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (user *model.User, flag bool, err error) {
	args := u.Called()
	return args.Get(0).(*model.User), args.Bool(1), args.Error(2)
}

// BeginTx implements User.
func (u *userMock) BeginTx() *gorm.DB {
	// return nil
	return nil
}

// CheckUsername implements User.
func (u *userMock) CheckUsername(ctx context.Context, tx *gorm.DB, username string) (user *model.User, flag bool, err error) {
	args := u.Called()
	return args.Get(0).(*model.User), args.Bool(1), args.Error(2)
}

// Register implements User.
func (u *userMock) Register(ctx context.Context, tx *gorm.DB, user model.User) (*model.User, error) {
	args := u.Called()
	return args.Get(0).(*model.User), args.Error(1)
}
func NewUserRepositoryMock() *userMock {
	return &userMock{}
}
