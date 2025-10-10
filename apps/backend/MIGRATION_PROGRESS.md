# Migration Progress: GORM to sqlx

**Last Updated**: October 9, 2025 (Products completed - 80% progress)
**Branch**: `migrate-sqlx`
**Status**: ✅ Complete (All Core Repositories Migrated - 80%)

## Overview

This document tracks the progress of migrating from GORM to sqlx + Goose for database operations.

## ✅ Completed Components

### Infrastructure
- [x] **Dependencies**: Added sqlx, lib/pq, and Goose to go.mod
- [x] **Database Module**: Created `database/sqlx.database.go` for sqlx connections
- [x] **Migration System**: Implemented Goose with embedded migrations (`database/migrate.go`)
- [x] **Initial Schema**: Complete migration in `database/migrations/20251009090803_init_schema.sql`
- [x] **Model Tags**: All structs updated from `gorm:"..."` to `db:"..."`
- [x] **Build Tools**: Updated makefile with migration commands

### Migrated Repositories

| Repository | Status | File | Notes |
|------------|--------|------|-------|
| User | ✅ Complete | `repository/user.go` | Reference implementation with all CRUD patterns |
| DatabaseTransaction | ✅ Complete | `repository/database_transaction.go` | Includes `WithTransaction` helper |
| Image | ✅ Complete | `repository/image.go` | Simple CRUD operations |
| AutoIncrementSequence | ✅ Complete | `repository/auto_increment_sequence.go` | Atomic counter with FOR UPDATE |
| Report | ✅ Complete | `repository/report.go` | CRUD with pagination |
| Team | ✅ Complete | `repository/team_repository_db.go` | Complex with JOIN queries and member relations |
| ProductCategory | ✅ Complete | `repository/product_category.go` | Simple CRUD with team filtering |
| Products | ✅ Complete | `repository/products.go` | Complex with multiple relations and junction tables |
| Mail | ⚠️ No change needed | `repository/mail.repository.smtp.go` | Uses external SMTP, no DB ops |

### Backup Files
All old GORM code backed up with `.gorm.go.backup` extension:
- `repository/user.gorm.go.backup`
- `repository/database_transaction.gorm.go.backup`
- `repository/image.gorm.go.backup`
- `repository/auto_increment_sequence.gorm.go.backup`
- `repository/report.gorm.go.backup`
- `repository/team_repository_db.gorm.go.backup`
- `repository/product_category.gorm.go.backup`
- `repository/products.gorm.go.backup`

## ⏳ Pending Repositories

### High Priority (Core Features)
All high-priority repositories have been migrated! ✅

### Low Priority (Optional)
| Repository | Complexity | Dependencies | Estimated Effort |
|------------|------------|--------------|------------------|
| **ODT** | 🟢 Low | External API | 1-2 hours |

## 📊 Progress Metrics

```
Total Repositories: 10
✅ Completed: 8 (80%)
⏳ Pending: 1 (10%)
⚠️  No Change Needed: 1 (10%)
```

### Completed (80%)
- ✅ User (Complex - with relations)
- ✅ DatabaseTransaction (Helper)
- ✅ Image (Simple)
- ✅ AutoIncrementSequence (Complex - with locking)
- ✅ Report (With pagination)
- ✅ Team (Complex - with JOIN queries and member relations)
- ✅ ProductCategory (Simple CRUD with team filtering)
- ✅ Products (Complex - with multiple relations and junction tables)

### Pending (10%)
- ⏳ ODT (Low Priority - external API, optional)

## 🔄 Migration Pattern

Each repository follows this pattern:

```go
// 1. Update interface to use *sqlx.Tx
type MyRepo interface {
    Method(ctx context.Context, tx *sqlx.Tx, ...) error
}

// 2. Update struct
type myRepoSqlx struct {
    db *sqlx.DB
}

// 3. Add getDB helper
func (r *myRepoSqlx) getDB(tx *sqlx.Tx) sqlx.ExtContext {
    if tx != nil {
        return tx
    }
    return r.db
}

// 4. Rewrite methods with raw SQL
func (r *myRepoSqlx) Method(ctx context.Context, tx *sqlx.Tx, ...) error {
    db := r.getDB(tx)
    query := `SELECT ... WHERE ... AND deleted_at IS NULL`
    return sqlx.GetContext(ctx, db, &result, query, args...)
}
```

## 🚧 Known Issues

### Services Layer ⚠️ PARTIALLY RESOLVED
Services that use DatabaseTransaction have been updated to use the new interface:
- **Before**: `tx := db.BeginTx()`
- **After**: `tx, err := db.BeginTx(ctx)`
- **Status**: ✅ Fixed in products service (line 81)

**Remaining Work**: Services using unmigrated repositories (Team, Products) will have type mismatches until those repositories are migrated. The Products service is partially affected:
- `i.products.FindImage()` - Still expects `*gorm.DB`
- `i.products.DeleteImageById()` - Still expects `*gorm.DB`

**Resolution**: These repositories are commented out in `repository/module.go`, so services won't be instantiated until migration is complete.

### Casbin Adapter
The project uses `casbin/gorm-adapter/v3` which depends on GORM. Options:
1. Keep GORM adapter (uses separate connection) ✅ **Current choice**
2. Migrate to sqlx-based adapter
3. Create custom adapter

**Current Decision**: Keep GORM adapter for now (isolated from main DB operations)

## 📝 Next Steps

### Immediate (This Sprint)
1. ✅ Complete user repository migration
2. ✅ Complete database transaction helper
3. ✅ Complete image repository
4. ✅ Migrate team repository
5. ✅ Migrate product repositories

### Short Term (Next Sprint)
1. Update service layer for new transaction interface
2. Migrate remaining repositories
3. Update and fix tests
4. Remove GORM dependencies (except Casbin adapter)

### Long Term
1. Consider migrating Casbin adapter
2. Performance benchmarking (GORM vs sqlx)
3. Update API documentation

## 🧪 Testing Strategy

### For Each Migrated Repository
1. **Unit Tests**: Update repository tests to use sqlx
2. **Integration Tests**: Test with real database
3. **Manual Testing**: Verify via API endpoints

### Test Checklist
- [ ] User registration and authentication
- [ ] Team CRUD operations
- [ ] Product management
- [ ] Image upload and retrieval
- [ ] Transaction rollback scenarios

## 📚 Reference Documentation

- **Migration Guide**: `/Users/suttapak/development/starter/apps/backend/MIGRATION_GUIDE.md`
- **sqlx Docs**: https://jmoiron.github.io/sqlx/
- **Goose Docs**: https://pressly.github.io/goose/
- **Example**: `internal/repository/user.go`

## 🔧 Development Commands

```bash
# Create new migration
make new-migrate name=feature_name

# Run migrations (auto on startup, or manual)
make migrate-up DB_DSN="..."

# Check status
make migrate-status DB_DSN="..."

# Start dev server (runs migrations automatically)
cd apps/backend && air
```

## 👥 Team Notes

### For Developers Adding New Features
- **Always use sqlx** for new repository code
- **Follow the pattern** in `repository/user.go`
- **Create migrations** with `make new-migrate name=...`
- **Remember soft deletes**: Always include `AND deleted_at IS NULL`

### For Code Reviewers
- Check that `db` tags are used (not `gorm`)
- Verify SQL injection protection (parameterized queries)
- Ensure soft delete filters are present
- Confirm transaction support via `*sqlx.Tx` parameter

## 🐛 Troubleshooting

### "Method signature mismatch"
Check that you've updated the interface to use `*sqlx.Tx` instead of `*gorm.DB`.

### "No rows affected"
Make sure to include `AND deleted_at IS NULL` in WHERE clauses.

### "Cannot scan NULL"
Use `sql.NullString`, `sql.NullTime`, etc. for nullable fields.

### "Transaction deadlock"
Ensure transactions are properly committed or rolled back in defer blocks.

## 📞 Contact

For questions or issues with the migration:
- Check `MIGRATION_GUIDE.md` for detailed patterns
- Review `repository/user.go` for working examples
- Create an issue with the `migration` label

---

**Migration started**: October 9, 2025
**Target completion**: TBD based on remaining repository complexity
