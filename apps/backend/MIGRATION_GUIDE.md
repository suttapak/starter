# GORM to sqlx Migration Guide

This guide documents the migration from GORM to sqlx + Goose for database operations in this project.

## What Changed

### Dependencies
- **Added**: `github.com/jmoiron/sqlx`, `github.com/lib/pq`, `github.com/pressly/goose/v3`
- **Kept (for now)**: GORM dependencies (will be removed after full migration)

### Database Module
- **New**: `database/sqlx.database.go` - sqlx connection setup
- **New**: `database/migrate.go` - Goose migration runner with embedded migrations
- **Updated**: `database/module.database.go` - now provides `*sqlx.DB` instead of `*gorm.DB`
- **Backup**: `database/gorm.database.go` - kept for reference

### Migrations
- **Location**: `apps/backend/database/migrations/`
- **Initial migration**: `20251009090803_init_schema.sql`
- **Features**:
  - Complete schema with all tables
  - Proper foreign keys and indexes
  - Seed data for roles, team_roles, report_json_schema_types
  - Embedded in binary and run automatically on startup

### Models
All model structs updated:
- Changed from `gorm:"tag"` to `db:"tag"`
- Changed `gorm.DeletedAt` to `sql.NullTime`
- Relations marked with `db:"-"` (must be loaded manually)
- Kept JSON tags for API responses

### Repositories
- **Interface changes**: Methods now use `*sqlx.Tx` instead of `*gorm.DB`
- **Migrated**: `internal/repository/user.go` (complete reference implementation)
- **Pending**: team, product, image, report, and other repositories

## How to Use sqlx

### Basic Query Patterns

#### Get Single Row
```go
func (r *repo) FindByID(ctx context.Context, tx *sqlx.Tx, id uint) (*model.Entity, error) {
    db := r.getDB(tx) // Helper to use tx or main db

    var entity model.Entity
    query := `SELECT id, name, created_at FROM entities WHERE id = $1 AND deleted_at IS NULL`

    err := sqlx.GetContext(ctx, db, &entity, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("entity not found")
        }
        return nil, fmt.Errorf("failed to find entity: %w", err)
    }

    return &entity, nil
}
```

#### Get Multiple Rows
```go
func (r *repo) FindAll(ctx context.Context, tx *sqlx.Tx) ([]model.Entity, error) {
    db := r.getDB(tx)

    var entities []model.Entity
    query := `SELECT id, name, created_at FROM entities WHERE deleted_at IS NULL ORDER BY created_at DESC`

    err := sqlx.SelectContext(ctx, db, &entities, query)
    if err != nil {
        return nil, fmt.Errorf("failed to find entities: %w", err)
    }

    return entities, nil
}
```

#### Insert with RETURNING
```go
func (r *repo) Create(ctx context.Context, tx *sqlx.Tx, entity model.Entity) (*model.Entity, error) {
    db := r.getDB(tx)

    query := `
        INSERT INTO entities (name, description, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        RETURNING id, created_at, updated_at
    `

    err := sqlx.GetContext(ctx, db, &entity, query, entity.Name, entity.Description)
    if err != nil {
        return nil, fmt.Errorf("failed to create entity: %w", err)
    }

    return &entity, nil
}
```

#### Update
```go
func (r *repo) Update(ctx context.Context, tx *sqlx.Tx, id uint, entity model.Entity) error {
    db := r.getDB(tx)

    query := `
        UPDATE entities
        SET name = $1, description = $2, updated_at = NOW()
        WHERE id = $3 AND deleted_at IS NULL
    `

    result, err := db.ExecContext(ctx, query, entity.Name, entity.Description, id)
    if err != nil {
        return fmt.Errorf("failed to update entity: %w", err)
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("entity not found")
    }

    return nil
}
```

#### Soft Delete
```go
func (r *repo) Delete(ctx context.Context, tx *sqlx.Tx, id uint) error {
    db := r.getDB(tx)

    query := `UPDATE entities SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

    result, err := db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete entity: %w", err)
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("entity not found")
    }

    return nil
}
```

### Transaction Helper
```go
// Add this helper to your repository struct
func (r *repo) getDB(tx *sqlx.Tx) sqlx.ExtContext {
    if tx != nil {
        return tx
    }
    return r.db
}
```

### Loading Relations
Since sqlx doesn't have eager loading, load relations manually:

```go
func (r *repo) FindByIDWithRelations(ctx context.Context, tx *sqlx.Tx, id uint) (*model.Entity, error) {
    db := r.getDB(tx)

    // 1. Load main entity
    var entity model.Entity
    mainQuery := `SELECT id, name, category_id FROM entities WHERE id = $1 AND deleted_at IS NULL`
    if err := sqlx.GetContext(ctx, db, &entity, mainQuery, id); err != nil {
        return nil, err
    }

    // 2. Load related category
    categoryQuery := `SELECT id, name FROM categories WHERE id = $1 AND deleted_at IS NULL`
    var category model.Category
    if err := sqlx.GetContext(ctx, db, &category, categoryQuery, entity.CategoryID); err == nil {
        entity.Category = &category
    }

    // 3. Load has-many relations
    itemsQuery := `SELECT id, entity_id, name FROM items WHERE entity_id = $1 AND deleted_at IS NULL`
    var items []model.Item
    if err := sqlx.SelectContext(ctx, db, &items, itemsQuery, id); err == nil {
        entity.Items = items
    }

    return &entity, nil
}
```

## Migration Workflow

### 1. Create Migration
```bash
make new-migrate name=add_new_feature
```

### 2. Edit Migration File
Edit `database/migrations/XXXXXX_add_new_feature.sql`:
```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE new_table (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_new_table_deleted_at ON new_table(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS new_table;
-- +goose StatementEnd
```

### 3. Run Migration
Migrations run automatically on app startup, or manually:
```bash
make migrate-up DB_DSN="host=localhost user=... password=... dbname=..."
```

## Migrating a Repository

### Step 1: Update Interface
Change from:
```go
type MyRepo interface {
    FindByID(ctx context.Context, tx *gorm.DB, id uint) (*model.Entity, error)
}
```

To:
```go
type MyRepo interface {
    FindByID(ctx context.Context, tx *sqlx.Tx, id uint) (*model.Entity, error)
}
```

### Step 2: Update Implementation
Change struct:
```go
type myRepo struct {
    db *sqlx.DB // Changed from *gorm.DB
}

func NewMyRepo(db *sqlx.DB) MyRepo {
    return &myRepo{db: db}
}
```

### Step 3: Rewrite Methods
Convert GORM queries to raw SQL:

**Before (GORM)**:
```go
func (r *myRepo) FindByID(ctx context.Context, tx *gorm.DB, id uint) (*model.Entity, error) {
    if tx == nil {
        tx = r.db
    }
    var entity model.Entity
    err := tx.WithContext(ctx).Preload("Category").First(&entity, id).Error
    return &entity, err
}
```

**After (sqlx)**:
```go
func (r *myRepo) FindByID(ctx context.Context, tx *sqlx.Tx, id uint) (*model.Entity, error) {
    db := r.getDB(tx)

    var entity model.Entity
    query := `SELECT id, name, category_id FROM entities WHERE id = $1 AND deleted_at IS NULL`

    err := sqlx.GetContext(ctx, db, &entity, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("entity not found")
        }
        return nil, fmt.Errorf("failed to find entity: %w", err)
    }

    // Load category manually
    if entity.CategoryID > 0 {
        var category model.Category
        catQuery := `SELECT id, name FROM categories WHERE id = $1 AND deleted_at IS NULL`
        if err := sqlx.GetContext(ctx, db, &category, catQuery, entity.CategoryID); err == nil {
            entity.Category = &category
        }
    }

    return &entity, nil
}
```

### Step 4: Update Module
Update `internal/repository/module.go`:
```go
fx.Provide(func(db *sqlx.DB) MyRepo {
    return NewMyRepo(db)
}),
```

## Testing

### Run Application
```bash
cd apps/backend
air
```

Check logs for:
- `Migrations completed successfully`
- No GORM-related errors

### Test Queries
Use the migrated user repository as reference for testing patterns.

## Common Issues

### 1. Missing Soft Delete Filter
Always include `AND deleted_at IS NULL` in WHERE clauses.

### 2. Wrong Placeholder
Use `$1, $2, $3` (PostgreSQL) not `?` (MySQL).

### 3. Relations Not Loading
Remember to manually load relations - sqlx doesn't auto-eager load.

### 4. Transaction Type Mismatch
Use `*sqlx.Tx` not `*gorm.DB` in method signatures.

## Rollback Plan

If issues arise:
1. The old GORM code is backed up (`.gorm.go.backup` files)
2. Restore GORM module in `database/module.database.go`
3. Revert repository module changes
4. Rebuild and deploy

## Next Steps

Priority order for migrating remaining repositories:
1. ✅ User repository (completed)
2. Team repository (high priority - core feature)
3. Product repository (high priority - core feature)
4. Image repository
5. Report repository
6. Auto-increment sequence repository
7. ODT repository

## Resources

- [sqlx Documentation](https://jmoiron.github.io/sqlx/)
- [Goose Documentation](https://pressly.github.io/goose/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- Example implementation: `internal/repository/user.go`
