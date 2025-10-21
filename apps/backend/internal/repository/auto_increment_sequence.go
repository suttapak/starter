package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/suttapak/starter/internal/model"
	"gorm.io/gorm"
)

type (
	AutoIncrementSequence interface {
		GetNextSequence(ctx context.Context, tx *gorm.DB, entityType model.EntityType, teamId uint, entityId uint) (uint, error)
		ResetSequence(ctx context.Context, tx *gorm.DB, entityType model.EntityType, teamId uint, entityId uint) error
	}

	autoIncrementSequence struct {
		db *gorm.DB
	}
)

// GetNextSequence atomically increments and returns the next sequence number
func (a *autoIncrementSequence) GetNextSequence(ctx context.Context, tx *gorm.DB, entityType model.EntityType, teamId uint, entityId uint) (uint, error) {
	if tx == nil {
		tx = a.db
	}

	var sequence model.AutoIncrementSequence

	// Use SELECT FOR UPDATE to prevent race conditions
	err := tx.WithContext(ctx).
		Set("gorm:query_option", "FOR UPDATE").
		Where("entity_type = ? AND team_id = ? AND entity_id = ?", entityType, teamId, entityId).
		First(&sequence).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new sequence starting at 1
			sequence = model.AutoIncrementSequence{
				EntityType: entityType,
				TeamID:     teamId,
				EntityID:   entityId,
				Sequence:   1,
			}
			if err := tx.Create(&sequence).Error; err != nil {
				return 0, fmt.Errorf("failed to create sequence: %w", err)
			}
			return 1, nil
		}
		return 0, fmt.Errorf("failed to get sequence: %w", err)
	}
	if time.Now().Truncate(24 * time.Hour).Equal(sequence.UpdatedAt.Truncate(24 * time.Hour)) {
		// Increment sequence
		sequence.Sequence++
	} else {
		if entityType == model.EntityTypeProduct || entityType == model.EntityTypeLot {
			sequence.Sequence++
		} else {
			// reset to 1 in this date
			sequence.Sequence = 1
		}
	}

	if err := tx.Save(&sequence).Error; err != nil {
		return 0, fmt.Errorf("failed to update sequence: %w", err)
	}
	return sequence.Sequence, nil
}

// ResetSequence resets the sequence counter to 0
func (a *autoIncrementSequence) ResetSequence(ctx context.Context, tx *gorm.DB, entityType model.EntityType, teamId uint, entityId uint) error {
	if tx == nil {
		tx = a.db
	}

	return tx.WithContext(ctx).
		Model(&model.AutoIncrementSequence{}).
		Where("entity_type = ? AND team_id = ? AND entity_id = ?", entityType, teamId, entityId).
		Update("sequence", 0).Error
}

func NewAutoIncrementSequence(db *gorm.DB) AutoIncrementSequence {
	return &autoIncrementSequence{
		db: db,
	}
}
