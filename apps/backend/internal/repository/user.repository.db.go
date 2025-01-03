package repository

import (
	"context"
	"errors"
	"github.com/suttapak/starter/internal/model"

	"gorm.io/gorm"
)

type (
	User interface {
		BeginTx() *gorm.DB
		Register(ctx context.Context, tx *gorm.DB, user model.User) (*model.User, error)
		CheckUsername(ctx context.Context, tx *gorm.DB, username string) (user *model.User, flag bool, err error)
		CheckEmail(ctx context.Context, tx *gorm.DB, email string) (user *model.User, flag bool, err error)
		GetUserByEmailOrUsername(ctx context.Context, tx *gorm.DB, emailOrUsername string) (user *model.User, err error)
		GetUserByUserId(ctx context.Context, tx *gorm.DB, uId uint) (user *model.User, err error)
		VerifyEmail(ctx context.Context, tx *gorm.DB, email string) (user *model.User, err error)
	}
	user struct {
		db *gorm.DB
	}
)

// GetUserByUserId implements User.
func (u user) GetUserByUserId(ctx context.Context, tx *gorm.DB, uId uint) (user *model.User, err error) {
	if tx == nil {
		tx = u.db
	}
	err = tx.WithContext(ctx).First(&user, uId).Error
	return
}

func (u user) VerifyEmail(ctx context.Context, tx *gorm.DB, email string) (user *model.User, err error) {
	if tx == nil {
		tx = u.db
	}
	err = tx.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Update("email_verifyed", true).Error
	if err != nil {
		return nil, err
	}

	// Fetch the updated user
	err = tx.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return
}

func (u user) GetUserByEmailOrUsername(ctx context.Context, tx *gorm.DB, emailOrUsername string) (user *model.User, err error) {
	//TODO implement me
	if tx == nil {
		tx = u.db
	}
	err = tx.WithContext(ctx).Where("email = ? or username = ?", emailOrUsername, emailOrUsername).First(&user).Error
	return
}

// CheckEmail implements User.
func (u user) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (user *model.User, flag bool, err error) {
	if tx == nil {
		tx = u.db
	}
	err = tx.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err == nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, true, err
	}
	return user, false, nil
}

func (u user) CheckUsername(ctx context.Context, tx *gorm.DB, username string) (user *model.User, flag bool, err error) {
	if tx == nil {
		tx = u.db
	}
	err = tx.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err == nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, true, err
	}
	return user, false, nil
}

func (u user) BeginTx() *gorm.DB {
	return u.db.Begin()
}

func (u user) Register(ctx context.Context, tx *gorm.DB, user model.User) (*model.User, error) {
	if tx == nil {
		tx = u.db
	}

	err := tx.WithContext(ctx).Create(&user).Error
	return &user, err
}

func newUserRepository(db *gorm.DB) User {
	return user{db: db}
}
