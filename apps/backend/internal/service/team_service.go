package service

import (
	"context"

	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/dto"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/internal/response"
	"github.com/suttapak/starter/logger"
)

type (
	Team interface {
		Create(ctx context.Context, ownerId uint, input dto.CreateTeamDto) (res *response.CreateTeamResponse, err error)
		GetTeamsMe(ctx context.Context, userId uint) (res []response.TeamResponse, err error)
	}
	team struct {
		logger         logger.AppLogger
		teamRepository repository.Team
		helper         helpers.Helper
	}
)

// GetTeamsMe implements Team.
func (t team) GetTeamsMe(ctx context.Context, userId uint) (res []response.TeamResponse, err error) {
	models, err := t.teamRepository.GetTeamsMe(ctx, nil, userId)
	if err != nil {
		t.logger.Error(err)
		return nil, errs.HandeGorm(err)
	}
	if err := t.helper.ParseJson(models, &res); err != nil {
		t.logger.Error(err)
		return nil, errs.ErrInvalid
	}
	return
}

// Create implements Team.
func (t team) Create(ctx context.Context, ownerId uint, input dto.CreateTeamDto) (res *response.CreateTeamResponse, err error) {
	// check with username is exist
	exists, err := t.teamRepository.CheckTeamUsernameIsExist(ctx, nil, input.Username)
	if err != nil {
		t.logger.Error(err)
		return nil, errs.HandeGorm(err)
	}
	if exists {
		t.logger.Error(err)
		return nil, errs.ErrTeamUsernameIsUsed
	}
	// create team with role id

	model, err := t.teamRepository.Create(ctx, nil, ownerId, repository.CreateTeamParams(input))
	if err != nil {
		t.logger.Error(err)
		return nil, errs.HandeGorm(err)
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
) Team {
	return team{logger: logger, teamRepository: teamRepository, helper: helper}
}
