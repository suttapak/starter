package database

import (
	"fmt"
	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/internal/idx"
	"github.com/suttapak/starter/internal/model"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newGorm(conf *config.Config) (*gorm.DB, error) {
	if conf.DB.DSN == "" {
		return nil, fmt.Errorf("dsn is empty")

	}
	db, err := gorm.Open(postgres.Open(conf.DB.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func migrateDb(db *gorm.DB) error {
	err := db.AutoMigrate(&model.User{}, &model.Post{}, &gormadapter.CasbinRule{})
	return err
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
