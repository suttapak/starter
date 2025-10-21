package database

import (
	"github.com/suttapak/starter/domain/database/seed"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewGorm),
	fx.Invoke(MigrateDb),
	fx.Invoke(seed.SeedRole),
)
