package app

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
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
	accountService command.AccountService,
	accountReadService query.AccountService,
	ledgerService command.LedgerService,
	counterService command.CounterService,
) {
	a.Queries = Queries{
		ReadVouchers: query.NewReadVouchersHandler(readModel, accountReadService),
	}
	a.Commands = Commands{
		CreateVoucher: command.NewCreateVoucherHandler(repo, accountService, counterService),
		AuditVoucher:  command.NewAuditVoucherHandler(repo),
		ReviewVoucher: command.NewReviewVoucherHandler(repo),
		UpdateVoucher: command.NewUpdateVoucherHandler(repo, accountService),
		PostVoucher:   command.NewPostVoucherHandler(readModel, repo, ledgerService),
		Migrate:       command.NewMigrationHandler(repo),
	}
}
