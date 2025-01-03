package controller

import (
	"github.com/gin-gonic/gin"
	"labostack/internal/dto"
	"labostack/internal/service"
	"net/http"
)

type (
	Auth interface {
		Register(c *gin.Context)
		Login(c *gin.Context)
		VerifyEmail(c *gin.Context)
	}

	auth struct {
		userService service.Auth
	}
)

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
func (a auth) VerifyEmail(c *gin.Context) {
	var (
		body dto.VerifyEmailDto
	)
	if err := c.ShouldBindQuery(&body); err != nil {
		handlerError(c, err)
		return
	}
	_, err := a.userService.VerifyEmail(c, body)
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
func (a auth) Login(c *gin.Context) {
	//TODO implement me
	var body dto.LoginDto

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		handlerError(c, err)
		return
	}
	res, err := a.userService.Login(c, body)
	if err != nil {
		handlerError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
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
func (a auth) Register(c *gin.Context) {
	var user dto.UserRegisterDto
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		handlerError(c, err)
		return
	}
	res, err := a.userService.Register(c, user)
	if err != nil {
		handlerError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func newAuth(userService service.Auth) Auth {
	return auth{userService: userService}
}
