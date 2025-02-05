package route

import (
	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"
)

type (
	team struct {
		r              *gin.Engine
		teamController controller.Team
		authMiddleware middleware.AuthGuardMiddleware
	}
)

func newTeam(r *gin.Engine, teamController controller.Team, authMiddleware middleware.AuthGuardMiddleware) (t *team) {
	return &team{
		r:              r,
		teamController: teamController,
		authMiddleware: authMiddleware,
	}
}

func useTeam(t *team) {
	r := t.r.Group("teams", t.authMiddleware.Protect)
	{
		r.GET("/me", t.teamController.GetTeamsMe)
		r.POST("/", t.teamController.Create)
	}
}
