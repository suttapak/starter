package bootstrap

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(newGin),
	fx.Provide(NewCarbin),
	fx.Invoke(useGin),
)
