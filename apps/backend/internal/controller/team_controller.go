package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/dto"
	"github.com/suttapak/starter/internal/service"
)

type (
	Team interface {
		Create(c *gin.Context)
		GetTeamsMe(c *gin.Context)
	}
	team struct {
		teamService service.Team
	}
)

// Create implements Team.
func (t *team) Create(c *gin.Context) {
	var (
		input dto.CreateTeamDto
	)
	userId, err := helpers.GetProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		handlerError(c, err)
		return
	}
	res, err := t.teamService.Create(c, userId, input)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

// GetTeamsMe implements Team.
func (t *team) GetTeamsMe(c *gin.Context) {
	userId, err := helpers.GetProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := t.teamService.GetTeamsMe(c, userId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

func NewTeam(teamService service.Team) Team {
	return &team{
		teamService: teamService,
	}
}
