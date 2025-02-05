package repository

import (
	"context"

	"github.com/suttapak/starter/internal/idx"
	"github.com/suttapak/starter/internal/model"
	"gorm.io/gorm"
)

type (
	Team interface {
		CommonDB
		Create(ctx context.Context, tx *gorm.DB, ownerId uint, params CreateTeamParams) (res *model.Team, err error)
		GetTeamsMe(ctx context.Context, tx *gorm.DB, userId uint) (res []model.Team, err error)
		CheckTeamUsernameIsExist(ctx context.Context, tx *gorm.DB, teamUsername string) (isExist bool, err error)
	}
	team struct {
		db *gorm.DB
	}
	CreateTeamParams struct {
		Name        string `json:"name"`
		Username    string `json:"username"`
		Description string `json:"description"`
	}
)

// CheckTeamUsernameIsExist implements Team.
func (t team) CheckTeamUsernameIsExist(ctx context.Context, tx *gorm.DB, teamUsername string) (exists bool, err error) {
	if tx == nil {
		tx = t.db
	}

	err = tx.Raw("SELECT EXISTS (SELECT 1 FROM teams WHERE username = ? LIMIT 1)", teamUsername).
		Scan(&exists).Error

	return exists, err
}

// GetTeamsMe implements Team.
func (t team) GetTeamsMe(ctx context.Context, tx *gorm.DB, userId uint) (res []model.Team, err error) {
	if tx == nil {
		tx = t.db
	}

	err = tx.Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ?", userId).
		Find(&res).Error

	return res, err
}

// Create implements Team.
func (t team) Create(ctx context.Context, tx *gorm.DB, ownerId uint, params CreateTeamParams) (res *model.Team, err error) {
	m := model.Team{
		Name:        params.Name,
		Username:    params.Username,
		Description: params.Description,
		TeamMembers: []model.TeamMember{{
			UserID:     ownerId,
			TeamRoleID: idx.TeamRoleOwnerID,
		}},
	}
	if tx == nil {
		tx = t.db
	}
	err = tx.Create(&m).Error
	return
}

// BeginTx implements Team.
func (t team) BeginTx() *gorm.DB {
	return t.db.Begin()
}

// CommitTx implements Team.
func (t team) CommitTx(tx *gorm.DB) {
	if tx == nil {
		return
	}
	tx.Commit()
}

// RollbackTx implements Team.
func (t team) RollbackTx(tx *gorm.DB) {
	if tx == nil {
		return
	}
	tx.Rollback()
}

func newTeam(db *gorm.DB) Team {
	return team{db: db}
}
