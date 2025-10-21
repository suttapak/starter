package database

import (
	"fmt"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/suttapak/starter/domain/config"
	"github.com/suttapak/starter/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGorm(conf *config.Config) (*gorm.DB, error) {
	if conf.DB.DSN == "" {
		return nil, fmt.Errorf("dsn is empty")

	}
	db, err := gorm.Open(postgres.Open(conf.DB.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MigrateDb(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.AutoIncrementSequence{},
		&model.Image{},
		&model.Role{},
		&model.User{},
		&model.ProfileImage{},
		&gormadapter.CasbinRule{},
	)
	return err
}
