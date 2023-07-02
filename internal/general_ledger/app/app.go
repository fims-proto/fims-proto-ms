package app

import (
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
)

type Queries struct {
	PagingAccounts            query.PagingAccountsHandler
	PagingAuxiliaryCategories query.PagingAuxiliaryCategoriesHandler
	PagingAuxiliaryAccounts   query.PagingAuxiliaryAccountsHandler
	CurrentPeriod             query.CurrentPeriodHandler
	PagingPeriods             query.PagingPeriodsHandler
	PagingLedgersByPeriod     query.PagingLedgersByPeriodHandler
	PagingAuxiliaryLedgers    query.PagingAuxiliaryLedgersHandler
	VoucherById               query.VoucherByIdHandler
	PagingVouchers            query.PagingVouchersHandler
}

type Commands struct {
	Initialize command.InitializeHandler

	AssignAuxiliaryCategory command.AssignAuxiliaryCategoryHandler

	ClosePeriod command.ClosePeriodHandler

	CreateVoucher       command.CreateVoucherHandler
	AuditVoucher        command.AuditVoucherHandler
	CancelAuditVoucher  command.CancelAuditVoucherHandler
	ReviewVoucher       command.ReviewVoucherHandler
	CancelReviewVoucher command.CancelReviewVoucherHandler
	UpdateVoucher       command.UpdateVoucherHandler
	PostVoucher         command.PostVoucherHandler

	CreateAuxiliaryCategory command.CreateAuxiliaryCategoryHandler
	CreateAuxiliaryAccount  command.CreateAuxiliaryAccountHandler

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
		PagingAccounts:            query.NewPagingAccountsHandler(readModel),
		PagingAuxiliaryCategories: query.NewPagingAuxiliaryCategoriesHandler(readModel),
		PagingAuxiliaryAccounts:   query.NewPagingAuxiliaryAccountsHandler(readModel),
		CurrentPeriod:             query.NewCurrentPeriodHandler(readModel),
		PagingPeriods:             query.NewPagingPeriodsHandler(readModel),
		PagingLedgersByPeriod:     query.NewPagingLedgersByPeriodHandler(readModel),
		PagingAuxiliaryLedgers:    query.NewPagingAuxiliaryLedgersHandler(readModel),
		VoucherById:               query.NewVoucherByIdHandler(readModel, userService),
		PagingVouchers:            query.NewPagingVouchersHandler(readModel, userService),
	}
	a.Commands = Commands{
		Initialize: command.NewInitializeHandler(repo, sobService, numberingService),

		AssignAuxiliaryCategory: command.NewAssignAuxiliaryCategoryHandler(repo),

		ClosePeriod: command.NewClosePeriodHandler(repo, numberingService),

		CreateVoucher:       command.NewCreateVoucherHandler(repo, numberingService),
		AuditVoucher:        command.NewAuditVoucherHandler(repo),
		CancelAuditVoucher:  command.NewCancelAuditVoucherHandler(repo),
		ReviewVoucher:       command.NewReviewVoucherHandler(repo),
		CancelReviewVoucher: command.NewCancelReviewVoucherHandler(repo),
		UpdateVoucher:       command.NewUpdateVoucherHandler(repo, numberingService),
		PostVoucher:         command.NewPostVoucherHandler(repo),

		CreateAuxiliaryCategory: command.NewCreateAuxiliaryCategoryHandler(repo),
		CreateAuxiliaryAccount:  command.NewCreateAuxiliaryAccountHandler(repo),

		Migrate: command.NewMigrationHandler(repo),
	}
}
