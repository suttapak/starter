package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/suttapak/starter/internal/model"
)

type (
	AutoIncrementSequence interface {
		GetNextSequence(ctx context.Context, tx *sqlx.Tx, entityType model.EntityType, teamId uint, entityId uint) (uint, error)
		ResetSequence(ctx context.Context, tx *sqlx.Tx, entityType model.EntityType, teamId uint, entityId uint) error
	}

	autoIncrementSequenceSqlx struct {
		db *sqlx.DB
	}
)

func NewAutoIncrementSequence(db *sqlx.DB) AutoIncrementSequence {
	return &autoIncrementSequenceSqlx{db: db}
}

// getDB returns the appropriate database connection (transaction or main DB)
func (a *autoIncrementSequenceSqlx) getDB(tx *sqlx.Tx) sqlx.ExtContext {
	if tx != nil {
		return tx
	}
	return a.db
}

// GetNextSequence atomically increments and returns the next sequence number
func (a *autoIncrementSequenceSqlx) GetNextSequence(ctx context.Context, tx *sqlx.Tx, entityType model.EntityType, teamId uint, entityId uint) (uint, error) {
	db := a.getDB(tx)

	// Use SELECT FOR UPDATE to prevent race conditions
	selectQuery := `
		SELECT id, entity_type, team_id, entity_id, sequence, created_at, updated_at
		FROM auto_increment_sequences
		WHERE entity_type = $1 AND team_id = $2 AND entity_id = $3 AND deleted_at IS NULL
		FOR UPDATE
	`

	var sequence model.AutoIncrementSequence
	err := sqlx.GetContext(ctx, db, &sequence, selectQuery, entityType, teamId, entityId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Create new sequence starting at 1
			insertQuery := `
				INSERT INTO auto_increment_sequences (entity_type, team_id, entity_id, sequence, created_at, updated_at)
				VALUES ($1, $2, $3, 1, NOW(), NOW())
				RETURNING id, sequence, created_at, updated_at
			`

			err := sqlx.GetContext(ctx, db, &sequence, insertQuery, entityType, teamId, entityId)
			if err != nil {
				return 0, fmt.Errorf("failed to create sequence: %w", err)
			}
			return 1, nil
		}
		return 0, fmt.Errorf("failed to get sequence: %w", err)
	}

	// Determine new sequence value based on date and entity type
	var newSequence uint
	today := time.Now().Truncate(24 * time.Hour)
	sequenceDate := sequence.UpdatedAt.Truncate(24 * time.Hour)

	if today.Equal(sequenceDate) {
		// Same day - increment sequence
		newSequence = sequence.Sequence + 1
	} else {
		// Different day
		if entityType == model.EntityTypeProduct || entityType == model.EntityTypeLot {
			// Products and lots continue incrementing across days
			newSequence = sequence.Sequence + 1
		} else {
			// Other types reset to 1 each day (transactions)
			newSequence = 1
		}
	}

	// Update the sequence
	updateQuery := `
		UPDATE auto_increment_sequences
		SET sequence = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`

	_, err = db.ExecContext(ctx, updateQuery, newSequence, sequence.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to update sequence: %w", err)
	}

	return newSequence, nil
}

// ResetSequence resets the sequence counter to 0
func (a *autoIncrementSequenceSqlx) ResetSequence(ctx context.Context, tx *sqlx.Tx, entityType model.EntityType, teamId uint, entityId uint) error {
	db := a.getDB(tx)

	query := `
		UPDATE auto_increment_sequences
		SET sequence = 0, updated_at = NOW()
		WHERE entity_type = $1 AND team_id = $2 AND entity_id = $3 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, query, entityType, teamId, entityId)
	if err != nil {
		return fmt.Errorf("failed to reset sequence: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("sequence not found")
	}

	return nil
}
