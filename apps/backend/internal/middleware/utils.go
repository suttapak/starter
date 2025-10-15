package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/internal/controller"
)

func handlerError(c *gin.Context, err error) {
	message := "Something went wrong"
	status := http.StatusBadRequest
	var appErr errs.AppError
	if errors.As(err, &appErr) {
		message = appErr.Message
		status = appErr.Code
	}
	// Create response object
	response := controller.Response[any]{
		Message: message,
		Status:  status,
		Data:    nil,
	}

	c.AbortWithStatusJSON(status, response)
}
