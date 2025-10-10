package integration_test

import (
	"os"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/suttapak/starter/bootstrap"
	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/i18n"
	"github.com/suttapak/starter/internal/controller"
	"github.com/suttapak/starter/internal/middleware"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/internal/service"
	"github.com/suttapak/starter/internal/testutil"
	"github.com/suttapak/starter/logger"
)

// TestServer holds the test server dependencies
type TestServer struct {
	Router         *gin.Engine
	DB             *sqlx.DB
	Config         *config.Config
	Logger         logger.AppLogger
	AuthController controller.Auth
}

// SetupTestServer creates a new test server with all dependencies
func SetupTestServer(t *testing.T) *TestServer {
	// Set test mode
	testutil.SetupTestMode()
	os.Setenv("CONFIG_PATH", "../../configs.test.toml")

	// Setup database
	db := testutil.SetupTestDB(t)
	testutil.SeedBasicData(t, db)

	// Load config
	// Create a test config
	conf := &config.Config{
		DB: config.DB{
			DSN: "host=localhost user=test_user password=test_password dbname=test_db port=5433 sslmode=disable TimeZone=Asia/Bangkok",
		},
		JWT: config.JWT{
			SECRET:         "test_secret",
			REFRESH_SECRET: "test_refresh_secret",
			EMAIL_SECRET:   "test_email_secret",
		},
		SERVER: config.SERVER{
			HOST:      "localhost",
			HOST_NAME: "http://localhost:8080",
			PORT:      "3000",
		},
		CORS: config.CORS{
			ALLOW_ORIGIN:      []string{"http://localhost:8080"},
			ALLOW_METHODS:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			ALLOW_HEADERS:     []string{"Authorization"},
			ALLOW_CREDENTIALS: true,
		},
		CARBIN: config.CARBIN{
			MODEL: "./carbin/authz_model.conf",
		},
	}

	// Setup logger - create simple test logger
	log := logger.NewLoggerMock()

	// Setup helpers
	helper := helpers.NewHelper()

	cas, err := bootstrap.NewCarbin(conf, db)
	require.NoError(t, err, "Failed to Init Casbin")

	// Setup repositories
	userRepo := repository.NewUser(db)
	dbTx := repository.NewDatabaseTransaction(db)
	mail := repository.NewMailRepositoryMock()
	imageRepo := repository.NewImage(db)
	teamRepo := repository.NewTeam(db)

	mail.On("Send").Return(nil)

	// Setup services
	jwtService := service.NewJWT(log, conf, helper)
	mailService := service.NewEmail(mail)
	authService := service.NewAuth(log, userRepo, jwtService, helper, mailService, conf)
	imageService := service.NewImageFileService()
	userService := service.NewUser(userRepo, imageRepo, imageService, log, helper)
	teamService := service.NewTeam(log, teamRepo, helper, mailService, jwtService, conf, userRepo)

	// Setup controllers
	authController := controller.NewAuth(authService, conf)

	// Setup middleware
	authGuard := middleware.NewAuthGuardMiddleware(jwtService, cas, log, userService, teamService, teamRepo)

	// Setup router (simple Gin setup instead of using bootstrap)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     conf.CORS.ALLOW_ORIGIN,
		AllowMethods:     conf.CORS.ALLOW_METHODS,
		AllowHeaders:     conf.CORS.ALLOW_HEADERS,
		AllowCredentials: conf.CORS.ALLOW_CREDENTIALS,
	}))
	router.Use(i18n.SetLocal)

	// Setup routes
	authGroup := router.Group("auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/refresh", authGuard.ProtectRefreshToken, authController.RefreshToken)
		authGroup.POST("/logout", authController.Logout)
		authGroup.GET("/email/verify", authController.VerifyEmail)
		authGroup.POST("/email/send-verify", authGuard.Protect, authController.SendVerifyEmail)
	}

	_ = dbTx // Keep for transaction tests

	return &TestServer{
		Router:         router,
		DB:             db,
		Config:         conf,
		Logger:         log,
		AuthController: authController,
	}
}

// Teardown cleans up test server resources
func (ts *TestServer) Teardown(t *testing.T) {
	testutil.TeardownTestDB(t, ts.DB)
}
