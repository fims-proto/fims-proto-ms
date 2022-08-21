package app

import (
	"github/fims-proto/fims-proto-ms/internal/journal/app/command"
	"github/fims-proto/fims-proto-ms/internal/journal/app/query"
	"github/fims-proto/fims-proto-ms/internal/journal/app/service"
	"github/fims-proto/fims-proto-ms/internal/journal/domain"
)

type Queries struct {
	JournalEntryById     query.JournalEntryByIdHandler
	PagingJournalEntries query.PagingJournalEntriesHandler
}

type Commands struct {
	CreateJournalEntry       command.CreateJournalEntryHandler
	AuditJournalEntry        command.AuditJournalEntryHandler
	CancelAuditJournalEntry  command.CancelAuditJournalEntryHandler
	ReviewJournalEntry       command.ReviewJournalEntryHandler
	CancelReviewJournalEntry command.CancelReviewJournalEntryHandler
	UpdateJournalEntry       command.UpdateJournalEntryHandler
	PostJournalEntry         command.PostJournalEntryHandler
	Migrate                  command.MigrationHandler
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
	readModel query.JournalReadModel,
	accountService service.AccountService,
	userService service.UserService,
	numberingService service.NumberingService,
) {
	a.Queries = Queries{
		JournalEntryById:     query.NewJournalEntryByIdHandler(readModel, accountService, userService),
		PagingJournalEntries: query.NewPagingJournalEntriesHandler(readModel, userService),
	}
	a.Commands = Commands{
		CreateJournalEntry:       command.NewCreateJournalEntryHandler(repo, accountService, numberingService),
		AuditJournalEntry:        command.NewAuditJournalEntryHandler(repo),
		CancelAuditJournalEntry:  command.NewCancelAuditJournalEntryHandler(repo),
		ReviewJournalEntry:       command.NewReviewJournalEntryHandler(repo),
		CancelReviewJournalEntry: command.NewCancelReviewJournalEntryHandler(repo),
		UpdateJournalEntry:       command.NewUpdateJournalEntryHandler(repo, accountService),
		PostJournalEntry:         command.NewPostJournalEntryHandler(repo, accountService),
		Migrate:                  command.NewMigrationHandler(repo),
	}
}
