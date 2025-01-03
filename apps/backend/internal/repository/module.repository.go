package repository

import "go.uber.org/fx"

var (
	Module = fx.Options(
		fx.Provide(newPost),
		fx.Provide(newUserRepository),
		fx.Provide(newMailRepository),
	)
)
