package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/dto"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/service"
)

type (
	Team interface {
		Create(c *gin.Context)
		GetTeamsMe(c *gin.Context)
		GetTeamMemberCount(c *gin.Context)
		GetTeamMembers(c *gin.Context)
		GetTeamByTeamId(c *gin.Context)
		GetPendingTeamMemberCount(c *gin.Context)
		GetPendingTeamMembers(c *gin.Context)
		GetTeamUserMe(c *gin.Context)
		UpdateMemberRole(c *gin.Context)
		SendInviteTeamMember(c *gin.Context)
		JoinTeamWithToken(c *gin.Context)
		CreateShearLink(c *gin.Context)
		JoinWithShearLink(c *gin.Context)
		GetTeamsFilter(c *gin.Context)
		CreateTeamPendingTeamMember(c *gin.Context)
		AcceptTeamMember(c *gin.Context)
		UpdateTeamInfo(c *gin.Context)
	}
	team struct {
		teamService service.Team
	}
)

// UpdateTeamInfo implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		team_id	path		integer							true	"Team ID".
//	@Param		data	body		service.UpdateTeamInfoRequest	true	"UpdateTeamInfoRequest".
//	@Success	201		{object}	Response[any]
//	@Failure	400		{object}	Response[any]
//	@Failure	404		{object}	Response[any]
//	@Failure	500		{object}	Response[any]
//	@Router		/teams/{team_id} [put]
func (t *team) UpdateTeamInfo(c *gin.Context) {
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	input := service.UpdateTeamInfoRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		handlerError(c, err)
		return
	}
	if err := t.teamService.UpdateTeamInfo(c, teamId, input); err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

// AcceptTeamMember implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		team_id	path		integer					true	"Team ID".
//	@Param		data	body		dto.AcceptTeamMemberDto	true	"AcceptTeamMemberDto".
//	@Success	201		{object}	Response[any]
//	@Failure	400		{object}	Response[any]
//	@Failure	404		{object}	Response[any]
//	@Failure	500		{object}	Response[any]
//	@Router		/teams/{team_id}/accept [post]
func (t *team) AcceptTeamMember(c *gin.Context) {
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	input := dto.AcceptTeamMemberDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		handlerError(c, err)
		return
	}
	err = t.teamService.AcceptTeamMember(c, teamId, input)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

// CreateTeamPendingTeamMember implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		team_id							path		integer	true	"Team ID".
//	@Success	201								{object}	Response[any]
//	@Failure	400								{object}	Response[any]
//	@Failure	404								{object}	Response[any]
//	@Failure	500								{object}	Response[any]
//	@Router		/teams/{team_id}/request-join	[post]
func (t *team) CreateTeamPendingTeamMember(c *gin.Context) {
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	userId, err := getProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	err = t.teamService.CreateTeamPendingTeamMember(c, teamId, userId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

// GetTeamsFilter implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		name	query		string	false	"Team Name".
//	@Success	200		{object}	ResponsePagination[[]service.TeamResponse]
//	@Failure	400		{object}	Response[any]
//	@Failure	404		{object}	Response[any]
//	@Failure	500		{object}	Response[any]
//	@Router		/teams/	[get]
func (t *team) GetTeamsFilter(c *gin.Context) {
	pg, err := helpers.NewPaginate(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	f, err := filter.New[filter.TeamFilter](c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := t.teamService.GetTeamsFilter(c, pg, f)
	if err != nil {
		handlerError(c, err)
		return
	}
	handlePaginationJsonResponse(c, res, pg)
}

// CreateShearLink implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		team_id							path		integer	true	"Team ID".
//	@Success	200								{object}	Response[string]
//	@Failure	400								{object}	Response[any]
//	@Failure	404								{object}	Response[any]
//	@Failure	500								{object}	Response[any]
//	@Router		/teams/{team_id}/shared-link	[post]
func (t *team) CreateShearLink(c *gin.Context) {
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := t.teamService.CreateShearLink(c, teamId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

// JoinWithShearLink implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		token				query		string	true	"TOKEN"
//	@Success	200					{object}	Response[any]
//	@Failure	400					{object}	Response[any]
//	@Failure	404					{object}	Response[any]
//	@Failure	500					{object}	Response[any]
//	@Router		/teams/join/link	[post]
func (t *team) JoinWithShearLink(c *gin.Context) {
	userId, err := getProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	token := c.Query("token")
	err = t.teamService.JoinWithShearLink(c, token, userId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

// JoinTeamWithToken implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		token		query		string	true	"TOKEN"
//	@Success	200			{object}	Response[any]
//	@Failure	400			{object}	Response[any]
//	@Failure	404			{object}	Response[any]
//	@Failure	500			{object}	Response[any]
//	@Router		/teams/join	[get]
func (t *team) JoinTeamWithToken(c *gin.Context) {
	token := c.Query("token")
	err := t.teamService.JoinTeamWithToken(c, token)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

// CreateTeamPendingTeamMember implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		data							body		dto.CreateTeamPendingTeamMemberDto	true	"CreateTeamPendingTeamMemberDto"
//	@Param		team_id							path		integer								true	"Team ID"
//	@Success	201								{object}	Response[any]
//	@Failure	400								{object}	Response[any]
//	@Failure	404								{object}	Response[any]
//	@Failure	500								{object}	Response[any]
//	@Router		/teams/{team_id}/pending-member	[post]
func (t *team) SendInviteTeamMember(c *gin.Context) {
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	input := dto.CreateTeamPendingTeamMemberDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		handlerError(c, err)
		return
	}
	err = t.teamService.SendInviteTeamMember(c, teamId, input)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

// UpdateMemberRole implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		data							body		dto.UpdateMemberRoleDto	true	"UpdateMemberRoleDto"
//	@Param		team_id							path		integer					true	"Team ID"
//	@Success	201								{object}	Response[any]
//	@Failure	400								{object}	Response[any]
//	@Failure	404								{object}	Response[any]
//	@Failure	500								{object}	Response[any]
//	@Router		/teams/{team_id}/member-role	[put]
func (t *team) UpdateMemberRole(c *gin.Context) {
	input := dto.UpdateMemberRoleDto{}
	if err := c.ShouldBindJSON(&input); err != nil {
		handlerError(c, err)
		return
	}
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	err = t.teamService.UpdateMemberRole(c, teamId, input)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

// GetTeamUserMe implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		team_id						path		integer	true	"Team ID"
//	@Success	200							{object}	Response[service.TeamMemberResponse]
//	@Failure	400							{object}	Response[any]
//	@Failure	404							{object}	Response[any]
//	@Failure	500							{object}	Response[any]
//	@Router		/teams/{team_id}/user-me	[get]
func (t *team) GetTeamUserMe(c *gin.Context) {
	userId, err := getProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := t.teamService.GetTeamUserMe(c, teamId, userId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

// GetPendingTeamMemberCount implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		team_id									path		integer	true	"Team ID"
//	@Success	200										{object}	Response[int64]
//	@Failure	400										{object}	Response[any]
//	@Failure	404										{object}	Response[any]
//	@Failure	500										{object}	Response[any]
//	@Router		/teams/{team_id}/pending-member-count	[get]
func (t *team) GetPendingTeamMemberCount(c *gin.Context) {
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := t.teamService.GetPendingTeamMemberCount(c, teamId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

// GetPendingTeamMembers implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		team_id								path		integer	true	"Team ID"
//	@Param		username							query		string	false	"Search"
//	@Param		page								query		int		false	"Page"
//	@Param		limit								query		int		false	"Limit"
//	@Success	200									{object}	ResponsePagination[[]service.TeamMemberResponse]
//	@Failure	400									{object}	Response[any]
//	@Failure	404									{object}	Response[any]
//	@Failure	500									{object}	Response[any]
//	@Router		/teams/{team_id}/pending-members	[get]
func (t *team) GetPendingTeamMembers(c *gin.Context) {
	pg, err := helpers.NewPaginate(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	f, err := filter.New[filter.TeamMemberFilter](c)
	if err != nil {
		handlerError(c, err)
		return
	}
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := t.teamService.GetPendingTeamMembers(c, pg, f, teamId)

	if err != nil {
		handlerError(c, err)
		return
	}
	handlePaginationJsonResponse(c, res, pg)
}

// GetTeamByTeamId implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		team_id				path		integer	true	"Team ID"
//	@Success	200					{object}	Response[service.TeamResponse]
//	@Failure	400					{object}	Response[any]
//	@Failure	404					{object}	Response[any]
//	@Failure	500					{object}	Response[any]
//	@Router		/teams/{team_id}	[get]
func (t *team) GetTeamByTeamId(c *gin.Context) {
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := t.teamService.GetTeamByTeamId(c, teamId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

// GetTeamMembers implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		team_id						path		integer	true	"Team ID"
//	@Param		username					query		string	false	"Search"
//	@Param		page						query		int		false	"Page"
//	@Param		limit						query		int		false	"Limit"
//	@Success	200							{object}	ResponsePagination[[]service.TeamMemberResponse]
//	@Failure	400							{object}	Response[any]
//	@Failure	404							{object}	Response[any]
//	@Failure	500							{object}	Response[any]
//	@Router		/teams/{team_id}/members	[get]
func (t *team) GetTeamMembers(c *gin.Context) {
	pg, err := helpers.NewPaginate(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	f, err := filter.New[filter.TeamMemberFilter](c)
	if err != nil {
		handlerError(c, err)
		return
	}
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := t.teamService.GetTeamMembers(c, pg, f, teamId)

	if err != nil {
		handlerError(c, err)
		return
	}
	handlePaginationJsonResponse(c, res, pg)
}

// GetTeamMemberCount implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		team_id							path		integer	true	"Team ID"
//	@Success	200								{object}	Response[int64]
//	@Failure	400								{object}	Response[any]
//	@Failure	404								{object}	Response[any]
//	@Failure	500								{object}	Response[any]
//	@Router		/teams/{team_id}/member-count	[get]
func (t *team) GetTeamMemberCount(c *gin.Context) {
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := t.teamService.GetTeamMemberCount(c, teamId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

// Create implements Team.
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		data	body		dto.CreateTeamDto	true	"CreateTeamDto"
//	@Success	201		{object}	Response[service.CreateTeamResponse]
//	@Failure	400		{object}	Response[any]
//	@Failure	404		{object}	Response[any]
//	@Failure	500		{object}	Response[any]
//	@Router		/teams/	[post]
func (t *team) Create(c *gin.Context) {
	var (
		input dto.CreateTeamDto
	)
	userId, err := getProtectUserId(c)
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
//
//	@Tags		teams
//	@Accept		json
//	@Produce	json
//	@Param		page		query		int	false	"Page"
//	@Param		limit		query		int	false	"Limit"
//	@Success	200			{object}	ResponsePagination[[]service.TeamResponse]
//	@Failure	400			{object}	Response[any]
//	@Failure	404			{object}	Response[any]
//	@Failure	500			{object}	Response[any]
//	@Router		/teams/me	[get]
func (t *team) GetTeamsMe(c *gin.Context) {
	userId, err := getProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	pg, err := helpers.NewPaginate(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := t.teamService.GetTeamsMe(c, pg, userId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handlePaginationJsonResponse(c, res, pg)
}

func NewTeam(teamService service.Team) Team {
	return &team{
		teamService: teamService,
	}
}
