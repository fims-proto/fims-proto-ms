package main

import (
	accountadapter "github/fims-proto/fims-proto-ms/internal/account/adapter"
	accountapp "github/fims-proto/fims-proto-ms/internal/account/app"
	accountquery "github/fims-proto/fims-proto-ms/internal/account/app/query"
	accountintraport "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
	ledgeradapter "github/fims-proto/fims-proto-ms/internal/ledger/adapter"
	ledgeraccountadapter "github/fims-proto/fims-proto-ms/internal/ledger/adapter/account"
	ledgervoucheradapter "github/fims-proto/fims-proto-ms/internal/ledger/adapter/voucher"
	ledgerapp "github/fims-proto/fims-proto-ms/internal/ledger/app"
	ledgercommand "github/fims-proto/fims-proto-ms/internal/ledger/app/command"
	ledgerintraport "github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"
	voucheradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter"
	voucheraccountadapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/account"
	voucherledgeradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/ledger"
	voucherapp "github/fims-proto/fims-proto-ms/internal/voucher/app"
	vouchercommand "github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	voucherquery "github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	voucherintraport "github/fims-proto/fims-proto-ms/internal/voucher/port/private/intraprocess"
	voucherhttpport "github/fims-proto/fims-proto-ms/internal/voucher/port/public/http"

	"github.com/gin-gonic/gin"
)

func main() {
	accountApplication := newAccountApplication()
	voucherApplication := newVoucherApplication(accountInterface, ledgerInterface)
	ledgerApplication := newLedgerApplication(accountInterface, voucherInterface)

	accountInterface := accountintraport.NewAccountInterface(accountApplication)
	voucherInterface := voucherintraport.NewVoucherInterface(voucherApplication)
	ledgerInterface := ledgerintraport.NewLedgerInterface(ledgerApplication)

	router := gin.Default()
	voucherhttpport.InitRouter(voucherhttpport.NewHandler(voucherApplication), router)

	if err := router.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

func newAccountApplication() accountapp.Application {
	memoryRepository := accountadapter.NewAccountMemoryRepository()

	return accountapp.Application{
		Queries: accountapp.Queries{
			ReadAccounts:     accountquery.NewReadAccountsHandler(memoryRepository),
			ValidateAccounts: accountquery.NewValidateAccountsHandler(memoryRepository),
		},
		Commands: accountapp.Commands{},
	}
}

func newVoucherApplication(accountInterface accountintraport.AccountInterface, ledgerInterface ledgerintraport.LedgerInterface) voucherapp.Application {
	memoryRepository := voucheradapter.NewVoucherMemoryRepository()
	accountService := voucheraccountadapter.NewIntraprocessAdapter(accountInterface)
	ledgerService := voucherledgeradapter.NewIntraprocessAdapter(ledgerInterface)

	return voucherapp.Application{
		Queries: voucherapp.Queries{
			ReadVouchers: voucherquery.NewReadVouchersHandler(memoryRepository),
		},
		Commands: voucherapp.Commands{
			RecordVoucher: vouchercommand.NewRecordVoucherHandler(memoryRepository, accountService),
			AuditVoucher:  vouchercommand.NewAuditVoucherHandler(memoryRepository),
			ReviewVoucher: vouchercommand.NewReviewVoucherHandler(memoryRepository),
			UpdateVoucher: vouchercommand.NewUpdateVoucherHandler(memoryRepository, accountService),
			PostVoucher:   vouchercommand.NewPostVoucherHandler(memoryRepository, memoryRepository, ledgerService),
		},
	}
}

func newLedgerApplication(accountInterface accountintraport.AccountInterface, voucherInterface voucherintraport.VoucherInterface) ledgerapp.Application {
	memoryRepository := ledgeradapter.NewLedgerMemoryRepository()
	accounterService := ledgeraccountadapter.NewIntraprocessAdapter(accountInterface)
	voucherService := ledgervoucheradapter.NewIntraprocessAdapter(voucherInterface)

	return ledgerapp.Application{
		Queries: ledgerapp.Queries{},
		Commands: ledgerapp.Commands{
			UpdateLedgerBalance: ledgercommand.NewUpdateLedgerBalanceHandler(memoryRepository, accounterService, voucherService),
		},
	}
}
