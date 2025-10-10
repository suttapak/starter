package repository

import (
	"context"

	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/idx"
	"github.com/suttapak/starter/internal/model"
	"gorm.io/gorm"
)

type (
	Team interface {
		Create(ctx context.Context, tx *gorm.DB, ownerId uint, params CreateTeamParams) (res *model.Team, err error)
		FindByUserId(ctx context.Context, tx *gorm.DB, pg *helpers.Pagination, userId uint) (res []model.Team, err error)
		ExistByUsername(ctx context.Context, tx *gorm.DB, teamUsername string) (isExist bool, err error)
		CountMemberByTeamId(ctx context.Context, tx *gorm.DB, teamId uint) (count int64, err error)
		CountPendingMemberByTeamId(ctx context.Context, tx *gorm.DB, teamId uint) (count int64, err error)
		FindAllMemberByTeamId(ctx context.Context, tx *gorm.DB, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) (res []model.TeamMember, err error)
		FindAllPendingMemberByTeamId(ctx context.Context, tx *gorm.DB, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) (res []model.TeamMember, err error)
		FindByTeamId(ctx context.Context, tx *gorm.DB, teamId uint) (res *model.Team, err error)
		FindByTeamIdAndUserId(ctx context.Context, tx *gorm.DB, teamId, userId uint) (res *model.TeamMember, err error)
		UpdateMemberRole(ctx context.Context, tx *gorm.DB, teamId, userId, roleId uint) (err error)
		CreatePendingMember(ctx context.Context, tx *gorm.DB, teamId, userId uint) (res *model.TeamMember, err error)
		ExistUserInTeamByTeamId(ctx context.Context, tx *gorm.DB, teamId, userId uint) (isExist bool, err error)
		CreateMemberByUserId(ctx context.Context, tx *gorm.DB, teamId, userId uint) (err error)
		FindAll(ctx context.Context, tx *gorm.DB, pg *helpers.Pagination, f *filter.TeamFilter) (res []model.Team, err error)
		FindMemberByTeamIdAndUserId(ctx context.Context, tx *gorm.DB, teamId, userId uint) (res *model.TeamMember, err error)
		AcceptTeamMember(ctx context.Context, tx *gorm.DB, teamId, userId, roleId uint) (err error)

		Update(ctx context.Context, tx *gorm.DB, teamId uint, params UpdateTeamInfoParams) (err error)
		FindAllEmailOfTeamAdminAndOwner(ctx context.Context, tx *gorm.DB, teamId uint) ([]FindAllTeamAdminAndOwnerResponse, error)
	}
	team struct {
		db *gorm.DB
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
		Email string `database:"email"`
	}
)

const sqlFindAllTeamAdminAndOwner = `
SELECT
  u.email as email
FROM
  users AS u
WHERE
  u.id in (
    SELECT
      user_id
    FROM
      team_members AS tm
    WHERE
      tm.team_role_id IN (?,?) AND tm.team_id = ?
  )
`

// FindAllTeamAdminAndOwner implements Transaction.
func (i *team) FindAllEmailOfTeamAdminAndOwner(ctx context.Context, tx *gorm.DB, teamId uint) ([]FindAllTeamAdminAndOwnerResponse, error) {
	if tx == nil {
		tx = i.db
	}
	res := []FindAllTeamAdminAndOwnerResponse{}
	tx = tx.WithContext(ctx)
	err := tx.Raw(
		sqlFindAllTeamAdminAndOwner,
		idx.TeamRoleOwnerID,
		idx.TeamRoleAdminID,
		teamId,
	).
		Scan(&res).Error
	return res, err
}

// Update implements Team.
func (t *team) Update(ctx context.Context, tx *gorm.DB, teamId uint, params UpdateTeamInfoParams) (err error) {
	team := model.Team{
		Name:        params.Name,
		Username:    params.Username,
		Address:     params.Address,
		Phone:       params.Phone,
		Email:       params.Email,
		Description: params.Description,
	}
	if tx == nil {
		tx = t.db
	}
	return tx.Select("name", "username", "address", "phone", "email", "description").Where("id = ?", teamId).Updates(&team).Error
}

// AcceptTeamMember implements Team.
func (t *team) AcceptTeamMember(ctx context.Context, tx *gorm.DB, teamId uint, userId uint, roleId uint) (err error) {
	if tx == nil {
		tx = t.db
	}
	member := model.TeamMember{
		TeamID:     teamId,
		UserID:     userId,
		TeamRoleID: roleId,
		IsActive:   true,
	}
	return tx.Where("team_id = ? and user_id = ?", teamId, userId).Updates(&member).Error
}

// FindMemberByTeamIdAndUserId implements Team.
func (t *team) FindMemberByTeamIdAndUserId(ctx context.Context, tx *gorm.DB, teamId uint, userId uint) (res *model.TeamMember, err error) {
	if tx == nil {
		tx = t.db
	}
	err = tx.Model(&model.TeamMember{}).
		Where("team_id = ? and user_id = ?", teamId, userId).
		First(&res).Error
	return
}

// FindAll implements Team.
func (t *team) FindAll(ctx context.Context, tx *gorm.DB, pg *helpers.Pagination, f *filter.TeamFilter) (res []model.Team, err error) {
	if tx == nil {
		tx = t.db
	}
	query := tx.Model(&model.Team{})
	if f != nil && f.Name != "" {

		query = query.Where("name ILIKE ?", f.Name+"%")
	}
	query.Count(&pg.Count)
	helpers.Paging(pg)
	err = query.Offset(pg.Offset).
		Limit(pg.Limit).
		Find(&res).Error
	return

}

// CreateMemberByUserId implements Team.
func (t *team) CreateMemberByUserId(ctx context.Context, tx *gorm.DB, teamId uint, userId uint) (err error) {
	if tx == nil {
		tx = t.db
	}
	member := model.TeamMember{
		TeamID:     teamId,
		UserID:     userId,
		TeamRoleID: idx.TeamRoleMemberID,
		IsActive:   true,
	}
	return tx.Create(&member).Error
}

func (t *team) ExistUserInTeamByTeamId(ctx context.Context, tx *gorm.DB, teamId uint, userId uint) (bool, error) {
	if tx == nil {
		tx = t.db
	}

	var exists bool
	err := tx.Model(&model.TeamMember{}).
		Select("COUNT(1) > 0").
		Where("team_id = ? AND user_id = ?", teamId, userId).
		Find(&exists).Error

	return exists, err
}

// CreatePendingMember adds a user to a team as a pending member.
func (t *team) CreatePendingMember(ctx context.Context, tx *gorm.DB, teamId uint, userId uint) (*model.TeamMember, error) {
	if tx == nil {
		tx = t.db
	}

	member := &model.TeamMember{
		TeamID:     teamId,
		UserID:     userId,
		TeamRoleID: idx.TeamRoleMemberID, // Corrected variable name
		IsActive:   false,
	}

	if err := tx.Create(member).Error; err != nil {
		return nil, err
	}

	return member, nil
}

// UpdateMemberRole implements Team.
func (t *team) UpdateMemberRole(ctx context.Context, tx *gorm.DB, teamId uint, userId uint, roleId uint) (err error) {
	if tx == nil {
		tx = t.db
	}
	err = tx.Model(&model.TeamMember{}).
		Where("team_id = ? and user_id = ?", teamId, userId).
		Update("team_role_id", roleId).Error
	return
}

// FindByTeamIdAndUserId implements Team.
func (t *team) FindByTeamIdAndUserId(ctx context.Context, tx *gorm.DB, teamId, userId uint) (res *model.TeamMember, err error) {
	if tx == nil {
		tx = t.db
	}
	err = tx.Model(&model.TeamMember{}).Where("team_id = ? and user_id = ?", teamId, userId).Preload("TeamRole").Preload("User").First(&res).Error
	return
}

// FindAllPendingMemberByTeamId implements Team.
func (t *team) FindAllPendingMemberByTeamId(ctx context.Context, tx *gorm.DB, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) (res []model.TeamMember, err error) {
	if tx == nil {
		tx = t.db
	}
	query := tx.Model(&model.TeamMember{}).
		Where("team_id = ? and is_active = false", teamId)
	if f != nil && f.Username != "" {
		query = query.Where("user_id IN (SELECT id FROM users WHERE username LIKE ?)", f.Username+"%")
	}
	query.Count(&pg.Count)
	helpers.Paging(pg)
	err = query.Offset(pg.Offset).
		Limit(pg.Limit).
		Preload("User").
		Preload("TeamRole").
		Find(&res).Error
	return
}

// CountPendingMemberByTeamId implements Team.
func (t *team) CountPendingMemberByTeamId(ctx context.Context, tx *gorm.DB, teamId uint) (count int64, err error) {
	if tx == nil {
		tx = t.db
	}
	err = tx.Model(&model.TeamMember{}).
		Where("team_id = ? and is_active = false", teamId).
		Count(&count).Error
	return
}

// FindByTeamId implements Team.
func (t *team) FindByTeamId(ctx context.Context, tx *gorm.DB, teamId uint) (res *model.Team, err error) {
	if tx == nil {
		tx = t.db
	}
	err = tx.WithContext(ctx).Where("id = ?", teamId).First(&res).Error
	return
}

// FindAllMemberByTeamId implements Team.
func (t *team) FindAllMemberByTeamId(ctx context.Context, tx *gorm.DB, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) (res []model.TeamMember, err error) {
	if tx == nil {
		tx = t.db
	}
	query := tx.Model(&model.TeamMember{}).
		Where("team_id = ? and is_active = true", teamId)
	if f != nil && f.Username != "" {
		query = query.Where("user_id IN (SELECT id FROM users WHERE username LIKE ?)", f.Username+"%")
	}
	query.Count(&pg.Count)
	helpers.Paging(pg)
	err = query.Offset(pg.Offset).
		Limit(pg.Limit).
		Preload("User").
		Preload("TeamRole").
		Find(&res).Error
	return
}

// CountMemberByTeamId implements Team.
func (t *team) CountMemberByTeamId(ctx context.Context, tx *gorm.DB, teamId uint) (count int64, err error) {
	if tx == nil {
		tx = t.db
	}
	err = tx.Model(&model.TeamMember{}).
		Where("team_id = ? and is_active = true", teamId).
		Count(&count).Error
	return
}

// ExistByUsername implements Team.
func (t *team) ExistByUsername(ctx context.Context, tx *gorm.DB, teamUsername string) (exists bool, err error) {
	if tx == nil {
		tx = t.db
	}

	err = tx.Raw("SELECT EXISTS (SELECT 1 FROM teams WHERE username = ? LIMIT 1)", teamUsername).
		Scan(&exists).Error

	return exists, err
}

// FindByUserId implements Team.
func (t *team) FindByUserId(ctx context.Context, tx *gorm.DB, pg *helpers.Pagination, userId uint) (res []model.Team, err error) {
	if tx == nil {
		tx = t.db
	}

	tx = tx.Model(&model.Team{}).Joins("JOIN team_members ON team_members.team_id = teams.id").
		Where("team_members.user_id = ? and team_members.is_active = true ", userId).
		Count(&pg.Count)
	helpers.Paging(pg)
	err = tx.Find(&res).Error

	return res, err
}

// Create implements Team.
func (t *team) Create(ctx context.Context, tx *gorm.DB, ownerId uint, params CreateTeamParams) (res *model.Team, err error) {
	m := model.Team{
		Name:        params.Name,
		Username:    params.Username,
		Description: params.Description,
		Address:     params.Address,
		Phone:       params.Phone,
		Email:       params.Email,
		TeamMembers: []model.TeamMember{{
			UserID:     ownerId,
			TeamRoleID: idx.TeamRoleOwnerID,
			IsActive:   true,
		}},
	}
	if tx == nil {
		tx = t.db
	}
	err = tx.Create(&m).Error
	return
}

func newTeam(db *gorm.DB) Team {
	return &team{db: db}
}
