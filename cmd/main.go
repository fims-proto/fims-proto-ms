package main

import (
	"flag"
	accountadapter "github/fims-proto/fims-proto-ms/internal/account/adapter"
	accountledgeradapter "github/fims-proto/fims-proto-ms/internal/account/adapter/ledger"
	accountapp "github/fims-proto/fims-proto-ms/internal/account/app"
	accountprivatehttpport "github/fims-proto/fims-proto-ms/internal/account/port/private/http"
	accountintraport "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
	"github/fims-proto/fims-proto-ms/internal/common/db"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	counteradapter "github/fims-proto/fims-proto-ms/internal/counter/adapter"
	counterapp "github/fims-proto/fims-proto-ms/internal/counter/app"
	counterprivatehttpport "github/fims-proto/fims-proto-ms/internal/counter/port/private/http"
	counterintraport "github/fims-proto/fims-proto-ms/internal/counter/port/private/intraprocess"
	ledgeradapter "github/fims-proto/fims-proto-ms/internal/ledger/adapter"
	ledgeraccountadapter "github/fims-proto/fims-proto-ms/internal/ledger/adapter/account"
	ledgervoucheradapter "github/fims-proto/fims-proto-ms/internal/ledger/adapter/voucher"
	ledgerapp "github/fims-proto/fims-proto-ms/internal/ledger/app"
	ledgertesthttpport "github/fims-proto/fims-proto-ms/internal/ledger/port/private/http"
	ledgerintraport "github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"
	sobadapter "github/fims-proto/fims-proto-ms/internal/sob/adapter"
	sobapp "github/fims-proto/fims-proto-ms/internal/sob/app"
	sobpublichttpport "github/fims-proto/fims-proto-ms/internal/sob/port/public/http"
	tenantdb "github/fims-proto/fims-proto-ms/internal/tenant/adapter/db"
	tenantapp "github/fims-proto/fims-proto-ms/internal/tenant/app"
	ginmiddleware "github/fims-proto/fims-proto-ms/internal/tenant/lib/gin-middleware"
	tenantmanager "github/fims-proto/fims-proto-ms/internal/tenant/lib/tenant-manager"
	tenantservice "github/fims-proto/fims-proto-ms/internal/tenant/lib/tenant-service"
	tenantintraport "github/fims-proto/fims-proto-ms/internal/tenant/port/private/intraprocess"
	"github/fims-proto/fims-proto-ms/internal/user/lib/authentication"
	voucheradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter"
	voucheraccountadapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/account"
	vouchercounteradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/counter"
	voucherledgeradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/ledger"
	voucherapp "github/fims-proto/fims-proto-ms/internal/voucher/app"
	voucherintraport "github/fims-proto/fims-proto-ms/internal/voucher/port/private/intraprocess"
	voucherpublichttpport "github/fims-proto/fims-proto-ms/internal/voucher/port/public/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func main() {
	flag.Parse()
	log.InitLoggers(log.NewStdLogEnablerAdapter(), log.NewStdLoggerAdapter())

	dbConnector := db.NewDBConnector()

	// >> TODO read from config file
	db, err := dbConnector.Open("fims-tenant-manager", "fims-tenant-manager")
	if err != nil {
		panic(errors.Wrap(err, "open fims-tenant-manager db connection failed"))
	}
	// << TODO
	tenantPostgresRepository := tenantdb.NewTenantPostgresRepository(db)
	tenantApplication := tenantapp.NewApplication(tenantPostgresRepository)
	tenantInterface := tenantintraport.NewTenantInterface(&tenantApplication)
	tenantService := tenantservice.NewTenantService(tenantInterface)

	tenantManager := tenantmanager.NewTenantManager(tenantService, dbConnector)

	// repositories
	sobRepository := sobadapter.NewSobMemoryRepository()
	accountRepository := accountadapter.NewAccountMemoryRepository()
	voucherRepository := voucheradapter.NewVoucherMemoryRepository()
	ledgerRepository := ledgeradapter.NewLedgerMemoryRepository()
	counterRepository := counteradapter.NewCounterMemoryRepository()

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
		ledgeraccountadapter.NewIntraprocessAdapter(accountInterface),
		ledgervoucheradapter.NewIntraprocessAdapter(voucherInterface),
	)

	counterApplication.Inject(
		&counterRepository, // because of pinter receiver
		&counterRepository,
	)

	router := gin.Default()
	router.Use(ginmiddleware.ResolveTenantBySubdomain(tenantManager))
	router.Use(authentication.Authn())
	sobpublichttpport.InitRouter(sobpublichttpport.NewHandler(&sobApplication), router)
	voucherpublichttpport.InitRouter(voucherpublichttpport.NewHandler(&voucherApplication), router)
	// below 2 are for dataload, can be integrated into onboarding procedure
	accountprivatehttpport.InitRouter(accountprivatehttpport.NewHandler(&accountApplication), router)
	counterprivatehttpport.InitRouter(counterprivatehttpport.NewHandler(&counterApplication), router)
	// TODO remove, test prupose
	ledgertesthttpport.InitRouter(ledgertesthttpport.NewHandler(ledgerRepository), router)

	if err := router.Run(":8080"); err != nil {
		panic(err.Error())
	}
}
