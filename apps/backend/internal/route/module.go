package route

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(newAuth),
	fx.Invoke(useAuth),
	fx.Provide(newUser),
	fx.Invoke(useUser),
	fx.Provide(newHealthCheck),
	fx.Invoke(useHealthCheck),
	fx.Provide(newTeam),
	fx.Invoke(useTeam),
	fx.Invoke(seedTeamPermission),
	fx.Invoke(seedProductPermission),
	fx.Provide(newProducts),
	fx.Invoke(useProducts),
	fx.Provide(newProductCategory),
	fx.Invoke(useProductCategory),
	fx.Invoke(seedProductCategoryPermission),
	fx.Provide(newReport),
	fx.Invoke(useReport),
)
