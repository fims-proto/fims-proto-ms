package main

import (
	"github.com/gin-gonic/gin"
	voucherAdapter "github/fims-proto/fims-proto-ms/internal/voucher/adapter"
	voucherApp "github/fims-proto/fims-proto-ms/internal/voucher/app"
	voucherCommand "github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	voucherQuery "github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	voucherHttpPort "github/fims-proto/fims-proto-ms/internal/voucher/port/public/http"
)

func main() {
	voucherApplication := newVoucherApplication()

	router := gin.Default()
	voucherHttpPort.InitRouter(voucherHttpPort.NewHandler(voucherApplication), router)

	if err := router.Run(":8080"); err != nil {
		panic(err.Error())
	}
}

func newVoucherApplication() voucherApp.Application {
	memoryRepository := voucherAdapter.NewVoucherMemoryRepository()

	return voucherApp.Application{
		Queries:  voucherApp.Queries{
			ReadVouchers: voucherQuery.NewAllVouchersHandler(memoryRepository),
		},
		Commands: voucherApp.Commands{
			RecordVoucher: voucherCommand.NewRecordVoucherHandler(&memoryRepository),
			AuditVoucher:  voucherCommand.NewAuditVoucherHandler(&memoryRepository),
			ReviewVoucher: voucherCommand.NewReviewVoucherHandler(&memoryRepository),
			UpdateVoucher: voucherCommand.NewUpdateVoucherHandler(&memoryRepository),
		},
	}
}
