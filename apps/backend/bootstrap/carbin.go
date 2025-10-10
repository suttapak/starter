package bootstrap

import (
	_ "embed"

	"github.com/jmoiron/sqlx"
	"github.com/suttapak/starter/config"

	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

//go:embed carbin/authz_model.conf
var carbinModel string

func NewCarbin(cfg *config.Config, db *sqlx.DB) (*casbin.Enforcer, error) {
	model, err := model.NewModelFromString(carbinModel)
	if err != nil {
		return nil, err
	}
	a, err := sqladapter.NewAdapter(db.DB, "postgres", "casbin_rule")
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(model, a)
	if err != nil {
		return nil, err
	}
	return e, nil
}
