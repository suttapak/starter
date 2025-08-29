package repository

import "go.uber.org/fx"

var (
	Module = fx.Options(
		fx.Provide(newUserRepository),
		fx.Provide(newMailRepository),
		fx.Provide(newTeam),
		fx.Provide(NewDatabaseTransaction),
		fx.Provide(NewProducts),
		fx.Provide(NewProductCategory),
		fx.Provide(NewAutoIncrementSequence),
		fx.Provide(NewODT),
		fx.Provide(NewReport),
		fx.Provide(NewImage),
	)
)
