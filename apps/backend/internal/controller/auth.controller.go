package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/internal/dto"
	"github.com/suttapak/starter/internal/service"
)

type (
	Auth interface {
		Register(c *gin.Context)
		Login(c *gin.Context)
		Logout(c *gin.Context)
		VerifyEmail(c *gin.Context)
	}

	auth struct {
		authService service.Auth
	}
)

// Logout implements Auth.
func (a *auth) Logout(c *gin.Context) {
	c.SetCookie("session", "", -1, "/", "localhost", false, true)
}

// VerifyEmail
// @BasePath /auth
// PingExample godoc
// @Summary verify email
// @Schemes
// @Description verify email
// @Tags auth
// @Accept json
// @Produce json
// @Param        token    query     string  false  "name search by q"  Format(email)
// @Success	200		{object}  response.UserResponse
// @Router	/auth/email/verify [post]
func (a *auth) VerifyEmail(c *gin.Context) {
	var (
		body dto.VerifyEmailDto
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
	c.Redirect(http.StatusMovedPermanently, "/auth/email/success")
}

// Login
// @BasePath /auth
// PingExample godoc
// @Summary user login
// @Schemes
// @Description user login
// @Tags auth
// @Accept json
// @Produce json
// @Param data body dto.LoginDto true "body data".
// @Success	200		{object}  response.AuthResponse
// @Router	/auth/login 	[post]
func (a *auth) Login(c *gin.Context) {
	//TODO implement me
	var body dto.LoginDto

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		handlerError(c, err)
		return
	}
	res, err := a.authService.Login(c, body)
	if err != nil {
		handlerError(c, err)
		return
	}
	c.SetCookie("session", res.Token, 0, "/", "localhost", false, true)
	handleJsonResponse(c, res)
}

// Register
// @BasePath /auth
// PingExample godoc
// @Summary register user
// @Schemes
// @Description register user
// @Tags auth
// @Accept json
// @Produce json
// @Param data body dto.UserRegisterDto true "body data".
// @Success	200		{object}  response.AuthResponse
// @Router	/auth/register 	[post]
func (a *auth) Register(c *gin.Context) {
	var user dto.UserRegisterDto
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

func NewAuth(authService service.Auth) Auth {
	return &auth{authService: authService}
}
