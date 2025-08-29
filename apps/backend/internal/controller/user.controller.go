package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/internal/service"
)

type (
	User interface {
		GetUserMe(c *gin.Context)
		GetUserById(c *gin.Context)
		FindUserByUsername(c *gin.Context)
		CheckUserIsVerifyEmail(c *gin.Context)
		CreateProfileImage(c *gin.Context)
	}
	user struct {
		userService service.UserService
	}
)

// CreateProfileImage implements User.
//
// @Tags      users
// @Accept    multipart/form-data
// @Produce   json
// @Param     file formData file true "Profile image file"
// @Success   201 {object} Response[service.UserResponse]
// @Failure   400 {object} Response[any]
// @Failure   404 {object} Response[any]
// @Failure   500 {object} Response[any]
// @Router    /users/profile-image [post]
func (a *user) CreateProfileImage(c *gin.Context) {
	userId, err := getProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := a.userService.CreateProfileImage(c, userId, file)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

// CheckUserIsVerifyEmail implements Auth.
func (a *user) CheckUserIsVerifyEmail(c *gin.Context) {
	userId, err := getProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := a.userService.CheckUserIsVerifyEmail(c, userId)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)

}

// FindUserByUsername implements User.
func (u *user) FindUserByUsername(c *gin.Context) {
	username := c.Query("username")
	res, err := u.userService.FindUserByUsername(c, username)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

// GetUserById implements User.
func (u *user) GetUserById(c *gin.Context) {
	// get user id form middleware
	uId, err := getUserIdFromParam(c)
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
func (u *user) GetUserMe(c *gin.Context) {
	uId, err := getProtectUserId(c)
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
	return &user{userService: userService}
}
