package route

import (
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"

	"github.com/gin-gonic/gin"
)

type (
	auth struct {
		r              *gin.Engine
		authController controller.Auth
		guard          middleware.AuthGuardMiddleware
	}
)

func newAuth(r *gin.Engine, authController controller.Auth, guard middleware.AuthGuardMiddleware) (a *auth) {
	a = &auth{
		r:              r,
		authController: authController,
		guard:          guard,
	}
	return
}

func useAuth(a *auth) {
	r := a.r.Group("auth")
	{
		r.POST("/login", a.authController.Login)
		r.POST("/register", a.authController.Register)
		r.POST("/refresh", a.guard.ProtectRefreshToken, a.authController.RefreshToken)
		r.POST("/logout", a.authController.Logout)
		r.GET("/email/verify", a.authController.VerifyEmail)
		r.POST("/email/send-verify", a.guard.Protect, a.authController.SendVerifyEmail)
	}
}
