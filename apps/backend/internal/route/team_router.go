package route

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	u := t.r.Group("teams", t.authMiddleware.Protect)
	{
		u.GET("/me", t.teamController.GetTeamsMe)
		u.GET("/", t.teamController.GetTeamsFilter)
		u.POST("/", t.teamController.Create)
		u.GET("/join", t.teamController.JoinTeamWithToken)
		u.POST("/join/link", t.teamController.JoinWithShearLink)
		u.POST("/:team_id/request-join", t.teamController.CreateTeamPendingTeamMember)
	}
	r := t.r.Group("teams", t.authMiddleware.Protect, t.authMiddleware.TeamPermission)
	{
		r.GET("/:team_id", t.teamController.GetTeamByTeamId)
		r.GET("/:team_id/member-count", t.teamController.GetTeamMemberCount)
		r.GET("/:team_id/members", t.teamController.GetTeamMembers)
		r.GET("/:team_id/pending-member-count", t.teamController.GetPendingTeamMemberCount)
		r.GET("/:team_id/pending-members", t.teamController.GetPendingTeamMembers)
		r.GET("/:team_id/user-me", t.teamController.GetTeamUserMe)
		r.PUT("/:team_id/member-role", t.teamController.UpdateMemberRole)
		r.POST("/:team_id/pending-member", t.teamController.SendInviteTeamMember)
		r.POST("/:team_id/shared-link", t.teamController.CreateShearLink)
		r.POST("/:team_id/accept", t.teamController.AcceptTeamMember)
		r.PUT("/:team_id", t.teamController.UpdateTeamInfo)
	}
}

func seedTeamPermission(db *gorm.DB) {
	/*
		TeamRoleOwnerID = iota + 1
		TeamRoleAdminID
		TeamRoleMemberID
	*/
	// v1 for role id, v2 for path, v3 for method
	var permission = []gormadapter.CasbinRule{
		{
			Ptype: "p",
			V0:    "1",
			V1:    "*",
			V2:    "*",
		},
		{
			Ptype: "p",
			V0:    "2",
			V1:    "*",
			V2:    "*",
		},
		{
			Ptype: "p",
			V0:    "3",
			V1:    "/teams/*",
			V2:    "GET",
		},
		{
			Ptype: "p",
			V0:    "3",
			V1:    "/teams",
			V2:    "POST",
		},
		{
			Ptype: "p",
			V0:    "3",
			V1:    "/teams/join/link",
			V2:    "POST",
		},
		{
			Ptype: "p",
			V0:    "3",
			V1:    "/teams/{id}/request-join",
			V2:    "POST",
		},
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&permission)

}
