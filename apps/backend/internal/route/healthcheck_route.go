package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UseHealthCheck(r *gin.Engine) {
	group := r.Group("health")
	{
		group.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"OK": true})
		})
	}
}
