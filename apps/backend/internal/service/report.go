package service

import (
	"bytes"
	"context"
	"mime/multipart"
	"strings"

	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/idx"
	"github.com/suttapak/starter/internal/model"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/logger"
)

type (
	Report interface {
		GetReport(ctx context.Context, id uint) (*ReportResponse, error)
		GetReports(ctx context.Context, pg *helpers.Pagination, f *filter.ReportFilter) ([]ReportResponse, error)
		UpdateReport(ctx context.Context, id uint, input *UpdateReportRequest, file *multipart.FileHeader) error
		DeleteReport(ctx context.Context, id uint) error
		CreateReport(ctx context.Context, req *CreateReportRequest, file *multipart.FileHeader) error
	}
	report struct {
		logger logger.AppLogger
		helper helpers.Helper
		report repository.Report
		odt    repository.Odt
		conf   *config.Config
	}

	UpdateReportRequest struct {
		Name string `json:"name" form:"name"`
	}

	ReportResponse struct {
		CommonModel
		Code                   string                   `json:"code"`
		Name                   string                   `json:"name"`
		DisplayName            string                   `json:"display_name"`
		Icon                   string                   `json:"icon"`
		ReportJsonSchemaTypeID idx.ReportJsonSchemaType `json:"report_json_schema_type_id"`
		ReportJsonSchemaType   ReportJsonSchemaType     `json:"report_json_schema_type"`
	}

	ReportJsonSchemaType struct {
		CommonModel
		Name string `json:"name"`
	}

	CreateReportRequest struct {
		Name string `json:"name" form:"name"`
	}
	RenderResponse struct {
		Url     string `json:"url"`
		Message string `json:"message"`
	}
)

func (i *report) GetReport(ctx context.Context, id uint) (*ReportResponse, error) {
	model, err := i.report.GetReport(ctx, nil, id)
	if err != nil {
		i.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	var res ReportResponse
	if err := i.helper.ParseJson(model, &res); err != nil {
		return nil, errs.ErrInvalid
	}
	return &res, nil
}

func (i *report) GetReports(ctx context.Context, pg *helpers.Pagination, f *filter.ReportFilter) ([]ReportResponse, error) {
	models, err := i.report.GetReports(ctx, nil, pg, f)
	if err != nil {
		i.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	var res []ReportResponse
	if err := i.helper.ParseJson(models, &res); err != nil {
		return nil, errs.ErrInvalid
	}
	return res, nil
}

func (i *report) CreateReport(ctx context.Context, req *CreateReportRequest, file *multipart.FileHeader) error {
	f, err := file.Open()
	if err != nil {
		i.logger.Error(err)
		return errs.ErrBadRequest
	}
	defer f.Close()
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(f); err != nil {
		i.logger.Error(err)
		return errs.ErrBadRequest
	}

	reader := bytes.NewReader(buf.Bytes())
	id, err := i.odt.UploadTemplate(reader)
	if err != nil {
		i.logger.Error(err)
		return errs.ErrInternal
	}
	modelData := model.ReportTemplate{
		Code:                   id,
		Name:                   req.Name,
		DisplayName:            req.Name,
		Icon:                   "file",
		ReportJsonSchemaTypeID: idx.ReportJsonSchemaTypeCommonId,
	}
	if err := i.report.CreateReport(ctx, nil, &modelData); err != nil {
		i.logger.Error(err)
	}
	return nil
}
func (i *report) UpdateReport(ctx context.Context, id uint, input *UpdateReportRequest, file *multipart.FileHeader) error {
	if file != nil {
		report, err := i.report.GetReport(ctx, nil, id)
		if err != nil {
			i.logger.Error(err)
			return errs.HandleGorm(err)
		}
		f, err := file.Open()
		if err != nil {
			i.logger.Error(err)
			return errs.ErrBadRequest
		}
		defer f.Close()
		buf := bytes.NewBuffer(nil)
		if _, err := buf.ReadFrom(f); err != nil {
			i.logger.Error(err)
			return errs.ErrBadRequest
		}

		reader := bytes.NewReader(buf.Bytes())
		if err := i.odt.UpdateTemplate(report.Code, reader); err != nil {
			i.logger.Error(err)
			return errs.HandleGorm(err)
		}
	}
	body := model.ReportTemplate{
		Name:        input.Name,
		DisplayName: strings.TrimSpace(input.Name),
	}
	if err := i.report.UpdateReport(ctx, nil, id, &body); err != nil {
		i.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}
func (i *report) DeleteReport(ctx context.Context, id uint) error {
	if err := i.report.DeleteReport(ctx, nil, id); err != nil {
		i.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}

func NewReport(
	logger logger.AppLogger,
	helper helpers.Helper,
	reportRepo repository.Report,
	odt repository.Odt,
	conf *config.Config,
) Report {
	return &report{
		logger: logger,
		helper: helper,
		report: reportRepo,
		odt:    odt,
		conf:   conf,
	}
}
