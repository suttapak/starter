package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/suttapak/starter/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T, db *sqlx.DB, email, username, password string) *model.User {
	ctx := context.Background()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	// Get default role ID
	var roleID int
	err = db.GetContext(ctx, &roleID, "SELECT id FROM roles WHERE name = $1", "User")
	require.NoError(t, err)

	// Create user
	user := &model.User{
		Email:         email,
		Username:      username,
		Password:      string(hashedPassword),
		FullName:      "Test User",
		RoleID:        uint(roleID),
		EmailVerifyed: false,
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (email, username, password, full_name, role_id, email_verifyed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err = db.GetContext(ctx, &user.ID, query,
		user.Email,
		user.Username,
		user.Password,
		user.FullName,
		user.RoleID,
		user.EmailVerifyed,
		user.CreatedAt,
		user.UpdatedAt,
	)
	require.NoError(t, err)

	return user
}

// CreateTestTeam creates a test team in the database
func CreateTestTeam(t *testing.T, db *sqlx.DB, name, username, description string, ownerID int) *model.Team {
	ctx := context.Background()

	team := &model.Team{
		Name:        name,
		Username:    username,
		Description: &description,
	}
	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()

	query := `
		INSERT INTO teams (name, username, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := db.GetContext(ctx, &team.ID, query,
		team.Name,
		team.Username,
		team.Description,
		team.CreatedAt,
		team.UpdatedAt,
	)
	require.NoError(t, err)

	// Get team role ID for owner
	var roleID uint
	err = db.GetContext(ctx, &roleID, "SELECT id FROM team_roles WHERE name = $1", "Owner")
	require.NoError(t, err)

	// Add owner as team member
	memberQuery := `
		INSERT INTO team_members (team_id, user_id, team_role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = db.ExecContext(ctx, memberQuery, team.ID, ownerID, roleID, true, time.Now(), time.Now())
	require.NoError(t, err)

	return team
}

// CreateTestProduct creates a test product in the database
func CreateTestProduct(t *testing.T, db *sqlx.DB, code, name, description string, price int, teamID int) *model.Product {
	ctx := context.Background()

	product := &model.Product{
		Code:        code,
		Name:        name,
		Description: description,
		Price:       int64(price * 100), // Convert to cents
		TeamID:      uint(teamID),
		UOM:         "ชิ้น",
	}
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	query := `
		INSERT INTO products (code, name, description, price, uom, team_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := db.GetContext(ctx, &product.ID, query,
		product.Code,
		product.Name,
		product.Description,
		product.Price,
		product.UOM,
		product.TeamID,
		product.CreatedAt,
		product.UpdatedAt,
	)
	require.NoError(t, err)

	return product
}

// CreateTestProductCategory creates a test product category in the database
func CreateTestProductCategory(t *testing.T, db *sqlx.DB, name string, teamID int) *model.ProductCategory {
	ctx := context.Background()

	category := &model.ProductCategory{
		Name:   name,
		TeamID: uint(teamID),
	}
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	query := `
		INSERT INTO product_categories (name, team_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := db.GetContext(ctx, &category.ID, query,
		category.Name,
		category.TeamID,
		category.CreatedAt,
		category.UpdatedAt,
	)
	require.NoError(t, err)

	return category
}

// CleanupUser deletes a test user from the database
func CleanupUser(t *testing.T, db *sqlx.DB, userID int) {
	ctx := context.Background()
	_, err := db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		t.Logf("Warning: failed to cleanup user %d: %v", userID, err)
	}
}

// CleanupTeam deletes a test team from the database
func CleanupTeam(t *testing.T, db *sqlx.DB, teamID int) {
	ctx := context.Background()
	_, err := db.ExecContext(ctx, "DELETE FROM teams WHERE id = $1", teamID)
	if err != nil {
		t.Logf("Warning: failed to cleanup team %d: %v", teamID, err)
	}
}
