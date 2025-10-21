package bootstrap

import (
	_ "embed"

	"gorm.io/gorm"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/suttapak/starter/domain/config"
)

//go:embed carbin/authz_model.conf
var carbinModel string

func NewCarbin(cfg *config.Config, db *gorm.DB) (*casbin.Enforcer, error) {
	model, err := model.NewModelFromString(carbinModel)
	if err != nil {
		return nil, err
	}
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(model, a)
	if err != nil {
		return nil, err
	}
	return e, nil
}
