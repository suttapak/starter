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

// GetUserByUserId implements User.
func (u *userMock) GetUserByUserId(ctx context.Context, tx *gorm.DB, uId uint) (user *model.User, err error) {
	args := u.Called(ctx, tx, uId)
	return args.Get(0).(*model.User), args.Error(1)
}

func (u *userMock) VerifyEmail(ctx context.Context, tx *gorm.DB, email string) (user *model.User, err error) {
	args := u.Called(ctx, tx, email)
	return args.Get(0).(*model.User), args.Error(1)
}

// GetUserByEmailOrUsername implements User.
func (u *userMock) GetUserByEmailOrUsername(ctx context.Context, tx *gorm.DB, emailOrUsername string) (user *model.User, err error) {
	//TODO implement me
	args := u.Called(ctx, tx, emailOrUsername)
	return args.Get(0).(*model.User), args.Error(1)
}

// CheckEmail implements User.
func (u *userMock) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (user *model.User, flag bool, err error) {
	args := u.Called(ctx, tx, email)
	return args.Get(0).(*model.User), args.Bool(1), args.Error(2)
}

// BeginTx implements User.
func (u *userMock) BeginTx() *gorm.DB {
	// return nil
	return nil
}

// CheckUsername implements User.
func (u *userMock) CheckUsername(ctx context.Context, tx *gorm.DB, username string) (user *model.User, flag bool, err error) {
	args := u.Called(ctx, tx, username)
	return args.Get(0).(*model.User), args.Bool(1), args.Error(2)
}

// Register implements User.
func (u *userMock) Register(ctx context.Context, tx *gorm.DB, user model.User) (*model.User, error) {
	args := u.Called(ctx, tx, user)
	return args.Get(0).(*model.User), args.Error(1)
}
func NewUserRepositoryMock() *userMock {
	return &userMock{}
}
