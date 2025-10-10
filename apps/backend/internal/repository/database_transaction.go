package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type (
	DatabaseTransaction interface {
		BeginTx(ctx context.Context) (*sqlx.Tx, error)
		CommitTx(tx *sqlx.Tx) error
		RollbackTx(tx *sqlx.Tx) error
		// WithTransaction executes a function within a transaction
		WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error
	}

	databaseTransactionSqlx struct {
		db *sqlx.DB
	}
)

// NewDatabaseTransaction creates a new database transaction helper
func NewDatabaseTransaction(db *sqlx.DB) DatabaseTransaction {
	return &databaseTransactionSqlx{db: db}
}

// BeginTx starts a new transaction
func (d *databaseTransactionSqlx) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

// CommitTx commits a transaction
func (d *databaseTransactionSqlx) CommitTx(tx *sqlx.Tx) error {
	if tx == nil {
		return nil
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RollbackTx rolls back a transaction
func (d *databaseTransactionSqlx) RollbackTx(tx *sqlx.Tx) error {
	if tx == nil {
		return nil
	}

	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	return nil
}

// WithTransaction executes a function within a transaction with automatic commit/rollback
func (d *databaseTransactionSqlx) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := d.BeginTx(ctx)
	if err != nil {
		return err
	}

	// Defer rollback in case of panic
	defer func() {
		if p := recover(); p != nil {
			_ = d.RollbackTx(tx)
			panic(p) // re-throw panic after rollback
		}
	}()

	// Execute the function
	if err := fn(tx); err != nil {
		if rbErr := d.RollbackTx(tx); rbErr != nil {
			return fmt.Errorf("failed to rollback after error: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	// Commit the transaction
	if err := d.CommitTx(tx); err != nil {
		return err
	}

	return nil
}
