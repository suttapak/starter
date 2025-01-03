package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"github.com/suttapak/starter/errs"
	"net/http"
)

var Module = fx.Options(
	fx.Provide(newPost),
	fx.Provide(newAuth),
)

func handlerError(c *gin.Context, err error) {
	var appErr errs.AppError
	if errors.As(err, &appErr) {
		c.AbortWithStatusJSON(appErr.Code, appErr.Message)
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
}
