package main

import (
	"bytes"
	"flag"
	"fmt"
	_ "github/fims-proto/fims-proto-ms/docs"
	accountAdapter "github/fims-proto/fims-proto-ms/internal/account/adapter/db"
	accountSobAdapter "github/fims-proto/fims-proto-ms/internal/account/adapter/sob"
	accountApp "github/fims-proto/fims-proto-ms/internal/account/app"
	accountPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/account/port/private/http"
	accountIntraPort "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
	"github/fims-proto/fims-proto-ms/internal/common/db"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	counterAdapter "github/fims-proto/fims-proto-ms/internal/counter/adapter/db"
	counterApp "github/fims-proto/fims-proto-ms/internal/counter/app"
	counterPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/counter/port/private/http"
	counterIntraPort "github/fims-proto/fims-proto-ms/internal/counter/port/private/intraprocess"
	"github/fims-proto/fims-proto-ms/internal/devops"
	ledgerAccountAdapter "github/fims-proto/fims-proto-ms/internal/ledger/adapter/account"
	ledgerAdapter "github/fims-proto/fims-proto-ms/internal/ledger/adapter/db"
	ledgerApp "github/fims-proto/fims-proto-ms/internal/ledger/app"
	ledgerPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/ledger/port/private/http"
	ledgerIntraPort "github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"
	ledgerPublicHttpPort "github/fims-proto/fims-proto-ms/internal/ledger/port/public/http"
	sobAdapter "github/fims-proto/fims-proto-ms/internal/sob/adapter/db"
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
	voucherAccountAdapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/account"
	voucherCounterAdapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/counter"
	voucherAdapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/db"
	voucherLedgerAdapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/ledger"
	voucherApp "github/fims-proto/fims-proto-ms/internal/voucher/app"
	voucherPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/voucher/port/private/http"
	voucherPublicHttpPort "github/fims-proto/fims-proto-ms/internal/voucher/port/public/http"
	"strings"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func main() {
	flag.Parse()

	loadConfig()

	log.InitLogger()
	defer cleanup()

	dbConnector := db.NewConnector()

	dbConnection, err := dbConnector.Open(viper.GetString("postgres.dsn"))
	if err != nil {
		panic(errors.Wrapf(err, "failed to initialize dbConnection connection"))
	}

	tenantPostgresRepository := tenantDb.NewTenantPostgresRepository(dbConnection)
	tenantApplication := tenantApp.NewApplication(tenantPostgresRepository)
	tenantInterface := tenantIntraPort.NewTenantInterface(&tenantApplication)
	tenantServiceImpl := tenantService.NewTenantService(tenantInterface)

	tenantManagerImpl := tenantManager.NewTenantManager(tenantServiceImpl, dbConnector)

	// repositories
	sobRepository := sobAdapter.NewSobPostgresRepository()
	accountRepository := accountAdapter.NewAccountPostgresRepository()
	voucherRepository := voucherAdapter.NewVoucherPostgresRepository()
	ledgerRepository := ledgerAdapter.NewLedgerPostgresRepository()
	counterRepository := counterAdapter.NewCounterPostgresRepository()

	// application - will be passed by reference, in order to make injection work
	sobApplication := sobApp.NewApplication()
	accountApplication := accountApp.NewApplication()
	voucherApplication := voucherApp.NewApplication()
	ledgerApplication := ledgerApp.NewApplication()
	counterApplication := counterApp.NewApplication()

	// intra process interfaces
	sobInterface := sobIntraPort.NewSobInterface(&sobApplication)
	accountInterface := accountIntraPort.NewAccountInterface(&accountApplication)
	ledgerInterface := ledgerIntraPort.NewLedgerInterface(&ledgerApplication)
	counterInterface := counterIntraPort.NewCounterInterface(&counterApplication)

	// application dependencies injection
	sobApplication.Inject(sobRepository, sobRepository)

	accountApplication.Inject(
		accountRepository,
		accountRepository,
		accountSobAdapter.NewIntraProcessAdapter(sobInterface),
	)

	voucherApplication.Inject(
		voucherRepository,
		voucherRepository,
		voucherAccountAdapter.NewIntraProcessAdapter(accountInterface),
		voucherLedgerAdapter.NewIntraProcessAdapter(ledgerInterface),
		voucherCounterAdapter.NewIntraProcessAdapter(counterInterface),
	)

	ledgerApplication.Inject(
		ledgerRepository,
		ledgerRepository,
		ledgerAccountAdapter.NewIntraProcessAdapter(accountInterface),
	)

	counterApplication.Inject(
		counterRepository, // because of pinter receiver
		counterRepository,
	)

	log.InfoWithoutCxt("All module applications initiated")

	router := gin.Default()
	router.Use(ginMiddleware.ResolveTenantBySubdomain(tenantManagerImpl))

	// public http API
	publicApiRouter := router.Group("/api/v1")
	sobPublicHttpPort.InitRouter(sobPublicHttpPort.NewHandler(&sobApplication), publicApiRouter)
	voucherPublicHttpPort.InitRouter(voucherPublicHttpPort.NewHandler(&voucherApplication), publicApiRouter)
	ledgerPublicHttpPort.InitRouter(ledgerPublicHttpPort.NewHandler(&ledgerApplication), publicApiRouter)

	// private http API, should have different authentication method then public API
	privateApiRouter := router.Group("/internal")
	sobPrivateHttpPort.InitRouter(sobPrivateHttpPort.NewHandler(&sobApplication), privateApiRouter)
	counterPrivateHttpPort.InitRouter(counterPrivateHttpPort.NewHandler(&counterApplication), privateApiRouter)
	accountPrivateHttpPort.InitRouter(accountPrivateHttpPort.NewHandler(&accountApplication), privateApiRouter)
	ledgerPrivateHttpPort.InitRouter(ledgerPrivateHttpPort.NewHandler(&ledgerApplication), privateApiRouter)
	voucherPrivateHttpPort.InitRouter(voucherPrivateHttpPort.NewHandler(&voucherApplication), privateApiRouter)

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
