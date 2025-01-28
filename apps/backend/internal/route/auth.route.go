package route

import (
	"github.com/suttapak/starter/internal/controller"

	"github.com/gin-gonic/gin"
)

type (
	auth struct {
		r              *gin.Engine
		authController controller.Auth
	}
)

func newAuth(r *gin.Engine, authController controller.Auth) (a auth) {
	a = auth{
		r:              r,
		authController: authController,
	}
	return
}

func useAuth(a auth) {
	r := a.r.Group("auth")
	{
		r.POST("/login", a.authController.Login)
		r.POST("/register", a.authController.Register)
		r.POST("/logout", a.authController.Logout)
	}
}
