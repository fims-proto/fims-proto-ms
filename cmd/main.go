package main

import (
	accountadapter "github/fims-proto/fims-proto-ms/internal/account/adapter"
	accountledgeradapter "github/fims-proto/fims-proto-ms/internal/account/adapter/ledger"
	accountapp "github/fims-proto/fims-proto-ms/internal/account/app"
	accountprivatehttpport "github/fims-proto/fims-proto-ms/internal/account/port/private/http"
	accountintraport "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
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
	voucheradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter"
	voucheraccountadapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/account"
	vouchercounteradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/counter"
	voucherledgeradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/ledger"
	voucherapp "github/fims-proto/fims-proto-ms/internal/voucher/app"
	voucherintraport "github/fims-proto/fims-proto-ms/internal/voucher/port/private/intraprocess"
	voucherpublichttpport "github/fims-proto/fims-proto-ms/internal/voucher/port/public/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// repositories
	accountRepository := accountadapter.NewAccountMemoryRepository()
	voucherRepository := voucheradapter.NewVoucherMemoryRepository()
	ledgerRepository := ledgeradapter.NewLedgerMemoryRepository()
	counterRepository := counteradapter.NewCounterMemoryRepository()

	// application - will be passed by reference, in order to make injectinon work
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
