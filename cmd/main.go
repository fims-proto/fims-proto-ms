package main

import (
	"github.com/gin-gonic/gin"
	accountAdapter "github/fims-proto/fims-proto-ms/internal/account/adapter"
	accountApp "github/fims-proto/fims-proto-ms/internal/account/app"
	accountQuery "github/fims-proto/fims-proto-ms/internal/account/app/query"
	accountIntraPort "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
	voucherAdapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter"
	voucherAccountAdapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/account"
	voucherApp "github/fims-proto/fims-proto-ms/internal/voucher/app"
	voucherCommand "github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	voucherQuery "github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	voucherHttpPort "github/fims-proto/fims-proto-ms/internal/voucher/port/public/http"
)

func main() {
	_, accountInterface := newAccountApplication()
	voucherApplication := newVoucherApplication(accountInterface)

	router := gin.Default()
	voucherHttpPort.InitRouter(voucherHttpPort.NewHandler(voucherApplication), router)
	if err := router.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

func newVoucherApplication(accountInterface accountIntraPort.AccountInterface) voucherApp.Application {
	memoryRepository := voucherAdapter.NewVoucherMemoryRepository()
	accountService := voucherAccountAdapter.NewIntraprocessService(accountInterface)

	return voucherApp.Application{
		Queries: voucherApp.Queries{
			ReadVouchers: voucherQuery.NewAllVouchersHandler(memoryRepository),
		},
		Commands: voucherApp.Commands{
			RecordVoucher: voucherCommand.NewRecordVoucherHandler(&memoryRepository, accountService),
			AuditVoucher:  voucherCommand.NewAuditVoucherHandler(&memoryRepository),
			ReviewVoucher: voucherCommand.NewReviewVoucherHandler(&memoryRepository),
			UpdateVoucher: voucherCommand.NewUpdateVoucherHandler(&memoryRepository, accountService),
		},
	}
}

func newAccountApplication() (accountApp.Application, accountIntraPort.AccountInterface) {
	memoryRepository := accountAdapter.NewAccountMemoryRepository()

	application := accountApp.Application{
		Queries: accountApp.Queries{
			ValidateAccounts: accountQuery.NewValidateAccountsHandler(memoryRepository),
		},
		Commands: accountApp.Commands{},
	}

	accountInterface := accountIntraPort.NewHandler(application)

	return application, accountInterface
}
