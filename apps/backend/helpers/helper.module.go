package helpers

import "go.uber.org/fx"

var (
	Module = fx.Options(
		fx.Provide(NewHelper),
	)
)
