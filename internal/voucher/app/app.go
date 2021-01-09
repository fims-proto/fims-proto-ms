package app

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
)

type Queries struct {
	AllVouchers query.AllVouchersHandler
}

type Commands struct {
	RecordVoucher command.RecordVoucherHandler
	AuditVoucher command.AuditVoucherHandler
	ReviewVoucher command.ReviewVoucherHandler
}

type Application struct {
	Queries Queries
	Commands Commands
}
