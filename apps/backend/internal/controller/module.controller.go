package controller

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(newPost),
	fx.Provide(NewAuth),
	fx.Provide(NewUser),
	fx.Provide(NewTeam),
)

type (
	Response[T any] struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
		Data    T      `json:"data"`
		Meta    any    `json:"meta"`
	}
)
