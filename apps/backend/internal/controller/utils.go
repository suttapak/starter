package controller

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
)

func handlePaginationJsonResponse(c *gin.Context, json any, pg *helpers.Pagination, msg ...string) {
	// Default message if none is provided
	message := "Success"
	if len(msg) > 0 {
		message = msg[0]
	}

	// Determine HTTP status based on request method
	status := http.StatusOK
	switch c.Request.Method {
	case http.MethodPost:
		status = http.StatusCreated
	case http.MethodPut, http.MethodPatch:
		status = http.StatusAccepted
	case http.MethodDelete:
		status = http.StatusNoContent
	}

	// Create response object
	response := Response[any]{
		Message: message,
		Status:  status,
		Data:    json,
		Meta:    pg,
	}

	// Send response with the determined status
	c.JSON(status, response)
}

func handleJsonResponse(c *gin.Context, json any, msg ...string) {
	// Default message if none is provided
	message := "Success"
	if len(msg) > 0 {
		message = msg[0]
	}

	// Determine HTTP status based on request method
	status := http.StatusOK
	switch c.Request.Method {
	case http.MethodPost:
		status = http.StatusCreated
	case http.MethodPut, http.MethodPatch:
		status = http.StatusAccepted
	case http.MethodDelete:
		status = http.StatusNoContent
	}

	// Create response object
	response := Response[any]{
		Message: message,
		Status:  status,
		Data:    json,
	}

	// Send response with the determined status
	c.JSON(status, response)
}

func handlerError(c *gin.Context, err error) {
	message := "Something went wrong"
	status := http.StatusBadRequest
	var appErr errs.AppError
	if errors.As(err, &appErr) {
		message = appErr.Message
		status = appErr.Code
	}
	// Create response object
	response := Response[any]{
		Message: message,
		Status:  status,
		Data:    nil,
	}

	c.AbortWithStatusJSON(status, response)
}

func getTeamId(c *gin.Context) (teamId uint, err error) {
	teamIdStr := c.Param("team_id")
	if teamIdStr == "" {
		return 0, errs.ErrNotActiveTeamId
	}
	teamIdInt, err := strconv.Atoi(teamIdStr)
	if err != nil {
		return 0, errs.ErrNotActiveTeamId
	}
	return uint(teamIdInt), nil
}

func getFiles(c *gin.Context) ([]*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, errs.ErrFileUploadNotFound
	}
	files := form.File["files"]
	if len(files) == 0 {
		return nil, errs.ErrFileUploadNoFile
	}
	return files, nil
}

func getProductImageId(c *gin.Context) (uint, error) {
	productImageIdStr := c.Param("product_image_id")

	productImageIdInt, err := strconv.Atoi(productImageIdStr)
	if err != nil {
		return 0, errs.ErrBadRequest
	}
	return uint(productImageIdInt), nil
}

func getProtectUserId(c *gin.Context) (uId uint, err error) {
	uIdStr, ok := c.Get("user_id")
	if !ok {
		return 0, errs.ErrUnauthorized
	}
	uId, ok = uIdStr.(uint)
	if !ok {
		return 0, errs.ErrUnauthorized
	}
	return uId, nil
}

func getUserIdFromParam(c *gin.Context) (uId uint, err error) {
	uIdStr := c.Param("id")
	u, err := strconv.Atoi(uIdStr)
	if err != nil {
		return 0, errs.ErrUnauthorized
	}
	return uint(u), nil
}
