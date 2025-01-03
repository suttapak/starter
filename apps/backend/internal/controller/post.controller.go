package controller

import (
	"labostack/internal/service"

	"github.com/gin-gonic/gin"
)

type (
	Post interface {
		Fist(c *gin.Context)
		Find(c *gin.Context)
		Create(c *gin.Context)
		Update(c *gin.Context)
		Delete(c *gin.Context)
	}
	post struct {
		postService service.Post
	}
)

// Create
// @BasePath /post
// PingExample godoc
// @Summary create post
// @Schemes
// @Description do ping
// @Tags post
// @Accept json
// @Produce json
// @Success	200		{object}  response.PostResponse
// @Router	/post 	[post]
func (p *post) Create(c *gin.Context) {
	panic("unimplemented")
}

// Delete implements Post.
func (p *post) Delete(c *gin.Context) {
	panic("unimplemented")
}

// Find implements Post.
func (p *post) Find(c *gin.Context) {
	panic("unimplemented")
}

// Fist implements Post.
func (p *post) Fist(c *gin.Context) {
	panic("unimplemented")
}

// Update implements Post.
func (p *post) Update(c *gin.Context) {
	panic("unimplemented")
}

func newPost(postService service.Post) Post {
	return &post{
		postService: postService,
	}
}
