package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/service"
)

type (
	User interface {
		GetUserMe(c *gin.Context)
		GetUserById(c *gin.Context)
	}
	user struct {
		userService service.UserService
	}
)

// GetUserById implements User.
func (u user) GetUserById(c *gin.Context) {
	// get user id form middleware
	uId, err := helpers.GetUserIdFromParam(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := u.userService.GetUserByUserId(c, uId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

// GetUserMe implements User.
func (u user) GetUserMe(c *gin.Context) {
	uId, err := helpers.GetProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := u.userService.GetUserByUserId(c, uId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

func NewUser(userService service.UserService) User {
	return user{userService: userService}
}
