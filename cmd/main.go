package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	accountNumberingAdapter "github/fims-proto/fims-proto-ms/internal/account/adapter/numbering"

	"golang.org/x/text/language"

	"github/fims-proto/fims-proto-ms/internal/common/db"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/devops"

	_ "github/fims-proto/fims-proto-ms/docs"
	accountAdapter "github/fims-proto/fims-proto-ms/internal/account/adapter/db"
	accountSobAdapter "github/fims-proto/fims-proto-ms/internal/account/adapter/sob"
	accountApp "github/fims-proto/fims-proto-ms/internal/account/app"
	accountPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/account/port/private/http"
	accountIntraPort "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
	accountPublicHttpPort "github/fims-proto/fims-proto-ms/internal/account/port/public/http"

	numberingAdapter "github/fims-proto/fims-proto-ms/internal/numbering/adapter/db"
	numberingApp "github/fims-proto/fims-proto-ms/internal/numbering/app"
	numberingPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/numbering/port/private/http"
	numberingIntraPort "github/fims-proto/fims-proto-ms/internal/numbering/port/private/intraprocess"

	userAdapter "github/fims-proto/fims-proto-ms/internal/user/adapter/db"
	userApp "github/fims-proto/fims-proto-ms/internal/user/app"
	userPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/user/port/private/http"
	userIntraPort "github/fims-proto/fims-proto-ms/internal/user/port/private/intraprocess"
	userPublicHttpPort "github/fims-proto/fims-proto-ms/internal/user/port/public/http"

	sobAccountAdapter "github/fims-proto/fims-proto-ms/internal/sob/adapter/account"
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

	journalAccountAdapter "github/fims-proto/fims-proto-ms/internal/journal/adapter/account"
	journalAdapter "github/fims-proto/fims-proto-ms/internal/journal/adapter/db"
	journalNumberingAdapter "github/fims-proto/fims-proto-ms/internal/journal/adapter/numbering"
	journalUserAdapter "github/fims-proto/fims-proto-ms/internal/journal/adapter/user"
	journalApp "github/fims-proto/fims-proto-ms/internal/journal/app"
	journalPrivateHttpPort "github/fims-proto/fims-proto-ms/internal/journal/port/private/http"
	journalPublicHttpPort "github/fims-proto/fims-proto-ms/internal/journal/port/public/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

	// i18n
	bundle := i18n.NewBundle(language.SimplifiedChinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("./i18n/zh-CN.json")

	// repositories
	sobRepository := sobAdapter.NewSobPostgresRepository()
	accountRepository := accountAdapter.NewAccountPostgresRepository()
	journalRepository := journalAdapter.NewJournalEntryPostgresRepository()
	numberingRepository := numberingAdapter.NewNumberingPostgresRepository()
	userRepository := userAdapter.NewUserPostgresRepository()

	// application - will be passed by reference, in order to make injection work
	sobApplication := sobApp.NewApplication()
	accountApplication := accountApp.NewApplication()
	journalApplication := journalApp.NewApplication()
	numberingApplication := numberingApp.NewApplication()
	userApplication := userApp.NewApplication()

	// intra process interfaces
	sobInterface := sobIntraPort.NewSobInterface(&sobApplication)
	accountInterface := accountIntraPort.NewAccountInterface(&accountApplication)
	numberingInterface := numberingIntraPort.NewNumberingInterface(&numberingApplication)
	userInterface := userIntraPort.NewUserInterface(&userApplication)

	// application dependencies injection
	accountServiceForSob := sobAccountAdapter.NewIntraProcessAdapter(accountInterface)
	sobApplication.Inject(
		sobRepository,
		sobRepository,
		accountServiceForSob,
	)

	sobServiceForAccount := accountSobAdapter.NewIntraProcessAdapter(sobInterface)
	numberingServiceForAccount := accountNumberingAdapter.NewIntraProcessAdapter(numberingInterface)
	accountApplication.Inject(
		accountRepository,
		accountRepository,
		sobServiceForAccount,
		numberingServiceForAccount,
	)

	accountServiceForJournal := journalAccountAdapter.NewIntraProcessAdapter(accountInterface)
	numberingServiceForJournal := journalNumberingAdapter.NewIntraProcessAdapter(numberingInterface)
	userServiceForJournal := journalUserAdapter.NewIntraProcessAdapter(userInterface)
	journalApplication.Inject(
		journalRepository,
		journalRepository,
		accountServiceForJournal,
		userServiceForJournal,
		numberingServiceForJournal,
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
	router.Use(ginMiddleware.ResolveTenantBySubdomain(tenantManagerImpl))
	router.Use(commonErrors.ErrorHandler(bundle))

	// public http API
	publicApiRouter := router.Group("/api/v1")
	sobPublicHttpPort.InitRouter(sobPublicHttpPort.NewHandler(&sobApplication), publicApiRouter)
	accountPublicHttpPort.InitRouter(accountPublicHttpPort.NewHandler(&accountApplication), publicApiRouter)
	journalPublicHttpPort.InitRouter(journalPublicHttpPort.NewHandler(&journalApplication), publicApiRouter)
	userPublicHttpPort.InitRouter(userPublicHttpPort.NewHandler(&userApplication), publicApiRouter)

	// private http API, should have different authentication method then public API
	privateApiRouter := router.Group("/internal")
	sobPrivateHttpPort.InitRouter(sobPrivateHttpPort.NewHandler(&sobApplication), privateApiRouter)
	numberingPrivateHttpPort.InitRouter(numberingPrivateHttpPort.NewHandler(&numberingApplication), privateApiRouter)
	accountPrivateHttpPort.InitRouter(accountPrivateHttpPort.NewHandler(&accountApplication), privateApiRouter)
	journalPrivateHttpPort.InitRouter(journalPrivateHttpPort.NewHandler(&journalApplication), privateApiRouter)
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
