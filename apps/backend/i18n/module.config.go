package i18n

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewI18N),
)
