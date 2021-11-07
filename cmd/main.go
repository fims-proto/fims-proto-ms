package main

import (
	"bytes"
	"flag"
	"fmt"
	_ "github/fims-proto/fims-proto-ms/docs"
	accountadapter "github/fims-proto/fims-proto-ms/internal/account/adapter/db"
	accountledgeradapter "github/fims-proto/fims-proto-ms/internal/account/adapter/ledger"
	accountapp "github/fims-proto/fims-proto-ms/internal/account/app"
	accountprivatehttpport "github/fims-proto/fims-proto-ms/internal/account/port/private/http"
	accountintraport "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
	"github/fims-proto/fims-proto-ms/internal/common/db"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	counteradapter "github/fims-proto/fims-proto-ms/internal/counter/adapter/db"
	counterapp "github/fims-proto/fims-proto-ms/internal/counter/app"
	counterprivatehttpport "github/fims-proto/fims-proto-ms/internal/counter/port/private/http"
	counterintraport "github/fims-proto/fims-proto-ms/internal/counter/port/private/intraprocess"
	"github/fims-proto/fims-proto-ms/internal/devops"
	ledgeraccountadapter "github/fims-proto/fims-proto-ms/internal/ledger/adapter/account"
	ledgeradapter "github/fims-proto/fims-proto-ms/internal/ledger/adapter/db"
	ledgervoucheradapter "github/fims-proto/fims-proto-ms/internal/ledger/adapter/voucher"
	ledgerapp "github/fims-proto/fims-proto-ms/internal/ledger/app"
	ledgerprivatehttpport "github/fims-proto/fims-proto-ms/internal/ledger/port/private/http"
	ledgerintraport "github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"
	sobadapter "github/fims-proto/fims-proto-ms/internal/sob/adapter/db"
	sobapp "github/fims-proto/fims-proto-ms/internal/sob/app"
	sobprivatehttpport "github/fims-proto/fims-proto-ms/internal/sob/port/private/http"
	sobpublichttpport "github/fims-proto/fims-proto-ms/internal/sob/port/public/http"
	tenantdb "github/fims-proto/fims-proto-ms/internal/tenant/adapter/db"
	tenantapp "github/fims-proto/fims-proto-ms/internal/tenant/app"
	ginmiddleware "github/fims-proto/fims-proto-ms/internal/tenant/lib/gin-middleware"
	tenantmanager "github/fims-proto/fims-proto-ms/internal/tenant/lib/tenant-manager"
	tenantservice "github/fims-proto/fims-proto-ms/internal/tenant/lib/tenant-service"
	tenantintraport "github/fims-proto/fims-proto-ms/internal/tenant/port/private/intraprocess"
	voucheraccountadapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/account"
	vouchercounteradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/counter"
	voucheradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/db"
	voucherledgeradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/ledger"
	voucherapp "github/fims-proto/fims-proto-ms/internal/voucher/app"
	voucherprivatehttpport "github/fims-proto/fims-proto-ms/internal/voucher/port/private/http"
	voucherintraport "github/fims-proto/fims-proto-ms/internal/voucher/port/private/intraprocess"
	voucherpublichttpport "github/fims-proto/fims-proto-ms/internal/voucher/port/public/http"
	"strings"

	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func main() {
	flag.Parse()

	loadConfig()

	log.InitLogger()
	defer cleanup()

	dbConnector := db.NewDBConnector(
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.dbName"),
		viper.GetString("postgres.timeZone"),
	)

	username := viper.GetString("postgres.username")

	db, err := dbConnector.Open(username, viper.GetString("postgres.password"))
	if err != nil {
		panic(errors.Wrapf(err, "open db connection for schema %s failed", username))
	}

	tenantPostgresRepository := tenantdb.NewTenantPostgresRepository(db)
	tenantApplication := tenantapp.NewApplication(tenantPostgresRepository)
	tenantInterface := tenantintraport.NewTenantInterface(&tenantApplication)
	tenantService := tenantservice.NewTenantService(tenantInterface)

	tenantManager := tenantmanager.NewTenantManager(tenantService, dbConnector)

	// repositories
	sobRepository := sobadapter.NewSobPostgresRepository()
	accountRepository := accountadapter.NewAccountPostgresRepository()
	voucherRepository := voucheradapter.NewVoucherPostgresRepository()
	ledgerRepository := ledgeradapter.NewLedgerPostgresRepository()
	counterRepository := counteradapter.NewCounterPostgresRepository()

	// application - will be passed by reference, in order to make injectinon work
	sobApplication := sobapp.NewApplication()
	accountApplication := accountapp.NewApplication()
	voucherApplication := voucherapp.NewApplication()
	ledgerApplication := ledgerapp.NewApplication()
	counterApplication := counterapp.NewApplication()

	// intrprocess interfaces
	accountInterface := accountintraport.NewAccountInterface(&accountApplication)
	voucherInterface := voucherintraport.NewVoucherInterface(&voucherApplication)
	ledgerInterface := ledgerintraport.NewLedgerInterface(&ledgerApplication)
	counterInterface := counterintraport.NewCounterInterface(&counterApplication)

	// application dependencies injection
	sobApplication.Inject(sobRepository, sobRepository)

	accountApplication.Inject(
		accountRepository,
		accountRepository,
		accountledgeradapter.NewIntraprocessAdapter(ledgerInterface),
	)

	voucherApplication.Inject(
		voucherRepository,
		voucherRepository,
		voucheraccountadapter.NewIntraprocessAdapter(accountInterface),
		voucherledgeradapter.NewIntraprocessAdapter(ledgerInterface),
		vouchercounteradapter.NewIntraprocessAdapter(counterInterface),
	)

	ledgerApplication.Inject(
		ledgerRepository,
		ledgerRepository,
		ledgeraccountadapter.NewIntraprocessAdapter(accountInterface),
		ledgervoucheradapter.NewIntraprocessAdapter(voucherInterface),
	)

	counterApplication.Inject(
		counterRepository, // because of pinter receiver
		counterRepository,
	)

	log.InfoWithoutCxt("All module applications initiated")

	router := gin.Default()
	router.Use(ginmiddleware.ResolveTenantBySubdomain(tenantManager))

	// public http API
	publicApiRouter := router.Group("/api/v1")
	sobpublichttpport.InitRouter(sobpublichttpport.NewHandler(&sobApplication), publicApiRouter)
	voucherpublichttpport.InitRouter(voucherpublichttpport.NewHandler(&voucherApplication), publicApiRouter)

	// private http API, should have different authentication method then public API
	privateApiRouter := router.Group("/internal")
	sobprivatehttpport.InitRouter(sobprivatehttpport.NewHandler(&sobApplication), privateApiRouter)
	counterprivatehttpport.InitRouter(counterprivatehttpport.NewHandler(&counterApplication), privateApiRouter)
	accountprivatehttpport.InitRouter(accountprivatehttpport.NewHandler(&accountApplication), privateApiRouter)
	ledgerprivatehttpport.InitRouter(ledgerprivatehttpport.NewHandler(&ledgerApplication), privateApiRouter)
	voucherprivatehttpport.InitRouter(voucherprivatehttpport.NewHandler(&voucherApplication), privateApiRouter)

	if strings.HasPrefix(viper.GetString("profile"), "dev") {
		// gin-swagger
		router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
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
	if !viper.IsSet("postgres.host") {
		checkResult.WriteString("postgres.host; ")
	}
	if !viper.IsSet("postgres.port") {
		checkResult.WriteString("postgres.port; ")
	}
	if !viper.IsSet("postgres.dbName") {
		checkResult.WriteString("postgres.dbName; ")
	}
	viper.SetDefault("postgres.timeZone", "UTC")
	if !viper.IsSet("postgres.username") {
		checkResult.WriteString("postgres.username; ")
	}
	if !viper.IsSet("postgres.password") {
		checkResult.WriteString("postgres.password; ")
	}
	// logger
	viper.SetDefault("logger.debug", false)
	viper.SetDefault("logger.jsonEncoding", true)

	if checkResult.Len() > 0 {
		panic("config missing: " + checkResult.String())
	}
}

func cleanup() {
	log.InfoWithoutCxt("fims terminating...")
	log.SyncLogger()
}
