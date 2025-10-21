package service

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewAuth),
	fx.Provide(NewJWT),
	fx.Provide(NewEmail),
	fx.Provide(NewUser),
	fx.Provide(NewCodeService),
	fx.Provide(NewExcelService),
	fx.Provide(NewImageFileService),
)
