package http

import (
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/localization"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	app       *app.Application
	localizer localization.Localizer
}

func NewHandler(app *app.Application, localizer localization.Localizer) Handler {
	if app == nil {
		panic("nil application")
	}
	return Handler{app: app, localizer: localizer}
}

func InitRouter(h Handler, api huma.API) {
	accountsApi := data.WithTag(api, "accounts")
	cashFlowItemsApi := data.WithTag(api, "cash-flow-items")
	ledgersApi := data.WithTag(api, "ledgers")
	periodsApi := data.WithTag(api, "periods")
	journalsApi := data.WithTag(api, "journals")

	accountDetailLink := map[string]*huma.Link{
		"readAccount": {
			OperationID: "readAccountById",
			Parameters: map[string]any{
				"sobId":     "$request.path.sobId",
				"accountId": "$response.body#/id",
			},
			Description: "Read account detail.",
		},
	}
	accountDetailLinkFromPath := map[string]*huma.Link{
		"readAccount": {
			OperationID: "readAccountById",
			Parameters: map[string]any{
				"sobId":     "$request.path.sobId",
				"accountId": "$request.path.accountId",
			},
			Description: "Read updated account.",
		},
	}
	journalDetailLinkFromBodyID := map[string]*huma.Link{
		"readJournal": {
			OperationID: "readJournalById",
			Parameters: map[string]any{
				"sobId":     "$request.path.sobId",
				"journalId": "$response.body#/id",
			},
			Description: "Read journal detail.",
		},
	}
	journalDetailLinkFromPath := map[string]*huma.Link{
		"readJournal": {
			OperationID: "readJournalById",
			Parameters: map[string]any{
				"sobId":     "$request.path.sobId",
				"journalId": "$request.path.journalId",
			},
			Description: "Read updated journal.",
		},
	}
	journalDetailLinkFromJournalID := map[string]*huma.Link{
		"readJournal": {
			OperationID: "readJournalById",
			Parameters: map[string]any{
				"sobId":     "$request.path.sobId",
				"journalId": "$response.body#/journalId",
			},
			Description: "Read created closing journal.",
		},
	}

	huma.Register(accountsApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/account-classes",
		Summary:     "List allowed account classes",
		OperationID: "readAccountClasses",
	}, h.ReadAccountClasses)

	huma.Register(accountsApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/accounts",
		Summary:     "List all accounts",
		OperationID: "readAllAccounts",
	}, h.ReadAllAccounts)

	huma.Register(accountsApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/account/{accountId}",
		Summary:     "Get account by ID",
		OperationID: "readAccountById",
	}, h.ReadAccountById)

	huma.Register(accountsApi, data.WithResponseLinks(huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/accounts",
		Summary:     "Create account",
		OperationID: "createAccount", DefaultStatus: 201,
	}, 201, accountDetailLink), h.CreateAccount)

	huma.Register(accountsApi, data.WithResponseLinks(huma.Operation{
		Method:      "PATCH",
		Path:        "/api/v1/sob/{sobId}/account/{accountId}",
		Summary:     "Update account",
		OperationID: "updateAccount", DefaultStatus: 204,
	}, 204, accountDetailLinkFromPath), h.UpdateAccount)

	huma.Register(accountsApi, huma.Operation{
		Method:      "DELETE",
		Path:        "/api/v1/sob/{sobId}/account/{accountId}",
		Summary:     "Delete account",
		OperationID: "deleteAccount", DefaultStatus: 204,
	}, h.DeleteAccount)

	huma.Register(cashFlowItemsApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/cash-flow-items",
		Summary:     "List cash flow items",
		OperationID: "readCashFlowItems",
	}, h.ReadCashFlowItems)

	huma.Register(ledgersApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/first-period/ledgers",
		Summary:     "List ledgers in first period",
		OperationID: "readFirstPeriodLedgers",
	}, h.ReadFirstPeriodLedgers)

	huma.Register(ledgersApi, huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/ledgers/initialize",
		Summary:     "Initialize ledgers balance",
		OperationID: "initializeLedgers", DefaultStatus: 204,
	}, h.InitializeLedgers)

	huma.Register(ledgersApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/ledgers",
		Summary:     "List ledgers by period range",
		OperationID: "readLedgersByPeriodRange",
	}, h.ReadLedgersByPeriodRange)

	huma.Register(ledgersApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/ledgers/transactions",
		Summary:     "Get ledger transactions",
		OperationID: "readLedgerTransactions",
	}, h.ReadLedgerTransactions)

	huma.Register(ledgersApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/ledgers/dimension-category/{dimensionCategoryId}/options",
		Summary:     "Get ledger by dimension category",
		OperationID: "readLedgerByDimensionCategory",
	}, h.ReadLedgerByDimensionCategory)

	huma.Register(periodsApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/periods",
		Summary:     "List all periods",
		OperationID: "readPeriods",
	}, h.ReadPeriods)

	huma.Register(periodsApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/period/{periodId}/pre-close-check",
		Summary:     "Pre-close check",
		OperationID: "preCloseCheck",
	}, h.PreCloseCheck)

	huma.Register(periodsApi, huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/period/{periodId}/close",
		Summary:     "Close period",
		OperationID: "closePeriod",
	}, h.ClosePeriod)

	huma.Register(periodsApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/periods/batch-pre-close-check",
		Summary:     "Batch pre-close check",
		OperationID: "batchPreCloseCheck",
	}, h.BatchPreCloseCheck)

	huma.Register(periodsApi, huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/periods/batch-close",
		Summary:     "Batch close periods",
		OperationID: "closePeriods", DefaultStatus: 204,
	}, h.ClosePeriods)

	huma.Register(journalsApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/journals",
		Summary:     "List journals",
		OperationID: "searchJournals",
	}, h.SearchJournals)

	huma.Register(journalsApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/journal/{journalId}",
		Summary:     "Get journal by ID",
		OperationID: "readJournalById",
	}, h.ReadJournalById)

	huma.Register(journalsApi, data.WithResponseLinks(huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/journals",
		Summary:     "Create journal",
		OperationID: "createJournal", DefaultStatus: 201,
	}, 201, journalDetailLinkFromBodyID), h.CreateJournal)

	huma.Register(journalsApi, data.WithResponseLinks(huma.Operation{
		Method:      "PATCH",
		Path:        "/api/v1/sob/{sobId}/journal/{journalId}",
		Summary:     "Update journal",
		OperationID: "updateJournal", DefaultStatus: 204,
	}, 204, journalDetailLinkFromPath), h.UpdateJournal)

	huma.Register(journalsApi, data.WithResponseLinks(huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/journal/{journalId}/audit",
		Summary:     "Audit journal",
		OperationID: "auditJournal", DefaultStatus: 204,
	}, 204, journalDetailLinkFromPath), h.AuditJournal)

	huma.Register(journalsApi, data.WithResponseLinks(huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/journal/{journalId}/cancel-audit",
		Summary:     "Cancel audit journal",
		OperationID: "cancelAuditJournal", DefaultStatus: 204,
	}, 204, journalDetailLinkFromPath), h.CancelAuditJournal)

	huma.Register(journalsApi, data.WithResponseLinks(huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/journal/{journalId}/review",
		Summary:     "Review journal",
		OperationID: "reviewJournal", DefaultStatus: 204,
	}, 204, journalDetailLinkFromPath), h.ReviewJournal)

	huma.Register(journalsApi, data.WithResponseLinks(huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/journal/{journalId}/cancel-review",
		Summary:     "Cancel review journal",
		OperationID: "cancelReviewJournal", DefaultStatus: 204,
	}, 204, journalDetailLinkFromPath), h.CancelReviewJournal)

	huma.Register(journalsApi, data.WithResponseLinks(huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/journal/{journalId}/post",
		Summary:     "Post journal",
		OperationID: "postJournal", DefaultStatus: 204,
	}, 204, journalDetailLinkFromPath), h.PostJournal)

	huma.Register(journalsApi, huma.Operation{
		Method:      "DELETE",
		Path:        "/api/v1/sob/{sobId}/journal/{journalId}",
		Summary:     "Delete system journal",
		OperationID: "deleteSystemJournal", DefaultStatus: 204,
	}, h.DeleteSystemJournal)

	huma.Register(journalsApi, data.WithResponseLinks(huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/journals/monthly-closing-journal",
		Summary:     "Create monthly closing journal",
		OperationID: "createMonthlyClosingJournal", DefaultStatus: 201,
	}, 201, journalDetailLinkFromJournalID), h.CreateMonthlyClosingJournal)

	huma.Register(journalsApi, data.WithResponseLinks(huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/journals/year-end-closing-journal",
		Summary:     "Create year-end closing journal",
		OperationID: "createYearEndClosingJournal", DefaultStatus: 201,
	}, 201, journalDetailLinkFromJournalID), h.CreateYearEndClosingJournal)

	huma.Register(journalsApi, data.WithResponseLinks(huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/journals/closing-journal",
		Summary:     "Get closing journal IDs",
		OperationID: "getClosingJournal",
	}, 200, map[string]*huma.Link{
		"readMonthlyClosingJournal": {
			OperationID: "readJournalById",
			Parameters: map[string]any{
				"sobId":     "$request.path.sobId",
				"journalId": "$response.body#/monthlyClosingJournalId",
			},
			Description: "Read monthly closing journal when monthlyClosingJournalId is present.",
		},
		"readYearEndClosingJournal": {
			OperationID: "readJournalById",
			Parameters: map[string]any{
				"sobId":     "$request.path.sobId",
				"journalId": "$response.body#/yearEndClosingJournalId",
			},
			Description: "Read year-end closing journal when yearEndClosingJournalId is present.",
		},
	}), h.GetClosingJournal)
}
