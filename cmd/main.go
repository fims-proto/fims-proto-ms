package main

import (
	accountadapter "github/fims-proto/fims-proto-ms/internal/account/adapter"
	accountapp "github/fims-proto/fims-proto-ms/internal/account/app"
	accountquery "github/fims-proto/fims-proto-ms/internal/account/app/query"
	accountintraport "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
	voucheradapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter"
	voucheraccountadapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter/account"
	voucherapp "github/fims-proto/fims-proto-ms/internal/voucher/app"
	vouchercommand "github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	voucherquery "github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	voucherhttpport "github/fims-proto/fims-proto-ms/internal/voucher/port/public/http"

	"github.com/gin-gonic/gin"
)

func main() {
	_, accountInterface := newAccountApplication()
	voucherApplication := newVoucherApplication(accountInterface)

	router := gin.Default()
<<<<<<< HEAD
	voucherhttpport.InitRouter(voucherhttpport.NewHandler(voucherApplication), router)

=======
	voucherHttpPort.InitRouter(voucherHttpPort.NewHandler(voucherApplication), router)
>>>>>>> master
	if err := router.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

func newVoucherApplication(accountInterface accountintraport.AccountInterface) voucherapp.Application {
	memoryRepository := voucheradapter.NewVoucherMemoryRepository()
	accountService := voucheraccountadapter.NewIntraprocessService(accountInterface)

	return voucherapp.Application{
		Queries: voucherapp.Queries{
			ReadVouchers: voucherquery.NewAllVouchersHandler(memoryRepository),
		},
		Commands: voucherapp.Commands{
			RecordVoucher: vouchercommand.NewRecordVoucherHandler(&memoryRepository, accountService),
			AuditVoucher:  vouchercommand.NewAuditVoucherHandler(&memoryRepository),
			ReviewVoucher: vouchercommand.NewReviewVoucherHandler(&memoryRepository),
			UpdateVoucher: vouchercommand.NewUpdateVoucherHandler(&memoryRepository, accountService),
		},
	}
}

func newAccountApplication() (accountapp.Application, accountintraport.AccountInterface) {
	memoryRepository := accountadapter.NewAccountMemoryRepository()

	application := accountapp.Application{
		Queries: accountapp.Queries{
			ValidateAccounts: accountquery.NewValidateAccountsHandler(memoryRepository),
		},
		Commands: accountapp.Commands{},
	}

	accountInterface := accountintraport.NewHandler(application)

	return application, accountInterface
}
