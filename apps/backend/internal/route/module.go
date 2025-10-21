package route

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Invoke(UseAuth),
	fx.Invoke(UseUser),
	fx.Invoke(UseHealthCheck),
)
