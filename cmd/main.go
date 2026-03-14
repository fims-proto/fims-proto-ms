package main

import (
	"net/http"
	"strings"

	_ "github/fims-proto/fims-proto-ms/docs/swagger_generated"
	"github/fims-proto/fims-proto-ms/internal/common/config"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	dedicatedDatasource "github/fims-proto/fims-proto-ms/internal/common/datasource/dedicated-datasource"
	multitenantDatasource "github/fims-proto/fims-proto-ms/internal/common/datasource/multitenant-datasource"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/common/localization"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/devops"
	dimensionAdapter "github/fims-proto/fims-proto-ms/internal/dimension/adapter/db"
	dimensionApp "github/fims-proto/fims-proto-ms/internal/dimension/app"
	dimensionPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/dimension/port/private/http"
	dimensionIntraPort "github/fims-proto/fims-proto-ms/internal/dimension/port/private/intraprocess"
	dimensionPublicHttpPort "github/fims-proto/fims-proto-ms/internal/dimension/port/public/http"
	generalLedgerAdapter "github/fims-proto/fims-proto-ms/internal/general_ledger/adapter/db"
	generalLedgerDimensionAdapter "github/fims-proto/fims-proto-ms/internal/general_ledger/adapter/dimension"
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
	reportAdapter "github/fims-proto/fims-proto-ms/internal/report/adapter/db"
	reportApp "github/fims-proto/fims-proto-ms/internal/report/app"
	reportPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/report/port/private/http"
	reportIntraPort "github/fims-proto/fims-proto-ms/internal/report/port/private/intraprocess"
	reportPublicHttpPort "github/fims-proto/fims-proto-ms/internal/report/port/public/http"
	sobAdapter "github/fims-proto/fims-proto-ms/internal/sob/adapter/db"
	sobGeneralLedgerAdapter "github/fims-proto/fims-proto-ms/internal/sob/adapter/general_ledger"
	sobReportAdapter "github/fims-proto/fims-proto-ms/internal/sob/adapter/report"
	sobApp "github/fims-proto/fims-proto-ms/internal/sob/app"
	sobPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/sob/port/private/http"
	sobIntraPort "github/fims-proto/fims-proto-ms/internal/sob/port/private/intraprocess"
	sobPublicHttpPort "github/fims-proto/fims-proto-ms/internal/sob/port/public/http"
	userAdapter "github/fims-proto/fims-proto-ms/internal/user/adapter/db"
	userApp "github/fims-proto/fims-proto-ms/internal/user/app"
	userPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/user/port/private/http"
	userIntraPort "github/fims-proto/fims-proto-ms/internal/user/port/private/intraprocess"
	userPublicHttpPort "github/fims-proto/fims-proto-ms/internal/user/port/public/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	defer cleanup()

	config.Initialize()

	log.Initialize()

	// i18n
	localizer := localization.NewLocalizer("./i18n", "zh-CN")

	var dataSource datasource.DataSource
	if config.GetBool("app.multiTenancy") {
		dataSource = multitenantDatasource.NewMultiTenantDataSource()
	} else {
		dataSource = dedicatedDatasource.NewDedicatedDataSource()
	}

	// repositories
	sobRepository := sobAdapter.NewSobPostgresRepository(dataSource)
	generalLedgerRepository := generalLedgerAdapter.NewGeneralLedgerPostgresRepository(dataSource)
	generalLedgerReadRepository := generalLedgerAdapter.NewGeneralLedgerPostgresReadRepository(dataSource)
	generalLedgerServiceRepository := reportAdapter.NewGeneralLedgerPostgresService(dataSource)
	reportRepository := reportAdapter.NewReportPostgresRepository(dataSource)
	reportReadRepository := reportAdapter.NewReportPostgresReadRepository(dataSource)
	numberingRepository := numberingAdapter.NewNumberingPostgresRepository(dataSource)
	userRepository := userAdapter.NewUserPostgresRepository(dataSource)
	dimensionRepository := dimensionAdapter.NewDimensionPostgresRepository(dataSource)
	dimensionReadRepository := dimensionAdapter.NewDimensionPostgresReadRepository(dataSource)

	// application - will be passed by reference, in order to make injection work
	sobApplication := sobApp.NewApplication()
	generalLedgerApplication := generalLedgerApp.NewApplication()
	numberingApplication := numberingApp.NewApplication()
	reportApplication := reportApp.NewApplication()
	userApplication := userApp.NewApplication()
	dimensionApplication := dimensionApp.NewApplication()

	// intra process interfaces
	sobInterface := sobIntraPort.NewSobInterface(&sobApplication)
	generalLedgerInterface := generalLedgerIntraPort.NewGeneralLedgerInterface(&generalLedgerApplication)
	numberingInterface := numberingIntraPort.NewNumberingInterface(&numberingApplication)
	reportInterface := reportIntraPort.NewReportInterface(&reportApplication)
	userInterface := userIntraPort.NewUserInterface(&userApplication)
	dimensionInterface := dimensionIntraPort.NewDimensionInterface(&dimensionApplication)

	// application dependencies injection
	generalLedgerServiceForSob := sobGeneralLedgerAdapter.NewIntraProcessAdapter(generalLedgerInterface)
	reportServiceForSob := sobReportAdapter.NewIntraProcessAdapter(reportInterface)
	sobApplication.Inject(
		sobRepository,
		sobRepository,
		generalLedgerServiceForSob,
		reportServiceForSob,
	)

	dimensionApplication.Inject(dimensionRepository, dimensionReadRepository)

	sobServiceForGeneralLedger := generalLedgerSobAdapter.NewIntraProcessAdapter(sobInterface)
	numberingServiceForGeneralLedger := generalLedgerNumberingAdapter.NewIntraProcessAdapter(numberingInterface)
	userServiceForGeneralLedger := generalLedgerUserAdapter.NewIntraProcessAdapter(userInterface)
	dimensionServiceForGeneralLedger := generalLedgerDimensionAdapter.NewIntraProcessAdapter(dimensionInterface)
	generalLedgerApplication.Inject(
		generalLedgerRepository,
		generalLedgerReadRepository,
		sobServiceForGeneralLedger,
		numberingServiceForGeneralLedger,
		userServiceForGeneralLedger,
		dimensionServiceForGeneralLedger,
	)

	reportApplication.Inject(
		reportRepository,
		reportReadRepository,
		generalLedgerServiceRepository,
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

	gin.SetMode(config.GetString("gin.releaseMode"))
	router := gin.Default()
	router.GET("/health/ping", func(c *gin.Context) { c.String(http.StatusOK, "Pong") })
	router.Use(datasource.ResolveSubdomain())
	router.Use(commonErrors.ErrorHandler(localizer))

	// public http API
	publicApiRouter := router.Group("/api/v1")
	sobPublicHttpPort.InitRouter(sobPublicHttpPort.NewHandler(&sobApplication), publicApiRouter)
	generalLedgerPublicHttpPort.InitRouter(
		generalLedgerPublicHttpPort.NewHandler(&generalLedgerApplication),
		publicApiRouter,
	)
	reportPublicHttpPort.InitRouter(reportPublicHttpPort.NewHandler(&reportApplication), publicApiRouter)
	userPublicHttpPort.InitRouter(userPublicHttpPort.NewHandler(&userApplication), publicApiRouter)
	dimensionPublicHttpPort.InitRouter(dimensionPublicHttpPort.NewHandler(&dimensionApplication), publicApiRouter)

	// private http API, should have different authentication method then public API
	privateApiRouter := router.Group("/internal")
	sobPrivateHttpPort.InitRouter(sobPrivateHttpPort.NewHandler(&sobApplication), privateApiRouter)
	numberingPrivateHttpPort.InitRouter(numberingPrivateHttpPort.NewHandler(&numberingApplication), privateApiRouter)
	generalLedgerPrivateHttpPort.InitRouter(
		generalLedgerPrivateHttpPort.NewHandler(&generalLedgerApplication),
		privateApiRouter,
	)
	reportPrivateHttpPort.InitRouter(reportPrivateHttpPort.NewHandler(&reportApplication), privateApiRouter)
	userPrivateHttpPort.InitRouter(userPrivateHttpPort.NewHandler(&userApplication), privateApiRouter)
	dimensionPrivateHttpPort.InitRouter(dimensionPrivateHttpPort.NewHandler(&dimensionApplication), privateApiRouter)

	if strings.HasPrefix(config.GetString("profile"), "dev") {
		// gin-swagger
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		// devops
		devopsApiRouter := router.Group("/devops/")
		devops.InitJwtHandler(devopsApiRouter)
	}

	log.InfoWithoutCxt("All module routers initiated")

	log.InfoWithoutCxt("Starting gin engine...")
	if err := router.Run(":" + config.GetString("app.port")); err != nil {
		panic(err.Error())
	}
}

func cleanup() {
	log.InfoWithoutCxt("fims terminating...")
	log.SyncLogger()
}
