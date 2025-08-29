package repository

import (
	"context"
	"errors"

	"github.com/suttapak/starter/internal/model"

	"gorm.io/gorm"
)

type (
	User interface {
		Register(ctx context.Context, tx *gorm.DB, user model.User) (*model.User, error)
		CheckUsername(ctx context.Context, tx *gorm.DB, username string) (user *model.User, flag bool, err error)
		CheckEmail(ctx context.Context, tx *gorm.DB, email string) (user *model.User, flag bool, err error)
		GetUserByEmailOrUsername(ctx context.Context, tx *gorm.DB, emailOrUsername string) (user *model.User, err error)
		FindById(ctx context.Context, tx *gorm.DB, uId uint) (user *model.User, err error)
		VerifyEmail(ctx context.Context, tx *gorm.DB, userId uint) (user *model.User, err error)
		FindByUsername(ctx context.Context, tx *gorm.DB, username string) (user []model.User, err error)
		CheckUserIsVerifyEmail(ctx context.Context, tx *gorm.DB, userId uint) (bool, error)
		// Create Image Profile From Image ID
		CreateImageProfile(ctx context.Context, tx *gorm.DB, userId uint, imageId uint) error
	}
	user struct {
		db *gorm.DB
	}
)

// CreateImageProfile implements User.
func (u *user) CreateImageProfile(ctx context.Context, tx *gorm.DB, userId uint, imageId uint) error {
	if tx == nil {
		tx = u.db
	}
	profileImage := &model.ProfileImage{
		UserID:  userId,
		ImageID: imageId,
	}
	return tx.WithContext(ctx).Create(profileImage).Error

}

// CheckUserIsVerifyEmail checks if a user's email is verified.
func (u *user) CheckUserIsVerifyEmail(ctx context.Context, tx *gorm.DB, userId uint) (bool, error) {
	if tx == nil {
		tx = u.db
	}

	var isVerified bool
	err := tx.WithContext(ctx).
		Model(&model.User{}).
		Select("email_verifyed").
		Where("id = ?", userId).
		Scan(&isVerified).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // User not found, treat as not verified
		}
		return false, err
	}

	return isVerified, nil
}

// FindByUsername implements User.
func (u *user) FindByUsername(ctx context.Context, tx *gorm.DB, username string) (user []model.User, err error) {
	if tx == nil {
		tx = u.db
	}
	err = tx.Model(&model.User{}).Where("username ILIKE ?", username+"%").Limit(10).Find(&user).Error
	return
}

// CommitTx implements User.
func (u *user) CommitTx(tx *gorm.DB) {
	if tx == nil {
		return
	}
	tx.Commit()
}

// RollbackTx implements User.
func (u *user) RollbackTx(tx *gorm.DB) {
	if tx == nil {
		return
	}
	tx.Rollback()
}

// FindById implements User.
func (u *user) FindById(ctx context.Context, tx *gorm.DB, uId uint) (user *model.User, err error) {
	if tx == nil {
		tx = u.db
	}
	err = tx.WithContext(ctx).Preload("ProfileImage.Image").Preload("Role").First(&user, uId).Error
	return
}

func (u *user) VerifyEmail(ctx context.Context, tx *gorm.DB, userId uint) (user *model.User, err error) {
	if tx == nil {
		tx = u.db
	}
	err = tx.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).Update("email_verifyed", true).Error
	if err != nil {
		return nil, err
	}

	// Fetch the updated user
	err = tx.WithContext(ctx).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return
}

func (u *user) GetUserByEmailOrUsername(ctx context.Context, tx *gorm.DB, emailOrUsername string) (user *model.User, err error) {
	//TODO implement me
	if tx == nil {
		tx = u.db
	}
	err = tx.WithContext(ctx).Where("email = ? or username = ?", emailOrUsername, emailOrUsername).First(&user).Error
	return
}

// CheckEmail implements User.
func (u *user) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (user *model.User, flag bool, err error) {
	if tx == nil {
		tx = u.db
	}
	err = tx.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err == nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, true, err
	}
	return user, false, nil
}

func (u *user) CheckUsername(ctx context.Context, tx *gorm.DB, username string) (user *model.User, flag bool, err error) {
	if tx == nil {
		tx = u.db
	}
	err = tx.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err == nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, true, err
	}
	return user, false, nil
}

func (u *user) Register(ctx context.Context, tx *gorm.DB, user model.User) (*model.User, error) {
	if tx == nil {
		tx = u.db
	}

	err := tx.WithContext(ctx).Create(&user).Error
	return &user, err
}

func newUserRepository(db *gorm.DB) User {
	return &user{db: db}
}
