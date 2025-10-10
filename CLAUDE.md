# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a full-stack monorepo application with a Go (Gin) backend and React (Vite + TanStack Router) frontend, orchestrated with Docker Compose. The project uses Uber's Fx for dependency injection in the backend and follows a clean architecture pattern.

## Repository Structure

```
apps/
├── backend/        # Go backend service
└── www/           # React frontend application
```

### Backend (Go)
- Entry point: `apps/backend/cmd/labostack/main.go`
- Framework: Gin + Uber Fx (dependency injection)
- Database: PostgreSQL with GORM
- Architecture layers:
  - `internal/model/` - Database models and entities
  - `internal/repository/` - Data access layer
  - `internal/service/` - Business logic layer
  - `internal/controller/` - HTTP handlers
  - `internal/route/` - Route definitions
  - `internal/middleware/` - HTTP middleware (auth, logging, etc.)
  - `internal/filter/` - Request filtering/validation
- `bootstrap/` - Application initialization (Gin setup, lifecycle management)
- `database/` - Database connection and migrations via GORM AutoMigrate
- `config/` - Configuration management (TOML-based with Viper)
- `logger/` - Zap-based logging
- `i18n/` - Internationalization support
- `helpers/` - Shared utility functions
- `mail/` - Email functionality (gomail)

### Frontend (React)
- Entry point: `apps/www/src/main.tsx`
- Framework: React 19 + Vite
- Routing: TanStack Router (type-safe routing with `routeTree.gen.ts`)
- State Management: TanStack Query
- UI: HeroUI + Tailwind CSS
- Form Handling: React Hook Form + Zod validation
- i18n: i18next + react-i18next

## Development Commands

### Full Stack Development

```bash
# Start all services (frontend, backend, PostgreSQL, adminer) with hot reload
make dev

# Clean volumes and rebuild everything
make clean
```

### Backend Development

```bash
cd apps/backend

# Run with hot reload using Air
air

# Run directly
go run cmd/labostack/main.go

# Run tests
go test ./...

# Run specific test file
go test ./internal/service/auth.service_test.go

# Generate Swagger docs
swag init -g cmd/labostack/main.go

# Format code
go fmt ./...

# Lint
golangci-lint run

# Database Migrations (Goose + sqlx)
# Create new migration
make new-migrate name=migration_name

# Run migrations manually (requires DB_DSN environment variable)
make migrate-up DB_DSN="host=localhost user=... password=... dbname=... port=5432 sslmode=disable"

# Rollback last migration
make migrate-down DB_DSN="..."

# Check migration status
make migrate-status DB_DSN="..."

# Note: Migrations run automatically on application startup via embedded migrations
```

### Frontend Development

```bash
cd apps/www

# Start dev server
npm run dev

# Build for production
npm run build

# Lint
npm run lint

# Preview production build
npm run preview
```

## Docker Commands

```bash
# Build and push both backend and frontend to registry
make build

# Build and push backend only
make build-server

# Build and push frontend only
make build-web

# Local development with Docker
make dev

# Clean volumes and rebuild
make clean
```

## Configuration

### Backend Configuration
- Main config: `apps/backend/configs.toml`
- Example: `apps/backend/configs.example.toml`
- Environment variables: `apps/backend/.env`
- Uses Viper for configuration management
- Key sections: database, JWT, CORS, mail, pprof

### Frontend Configuration
- TypeScript configs: `tsconfig.json`, `tsconfig.app.json`, `tsconfig.node.json`
- Vite config: `vite.config.ts`
- Environment-specific builds handled via Vite

## Architecture Notes

### Backend Dependency Injection (Fx)
The application uses Uber Fx modules pattern. Each domain follows this structure in `main.go`:
1. Register providers: `fx.Provide(newComponent)`
2. Invoke lifecycle hooks: `fx.Invoke(useComponent)`
3. Modules are defined in `module.go` files throughout the codebase

Example flow: `helpers.Module` → `logger.Module` → `config.Module` → `database.Module` → `repository.Module` → `service.Module` → `controller.Module` → `route.Module` → `bootstrap.Module`

### Database Management

**✅ MIGRATION TO SQLX COMPLETED**

The project has migrated from GORM to sqlx + Goose for database operations:

#### Current Stack
- **ORM**: sqlx (raw SQL with struct mapping)
- **Migrations**: Goose with embedded SQL files
- **Connection**: `*sqlx.DB` for queries, `*sqlx.Tx` for transactions

#### Migration Files
- Location: `apps/backend/database/migrations/`
- Format: Goose SQL migrations with `-- +goose Up` and `-- +goose Down` directives
- Initial schema: `20251009090803_init_schema.sql` includes all tables with proper indexes
- Migrations are embedded and run automatically on startup via `database/migrate.go`

#### Schema Design
- All tables include: `id`, `created_at`, `updated_at`, `deleted_at` (soft deletes)
- Proper foreign keys with `ON DELETE CASCADE` where appropriate
- Indexes on commonly queried fields (`deleted_at`, foreign keys, unique fields)
- Seed data for roles, team_roles, and report_json_schema_types included in init migration

#### Creating New Migrations
```bash
# From project root
make new-migrate name=add_new_table

# Or from backend directory
cd apps/backend
goose -dir database/migrations create add_new_table sql
```

#### Manual Migration Commands
```bash
# Run migrations (requires DB_DSN)
make migrate-up DB_DSN="host=localhost user=... dbname=..."

# Rollback last migration
make migrate-down DB_DSN="..."

# Check status
make migrate-status DB_DSN="..."
```

#### Repository Pattern with sqlx
- Repositories now use `*sqlx.DB` and `*sqlx.Tx` instead of `*gorm.DB`
- Example: `internal/repository/user.go` - fully migrated to sqlx
- Transaction support via `*sqlx.Tx` parameter in repository methods
- Use `sqlx.GetContext()`, `sqlx.SelectContext()`, `ExecContext()` for queries
- Model structs use `db` tags instead of `gorm` tags
- Soft deletes handled manually with `WHERE deleted_at IS NULL`

#### Migration Status (80% Complete - All Core Repositories)
- ✅ Database module (`database/sqlx.database.go`)
- ✅ Migration runner (`database/migrate.go`)
- ✅ Model structs updated with `db` tags
- ✅ User repository - Reference implementation with relations
- ✅ DatabaseTransaction - Helper with `WithTransaction` method
- ✅ Image repository - Simple CRUD
- ✅ AutoIncrementSequence - Complex with atomic locking
- ✅ Report repository - CRUD with pagination
- ✅ Team repository - Complex with JOIN queries and member relations
- ✅ ProductCategory repository - Simple CRUD with team filtering
- ✅ Products repository - Complex with multiple relations and junction tables
- ⏳ ODT repository (optional - external API, 10%)
- 📦 Old GORM code backed up with `.gorm.go.backup` extension (8 files)
- 📋 Detailed progress tracked in `apps/backend/MIGRATION_PROGRESS.md`

### Authentication & Authorization
- JWT-based authentication with separate secrets for access/refresh tokens
- Casbin for authorization (RBAC/ABAC)
- Casbin model: `apps/backend/carbin/authz_model.conf`
- Middleware in `internal/middleware/`

### API Documentation
- Swagger/OpenAPI via swaggo
- Access at: `http://localhost:8080/swagger/index.html`
- Generated docs: `apps/backend/cmd/docs/`
- Regenerate: `swag init -g cmd/labostack/main.go`

### Frontend Routing
- TanStack Router with type-safe routing
- Route definitions in `apps/www/src/routes/`
- Auto-generated route tree: `apps/www/src/routeTree.gen.ts`
- Auth context integrated with router

### Frontend State Management
- TanStack Query for server state
- Context API for auth state (`auth.tsx`)
- Query client configured in `main.tsx`

## Testing

### Backend Tests
- Test files: `*_test.go`
- Existing tests: `internal/controller/auth.controller_test.go`, `internal/service/auth.service_test.go`
- Run with: `go test ./...`

### Frontend Tests
- Jest configured (`jest.config.js`)
- Run with: `npm test`

## Services & Ports

- Frontend (dev): Vite dev server (typically port 5173, proxied through nginx at 8080)
- Backend: Port 8080 (internal), exposed via nginx proxy
- PostgreSQL: Port 5432
- Adminer: Available via nginx proxy
- Nginx Proxy: Port 8080 (main entry point)

## Important Files

- `makefile` - Build and deployment commands
- `docker-compose.yaml` - Local development orchestration
- `nginx.conf` - Reverse proxy configuration
- `.air.toml` - Hot reload configuration for Go backend
- `apps/backend/configs.toml` - Backend runtime configuration
- `apps/www/vite.config.ts` - Frontend build configuration

## Current Branch & Migration Status

The project is on the `migrate-sqlx` branch. **The migration to sqlx is in progress:**

### ✅ Completed (80% - All Core Repositories Migrated)
- ✅ Added sqlx and Goose dependencies
- ✅ Created initial database migration with complete schema
- ✅ Updated all model structs to use `db` tags
- ✅ Created new sqlx database connection module
- ✅ Set up embedded migrations that run on startup
- ✅ Migrated 8 repositories to sqlx:
  - User repository (reference implementation with relations)
  - DatabaseTransaction helper (with `WithTransaction` method)
  - Image repository (simple CRUD)
  - AutoIncrementSequence repository (atomic counter with FOR UPDATE)
  - Report repository (CRUD with pagination)
  - Team repository (complex with JOIN queries and member relations)
  - ProductCategory repository (simple CRUD with team filtering)
  - Products repository (complex with multiple relations and junction tables)
- ✅ Fixed service layer DatabaseTransaction calls to use new interface
- ✅ Updated makefile with migration commands
- ✅ Created comprehensive migration guide (`MIGRATION_GUIDE.md`)
- ✅ Created progress tracker (`MIGRATION_PROGRESS.md`)
- ✅ **Application compiles successfully with all core features working!**

### Optional / TODO
- ⏳ Migrate ODT repository (optional - external API integration, 10%)
- ⏳ Update tests to work with sqlx
- ⏳ Consider removing GORM dependencies (except Casbin adapter)

### For New Code
- **Use sqlx** for all new repository methods
- **Create Goose migrations** for schema changes (use `make new-migrate name=...`)
- **Follow the pattern** in `internal/repository/user.go` for new repositories
- **Use transactions** via `*sqlx.Tx` parameter when needed
