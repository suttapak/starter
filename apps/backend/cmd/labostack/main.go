package main

import (
	"labostack/boostrap"
	"labostack/config"
	"labostack/database"
	"labostack/internal/controller"
	"labostack/internal/middleware"
	"labostack/internal/repository"
	"labostack/internal/route"
	"labostack/internal/service"
	"labostack/logger"
	"labostack/mail"

	"go.uber.org/fx"
)

func main() {
	fx.
		New(
			logger.Module,
			config.Module,
			database.Module,
			mail.Module,
			repository.Module,
			service.Module,
			middleware.Module,
			controller.Module,
			route.Module,
			boostrap.Module,
		).
		Run()
}
