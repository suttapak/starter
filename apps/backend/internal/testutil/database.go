package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/suttapak/starter/database"
	"github.com/suttapak/starter/logger"
)

// SetupTestDB sets up a test database connection and runs migrations
func SetupTestDB(t *testing.T) *sqlx.DB {
	// Set config path for test environment
	os.Setenv("CONFIG_PATH", "../../configs.test.toml")

	// Database connection string for tests
	dsn := "host=localhost user=test_user password=test_password dbname=test_db port=5433 sslmode=disable TimeZone=Asia/Bangkok"

	// Wait for database to be ready
	var db *sqlx.DB
	var err error
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	require.NoError(t, err, "Failed to connect to test database")

	// Create a simple test logger
	testLogger := logger.NewLoggerMock()

	// Run migrations
	err = database.RunMigrations(db.DB, testLogger)
	require.NoError(t, err, "Failed to run migrations")

	return db
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(t *testing.T, db *sqlx.DB) {
	if db != nil {
		// Clean all tables
		tables := []string{
			"product_images",
			"products",
			"product_categories",
			"product_product_categories",
			"images",
			"report_templates",
			"report_json_schema_types",
			"team_members",
			"team_roles",
			"teams",
			"users",
			"roles",
			"auto_increment_sequences",
		}

		ctx := context.Background()
		for _, table := range tables {
			query := fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY CASCADE`, table)

			_, err := db.ExecContext(ctx, query)
			if err != nil {
				t.Logf("Warning: failed to clean table %s: %v", table, err)
			}
		}

		db.Close()
	}
}

// SeedRoles seeds the roles table with test data
func SeedRoles(t *testing.T, db *sqlx.DB) {
	ctx := context.Background()

	// Check if roles already exist
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM roles")
	require.NoError(t, err)

	if count > 0 {
		return // Roles already seeded
	}

	roles := []map[string]interface{}{
		{"name": "User"},
		{"name": "Moderator"},
		{"name": "Admin"},
		{"name": "SuperAdmin"},
	}

	for _, role := range roles {
		_, err := db.ExecContext(ctx,
			"INSERT INTO roles (name, created_at, updated_at) VALUES ($1, NOW(), NOW())",
			role["name"],
		)
		require.NoError(t, err)
	}
}

// SeedTeamRoles seeds the team_roles table with test data
func SeedTeamRoles(t *testing.T, db *sqlx.DB) {
	ctx := context.Background()

	// Check if team roles already exist
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM team_roles")
	require.NoError(t, err)

	if count > 0 {
		return // Team roles already seeded
	}

	teamRoles := []map[string]interface{}{
		{"name": "Owner"},
		{"name": "Admin"},
		{"name": "Member"},
	}

	for _, role := range teamRoles {
		_, err := db.ExecContext(ctx,
			"INSERT INTO team_roles (name, created_at, updated_at) VALUES ($1,  NOW(), NOW())",
			role["name"],
		)
		require.NoError(t, err)
	}
}

// SeedReportSchemaTypes seeds the report_json_schema_types table with test data
func SeedReportSchemaTypes(t *testing.T, db *sqlx.DB) {
	ctx := context.Background()

	// Check if schema types already exist
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM report_json_schema_types")
	require.NoError(t, err)

	if count > 0 {
		return // Schema types already seeded
	}

	schemaTypes := []map[string]interface{}{
		{"name": "text"},
		{"name": "number"},
		{"name": "date"},
		{"name": "select"},
	}

	for _, schemaType := range schemaTypes {
		_, err := db.ExecContext(ctx,
			"INSERT INTO report_json_schema_types (name, created_at, updated_at) VALUES ($1, NOW(), NOW())",
			schemaType["name"],
		)
		require.NoError(t, err)
	}
}

// SeedBasicData seeds roles, team_roles, and report_json_schema_types
func SeedBasicData(t *testing.T, db *sqlx.DB) {
	SeedRoles(t, db)
	SeedTeamRoles(t, db)
	SeedReportSchemaTypes(t, db)
}
