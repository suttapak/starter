package service

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewAuth),
	fx.Provide(newJWT),
	fx.Provide(newEmailService),
	fx.Provide(newUserService),
	fx.Provide(NewTeam),
	fx.Provide(NewProducts),
	fx.Provide(NewProductCategory),
	fx.Provide(NewCodeService),
	fx.Provide(NewReport),
	fx.Provide(NewExcelService),
	fx.Provide(NewImageFileService),
)
