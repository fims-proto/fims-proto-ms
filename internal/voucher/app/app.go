package app

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/service"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
)

type Queries struct {
	VoucherById    query.VoucherByIdHandler
	PagingVouchers query.PagingVouchersHandler
}

type Commands struct {
	CreateVoucher       command.CreateVoucherHandler
	AuditVoucher        command.AuditVoucherHandler
	CancelAuditVoucher  command.CancelAuditVoucherHandler
	ReviewVoucher       command.ReviewVoucherHandler
	CancelReviewVoucher command.CancelReviewVoucherHandler
	UpdateVoucher       command.UpdateVoucherHandler
	PostVoucher         command.PostVoucherHandler
	Migrate             command.MigrationHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(
	repo domain.Repository,
	readModel query.VoucherReadModel,
	accountService service.AccountService,
	userService service.UserService,
	numberingService service.NumberingService,
) {
	a.Queries = Queries{
		VoucherById:    query.NewVoucherByIdHandler(readModel, accountService, userService),
		PagingVouchers: query.NewPagingVouchersHandler(readModel, userService),
	}
	a.Commands = Commands{
		CreateVoucher:       command.NewCreateVoucherHandler(repo, accountService, numberingService),
		AuditVoucher:        command.NewAuditVoucherHandler(repo),
		CancelAuditVoucher:  command.NewCancelAuditVoucherHandler(repo),
		ReviewVoucher:       command.NewReviewVoucherHandler(repo),
		CancelReviewVoucher: command.NewCancelReviewVoucherHandler(repo),
		UpdateVoucher:       command.NewUpdateVoucherHandler(repo, accountService),
		PostVoucher:         command.NewPostVoucherHandler(repo, accountService),
		Migrate:             command.NewMigrationHandler(repo),
	}
}
