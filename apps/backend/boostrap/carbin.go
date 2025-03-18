package boostrap

import (
	"github.com/suttapak/starter/config"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func newCarbin(cfg *config.Config, db *gorm.DB) (*casbin.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(cfg.CARBIN.MODEL, a)
	if err != nil {
		return nil, err
	}
	return e, nil
}
