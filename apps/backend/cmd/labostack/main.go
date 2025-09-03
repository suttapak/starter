package main

import (
	"time"

	"github.com/suttapak/starter/bootstrap"
	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/database"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/i18n"
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/internal/route"
	"github.com/suttapak/starter/internal/service"
	"github.com/suttapak/starter/logger"
	"github.com/suttapak/starter/mail"

	"go.uber.org/fx"
)

func init() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server caller server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	fx.
		New(
			helpers.Module,
			logger.Module,
			config.Module,
			i18n.Module,
			database.Module,
			mail.Module,
			repository.Module,
			service.Module,
			middleware.Module,
			controller.Module,
			route.Module,
			bootstrap.Module,
		).
		Run()
}
