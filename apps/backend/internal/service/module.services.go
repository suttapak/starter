package service

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(newPost),
	fx.Provide(newAuth),
	fx.Provide(newJWT),
	fx.Provide(newEmailService),
	fx.Provide(newUserService),
)
