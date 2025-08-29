package route

import (
	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"
)

type (
	user struct {
		r              *gin.Engine
		userController controller.User
		guard          middleware.AuthGuardMiddleware
	}
)

func newUser(r *gin.Engine, userController controller.User, guard middleware.AuthGuardMiddleware) (u *user) {
	u = &user{
		r:              r,
		userController: userController,
		guard:          guard,
	}
	return
}

func useUser(u *user) {
	group := u.r.Group("users", u.guard.Protect)
	{
		group.GET("/:id", u.userController.GetUserById)
		group.GET("/me", u.userController.GetUserMe)
		group.GET("/by-username", u.userController.FindUserByUsername)
		group.GET("/verify-email", u.userController.CheckUserIsVerifyEmail)
		group.POST("/profile-image", u.userController.CreateProfileImage)
	}
}
