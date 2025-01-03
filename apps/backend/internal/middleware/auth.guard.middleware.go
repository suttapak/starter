package middleware

import (
	"labostack/internal/service"
	"labostack/logger"
	"net/http"
	"strconv"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type (
	AuthGuardMiddleware interface {
		Protect(c *gin.Context)
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
func (a authGuardMiddleware) Permission(c *gin.Context) {
	userIdStr, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}
	uId, ok := userIdStr.(uint)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}
	user, err := a.user.GetUserByUserId(c, uId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}
	if err := a.enforcer.LoadPolicy(); err != nil {
		a.logger.Error(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}
	allowed, err := a.enforcer.Enforce(strconv.Itoa(int(user.ID)), c.Request.URL.Path, c.Request.Method)
	if err != nil {
		a.logger.Error(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}
	if !allowed {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}
	c.Next()
}

// Protect implements AuthGuardMiddleware.
func (a authGuardMiddleware) Protect(c *gin.Context) {
	token := c.GetHeader("Authorization")
	splitedToken := strings.Split(token, " ")
	if len(splitedToken) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}
	if splitedToken[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}
	uId, err := a.jwt.GetUserIdFormToken(c, splitedToken[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		return
	}
	user, err := a.user.GetUserByUserId(c, uId)
	if err != nil {
		// form service not logger err
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
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
) AuthGuardMiddleware {
	return authGuardMiddleware{
		jwt:      jwt,
		enforcer: enforcer,
		logger:   logger,
		user:     user,
	}
}
