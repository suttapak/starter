package route

import (
	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"
)

func UseUser(
	r *gin.Engine,
	userController controller.User,
	guard middleware.AuthGuardMiddleware,
) {
	group := r.Group("users", guard.Protect)
	{
		group.GET("/:id", userController.GetUserById)
		group.GET("/me", userController.GetUserMe)
		group.GET("/by-username", userController.FindUserByUsername)
		group.GET("/verify-email", userController.CheckUserIsVerifyEmail)
		group.POST("/profile-image", userController.CreateProfileImage)
	}
}
