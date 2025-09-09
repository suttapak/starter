package middleware

import (
	"strconv"
	"strings"

	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/internal/service"
	"github.com/suttapak/starter/logger"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type (
	AuthGuardMiddleware interface {
		Protect(c *gin.Context)
		ProtectRefreshToken(c *gin.Context)
		Permission(c *gin.Context)
		TeamPermission(c *gin.Context)
	}
	authGuardMiddleware struct {
		jwt      service.JWTService
		enforcer *casbin.Enforcer
		logger   logger.AppLogger
		user     service.UserService
		team     service.Team
		teamRepo repository.Team
	}
)

// TeamPermission implements AuthGuardMiddleware.
func (a *authGuardMiddleware) TeamPermission(c *gin.Context) {
	userIdStr, ok := c.Get("user_id")
	if !ok {
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	uId, ok := userIdStr.(uint)
	if !ok {
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	teamIdStr := c.Param("team_id")
	teamId, err := strconv.Atoi(teamIdStr)
	if err != nil {
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	member, err := a.team.GetTeamUserMe(c, uint(teamId), uId)
	if err != nil {
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	isExist, err := a.teamRepo.ExistUserInTeamByTeamId(c, nil, uint(teamId), uId)
	if err != nil {
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	if !isExist {
		handlerError(c, errs.ErrUnauthorized)
		return
	}

	if err := a.enforcer.LoadPolicy(); err != nil {
		a.logger.Error(err)
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	allowed, err := a.enforcer.Enforce(strconv.Itoa(int(member.TeamRoleID)), c.Request.URL.Path, c.Request.Method)
	if err != nil {
		a.logger.Error(err)
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	if !allowed {
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	c.Next()
}

// Permission implements AuthGuardMiddleware.
func (a *authGuardMiddleware) Permission(c *gin.Context) {
	userIdStr, ok := c.Get("user_id")
	if !ok {
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	uId, ok := userIdStr.(uint)
	if !ok {
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	user, err := a.user.GetUserByUserId(c, uId)
	if err != nil {
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	if err := a.enforcer.LoadPolicy(); err != nil {
		a.logger.Error(err)
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	allowed, err := a.enforcer.Enforce(strconv.Itoa(int(user.ID)), c.Request.URL.Path, c.Request.Method)
	if err != nil {
		a.logger.Error(err)
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	if !allowed {
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	c.Next()
}

// Protect implements AuthGuardMiddleware.
func (a *authGuardMiddleware) Protect(c *gin.Context) {
	var (
		token string
	)
	token, err := c.Cookie("session")
	if token == "" {
		token = c.GetHeader("Authorization")
		splitToken := strings.Split(token, " ")

		if len(splitToken) != 2 {
			token = ""
		}
		if len(splitToken) == 2 {
			token = splitToken[1]
		}
	}

	uId, err := a.jwt.GetUserIdFormToken(c, token)
	if err != nil {
		a.logger.Error(err)
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	user, err := a.user.GetUserByUserId(c, uId)
	if err != nil {
		// form service not logger err
		a.logger.Error(err)
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	c.Set("user_id", user.ID)
	c.Next()
}

func (a *authGuardMiddleware) ProtectRefreshToken(c *gin.Context) {
	var (
		token string
	)
	token = c.GetHeader("Authorization")
	splitToken := strings.Split(token, " ")

	if len(splitToken) != 2 {
		token = ""
	}
	if len(splitToken) == 2 {
		token = splitToken[1]
	}

	uId, err := a.jwt.GetUserIdFormRefreshToken(c, token)
	if err != nil {
		a.logger.Error(err)
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	user, err := a.user.GetUserByUserId(c, uId)
	if err != nil {
		// form service not logger err
		a.logger.Error(err)
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	c.Set("user_id", user.ID)
	c.Next()
}

func newAuthGuardMiddleware(
	jwt service.JWTService,
	enforcer *casbin.Enforcer,
	logger logger.AppLogger,
	user service.UserService,
	team service.Team,
	teamRepo repository.Team,
) AuthGuardMiddleware {
	return &authGuardMiddleware{
		jwt:      jwt,
		enforcer: enforcer,
		logger:   logger,
		user:     user,
		team:     team,
		teamRepo: teamRepo,
	}
}
