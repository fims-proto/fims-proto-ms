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
	PagingLedgersByPeriod     query.LedgersByPeriodRangeHandler
	LedgerSummary             query.LedgerSummaryHandler
	AuxiliaryLedgerSummary    query.AuxiliaryLedgerSummaryHandler
	PagingLedgerEntries       query.PagingLedgerEntriesHandler
	JournalById               query.JournalByIdHandler
	PagingJournals            query.PagingJournalsHandler
}

type Commands struct {
	Initialize               command.InitializeHandler
	InitializeLedgersBalance command.InitializeLedgersBalanceHandler

	CreateAccount command.CreateAccountHandler
	UpdateAccount command.UpdateAccountHandler

	ClosePeriod command.ClosePeriodHandler

	CreateJournal       command.CreateJournalHandler
	AuditJournal        command.AuditJournalHandler
	CancelAuditJournal  command.CancelAuditJournalHandler
	ReviewJournal       command.ReviewJournalHandler
	CancelReviewJournal command.CancelReviewJournalHandler
	UpdateJournal       command.UpdateJournalHandler
	PostJournal         command.PostJournalHandler

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
		PagingLedgersByPeriod:     query.NewLedgersByPeriodRangeHandler(readModel),
		LedgerSummary:             query.NewLedgerSummaryHandler(readModel),
		AuxiliaryLedgerSummary:    query.NewAuxiliaryLedgerSummaryHandler(readModel),
		PagingLedgerEntries:       query.NewPagingLedgerEntriesHandler(readModel),
		JournalById:               query.NewJournalByIdHandler(readModel, userService),
		PagingJournals:            query.NewPagingJournalsHandler(readModel, userService),
	}
	a.Commands = Commands{
		Initialize:               command.NewInitializeHandler(repo, sobService, numberingService),
		InitializeLedgersBalance: command.NewInitializeLedgersBalanceHandler(repo, sobService),

		CreateAccount: command.NewCreateAccountHandler(repo, sobService),
		UpdateAccount: command.NewUpdateAccountHandler(repo, sobService),

		ClosePeriod: command.NewClosePeriodHandler(repo, numberingService),

		CreateJournal:       command.NewCreateJournalHandler(repo, numberingService),
		AuditJournal:        command.NewAuditJournalHandler(repo),
		CancelAuditJournal:  command.NewCancelAuditJournalHandler(repo),
		ReviewJournal:       command.NewReviewJournalHandler(repo),
		CancelReviewJournal: command.NewCancelReviewJournalHandler(repo),
		UpdateJournal:       command.NewUpdateJournalHandler(repo, numberingService),
		PostJournal:         command.NewPostJournalHandler(repo),

		CreateAuxiliaryCategory: command.NewCreateAuxiliaryCategoryHandler(repo),
		CreateAuxiliaryAccount:  command.NewCreateAuxiliaryAccountHandler(repo),

		Migrate: command.NewMigrationHandler(repo),
	}
}
