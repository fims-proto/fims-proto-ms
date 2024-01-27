package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github/fims-proto/fims-proto-ms/internal/common/database"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/common/localization"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	tenantDb "github/fims-proto/fims-proto-ms/internal/tenant/adapter/db"
	tenantApp "github/fims-proto/fims-proto-ms/internal/tenant/app"
	ginMiddleware "github/fims-proto/fims-proto-ms/internal/tenant/lib/gin-middleware"
	tenantManager "github/fims-proto/fims-proto-ms/internal/tenant/lib/tenant-manager"
	tenantService "github/fims-proto/fims-proto-ms/internal/tenant/lib/tenant-service"
	tenantIntraPort "github/fims-proto/fims-proto-ms/internal/tenant/port/private/intraprocess"
	"net/http"
)

func main() {
	flag.Parse()

	loadConfig()

	log.InitLogger()
	defer cleanup()

	dbConnector := database.NewConnector()

	dbConnection, err := dbConnector.Open(viper.GetString("postgres.dsn"))
	if err != nil {
		panic(fmt.Errorf("failed to initialize dbConnection connection: %w", err))
	}

	tenantPostgresRepository := tenantDb.NewTenantPostgresRepository(dbConnection)
	tenantApplication := tenantApp.NewApplication(tenantPostgresRepository)
	tenantInterface := tenantIntraPort.NewTenantInterface(&tenantApplication)
	tenantServiceImpl := tenantService.NewTenantService(tenantInterface)

	tenantManagerImpl := tenantManager.NewTenantManager(tenantServiceImpl, dbConnector)

	// i18n
	localizer := localization.NewLocalizer("./i18n", "zh-CN")

	router := gin.Default()
	router.GET("/health/ping", func(c *gin.Context) { c.String(http.StatusOK, "Pong") })
	router.Use(ginMiddleware.ResolveTenantBySubdomain(tenantManagerImpl))
	router.Use(commonErrors.ErrorHandler(localizer))

	// public http API
	_ = router.Group("/api/v1")
	_ = router.Group("/internal")

	log.InfoWithoutCxt("Starting gin engine...")
	if err := router.Run(":" + viper.GetString("app.port")); err != nil {
		panic(err.Error())
	}
}

func loadConfig() {
	// environment variables
	if err := viper.BindEnv("profile", "PROFILE"); err != nil {
		panic(fmt.Errorf("failed to bind ENV profile: %w", err))
	}
	viper.SetDefault("profile", "dev")

	_ = viper.BindEnv("postgres.dsn", "DSN")

	// read config
	profile := viper.GetString("profile")
	viper.SetConfigName(fmt.Sprintf("application-%s", profile))
	viper.AddConfigPath("./config/")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to load config file: %w", err))
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
