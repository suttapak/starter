package route

import (
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UseAuth(engine *gin.Engine, c controller.Auth, guard middleware.AuthGuardMiddleware) {
	r := engine.Group("auth")
	{
		r.POST("/login", c.Login)
		r.POST("/register", c.Register)
		r.POST("/refresh", guard.ProtectRefreshToken, c.RefreshToken)
		r.POST("/logout", c.Logout)
		r.GET("/email/verify", c.VerifyEmail)
		r.POST("/email/send-verify", guard.Protect, c.SendVerifyEmail)
	}
}
