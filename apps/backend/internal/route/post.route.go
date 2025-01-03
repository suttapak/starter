package route

import (
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"

	"github.com/gin-gonic/gin"
)

type (
	post struct {
		r           *gin.Engine
		postHandler controller.Post
		guard       middleware.AuthGuardMiddleware
	}
)

func newPost(r *gin.Engine, postHandler controller.Post, guard middleware.AuthGuardMiddleware) (p post) {
	p = post{
		r:           r,
		postHandler: postHandler,
		guard:       guard,
	}
	return
}
func usePost(post post) {
	group := post.r.Group("post", post.guard.Protect, post.guard.Permission)
	{
		group.GET("/", func(c *gin.Context) {

		})
		group.POST("/", post.postHandler.Create)
	}
}
