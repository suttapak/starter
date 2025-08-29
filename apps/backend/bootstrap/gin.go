package bootstrap

import (
	"context"
	"net"
	"net/http"

	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/i18n"
	"github.com/suttapak/starter/logger"

	_ "github.com/suttapak/starter/cmd/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

func newGin(conf *config.Config, log logger.AppLogger) *gin.Engine {
	r := gin.Default()

	if conf.PPROF.ENABLE {
		pprof.Register(r)
	}

	r.Use(logger.GinLogger(log))
	// set lang
	r.Use(i18n.SetLocal)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     conf.CORS.ALLOW_ORIGIN,
		AllowMethods:     conf.CORS.ALLOW_METHODS,
		AllowHeaders:     conf.CORS.ALLOW_HEADERS,
		ExposeHeaders:    conf.CORS.EXPOSE_HEADERS,
		AllowCredentials: conf.CORS.ALLOW_CREDENTIALS,
	}))

	r.Static("/static", "./public/static")
	r.Static("/public/static", "./public/static")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func useGin(
	lc fx.Lifecycle,
	app *gin.Engine,
) {
	srv := &http.Server{Addr: ":8080", Handler: app}
	lc.Append(fx.Hook{

		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			go func() {
				err := srv.Serve(ln)
				if err != nil {
					panic(err)
				}
			}()
			return nil

		},
		OnStop: func(ctx context.Context) error {
			err := srv.Shutdown(ctx)
			if err != nil {
				return err
			}
			return nil
		},
	})
}
