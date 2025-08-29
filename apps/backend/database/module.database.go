package database

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(newGorm),
	fx.Invoke(migrateDb),
	fx.Invoke(seedRole),
	fx.Invoke(seedTeamRole),
	fx.Invoke(seedReportJsonSchemaType),
)
