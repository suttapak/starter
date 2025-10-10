package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/model"
)

type (
	Report interface {
		FindById(ctx context.Context, tx *sqlx.Tx, id uint) (*model.ReportTemplate, error)
		FindAll(ctx context.Context, tx *sqlx.Tx, pg *helpers.Pagination, f *filter.ReportFilter) ([]model.ReportTemplate, error)
		Create(ctx context.Context, tx *sqlx.Tx, m *model.ReportTemplate) error
		Save(ctx context.Context, tx *sqlx.Tx, id uint, m *model.ReportTemplate) error
		DeleteById(ctx context.Context, tx *sqlx.Tx, id uint) error
	}

	reportSqlx struct {
		db *sqlx.DB
	}
)

func NewReport(db *sqlx.DB) Report {
	return &reportSqlx{db: db}
}

// getDB returns the appropriate database connection (transaction or main DB)
func (r *reportSqlx) getDB(tx *sqlx.Tx) sqlx.ExtContext {
	if tx != nil {
		return tx
	}
	return r.db
}

// FindById retrieves a report template by ID
func (r *reportSqlx) FindById(ctx context.Context, tx *sqlx.Tx, id uint) (*model.ReportTemplate, error) {
	db := r.getDB(tx)

	var result model.ReportTemplate
	query := `
		SELECT id, code, name, display_name, icon, report_json_schema_type_id, created_at, updated_at
		FROM report_templates
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := sqlx.GetContext(ctx, db, &result, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("report template not found")
		}
		return nil, fmt.Errorf("failed to find report template: %w", err)
	}

	return &result, nil
}

// FindAll retrieves all report templates with pagination
func (r *reportSqlx) FindAll(ctx context.Context, tx *sqlx.Tx, pg *helpers.Pagination, f *filter.ReportFilter) ([]model.ReportTemplate, error) {
	db := r.getDB(tx)

	// Count total records
	countQuery := `SELECT COUNT(*) FROM report_templates WHERE deleted_at IS NULL`
	if err := sqlx.GetContext(ctx, db, &pg.Count, countQuery); err != nil {
		return nil, fmt.Errorf("failed to count report templates: %w", err)
	}

	// Apply pagination
	helpers.Paging(pg)

	// Query with pagination
	query := `
		SELECT id, code, name, display_name, icon, report_json_schema_type_id, created_at, updated_at
		FROM report_templates
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var results []model.ReportTemplate
	err := sqlx.SelectContext(ctx, db, &results, query, pg.Limit, pg.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find report templates: %w", err)
	}

	return results, nil
}

// Create creates a new report template
func (r *reportSqlx) Create(ctx context.Context, tx *sqlx.Tx, m *model.ReportTemplate) error {
	db := r.getDB(tx)

	query := `
		INSERT INTO report_templates (code, name, display_name, icon, report_json_schema_type_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := sqlx.GetContext(ctx, db, m, query,
		m.Code,
		m.Name,
		m.DisplayName,
		m.Icon,
		m.ReportJsonSchemaTypeID,
	)

	if err != nil {
		return fmt.Errorf("failed to create report template: %w", err)
	}

	return nil
}

// Save updates an existing report template
func (r *reportSqlx) Save(ctx context.Context, tx *sqlx.Tx, id uint, m *model.ReportTemplate) error {
	db := r.getDB(tx)

	query := `
		UPDATE report_templates
		SET code = $1, name = $2, display_name = $3, icon = $4, report_json_schema_type_id = $5, updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, query,
		m.Code,
		m.Name,
		m.DisplayName,
		m.Icon,
		m.ReportJsonSchemaTypeID,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update report template: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("report template not found")
	}

	return nil
}

// DeleteById soft deletes a report template
func (r *reportSqlx) DeleteById(ctx context.Context, tx *sqlx.Tx, id uint) error {
	db := r.getDB(tx)

	query := `
		UPDATE report_templates
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete report template: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("report template not found")
	}

	return nil
}
