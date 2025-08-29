package route

import (
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"

	"github.com/gin-gonic/gin"
)

type (
	report struct {
		r      *gin.Engine
		guard  middleware.AuthGuardMiddleware
		report controller.Report
	}
)

func newReport(r *gin.Engine, reportController controller.Report, guard middleware.AuthGuardMiddleware) *report {
	return &report{
		r:      r,
		report: reportController,
		guard:  guard,
	}
}

func useReport(a *report) {
	r := a.r.Group("report", a.guard.Protect)
	{
		r.GET("/:report_id", a.report.GetReport)
		r.GET("", a.report.GetReports)
		r.POST("", a.report.CreateReport)
		r.PUT("/:report_id", a.report.UpdateReport)
		r.DELETE("/:report_id", a.report.DeleteReport)
	}
}
