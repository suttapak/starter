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
	// i18n
	i18nService, err := i18n.NewI18N(conf, log)
	require.NoError(t, err, "Failed to Init I18n")

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
	productRepo := repository.NewProducts(db)
	productCategoryRepo := repository.NewProductCategory(db)
	sequenceRepository := repository.NewAutoIncrementSequence(db)

	mail.On("Send").Return(nil)

	// Setup services
	jwtService := service.NewJWT(log, conf, helper)
	codeService := service.NewCodeService(log, sequenceRepository)
	excelService := service.NewExcelService()
	mailService := service.NewEmail(mail)
	authService := service.NewAuth(log, userRepo, jwtService, helper, mailService, conf)
	imageService := service.NewImageFileService()
	userService := service.NewUser(userRepo, imageRepo, imageService, log, helper)
	teamService := service.NewTeam(log, teamRepo, helper, mailService, jwtService, conf, userRepo)
	productService := service.NewProducts(log, helper, dbTx, productRepo, codeService, excelService, i18nService, imageRepo, imageService)
	productCategoryService := service.NewProductCategory(log, helper, productCategoryRepo)

	// Setup controllers
	authController := controller.NewAuth(authService, conf)
	userController := controller.NewUser(userService)
	teamController := controller.NewTeam(teamService)
	productController := controller.NewProducts(productService)
	productCategoryController := controller.NewProductCategory(productCategoryService)

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
	// Auth routes
	authGroup := router.Group("auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/refresh", authGuard.ProtectRefreshToken, authController.RefreshToken)
		authGroup.POST("/logout", authController.Logout)
		authGroup.GET("/email/verify", authController.VerifyEmail)
		authGroup.POST("/email/send-verify", authGuard.Protect, authController.SendVerifyEmail)
	}

	// User routes
	userGroup := router.Group("users", authGuard.Protect)
	{
		userGroup.GET("/:id", userController.GetUserById)
		userGroup.GET("/me", userController.GetUserMe)
		userGroup.GET("/by-username", userController.FindUserByUsername)
		userGroup.GET("/verify-email", userController.CheckUserIsVerifyEmail)
		userGroup.POST("/profile-image", userController.CreateProfileImage)
	}

	// Team routes
	teamGroup := router.Group("teams", authGuard.Protect)
	{
		teamGroup.GET("/me", teamController.GetTeamsMe)
		teamGroup.GET("/", teamController.GetTeamsFilter)
		teamGroup.POST("/", teamController.Create)
		teamGroup.GET("/join", teamController.JoinTeamWithToken)
		teamGroup.POST("/join/link", teamController.JoinWithShearLink)
		teamGroup.POST("/:team_id/request-join", teamController.CreateTeamPendingTeamMember)
	}

	// Team routes with permission check
	teamPermGroup := router.Group("teams", authGuard.Protect, authGuard.TeamPermission)
	{
		teamPermGroup.GET("/:team_id", teamController.GetTeamByTeamId)
		teamPermGroup.GET("/:team_id/member-count", teamController.GetTeamMemberCount)
		teamPermGroup.GET("/:team_id/members", teamController.GetTeamMembers)
		teamPermGroup.GET("/:team_id/pending-member-count", teamController.GetPendingTeamMemberCount)
		teamPermGroup.GET("/:team_id/pending-members", teamController.GetPendingTeamMembers)
		teamPermGroup.GET("/:team_id/user-me", teamController.GetTeamUserMe)
		teamPermGroup.PUT("/:team_id/member-role", teamController.UpdateMemberRole)
		teamPermGroup.POST("/:team_id/pending-member", teamController.SendInviteTeamMember)
		teamPermGroup.POST("/:team_id/shared-link", teamController.CreateShearLink)
		teamPermGroup.POST("/:team_id/accept", teamController.AcceptTeamMember)
		teamPermGroup.PUT("/:team_id", teamController.UpdateTeamInfo)
	}

	// Product routes
	productGroup := router.Group("teams/:team_id/products", authGuard.Protect, authGuard.TeamPermission)
	{
		productGroup.GET("/:products_id", productController.GetProduct)
		productGroup.GET("", productController.GetProducts)
		productGroup.POST("", productController.CreateProducts)
		productGroup.PUT("/:products_id", productController.UpdateProducts)
		productGroup.POST("/:products_id/upload_image", productController.UploadProductImages)
		productGroup.DELETE("/:products_id", productController.DeleteProducts)
		productGroup.DELETE("/:products_id/product_image/:product_image_id", productController.DeleteProductImage)
	}

	// Product Category routes
	productCategoryGroup := router.Group("teams/:team_id/product_category", authGuard.Protect, authGuard.TeamPermission)
	{
		productCategoryGroup.GET("/:product_category_id", productCategoryController.GetProductCategory)
		productCategoryGroup.GET("", productCategoryController.GetProductCategories)
		productCategoryGroup.POST("", productCategoryController.CreateProductCategory)
		productCategoryGroup.PUT("/:product_category_id", productCategoryController.UpdateProductCategory)
		productCategoryGroup.DELETE("/:product_category_id", productCategoryController.DeleteProductCategory)
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
