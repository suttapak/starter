package main

import (
	"github.com/suttapak/starter/boostrap"
	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/database"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/internal/route"
	"github.com/suttapak/starter/internal/service"
	"github.com/suttapak/starter/logger"
	"github.com/suttapak/starter/mail"

	"go.uber.org/fx"
)

func main() {
	fx.
		New(
			helpers.Module,
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
