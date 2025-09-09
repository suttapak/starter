package service

import (
	"context"
	"fmt"

	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/logger"
)

type (
	Team interface {
		Create(ctx context.Context, ownerId uint, input CreateTeamDto) (res *CreateTeamResponse, err error)
		GetTeamsMe(ctx context.Context, pg *helpers.Pagination, userId uint) (res []TeamResponse, err error)
		GetTeamMemberCount(ctx context.Context, teamId uint) (res int64, err error)
		GetTeamMembers(ctx context.Context, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) (res []TeamMemberResponse, err error)
		GetTeamByTeamId(ctx context.Context, teamId uint) (res *TeamResponse, err error)
		GetPendingTeamMemberCount(ctx context.Context, teamId uint) (res int64, err error)
		GetPendingTeamMembers(ctx context.Context, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) (res []TeamMemberResponse, err error)
		GetTeamUserMe(ctx context.Context, teamId, userId uint) (res *TeamMemberResponse, err error)
		UpdateMemberRole(ctx context.Context, teamId uint, input UpdateMemberRoleDto) (err error)
		SendInviteTeamMember(ctx context.Context, teamId uint, input CreateTeamPendingTeamMemberDto) (err error)
		JoinTeamWithToken(ctx context.Context, token string) (err error)
		CreateShearLink(ctx context.Context, teamId uint) (res string, err error)
		JoinWithShearLink(ctx context.Context, token string, userId uint) (err error)
		GetTeamsFilter(ctx context.Context, pg *helpers.Pagination, f *filter.TeamFilter) (res []TeamResponse, err error)
		CreateTeamPendingTeamMember(ctx context.Context, teamId uint, userId uint) (err error)
		AcceptTeamMember(ctx context.Context, teamId uint, input AcceptTeamMemberDto) (err error)

		UpdateTeamInfo(ctx context.Context, teamId uint, input UpdateTeamInfoRequest) error
	}
	team struct {
		logger         logger.AppLogger
		teamRepository repository.Team
		userRepository repository.User
		helper         helpers.Helper
		emailService   Email
		jwt            JWTService
		config         *config.Config
	}
	UpdateTeamInfoRequest struct {
		Name        string `json:"name"`
		Username    string `json:"username"`
		Address     string `json:"address"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Description string `json:"description"`
	}
	CreateTeamResponse struct {
		CommonModel
		Name        string `json:"name"`
		Username    string `json:"username"`
		Address     string `json:"address"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Description string `json:"description"`
	}
	TeamResponse struct {
		CommonModel
		Name        string `json:"name"`
		Username    string `json:"username"`
		Address     string `json:"address"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Description string `json:"description"`
	}
	TeamMemberResponse struct {
		CommonModel
		TeamID     uint             `json:"team_id" gorm:"primaryKey"`
		UserID     uint             `json:"user_id" gorm:"primaryKey"`
		TeamRoleID uint             `json:"team_role_id"`
		User       UserResponse     `json:"user"`
		TeamRole   TeamRoleResponse `json:"team_role"`
		IsActive   bool             `json:"is_active"`
	}
	TeamRoleResponse struct {
		CommonModel
		Name string `json:"name"`
	}
	CreateTeamDto struct {
		Name        string `json:"name"`
		Username    string `json:"username"`
		Address     string `json:"address"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Description string `json:"description"`
	}
	UpdateMemberRoleDto struct {
		UserId uint `json:"user_id"`
		RoleId uint `json:"role_id"`
	}
	CreateTeamPendingTeamMemberDto struct {
		Username string `json:"username"`
	}
	AcceptTeamMemberDto struct {
		UserID uint `json:"user_id"`
		RoleID uint `json:"role_id"`
	}
)

// UpdateTeamInfo implements Team.
func (t *team) UpdateTeamInfo(ctx context.Context, teamId uint, input UpdateTeamInfoRequest) error {
	if err := t.teamRepository.Update(ctx, nil, teamId, repository.UpdateTeamInfoParams(input)); err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}

// AcceptTeamMember implements Team.
func (t *team) AcceptTeamMember(ctx context.Context, teamId uint, input AcceptTeamMemberDto) (err error) {
	// check user is already in team
	isExist, err := t.teamRepository.ExistUserInTeamByTeamId(ctx, nil, teamId, input.UserID)
	if err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	if !isExist {
		return errs.ErrUserAlreadyInTeam
	}
	model, err := t.teamRepository.FindMemberByTeamIdAndUserId(ctx, nil, teamId, input.UserID)
	if err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	if model.IsActive {
		return errs.ErrUserAlreadyInTeam
	}
	if err := t.teamRepository.AcceptTeamMember(ctx, nil, teamId, input.UserID, input.RoleID); err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}

// CreateTeamPendingTeamMember implements Team.
func (t *team) CreateTeamPendingTeamMember(ctx context.Context, teamId uint, userId uint) (err error) {
	// check user is already in team
	isExist, err := t.teamRepository.ExistUserInTeamByTeamId(ctx, nil, teamId, userId)
	if err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	if isExist {
		return errs.ErrUserAlreadyInTeam
	}
	if _, err := t.teamRepository.CreatePendingMember(ctx, nil, teamId, userId); err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return
}

// GetTeamsFilter implements Team.
func (t *team) GetTeamsFilter(ctx context.Context, pg *helpers.Pagination, f *filter.TeamFilter) (res []TeamResponse, err error) {
	models, err := t.teamRepository.FindAll(ctx, nil, pg, f)
	if err != nil {
		t.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	if err := t.helper.ParseJson(models, &res); err != nil {
		t.logger.Error(err)
		return nil, errs.ErrInvalid
	}
	return
}

// CreateShearLink implements Team.
func (t *team) CreateShearLink(ctx context.Context, teamId uint) (res string, err error) {
	token, err := t.jwt.GenerateTeamToken(ctx, teamId, 0)
	if err != nil {
		t.logger.Error(err)
		return "", errs.ErrGenerateJWTFail
	}
	return fmt.Sprintf("%s/team/join-team?token=%s", t.config.SERVER.HOST_NAME, token), nil
}

// JoinWithShearLink implements Team.
func (t *team) JoinWithShearLink(ctx context.Context, token string, userId uint) (err error) {
	resToken, err := t.jwt.GetTeamFormToken(ctx, token)
	if err != nil {
		t.logger.Error(err)
		return errs.ErrGenerateJWTFail
	}
	teamId := resToken.TeamId
	// check user is already in team
	isExist, err := t.teamRepository.ExistUserInTeamByTeamId(ctx, nil, teamId, userId)
	if err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	if isExist {
		return errs.ErrUserAlreadyInTeam
	}
	if err := t.teamRepository.CreateMemberByUserId(ctx, nil, teamId, userId); err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}

// JoinTeamWithToken implements Team.
func (t *team) JoinTeamWithToken(ctx context.Context, token string) (err error) {
	resToken, err := t.jwt.GetTeamFormToken(ctx, token)
	if err != nil {
		t.logger.Error(err)
		return errs.ErrGenerateJWTFail
	}
	userId := resToken.UserId
	teamId := resToken.TeamId
	// check user is already in team
	isExist, err := t.teamRepository.ExistUserInTeamByTeamId(ctx, nil, teamId, userId)
	if err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	if isExist {
		return errs.ErrUserAlreadyInTeam
	}
	if err := t.teamRepository.CreateMemberByUserId(ctx, nil, teamId, userId); err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}

// CreateTeamPendingTeamMember implements Team.
func (t *team) SendInviteTeamMember(ctx context.Context, teamId uint, input CreateTeamPendingTeamMemberDto) (err error) {
	userModel, err := t.userRepository.GetUserByEmailOrUsername(ctx, nil, input.Username)
	if err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	userId := userModel.ID
	// check user is already in team
	isExist, err := t.teamRepository.ExistUserInTeamByTeamId(ctx, nil, teamId, userId)
	if err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	if isExist {
		return errs.ErrUserAlreadyInTeam
	}

	var (
		email    = userModel.Email
		isVerify = userModel.EmailVerifyed
	)
	if !isVerify {
		return errs.ErrEmailNotVerify
	}
	token, err := t.jwt.GenerateTeamToken(ctx, teamId, userId)
	if err != nil {
		t.logger.Error(err)
		return errs.ErrGenerateJWTFail
	}
	mailBody := InviteTeamMemberTemplateDataDto{
		TeamName:     "",
		JoinTeamLink: fmt.Sprintf("%s/api/v1/teams/join?token=%s", t.config.SERVER.HOST_NAME, token),
	}
	if err := t.emailService.NewRequest([]string{email}, "Join Team").ParseInviteTeamMemberTemplate(ctx, &mailBody).SendMail(ctx); err != nil {
		t.logger.Error(err)
		return errs.ErrSendEmail
	}
	return
}

// UpdateMemberRole implements Team.
func (t *team) UpdateMemberRole(ctx context.Context, teamId uint, input UpdateMemberRoleDto) (err error) {
	if err := t.teamRepository.UpdateMemberRole(ctx, nil, teamId, input.UserId, input.RoleId); err != nil {
		t.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}

// GetTeamUserMe implements Team.
func (t *team) GetTeamUserMe(ctx context.Context, teamId uint, userId uint) (res *TeamMemberResponse, err error) {
	model, err := t.teamRepository.FindByTeamIdAndUserId(ctx, nil, teamId, userId)
	if err != nil {
		t.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	if err := t.helper.ParseJson(model, &res); err != nil {
		t.logger.Error(err)
		return nil, err
	}
	return
}

// GetPendingTeamMemberCount implements Team.
func (t *team) GetPendingTeamMemberCount(ctx context.Context, teamId uint) (res int64, err error) {
	res, err = t.teamRepository.CountPendingMemberByTeamId(ctx, nil, teamId)
	if err != nil {
		t.logger.Error(err)
		return 0, errs.HandleGorm(err)
	}
	return
}

// GetPendingTeamMembers implements Team.
func (t *team) GetPendingTeamMembers(ctx context.Context, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) (res []TeamMemberResponse, err error) {
	model, err := t.teamRepository.FindAllPendingMemberByTeamId(ctx, nil, pg, f, teamId)
	if err != nil {
		t.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	if err := t.helper.ParseJson(model, &res); err != nil {
		t.logger.Error(err)
		return nil, err
	}
	return

}

// GetTeamByTeamId implements Team.
func (t *team) GetTeamByTeamId(ctx context.Context, teamId uint) (res *TeamResponse, err error) {
	model, err := t.teamRepository.FindByTeamId(ctx, nil, teamId)
	if err != nil {
		t.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	if err := t.helper.ParseJson(model, &res); err != nil {
		t.logger.Error(err)
		return nil, err
	}
	return

}

// GetTeamMembers implements Team.
func (t *team) GetTeamMembers(ctx context.Context, pg *helpers.Pagination, f *filter.TeamMemberFilter, teamId uint) (res []TeamMemberResponse, err error) {
	model, err := t.teamRepository.FindAllMemberByTeamId(ctx, nil, pg, f, teamId)
	if err != nil {
		t.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	if err := t.helper.ParseJson(model, &res); err != nil {
		t.logger.Error(err)
		return nil, err
	}
	return

}

// GetTeamMemberCount implements Team.
func (t *team) GetTeamMemberCount(ctx context.Context, teamId uint) (res int64, err error) {
	res, err = t.teamRepository.CountMemberByTeamId(ctx, nil, teamId)
	if err != nil {
		t.logger.Error(err)
		return 0, errs.HandleGorm(err)
	}
	return
}

// GetTeamsMe implements Team.
func (t *team) GetTeamsMe(ctx context.Context, pg *helpers.Pagination, userId uint) (res []TeamResponse, err error) {
	models, err := t.teamRepository.FindByUserId(ctx, nil, pg, userId)
	if err != nil {
		t.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	if err := t.helper.ParseJson(models, &res); err != nil {
		t.logger.Error(err)
		return nil, errs.ErrInvalid
	}
	return
}

// Create implements Team.
func (t *team) Create(ctx context.Context, ownerId uint, input CreateTeamDto) (res *CreateTeamResponse, err error) {
	// check with username is exist
	exists, err := t.teamRepository.ExistByUsername(ctx, nil, input.Username)
	if err != nil {
		t.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	if exists {
		t.logger.Error(err)
		return nil, errs.ErrTeamUsernameIsUsed
	}
	// create team with role id

	model, err := t.teamRepository.Create(ctx, nil, ownerId, repository.CreateTeamParams(input))
	if err != nil {
		t.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	if err := t.helper.ParseJson(model, &res); err != nil {
		t.logger.Error(err)
		return nil, errs.ErrInvalid
	}
	// make json response
	return
}

func NewTeam(
	logger logger.AppLogger,
	teamRepository repository.Team,
	helper helpers.Helper,
	emailService Email,
	jwt JWTService,
	config *config.Config,
	userRepository repository.User,
) Team {
	return &team{
		logger:         logger,
		teamRepository: teamRepository,
		helper:         helper,
		emailService:   emailService,
		jwt:            jwt,
		config:         config,
		userRepository: userRepository,
	}
}
