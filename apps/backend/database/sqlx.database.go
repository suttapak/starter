package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/suttapak/starter/config"
)

func newSqlx(conf *config.Config) (*sqlx.DB, error) {
	if conf.DB.DSN == "" {
		return nil, fmt.Errorf("dsn is empty")
	}

	// Parse DSN and convert to PostgreSQL connection string
	// GORM DSN format: "host=db user=suttapak password=p@ssw0rd dbname=starter port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	// sqlx needs standard PostgreSQL connection string

	db, err := sqlx.Connect("postgres", conf.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
