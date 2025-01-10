package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/errs"
)

func GetProtectUserId(c *gin.Context) (uId uint, err error) {
	uIdStr, ok := c.Get("userId")
	if !ok {
		return 0, errs.ErrUnauthorized
	}
	uId, ok = uIdStr.(uint)
	if !ok {
		return 0, errs.ErrUnauthorized
	}
	return uId, nil
}

func GetUserIdFromParam(c *gin.Context) (uId uint, err error) {
	uIdStr := c.Param("uId")
	u, err := strconv.Atoi(uIdStr)
	if err != nil {
		return 0, errs.ErrUnauthorized
	}
	return uint(u), nil
}
