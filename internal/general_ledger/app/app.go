package app

import (
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
)

type Queries struct {
	AllAccounts                query.AllAccountsHandler
	AccountById                query.AccountByIdHandler
	AllPeriods                 query.AllPeriodsHandler
	FirstPeriodLedgers         query.FirstPeriodLedgersHandler
	PagingLedgersByPeriod      query.LedgersByPeriodRangeHandler
	LedgersByDimensionCategory query.LedgersByDimensionCategoryHandler
	LedgerEntries              query.LedgerEntriesHandler
	JournalById                query.JournalByIdHandler
	PagingJournals             query.PagingJournalsHandler
	PeriodPreCloseCheck        query.PeriodPreCloseCheckHandler
	ClosingJournalIdsByPeriod  query.ClosingJournalIdsByPeriodHandler
}

type Commands struct {
	Initialize               command.InitializeHandler
	InitializeLedgersBalance command.InitializeLedgersBalanceHandler

	CreateAccount command.CreateAccountHandler
	UpdateAccount command.UpdateAccountHandler
	DeleteAccount command.DeleteAccountHandler

	ClosePeriod command.ClosePeriodHandler

	CreateJournal       command.CreateJournalHandler
	AuditJournal        command.AuditJournalHandler
	CancelAuditJournal  command.CancelAuditJournalHandler
	ReviewJournal       command.ReviewJournalHandler
	CancelReviewJournal command.CancelReviewJournalHandler
	UpdateJournal       command.UpdateJournalHandler
	PostJournal         command.PostJournalHandler

	CreateMonthlyClosingJournal command.CreateMonthlyClosingJournalHandler
	CreateYearEndClosingJournal command.CreateYearEndClosingJournalHandler
	DeleteSystemJournal         command.DeleteSystemJournalHandler

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
	dimensionService service.DimensionService,
) {
	a.Queries = Queries{
		AllAccounts:                query.NewAllAccountsHandler(readModel),
		AccountById:                query.NewAccountByIdHandler(readModel, dimensionService),
		AllPeriods:                 query.NewAllPeriodsHandler(readModel),
		FirstPeriodLedgers:         query.NewFirstPeriodLedgersHandler(readModel),
		PagingLedgersByPeriod:      query.NewLedgersByPeriodRangeHandler(readModel),
		LedgersByDimensionCategory: query.NewLedgersByDimensionCategoryHandler(readModel),
		LedgerEntries:              query.NewLedgerEntriesHandler(readModel),
		JournalById:                query.NewJournalByIdHandler(readModel, userService, dimensionService),
		PagingJournals:             query.NewPagingJournalsHandler(readModel, userService),
		PeriodPreCloseCheck:        query.NewPeriodPreCloseCheckHandler(readModel),
		ClosingJournalIdsByPeriod:  query.NewClosingJournalIdsByPeriodHandler(readModel),
	}
	a.Commands = Commands{
		Initialize:               command.NewInitializeHandler(repo, sobService, numberingService),
		InitializeLedgersBalance: command.NewInitializeLedgersBalanceHandler(repo),

		CreateAccount: command.NewCreateAccountHandler(repo, sobService),
		UpdateAccount: command.NewUpdateAccountHandler(repo, sobService),
		DeleteAccount: command.NewDeleteAccountHandler(repo),

		ClosePeriod: command.NewClosePeriodHandler(repo, numberingService),

		CreateJournal:       command.NewCreateJournalHandler(repo, numberingService, dimensionService, sobService),
		AuditJournal:        command.NewAuditJournalHandler(repo),
		CancelAuditJournal:  command.NewCancelAuditJournalHandler(repo),
		ReviewJournal:       command.NewReviewJournalHandler(repo),
		CancelReviewJournal: command.NewCancelReviewJournalHandler(repo),
		UpdateJournal:       command.NewUpdateJournalHandler(repo, numberingService, dimensionService, sobService),
		PostJournal:         command.NewPostJournalHandler(repo),

		CreateMonthlyClosingJournal: command.NewCreateMonthlyClosingJournalHandler(repo, numberingService, dimensionService, sobService),
		CreateYearEndClosingJournal: command.NewCreateYearEndClosingJournalHandler(repo, numberingService, dimensionService, sobService),
		DeleteSystemJournal:         command.NewDeleteSystemJournalHandler(repo),

		Migrate: command.NewMigrationHandler(repo),
	}
}
