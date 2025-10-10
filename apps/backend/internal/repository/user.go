package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/suttapak/starter/internal/model"
)

type (
	User interface {
		Register(ctx context.Context, tx *sqlx.Tx, user model.User) (*model.User, error)
		CheckUsername(ctx context.Context, tx *sqlx.Tx, username string) (user *model.User, flag bool, err error)
		CheckEmail(ctx context.Context, tx *sqlx.Tx, email string) (user *model.User, flag bool, err error)
		GetUserByEmailOrUsername(ctx context.Context, tx *sqlx.Tx, emailOrUsername string) (user *model.User, err error)
		FindById(ctx context.Context, tx *sqlx.Tx, uId uint) (user *model.User, err error)
		VerifyEmail(ctx context.Context, tx *sqlx.Tx, userId uint) (user *model.User, err error)
		FindByUsername(ctx context.Context, tx *sqlx.Tx, username string) (user []model.User, err error)
		IsVerifyEmailByUserId(ctx context.Context, tx *sqlx.Tx, userId uint) (bool, error)
		CreateImageProfile(ctx context.Context, tx *sqlx.Tx, userId uint, imageId uint) error
	}

	userSqlx struct {
		db *sqlx.DB
	}
)

func NewUser(db *sqlx.DB) User {
	return &userSqlx{db: db}
}

// getDB returns the appropriate database connection (transaction or main DB)
func (u *userSqlx) getDB(tx *sqlx.Tx) sqlx.ExtContext {
	if tx != nil {
		return tx
	}
	return u.db
}

// Register creates a new user
func (u *userSqlx) Register(ctx context.Context, tx *sqlx.Tx, user model.User) (*model.User, error) {
	db := u.getDB(tx)

	query := `
		INSERT INTO users (username, password, email, email_verifyed, full_name, role_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := sqlx.GetContext(ctx, db, &user, query,
		user.Username,
		user.Password,
		user.Email,
		user.EmailVerifyed,
		user.FullName,
		user.RoleID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return &user, nil
}

// CheckUsername checks if a username exists
func (u *userSqlx) CheckUsername(ctx context.Context, tx *sqlx.Tx, username string) (*model.User, bool, error) {
	db := u.getDB(tx)

	var user model.User
	query := `SELECT id, username, email FROM users WHERE username = $1 AND deleted_at IS NULL LIMIT 1`

	err := sqlx.GetContext(ctx, db, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to check username: %w", err)
	}

	return &user, true, nil
}

// CheckEmail checks if an email exists
func (u *userSqlx) CheckEmail(ctx context.Context, tx *sqlx.Tx, email string) (*model.User, bool, error) {
	db := u.getDB(tx)

	var user model.User
	query := `SELECT id, username, email FROM users WHERE email = $1 AND deleted_at IS NULL LIMIT 1`

	err := sqlx.GetContext(ctx, db, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to check email: %w", err)
	}

	return &user, true, nil
}

// GetUserByEmailOrUsername finds a user by email or username
func (u *userSqlx) GetUserByEmailOrUsername(ctx context.Context, tx *sqlx.Tx, emailOrUsername string) (*model.User, error) {
	db := u.getDB(tx)

	var user model.User
	query := `
		SELECT id, username, password, email, email_verifyed, full_name, role_id, created_at, updated_at
		FROM users
		WHERE (email = $1 OR username = $1) AND deleted_at IS NULL
		LIMIT 1
	`

	err := sqlx.GetContext(ctx, db, &user, query, emailOrUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// FindById finds a user by ID with profile image and role
func (u *userSqlx) FindById(ctx context.Context, tx *sqlx.Tx, uId uint) (*model.User, error) {
	db := u.getDB(tx)

	var user model.User
	query := `
		SELECT u.id, u.username, u.email, u.email_verifyed, u.full_name, u.role_id, u.created_at, u.updated_at
		FROM users u
		WHERE u.id = $1 AND u.deleted_at IS NULL
		LIMIT 1
	`

	err := sqlx.GetContext(ctx, db, &user, query, uId)
	if err != nil {
		return nil, err
	}

	// Load role
	roleQuery := `SELECT id, name FROM roles WHERE id = $1 AND deleted_at IS NULL`
	var role model.Role
	if err := sqlx.GetContext(ctx, db, &role, roleQuery, user.RoleID); err == nil {
		user.Role = &role
	}

	// Load profile images
	imagesQuery := `
		SELECT pi.id, pi.user_id, pi.image_id, pi.created_at, pi.updated_at,
		       i.path, i.url, i.size, i.width, i.height, i.type
		FROM profile_images pi
		INNER JOIN images i ON i.id = pi.image_id
		WHERE pi.user_id = $1 AND pi.deleted_at IS NULL AND i.deleted_at IS NULL
	`

	rows, err := db.QueryxContext(ctx, imagesQuery, uId)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var pi model.ProfileImage
			var img model.Image
			if err := rows.Scan(
				&pi.ID, &pi.UserID, &pi.ImageID, &pi.CreatedAt, &pi.UpdatedAt,
				&img.Path, &img.Url, &img.Size, &img.Width, &img.Height, &img.Type,
			); err == nil {
				pi.Image = &img
				user.ProfileImage = append(user.ProfileImage, pi)
			}
		}
	}

	return &user, nil
}

// VerifyEmail marks a user's email as verified
func (u *userSqlx) VerifyEmail(ctx context.Context, tx *sqlx.Tx, userId uint) (*model.User, error) {
	db := u.getDB(tx)

	query := `
		UPDATE users
		SET email_verifyed = true, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, username, email, email_verifyed, full_name, role_id, created_at, updated_at
	`

	var user model.User
	err := sqlx.GetContext(ctx, db, &user, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}

	return &user, nil
}

// FindByUsername searches for users by username prefix
func (u *userSqlx) FindByUsername(ctx context.Context, tx *sqlx.Tx, username string) ([]model.User, error) {
	db := u.getDB(tx)

	query := `
		SELECT id, username, email, full_name, role_id, created_at, updated_at
		FROM users
		WHERE username ILIKE $1 AND deleted_at IS NULL
		ORDER BY username
		LIMIT 10
	`

	var users []model.User
	err := sqlx.SelectContext(ctx, db, &users, query, username+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to find users by username: %w", err)
	}

	return users, nil
}

// IsVerifyEmailByUserId checks if a user's email is verified
func (u *userSqlx) IsVerifyEmailByUserId(ctx context.Context, tx *sqlx.Tx, userId uint) (bool, error) {
	db := u.getDB(tx)

	var isVerified bool
	query := `SELECT email_verifyed FROM users WHERE id = $1 AND deleted_at IS NULL`

	err := sqlx.GetContext(ctx, db, &isVerified, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check email verification: %w", err)
	}

	return isVerified, nil
}

// CreateImageProfile creates a profile image entry
func (u *userSqlx) CreateImageProfile(ctx context.Context, tx *sqlx.Tx, userId uint, imageId uint) error {
	db := u.getDB(tx)

	query := `
		INSERT INTO profile_images (user_id, image_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`

	_, err := db.ExecContext(ctx, query, userId, imageId)
	if err != nil {
		return fmt.Errorf("failed to create profile image: %w", err)
	}

	return nil
}
