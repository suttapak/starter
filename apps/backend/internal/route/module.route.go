package route

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(newPost),
	fx.Invoke(usePost),
	fx.Provide(newAuth),
	fx.Invoke(useAuth),
)
