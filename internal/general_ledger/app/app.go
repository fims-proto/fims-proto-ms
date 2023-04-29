package app

import (
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
)

type Queries struct {
	PagingAccounts        query.PagingAccountsHandler
	AccountsByIds         query.AccountsByIdsHandler
	AccountsByNumbers     query.AccountsByNumbersHandler
	CurrentPeriod         query.CurrentPeriodHandler
	PagingPeriods         query.PagingPeriodsHandler
	PeriodsByIds          query.PeriodsByIdsHandler
	PeriodById            query.PeriodByIdHandler
	PagingLedgersByPeriod query.PagingLedgersByPeriodHandler
	VoucherById           query.VoucherByIdHandler
	PagingVouchers        query.PagingVouchersHandler
}

type Commands struct {
	Initialize command.InitializeHandler

	CreatePeriod command.CreatePeriodHandler
	ClosePeriod  command.ClosePeriodHandler

	CreateVoucher       command.CreateVoucherHandler
	AuditVoucher        command.AuditVoucherHandler
	CancelAuditVoucher  command.CancelAuditVoucherHandler
	ReviewVoucher       command.ReviewVoucherHandler
	CancelReviewVoucher command.CancelReviewVoucherHandler
	UpdateVoucher       command.UpdateVoucherHandler
	PostVoucher         command.PostVoucherHandler

	Migrate command.MigrationHandler
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
	readModel query.GeneralLedgerReadModel,
	sobService service.SobService,
	numberingService service.NumberingService,
	userService service.UserService,
) {
	a.Queries = Queries{
		PagingAccounts:    query.NewPagingAccountsHandler(readModel),
		AccountsByNumbers: query.NewAccountsByNumbersHandler(readModel),
		AccountsByIds:     query.NewAccountsByIdsHandler(readModel),

		CurrentPeriod: query.NewCurrentPeriodHandler(readModel),
		PagingPeriods: query.NewPagingPeriodsHandler(readModel),
		PeriodsByIds:  query.NewPeriodsByIdsHandler(readModel),
		PeriodById:    query.NewPeriodByIdHandler(readModel),

		PagingLedgersByPeriod: query.NewPagingLedgersByPeriodHandler(readModel),

		VoucherById:    query.NewVoucherByIdHandler(readModel, userService),
		PagingVouchers: query.NewPagingVouchersHandler(readModel, userService),
	}
	a.Commands = Commands{
		Initialize: command.NewInitializeHandler(repo, readModel, sobService, numberingService),

		CreatePeriod: command.NewCreatePeriodHandler(repo, readModel, numberingService),
		ClosePeriod:  command.NewClosePeriodHandler(repo, readModel, numberingService),

		CreateVoucher:       command.NewCreateVoucherHandler(repo, readModel, numberingService),
		AuditVoucher:        command.NewAuditVoucherHandler(repo),
		CancelAuditVoucher:  command.NewCancelAuditVoucherHandler(repo),
		ReviewVoucher:       command.NewReviewVoucherHandler(repo),
		CancelReviewVoucher: command.NewCancelReviewVoucherHandler(repo),
		UpdateVoucher:       command.NewUpdateVoucherHandler(repo, readModel, numberingService),
		PostVoucher:         command.NewPostVoucherHandler(repo, readModel),

		Migrate: command.NewMigrationHandler(repo),
	}
}
