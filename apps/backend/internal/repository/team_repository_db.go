package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/idx"
	"github.com/suttapak/starter/internal/model"
)

type (
	Team interface {
		Create(ctx context.Context, tx *sqlx.Tx, ownerId uint, params CreateTeamParams) (res *model.Team, err error)
		FindByUserId(ctx context.Context, tx *sqlx.Tx, pg *helpers.Pagination, userId uint) (res []model.Team, err error)
		ExistByUsername(ctx context.Context, tx *sqlx.Tx, teamUsername string) (isExist bool, err error)
		CountMemberByTeamId(ctx context.Context, tx *sqlx.Tx, teamId uint) (count int64, err error)
		CountPendingMemberByTeamId(ctx context.Context, tx *sqlx.Tx, teamId uint) (count int64, err error)
		FindAllMemberByTeamId(ctx context.Context, tx *sqlx.Tx, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) (res []model.TeamMember, err error)
		FindAllPendingMemberByTeamId(ctx context.Context, tx *sqlx.Tx, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) (res []model.TeamMember, err error)
		FindByTeamId(ctx context.Context, tx *sqlx.Tx, teamId uint) (res *model.Team, err error)
		FindByTeamIdAndUserId(ctx context.Context, tx *sqlx.Tx, teamId, userId uint) (res *model.TeamMember, err error)
		UpdateMemberRole(ctx context.Context, tx *sqlx.Tx, teamId, userId, roleId uint) (err error)
		CreatePendingMember(ctx context.Context, tx *sqlx.Tx, teamId, userId uint) (res *model.TeamMember, err error)
		ExistUserInTeamByTeamId(ctx context.Context, tx *sqlx.Tx, teamId, userId uint) (isExist bool, err error)
		CreateMemberByUserId(ctx context.Context, tx *sqlx.Tx, teamId, userId uint) (err error)
		FindAll(ctx context.Context, tx *sqlx.Tx, pg *helpers.Pagination, f *filter.TeamFilter) (res []model.Team, err error)
		FindMemberByTeamIdAndUserId(ctx context.Context, tx *sqlx.Tx, teamId, userId uint) (res *model.TeamMember, err error)
		AcceptTeamMember(ctx context.Context, tx *sqlx.Tx, teamId, userId, roleId uint) (err error)
		Update(ctx context.Context, tx *sqlx.Tx, teamId uint, params UpdateTeamInfoParams) (err error)
		FindAllEmailOfTeamAdminAndOwner(ctx context.Context, tx *sqlx.Tx, teamId uint) ([]FindAllTeamAdminAndOwnerResponse, error)
	}

	teamSqlx struct {
		db *sqlx.DB
	}

	CreateTeamParams struct {
		Name        string `json:"name"`
		Username    string `json:"username"`
		Address     string `json:"address"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Description string `json:"description"`
	}

	UpdateTeamInfoParams struct {
		Name        string `json:"name"`
		Username    string `json:"username"`
		Address     string `json:"address"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Description string `json:"description"`
	}

	FindAllTeamAdminAndOwnerResponse struct {
		Email string `db:"email"`
	}
)

func NewTeam(db *sqlx.DB) Team {
	return &teamSqlx{db: db}
}

// getDB returns the appropriate database connection (transaction or main DB)
func (t *teamSqlx) getDB(tx *sqlx.Tx) sqlx.ExtContext {
	if tx != nil {
		return tx
	}
	return t.db
}

// Create creates a new team with the owner as the first member
func (t *teamSqlx) Create(ctx context.Context, tx *sqlx.Tx, ownerId uint, params CreateTeamParams) (*model.Team, error) {
	db := t.getDB(tx)

	// Insert team
	teamQuery := `
		INSERT INTO teams (name, username, address, phone, email, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	team := &model.Team{
		Name:        params.Name,
		Username:    params.Username,
		Address:     &params.Address,
		Phone:       &params.Phone,
		Email:       &params.Email,
		Description: &params.Description,
	}

	err := sqlx.GetContext(ctx, db, team, teamQuery,
		params.Name,
		params.Username,
		params.Address,
		params.Phone,
		params.Email,
		params.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	// Insert owner as first member
	memberQuery := `
		INSERT INTO team_members (team_id, user_id, team_role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, true, NOW(), NOW())
	`

	_, err = db.ExecContext(ctx, memberQuery, team.ID, ownerId, idx.TeamRoleOwnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to create team owner member: %w", err)
	}

	return team, nil
}

// FindByUserId retrieves all teams for a specific user
func (t *teamSqlx) FindByUserId(ctx context.Context, tx *sqlx.Tx, pg *helpers.Pagination, userId uint) ([]model.Team, error) {
	db := t.getDB(tx)

	// Count total records
	countQuery := `
		SELECT COUNT(DISTINCT teams.id)
		FROM teams
		JOIN team_members ON team_members.team_id = teams.id
		WHERE team_members.user_id = $1
			AND team_members.is_active = true
			AND teams.deleted_at IS NULL
	`
	if err := sqlx.GetContext(ctx, db, &pg.Count, countQuery, userId); err != nil {
		return nil, fmt.Errorf("failed to count teams: %w", err)
	}

	// Apply pagination
	helpers.Paging(pg)

	// Query with pagination
	query := `
		SELECT DISTINCT teams.id, teams.name, teams.address, teams.phone, teams.email,
			teams.username, teams.description, teams.created_at, teams.updated_at
		FROM teams
		JOIN team_members ON team_members.team_id = teams.id
		WHERE team_members.user_id = $1
			AND team_members.is_active = true
			AND teams.deleted_at IS NULL
		ORDER BY teams.created_at DESC
		LIMIT $2 OFFSET $3
	`

	var results []model.Team
	err := sqlx.SelectContext(ctx, db, &results, query, userId, pg.Limit, pg.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find teams by user: %w", err)
	}

	return results, nil
}

// ExistByUsername checks if a team with the given username exists
func (t *teamSqlx) ExistByUsername(ctx context.Context, tx *sqlx.Tx, teamUsername string) (bool, error) {
	db := t.getDB(tx)

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE username = $1 AND deleted_at IS NULL)`

	err := sqlx.GetContext(ctx, db, &exists, query, teamUsername)
	if err != nil {
		return false, fmt.Errorf("failed to check team username existence: %w", err)
	}

	return exists, nil
}

// CountMemberByTeamId counts active members in a team
func (t *teamSqlx) CountMemberByTeamId(ctx context.Context, tx *sqlx.Tx, teamId uint) (int64, error) {
	db := t.getDB(tx)

	var count int64
	query := `
		SELECT COUNT(*)
		FROM team_members
		WHERE team_id = $1 AND is_active = true AND deleted_at IS NULL
	`

	err := sqlx.GetContext(ctx, db, &count, query, teamId)
	if err != nil {
		return 0, fmt.Errorf("failed to count team members: %w", err)
	}

	return count, nil
}

// CountPendingMemberByTeamId counts pending members in a team
func (t *teamSqlx) CountPendingMemberByTeamId(ctx context.Context, tx *sqlx.Tx, teamId uint) (int64, error) {
	db := t.getDB(tx)

	var count int64
	query := `
		SELECT COUNT(*)
		FROM team_members
		WHERE team_id = $1 AND is_active = false AND deleted_at IS NULL
	`

	err := sqlx.GetContext(ctx, db, &count, query, teamId)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending team members: %w", err)
	}

	return count, nil
}

// FindAllMemberByTeamId retrieves all active members of a team with pagination
func (t *teamSqlx) FindAllMemberByTeamId(ctx context.Context, tx *sqlx.Tx, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) ([]model.TeamMember, error) {
	db := t.getDB(tx)

	// Build count query
	countQuery := `
		SELECT COUNT(*)
		FROM team_members
		WHERE team_id = $1 AND is_active = true AND deleted_at IS NULL
	`
	args := []interface{}{teamId}

	if f != nil && f.Username != "" {
		countQuery += ` AND user_id IN (SELECT id FROM users WHERE username LIKE $2 AND deleted_at IS NULL)`
		args = append(args, f.Username+"%")
	}

	if err := sqlx.GetContext(ctx, db, &pg.Count, countQuery, args...); err != nil {
		return nil, fmt.Errorf("failed to count team members: %w", err)
	}

	// Apply pagination
	helpers.Paging(pg)

	// Build main query
	query := `
		SELECT tm.id, tm.team_id, tm.user_id, tm.team_role_id, tm.is_active,
			tm.created_at, tm.updated_at
		FROM team_members tm
		WHERE tm.team_id = $1 AND tm.is_active = true AND tm.deleted_at IS NULL
	`

	if f != nil && f.Username != "" {
		if len(args) == 2 {
			query += ` AND tm.user_id IN (SELECT id FROM users WHERE username LIKE $2 AND deleted_at IS NULL)`
		}
	}

	query += ` ORDER BY tm.created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, pg.Limit, pg.Offset)

	var members []model.TeamMember
	err := sqlx.SelectContext(ctx, db, &members, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find team members: %w", err)
	}

	// Load relations
	for i := range members {
		if err := t.loadTeamMemberRelations(ctx, db, &members[i]); err != nil {
			return nil, err
		}
	}

	return members, nil
}

// FindAllPendingMemberByTeamId retrieves all pending members of a team with pagination
func (t *teamSqlx) FindAllPendingMemberByTeamId(ctx context.Context, tx *sqlx.Tx, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) ([]model.TeamMember, error) {
	db := t.getDB(tx)

	// Build count query
	countQuery := `
		SELECT COUNT(*)
		FROM team_members
		WHERE team_id = $1 AND is_active = false AND deleted_at IS NULL
	`
	args := []interface{}{teamId}

	if f != nil && f.Username != "" {
		countQuery += ` AND user_id IN (SELECT id FROM users WHERE username LIKE $2 AND deleted_at IS NULL)`
		args = append(args, f.Username+"%")
	}

	if err := sqlx.GetContext(ctx, db, &pg.Count, countQuery, args...); err != nil {
		return nil, fmt.Errorf("failed to count pending team members: %w", err)
	}

	// Apply pagination
	helpers.Paging(pg)

	// Build main query
	query := `
		SELECT tm.id, tm.team_id, tm.user_id, tm.team_role_id, tm.is_active,
			tm.created_at, tm.updated_at
		FROM team_members tm
		WHERE tm.team_id = $1 AND tm.is_active = false AND tm.deleted_at IS NULL
	`

	if f != nil && f.Username != "" {
		if len(args) == 2 {
			query += ` AND tm.user_id IN (SELECT id FROM users WHERE username LIKE $2 AND deleted_at IS NULL)`
		}
	}

	query += ` ORDER BY tm.created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, pg.Limit, pg.Offset)

	var members []model.TeamMember
	err := sqlx.SelectContext(ctx, db, &members, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find pending team members: %w", err)
	}

	// Load relations
	for i := range members {
		if err := t.loadTeamMemberRelations(ctx, db, &members[i]); err != nil {
			return nil, err
		}
	}

	return members, nil
}

// FindByTeamId retrieves a team by ID
func (t *teamSqlx) FindByTeamId(ctx context.Context, tx *sqlx.Tx, teamId uint) (*model.Team, error) {
	db := t.getDB(tx)

	var team model.Team
	query := `
		SELECT id, name, address, phone, email, username, description, created_at, updated_at
		FROM teams
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := sqlx.GetContext(ctx, db, &team, query, teamId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("team not found")
		}
		return nil, fmt.Errorf("failed to find team: %w", err)
	}

	return &team, nil
}

// FindByTeamIdAndUserId retrieves a team member by team ID and user ID with relations
func (t *teamSqlx) FindByTeamIdAndUserId(ctx context.Context, tx *sqlx.Tx, teamId, userId uint) (*model.TeamMember, error) {
	db := t.getDB(tx)

	var member model.TeamMember
	query := `
		SELECT id, team_id, user_id, team_role_id, is_active, created_at, updated_at
		FROM team_members
		WHERE team_id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	err := sqlx.GetContext(ctx, db, &member, query, teamId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("team member not found")
		}
		return nil, fmt.Errorf("failed to find team member: %w", err)
	}

	// Load relations
	if err := t.loadTeamMemberRelations(ctx, db, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// UpdateMemberRole updates a team member's role
func (t *teamSqlx) UpdateMemberRole(ctx context.Context, tx *sqlx.Tx, teamId, userId, roleId uint) error {
	db := t.getDB(tx)

	query := `
		UPDATE team_members
		SET team_role_id = $1, updated_at = NOW()
		WHERE team_id = $2 AND user_id = $3 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, query, roleId, teamId, userId)
	if err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("team member not found")
	}

	return nil
}

// CreatePendingMember adds a user to a team as a pending member
func (t *teamSqlx) CreatePendingMember(ctx context.Context, tx *sqlx.Tx, teamId, userId uint) (*model.TeamMember, error) {
	db := t.getDB(tx)

	query := `
		INSERT INTO team_members (team_id, user_id, team_role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, false, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	member := &model.TeamMember{
		TeamID:     teamId,
		UserID:     userId,
		TeamRoleID: idx.TeamRoleMemberID,
		IsActive:   false,
	}

	err := sqlx.GetContext(ctx, db, member, query, teamId, userId, idx.TeamRoleMemberID)
	if err != nil {
		return nil, fmt.Errorf("failed to create pending member: %w", err)
	}

	return member, nil
}

// ExistUserInTeamByTeamId checks if a user exists in a team
func (t *teamSqlx) ExistUserInTeamByTeamId(ctx context.Context, tx *sqlx.Tx, teamId, userId uint) (bool, error) {
	db := t.getDB(tx)

	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM team_members
			WHERE team_id = $1 AND user_id = $2 AND deleted_at IS NULL
		)
	`

	err := sqlx.GetContext(ctx, db, &exists, query, teamId, userId)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence in team: %w", err)
	}

	return exists, nil
}

// CreateMemberByUserId creates an active team member
func (t *teamSqlx) CreateMemberByUserId(ctx context.Context, tx *sqlx.Tx, teamId, userId uint) error {
	db := t.getDB(tx)

	query := `
		INSERT INTO team_members (team_id, user_id, team_role_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, true, NOW(), NOW())
	`

	_, err := db.ExecContext(ctx, query, teamId, userId, idx.TeamRoleMemberID)
	if err != nil {
		return fmt.Errorf("failed to create team member: %w", err)
	}

	return nil
}

// FindAll retrieves all teams with pagination
func (t *teamSqlx) FindAll(ctx context.Context, tx *sqlx.Tx, pg *helpers.Pagination, f *filter.TeamFilter) ([]model.Team, error) {
	db := t.getDB(tx)

	// Build count query
	countQuery := `SELECT COUNT(*) FROM teams WHERE deleted_at IS NULL`
	args := []interface{}{}

	if f != nil && f.Name != "" {
		countQuery += ` AND name ILIKE $1`
		args = append(args, f.Name+"%")
	}

	if err := sqlx.GetContext(ctx, db, &pg.Count, countQuery, args...); err != nil {
		return nil, fmt.Errorf("failed to count teams: %w", err)
	}

	// Apply pagination
	helpers.Paging(pg)

	// Build main query
	query := `
		SELECT id, name, address, phone, email, username, description, created_at, updated_at
		FROM teams
		WHERE deleted_at IS NULL
	`

	if f != nil && f.Name != "" {
		if len(args) == 1 {
			query += ` AND name ILIKE $1`
		}
	}

	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, pg.Limit, pg.Offset)

	var teams []model.Team
	err := sqlx.SelectContext(ctx, db, &teams, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find teams: %w", err)
	}

	return teams, nil
}

// FindMemberByTeamIdAndUserId retrieves a team member by team ID and user ID
func (t *teamSqlx) FindMemberByTeamIdAndUserId(ctx context.Context, tx *sqlx.Tx, teamId, userId uint) (*model.TeamMember, error) {
	db := t.getDB(tx)

	var member model.TeamMember
	query := `
		SELECT id, team_id, user_id, team_role_id, is_active, created_at, updated_at
		FROM team_members
		WHERE team_id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	err := sqlx.GetContext(ctx, db, &member, query, teamId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("team member not found")
		}
		return nil, fmt.Errorf("failed to find team member: %w", err)
	}

	return &member, nil
}

// AcceptTeamMember activates a pending team member
func (t *teamSqlx) AcceptTeamMember(ctx context.Context, tx *sqlx.Tx, teamId, userId, roleId uint) error {
	db := t.getDB(tx)

	query := `
		UPDATE team_members
		SET team_role_id = $1, is_active = true, updated_at = NOW()
		WHERE team_id = $2 AND user_id = $3 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, query, roleId, teamId, userId)
	if err != nil {
		return fmt.Errorf("failed to accept team member: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("team member not found")
	}

	return nil
}

// Update updates team information
func (t *teamSqlx) Update(ctx context.Context, tx *sqlx.Tx, teamId uint, params UpdateTeamInfoParams) error {
	db := t.getDB(tx)

	query := `
		UPDATE teams
		SET name = $1, username = $2, address = $3, phone = $4, email = $5, description = $6, updated_at = NOW()
		WHERE id = $7 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, query,
		params.Name,
		params.Username,
		params.Address,
		params.Phone,
		params.Email,
		params.Description,
		teamId,
	)
	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("team not found")
	}

	return nil
}

// FindAllEmailOfTeamAdminAndOwner retrieves all emails of team admins and owners
func (t *teamSqlx) FindAllEmailOfTeamAdminAndOwner(ctx context.Context, tx *sqlx.Tx, teamId uint) ([]FindAllTeamAdminAndOwnerResponse, error) {
	db := t.getDB(tx)

	query := `
		SELECT u.email as email
		FROM users AS u
		WHERE u.id IN (
			SELECT user_id
			FROM team_members AS tm
			WHERE tm.team_role_id IN ($1, $2)
				AND tm.team_id = $3
				AND tm.deleted_at IS NULL
		) AND u.deleted_at IS NULL
	`

	var results []FindAllTeamAdminAndOwnerResponse
	err := sqlx.SelectContext(ctx, db, &results, query, idx.TeamRoleOwnerID, idx.TeamRoleAdminID, teamId)
	if err != nil {
		return nil, fmt.Errorf("failed to find admin and owner emails: %w", err)
	}

	return results, nil
}

// loadTeamMemberRelations loads related User and TeamRole for a TeamMember
func (t *teamSqlx) loadTeamMemberRelations(ctx context.Context, db sqlx.ExtContext, member *model.TeamMember) error {
	// Load User
	userQuery := `
		SELECT id, username, password, email, email_verifyed, role_id, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	user := &model.User{}
	if err := sqlx.GetContext(ctx, db, user, userQuery, member.UserID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to load user: %w", err)
		}
	} else {
		member.User = user
	}

	// Load TeamRole
	roleQuery := `
		SELECT id, name, created_at, updated_at
		FROM team_roles
		WHERE id = $1 AND deleted_at IS NULL
	`
	role := &model.TeamRole{}
	if err := sqlx.GetContext(ctx, db, role, roleQuery, member.TeamRoleID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to load team role: %w", err)
		}
	} else {
		member.TeamRole = role
	}

	return nil
}
