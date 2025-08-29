package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/service"
)

type (
	Report interface {
		GetReport(c *gin.Context)
		GetReports(c *gin.Context)
		CreateReport(c *gin.Context)
		UpdateReport(c *gin.Context)
		DeleteReport(c *gin.Context)
		// Render(c *gin.Context)
	}

	report struct {
		report service.Report
	}
)

// // Render implements Report.
// func (i *report) Render(c *gin.Context) {
// 	reportTemplateId, err := getReportId(c)
// 	if err != nil {
// 		handlerError(c, err)
// 		return
// 	}
// 	transactionId, err := getTransactionId(c)
// 	if err != nil {
// 		handlerError(c, err)
// 		return
// 	}
// 	res, err := i.report.Render(c, reportTemplateId, transactionId)
// 	if err != nil {
// 		handlerError(c, err)
// 		return
// 	}
// 	handleJsonResponse(c, res)
// }

func getReportId(c *gin.Context) (uint, error) {
	idStr := c.Param("report_id")
	if idStr == "" {
		return 0, errs.ErrNotFound
	}
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errs.ErrBadRequest
	}
	return uint(idInt), nil
}

func (i *report) GetReport(c *gin.Context) {
	id, err := getReportId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := i.report.GetReport(c, id)
	if err != nil {
		handlerError(c, err)
		return
	}

	handleJsonResponse(c, res)
}
func (i *report) GetReports(c *gin.Context) {
	pg, err := helpers.NewPaginate(c)
	if err != nil {
		handlerError(c, err)
		return
	}

	f, err := filter.New[filter.ReportFilter](c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := i.report.GetReports(c, pg, f)
	if err != nil {
		handlerError(c, err)
		return
	}
	handlePaginationJsonResponse(c, res, pg)
}
func (i *report) CreateReport(c *gin.Context) {
	input := service.CreateReportRequest{}
	if err := c.Bind(&input); err != nil {
		handlerError(c, errs.ErrBadRequest)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		handlerError(c, err)
		return
	}
	if err := i.report.CreateReport(c, &input, file); err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}
func (i *report) UpdateReport(c *gin.Context) {
	id, err := getReportId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	input := service.UpdateReportRequest{}
	if err := c.Bind(&input); err != nil {
		handlerError(c, errs.ErrBadRequest)
		return
	}
	file, _ := c.FormFile("file")

	if err := i.report.UpdateReport(c, id, &input, file); err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

func (i *report) DeleteReport(c *gin.Context) {
	id, err := getReportId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	if err := i.report.DeleteReport(c, id); err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

func NewReport(
	reportService service.Report,
) Report {
	return &report{
		report: reportService,
	}
}
