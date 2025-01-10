package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/errs"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(newPost),
	fx.Provide(NewAuth),
)

type (
	Response[T any] struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
		Data    T      `json:"data"`
		Meta    any    `json:"meta"`
	}
)

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
