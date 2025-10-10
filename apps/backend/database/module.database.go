package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/suttapak/starter/logger"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(newSqlx),
	fx.Invoke(func(db *sqlx.DB, log logger.AppLogger) error {
		return RunMigrations(db.DB, log)
	}),
)
