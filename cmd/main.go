package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github/fims-proto/fims-proto-ms/docs"
	"github/fims-proto/fims-proto-ms/internal/common/database"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/common/localization"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/devops"
	generalLedgerAdapter "github/fims-proto/fims-proto-ms/internal/general_ledger/adapter/db"
	generalLedgerNumberingAdapter "github/fims-proto/fims-proto-ms/internal/general_ledger/adapter/numbering"
	generalLedgerSobAdapter "github/fims-proto/fims-proto-ms/internal/general_ledger/adapter/sob"
	generalLedgerUserAdapter "github/fims-proto/fims-proto-ms/internal/general_ledger/adapter/user"
	generalLedgerApp "github/fims-proto/fims-proto-ms/internal/general_ledger/app"
	generalLedgerPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/general_ledger/port/private/http"
	generalLedgerIntraPort "github/fims-proto/fims-proto-ms/internal/general_ledger/port/private/intraprocess"
	generalLedgerPublicHttpPort "github/fims-proto/fims-proto-ms/internal/general_ledger/port/public/http"
	numberingAdapter "github/fims-proto/fims-proto-ms/internal/numbering/adapter/db"
	numberingApp "github/fims-proto/fims-proto-ms/internal/numbering/app"
	numberingPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/numbering/port/private/http"
	numberingIntraPort "github/fims-proto/fims-proto-ms/internal/numbering/port/private/intraprocess"
	sobAdapter "github/fims-proto/fims-proto-ms/internal/sob/adapter/db"
	sobGeneralLedgerAdapter "github/fims-proto/fims-proto-ms/internal/sob/adapter/general_ledger"
	sobApp "github/fims-proto/fims-proto-ms/internal/sob/app"
	sobPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/sob/port/private/http"
	sobIntraPort "github/fims-proto/fims-proto-ms/internal/sob/port/private/intraprocess"
	sobPublicHttpPort "github/fims-proto/fims-proto-ms/internal/sob/port/public/http"
	tenantDb "github/fims-proto/fims-proto-ms/internal/tenant/adapter/db"
	tenantApp "github/fims-proto/fims-proto-ms/internal/tenant/app"
	ginMiddleware "github/fims-proto/fims-proto-ms/internal/tenant/lib/gin-middleware"
	tenantManager "github/fims-proto/fims-proto-ms/internal/tenant/lib/tenant-manager"
	tenantService "github/fims-proto/fims-proto-ms/internal/tenant/lib/tenant-service"
	tenantIntraPort "github/fims-proto/fims-proto-ms/internal/tenant/port/private/intraprocess"
	userAdapter "github/fims-proto/fims-proto-ms/internal/user/adapter/db"
	userApp "github/fims-proto/fims-proto-ms/internal/user/app"
	userPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/user/port/private/http"
	userIntraPort "github/fims-proto/fims-proto-ms/internal/user/port/private/intraprocess"
	userPublicHttpPort "github/fims-proto/fims-proto-ms/internal/user/port/public/http"
)

func main() {
	flag.Parse()

	loadConfig()

	log.InitLogger()
	defer cleanup()

	dbConnector := database.NewConnector()

	dbConnection, err := dbConnector.Open(viper.GetString("postgres.dsn"))
	if err != nil {
		panic(errors.Wrapf(err, "failed to initialize dbConnection connection"))
	}

	tenantPostgresRepository := tenantDb.NewTenantPostgresRepository(dbConnection)
	tenantApplication := tenantApp.NewApplication(tenantPostgresRepository)
	tenantInterface := tenantIntraPort.NewTenantInterface(&tenantApplication)
	tenantServiceImpl := tenantService.NewTenantService(tenantInterface)

	tenantManagerImpl := tenantManager.NewTenantManager(tenantServiceImpl, dbConnector)

	// i18n
	localizer := localization.NewLocalizer("./i18n", "zh-CN")

	// repositories
	sobRepository := sobAdapter.NewSobPostgresRepository()
	generalLedgerRepository := generalLedgerAdapter.NewGeneralLedgerPostgresRepository()
	numberingRepository := numberingAdapter.NewNumberingPostgresRepository()
	userRepository := userAdapter.NewUserPostgresRepository()

	// application - will be passed by reference, in order to make injection work
	sobApplication := sobApp.NewApplication()
	generalLedgerApplication := generalLedgerApp.NewApplication()
	numberingApplication := numberingApp.NewApplication()
	userApplication := userApp.NewApplication()

	// intra process interfaces
	sobInterface := sobIntraPort.NewSobInterface(&sobApplication)
	generalLedgerInterface := generalLedgerIntraPort.NewGeneralLedgerInterface(&generalLedgerApplication)
	numberingInterface := numberingIntraPort.NewNumberingInterface(&numberingApplication)
	userInterface := userIntraPort.NewUserInterface(&userApplication)

	// application dependencies injection
	generalLedgerServiceForSob := sobGeneralLedgerAdapter.NewIntraProcessAdapter(generalLedgerInterface)
	sobApplication.Inject(
		sobRepository,
		sobRepository,
		generalLedgerServiceForSob,
	)

	sobServiceForGeneralLedger := generalLedgerSobAdapter.NewIntraProcessAdapter(sobInterface)
	numberingServiceForGeneralLedger := generalLedgerNumberingAdapter.NewIntraProcessAdapter(numberingInterface)
	userServiceForGeneralLedger := generalLedgerUserAdapter.NewIntraProcessAdapter(userInterface)
	generalLedgerApplication.Inject(
		generalLedgerRepository,
		generalLedgerRepository,
		sobServiceForGeneralLedger,
		numberingServiceForGeneralLedger,
		userServiceForGeneralLedger,
	)

	numberingApplication.Inject(
		numberingRepository,
		numberingRepository,
	)

	userApplication.Inject(
		userRepository,
		userRepository,
	)

	log.InfoWithoutCxt("All module applications initiated")

	router := gin.Default()
	router.GET("/health/ping", func(c *gin.Context) { c.String(http.StatusOK, "Pong") })
	router.Use(ginMiddleware.ResolveTenantBySubdomain(tenantManagerImpl))
	router.Use(commonErrors.ErrorHandler(localizer))

	// public http API
	publicApiRouter := router.Group("/api/v1")
	sobPublicHttpPort.InitRouter(sobPublicHttpPort.NewHandler(&sobApplication), publicApiRouter)
	generalLedgerPublicHttpPort.InitRouter(generalLedgerPublicHttpPort.NewHandler(&generalLedgerApplication), publicApiRouter)
	userPublicHttpPort.InitRouter(userPublicHttpPort.NewHandler(&userApplication), publicApiRouter)

	// private http API, should have different authentication method then public API
	privateApiRouter := router.Group("/internal")
	sobPrivateHttpPort.InitRouter(sobPrivateHttpPort.NewHandler(&sobApplication), privateApiRouter)
	numberingPrivateHttpPort.InitRouter(numberingPrivateHttpPort.NewHandler(&numberingApplication), privateApiRouter)
	generalLedgerPrivateHttpPort.InitRouter(generalLedgerPrivateHttpPort.NewHandler(&generalLedgerApplication), privateApiRouter)
	userPrivateHttpPort.InitRouter(userPrivateHttpPort.NewHandler(&userApplication), privateApiRouter)

	if strings.HasPrefix(viper.GetString("profile"), "dev") {
		// gin-swagger
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		// devops
		devopsApiRouter := router.Group("/devops/")
		devops.InitJwtHandler(devopsApiRouter)
	}

	log.InfoWithoutCxt("All module routers initiated")

	log.InfoWithoutCxt("Starting gin engine...")
	if err := router.Run(":" + viper.GetString("app.port")); err != nil {
		panic(err.Error())
	}
}

func loadConfig() {
	// environment variables
	if err := viper.BindEnv("profile", "PROFILE"); err != nil {
		panic(errors.Wrap(err, "failed to bind ENV profile"))
	}
	viper.SetDefault("profile", "dev")

	_ = viper.BindEnv("postgres.dsn", "DSN")

	// read config
	profile := viper.GetString("profile")
	viper.SetConfigName(fmt.Sprintf("application-%s", profile))
	viper.AddConfigPath("./config/")
	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "failed to load config file"))
	}

	// check mandatory and set defaults:
	checkResult := bytes.Buffer{}
	// app
	viper.SetDefault("app.port", "5002")
	// postgres
	if !viper.IsSet("postgres.dsn") {
		checkResult.WriteString("postgres.dsn; ")
	}
	// logger
	viper.SetDefault("logger.debug", false)
	viper.SetDefault("logger.jsonEncoding", true)
	viper.SetDefault("logger.showSql", false)

	if checkResult.Len() > 0 {
		panic("config missing: " + checkResult.String())
	}
}

func cleanup() {
	log.InfoWithoutCxt("fims terminating...")
	log.SyncLogger()
}
