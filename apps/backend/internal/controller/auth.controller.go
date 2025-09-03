package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/internal/service"
)

type (
	Auth interface {
		Register(c *gin.Context)
		Login(c *gin.Context)
		Logout(c *gin.Context)
		VerifyEmail(c *gin.Context)
		SendVerifyEmail(c *gin.Context)
		RefreshToken(c *gin.Context)
	}

	auth struct {
		conf        *config.Config
		authService service.Auth
	}
)

// RefreshToken implements Auth.
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Success	201	{object}	Response[service.AuthResponse]
//	@Failure	400	{object}	Response[any]
//	@Failure	404	{object}	Response[any]
//	@Failure	500	{object}	Response[any]
//	@Router		/auth/refresh [post]
func (a *auth) RefreshToken(c *gin.Context) {
	userId, err := getProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := a.authService.RefreshToken(c, userId)
	if err != nil {
		handlerError(c, err)
		return
	}
	c.SetCookie("session", res.Token, 0, "/", a.conf.SERVER.HOST, false, true)
	handleJsonResponse(c, res)

}

// SendVerifyEmail implements Auth.
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Success	201	{object}	Response[any]
//	@Failure	400	{object}	Response[any]
//	@Failure	404	{object}	Response[any]
//	@Failure	500	{object}	Response[any]
//	@Router		/auth/email/send-verify [post]
func (a *auth) SendVerifyEmail(c *gin.Context) {
	userId, err := getProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	input := service.SendVerifyEmailDto{UserID: userId}
	err = a.authService.SendVerifyEmail(c, input)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

// Logout implements Auth.
func (a *auth) Logout(c *gin.Context) {
	c.SetCookie("session", "", -1, "/", a.conf.SERVER.HOST, false, true)
}

// VerifyEmail
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		data	body	service.VerifyEmailDto	true	"body data".
//	@Success	301
//	@Failure	400	{object}	Response[any]
//	@Failure	404	{object}	Response[any]
//	@Failure	500	{object}	Response[any]
//	@Router		/auth/email/verify [post]
func (a *auth) VerifyEmail(c *gin.Context) {
	var (
		body service.VerifyEmailDto
	)
	if err := c.ShouldBindQuery(&body); err != nil {
		handlerError(c, err)
		return
	}
	_, err := a.authService.VerifyEmail(c, body)
	if err != nil {
		handlerError(c, err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/verify-email-success")
}

// Login
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		data	body		service.LoginDto	true	"body data".
//	@Success	201		{object}	Response[service.AuthResponse]
//	@Failure	400		{object}	Response[any]
//	@Failure	404		{object}	Response[any]
//	@Failure	500		{object}	Response[any]
//	@Router		/auth/login [post]
func (a *auth) Login(c *gin.Context) {
	//TODO implement me
	var body service.LoginDto

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		handlerError(c, err)
		return
	}
	res, err := a.authService.Login(c, body)
	if err != nil {
		handlerError(c, err)
		return
	}
	c.SetCookie("session", res.Token, 0, "/", a.conf.SERVER.HOST, false, true)
	handleJsonResponse(c, res)
}

// Register
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		data	body		service.UserRegisterDto	true	"body data".
//	@Success	201		{object}	Response[service.AuthResponse]
//	@Failure	400		{object}	Response[any]
//	@Failure	404		{object}	Response[any]
//	@Failure	500		{object}	Response[any]
//	@Router		/auth/register [post]
func (a *auth) Register(c *gin.Context) {
	var user service.UserRegisterDto
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		handlerError(c, err)
		return
	}
	res, err := a.authService.Register(c, user)
	if err != nil {
		handlerError(c, err)
		return
	}

	handleJsonResponse(c, res)
}

func NewAuth(authService service.Auth, conf *config.Config) Auth {
	return &auth{authService: authService, conf: conf}
}
