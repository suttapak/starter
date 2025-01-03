package boostrap

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(newGin),
	fx.Provide(newCarbin),
	fx.Invoke(useGin),
)
