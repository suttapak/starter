package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	healthcheck struct {
		r *gin.Engine
	}
)

func newHealthCheck(r *gin.Engine) (p healthcheck) {
	p = healthcheck{
		r: r,
	}
	return
}
func useHealthCheck(post healthcheck) {
	group := post.r.Group("healthcheck")
	{
		group.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"OK": true})
		})
	}
}
