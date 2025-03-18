package filter

import "github.com/gin-gonic/gin"

func New[T any](c *gin.Context) (*T, error) {
	var f T

	if err := c.ShouldBindQuery(&f); err != nil {
		return nil, err
	}

	return &f, nil
}
