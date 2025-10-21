package repository

import (
	"go.uber.org/fx"
)

var (
	Module = fx.Options(
		// Migrated repositories using sqlx
		fx.Provide(NewUser),
		fx.Provide(NewDatabaseTransaction),
		fx.Provide(NewImage),
		fx.Provide(NewAutoIncrementSequence),
		fx.Provide(NewODT),
		fx.Provide(NewMailRepository),
	)
)
