package controller

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewAuth),
	fx.Provide(NewUser),
	fx.Provide(NewTeam),
	fx.Provide(NewProducts),
	fx.Provide(NewProductCategory),
	fx.Provide(NewReport),
)

type (
	Response[T any] struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
		Data    T      `json:"data"`
		Meta    any    `json:"meta"`
	}
)
