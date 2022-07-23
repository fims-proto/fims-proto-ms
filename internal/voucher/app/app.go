package app

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/service"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
)

type Queries struct {
	ReadVouchers query.ReadVouchersHandler
}

type Commands struct {
	CreateVoucher command.CreateVoucherHandler
	AuditVoucher  command.AuditVoucherHandler
	ReviewVoucher command.ReviewVoucherHandler
	UpdateVoucher command.UpdateVoucherHandler
	PostVoucher   command.PostVoucherHandler
	Migrate       command.MigrationHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(
	readModel query.VouchersReadModel,
	repo domain.Repository,
	accountService service.AccountService,
	ledgerService service.LedgerService,
	userService service.UserService,
	numberingService service.NumberingService,
) {
	a.Queries = Queries{
		ReadVouchers: query.NewReadVouchersHandler(readModel, accountService, userService, ledgerService),
	}
	a.Commands = Commands{
		CreateVoucher: command.NewCreateVoucherHandler(repo, accountService, numberingService, ledgerService),
		AuditVoucher:  command.NewAuditVoucherHandler(repo),
		ReviewVoucher: command.NewReviewVoucherHandler(repo),
		UpdateVoucher: command.NewUpdateVoucherHandler(repo, accountService, ledgerService),
		PostVoucher:   command.NewPostVoucherHandler(readModel, repo, ledgerService),
		Migrate:       command.NewMigrationHandler(repo),
	}
}
