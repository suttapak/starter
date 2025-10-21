package middleware

import (
	"strconv"
	"strings"

	"github.com/suttapak/starter/errs"
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
	}
	authGuardMiddleware struct {
		jwt      service.JWTService
		enforcer *casbin.Enforcer
		logger   logger.AppLogger
		user     service.UserService
	}
)

// Permission implements AuthGuardMiddleware.
func (a *authGuardMiddleware) Permission(c *gin.Context) {
	userIdStr, ok := c.Get("user_id")
	if !ok {
		a.logger.Error("can not get userId")
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	uId, ok := userIdStr.(uint)
	if !ok {
		a.logger.Error("user id not uint")
		handlerError(c, errs.ErrUnauthorized)
		return
	}
	user, err := a.user.GetUserByUserId(c, uId)
	if err != nil {
		a.logger.Error(err)
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
	token, _ = c.Cookie("session")
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

func NewAuthGuardMiddleware(
	jwt service.JWTService,
	enforcer *casbin.Enforcer,
	logger logger.AppLogger,
	user service.UserService,
) AuthGuardMiddleware {
	return &authGuardMiddleware{
		jwt:      jwt,
		enforcer: enforcer,
		logger:   logger,
		user:     user,
	}
}
