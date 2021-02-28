package app

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
)

type Queries struct {
	ReadVouchers query.ReadVouchersHandler
}

type Commands struct {
	RecordVoucher command.RecordVoucherHandler
	AuditVoucher  command.AuditVoucherHandler
	ReviewVoucher command.ReviewVoucherHandler
	UpdateVoucher command.UpdateVoucherHandler
	PostVoucher   command.PostVoucherHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}
