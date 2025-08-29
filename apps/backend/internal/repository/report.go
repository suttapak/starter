package repository

import (
	"context"

	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/model"
	"gorm.io/gorm"
)

type (
	Report interface {
		GetReport(ctx context.Context, tx *gorm.DB, id uint) (*model.ReportTemplate, error)
		GetReports(ctx context.Context, tx *gorm.DB, pg *helpers.Pagination, f *filter.ReportFilter) ([]model.ReportTemplate, error)
		CreateReport(ctx context.Context, tx *gorm.DB, m *model.ReportTemplate) error
		UpdateReport(ctx context.Context, tx *gorm.DB, id uint, m *model.ReportTemplate) error
		DeleteReport(ctx context.Context, tx *gorm.DB, id uint) error
	}

	report struct {
		db *gorm.DB
	}
)

func (i *report) GetReport(ctx context.Context, tx *gorm.DB, id uint) (*model.ReportTemplate, error) {
	if tx == nil {
		tx = i.db
	}
	var res model.ReportTemplate
	err := tx.Where("id = ?", id).First(&res).Error
	return &res, err
}

func (i *report) GetReports(ctx context.Context, tx *gorm.DB, pg *helpers.Pagination, f *filter.ReportFilter) ([]model.ReportTemplate, error) {
	if tx == nil {
		tx = i.db
	}
	var res []model.ReportTemplate
	tx = tx.Model(&model.ReportTemplate{}).Count(&pg.Count)
	helpers.Paging(pg)
	err := tx.Offset(pg.Offset).Limit(pg.Limit).Find(&res).Error
	return res, err
}

func (i *report) CreateReport(ctx context.Context, tx *gorm.DB, m *model.ReportTemplate) error {
	if tx == nil {
		tx = i.db
	}
	err := tx.Create(&m).Error
	return err
}

func (i *report) UpdateReport(ctx context.Context, tx *gorm.DB, id uint, m *model.ReportTemplate) error {
	if tx == nil {
		tx = i.db
	}
	err := tx.Model(&model.ReportTemplate{}).Where("id = ?", id).Updates(&m).Error
	return err
}

func (i *report) DeleteReport(ctx context.Context, tx *gorm.DB, id uint) error {
	if tx == nil {
		tx = i.db
	}
	input := model.ReportTemplate{}
	err := tx.Where("id = ?", id).Delete(&input).Error
	return err
}

func NewReport(db *gorm.DB) Report {
	return &report{
		db: db,
	}
}
