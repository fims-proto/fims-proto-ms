package app

import (
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
)

type Queries struct {
	AllAccounts               query.AllAccountsHandler
	PagingAccounts            query.PagingAccountsHandler
	AccountById               query.AccountByIdHandler
	PagingAuxiliaryCategories query.PagingAuxiliaryCategoriesHandler
	AuxiliaryCategoryByKey    query.AuxiliaryCategoryByKeyHandler
	PagingAuxiliaryAccounts   query.PagingAuxiliaryAccountsHandler
	CurrentPeriod             query.CurrentPeriodHandler
	AllPeriods                query.AllPeriodsHandler
	FirstPeriodLedgers        query.FirstPeriodLedgersHandler
	PagingLedgersByPeriod     query.PagingLedgersByPeriodHandler
	LedgerSummary             query.LedgerSummaryHandler
	AuxiliaryLedgerSummary    query.AuxiliaryLedgerSummaryHandler
	PagingLedgerEntries       query.PagingLedgerEntriesHandler
	VoucherById               query.VoucherByIdHandler
	PagingVouchers            query.PagingVouchersHandler
}

type Commands struct {
	Initialize               command.InitializeHandler
	InitializeLedgersBalance command.InitializeLedgersBalanceHandler

	CreateAccount command.CreateAccountHandler
	UpdateAccount command.UpdateAccountHandler

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
		AllAccounts:               query.NewAllAccountsHandler(readModel),
		PagingAccounts:            query.NewPagingAccountsHandler(readModel),
		AccountById:               query.NewAccountByIdHandler(readModel),
		PagingAuxiliaryCategories: query.NewPagingAuxiliaryCategoriesHandler(readModel),
		AuxiliaryCategoryByKey:    query.NewAuxiliaryCategoryByKeyHandler(readModel),
		PagingAuxiliaryAccounts:   query.NewPagingAuxiliaryAccountsHandler(readModel),
		CurrentPeriod:             query.NewCurrentPeriodHandler(readModel),
		AllPeriods:                query.NewAllPeriodsHandler(readModel),
		FirstPeriodLedgers:        query.NewFirstPeriodLedgersHandler(readModel),
		PagingLedgersByPeriod:     query.NewPagingLedgersByPeriodHandler(readModel),
		LedgerSummary:             query.NewLedgerSummaryHandler(readModel),
		AuxiliaryLedgerSummary:    query.NewAuxiliaryLedgerSummaryHandler(readModel),
		PagingLedgerEntries:       query.NewPagingLedgerEntriesHandler(readModel),
		VoucherById:               query.NewVoucherByIdHandler(readModel, userService),
		PagingVouchers:            query.NewPagingVouchersHandler(readModel, userService),
	}
	a.Commands = Commands{
		Initialize:               command.NewInitializeHandler(repo, sobService, numberingService),
		InitializeLedgersBalance: command.NewInitializeLedgersBalanceHandler(repo, sobService),

		CreateAccount: command.NewCreateAccountHandler(repo, sobService),
		UpdateAccount: command.NewUpdateAccountHandler(repo, sobService),

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
