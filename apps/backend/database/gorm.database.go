package database

import (
	"fmt"

	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/internal/idx"
	"github.com/suttapak/starter/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func newGorm(conf *config.Config) (*gorm.DB, error) {
	if conf.DB.DSN == "" {
		return nil, fmt.Errorf("dsn is empty")

	}
	db, err := gorm.Open(postgres.Open(conf.DB.DSN), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func migrateDb(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Image{},
		&model.User{},
		&model.ProfileImage{},
		&model.Team{},
		&model.TeamMember{},
		&model.TeamRole{},
		&model.Product{},
		&model.ProductCategory{},
		&model.ProductImage{},
		&model.ProductProductCategory{},
		&model.AutoIncrementSequence{},
		&model.ReportJsonSchemaType{},
		&model.ReportTemplate{},
	)

	return err
}

func seedReportJsonSchemaType(db *gorm.DB) error {
	var (
		count int64
	)

	if err := db.Model(&model.ReportJsonSchemaType{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		var types = []*model.ReportJsonSchemaType{{
			CommonModel: model.CommonModel{ID: uint(idx.ReportJsonSchemaTypeCommonId)},
			Name:        "Common",
		},
		}
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&types).Error; err != nil {
			return err
		}
	}
	return nil
}
func seedRole(db *gorm.DB) error {
	var (
		count int64
	)

	if err := db.Model(&model.Role{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		var roles = []*model.Role{{
			CommonModel: model.CommonModel{ID: idx.RoleUser},
			Name:        "User",
		}, {
			CommonModel: model.CommonModel{ID: idx.RoleModerator},
			Name:        "Moderator",
		}, {
			CommonModel: model.CommonModel{ID: idx.RoleAdmin},
			Name:        "Admin",
		}, {
			CommonModel: model.CommonModel{ID: idx.RoleSuperAdmin},
			Name:        "SuperAdmin",
		}}
		if err := db.Create(&roles).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedTeamRole(db *gorm.DB) error {
	var (
		count int64
	)
	if err := db.Model(&model.TeamRole{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		var roles = []*model.TeamRole{
			{
				CommonModel: model.CommonModel{ID: idx.TeamRoleOwnerID},
				Name:        "Owner",
			},
			{
				CommonModel: model.CommonModel{ID: idx.TeamRoleAdminID},
				Name:        "Admin",
			},
			{
				CommonModel: model.CommonModel{ID: idx.TeamRoleMemberID},
				Name:        "Member",
			},
		}
		if err := db.Create(&roles).Error; err != nil {
			return err
		}
	}
	return nil
}
