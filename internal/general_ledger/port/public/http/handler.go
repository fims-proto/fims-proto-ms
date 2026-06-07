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

	huma.Register(accountsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/account-classes", Summary: "List allowed account classes"}, h.ReadAccountClasses)
	huma.Register(accountsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/accounts", Summary: "List all accounts"}, h.ReadAllAccounts)
	huma.Register(accountsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/account/{accountId}", Summary: "Get account by ID"}, h.ReadAccountById)
	huma.Register(accountsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/accounts", Summary: "Create account", DefaultStatus: 201}, h.CreateAccount)
	huma.Register(accountsApi, huma.Operation{Method: "PATCH", Path: "/api/v1/sob/{sobId}/account/{accountId}", Summary: "Update account", DefaultStatus: 204}, h.UpdateAccount)
	huma.Register(accountsApi, huma.Operation{Method: "DELETE", Path: "/api/v1/sob/{sobId}/account/{accountId}", Summary: "Delete account", DefaultStatus: 204}, h.DeleteAccount)

	huma.Register(cashFlowItemsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/cash-flow-items", Summary: "List cash flow items"}, h.ReadCashFlowItems)

	huma.Register(ledgersApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/first-period/ledgers", Summary: "List ledgers in first period"}, h.ReadFirstPeriodLedgers)
	huma.Register(ledgersApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/ledgers/initialize", Summary: "Initialize ledgers balance", DefaultStatus: 204}, h.InitializeLedgers)
	huma.Register(ledgersApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/ledgers", Summary: "List ledgers by period range"}, h.ReadLedgersByPeriodRange)
	huma.Register(ledgersApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/ledgers/transactions", Summary: "Get ledger transactions"}, h.ReadLedgerTransactions)
	huma.Register(ledgersApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/ledgers/dimension-category/{dimensionCategoryId}/options", Summary: "Get ledger by dimension category"}, h.ReadLedgerByDimensionCategory)

	huma.Register(periodsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/periods", Summary: "List all periods"}, h.ReadPeriods)
	huma.Register(periodsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/period/{periodId}/pre-close-check", Summary: "Pre-close check"}, h.PreCloseCheck)
	huma.Register(periodsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/period/{periodId}/close", Summary: "Close period"}, h.ClosePeriod)
	huma.Register(periodsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/periods/batch-pre-close-check", Summary: "Batch pre-close check"}, h.BatchPreCloseCheck)
	huma.Register(periodsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/periods/batch-close", Summary: "Batch close periods", DefaultStatus: 204}, h.ClosePeriods)

	huma.Register(journalsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/journals", Summary: "List journals"}, h.SearchJournals)
	huma.Register(journalsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/journal/{journalId}", Summary: "Get journal by ID"}, h.ReadJournalById)
	huma.Register(journalsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/journals", Summary: "Create journal", DefaultStatus: 201}, h.CreateJournal)
	huma.Register(journalsApi, huma.Operation{Method: "PATCH", Path: "/api/v1/sob/{sobId}/journal/{journalId}", Summary: "Update journal", DefaultStatus: 204}, h.UpdateJournal)
	huma.Register(journalsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/journal/{journalId}/audit", Summary: "Audit journal", DefaultStatus: 204}, h.AuditJournal)
	huma.Register(journalsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/journal/{journalId}/cancel-audit", Summary: "Cancel audit journal", DefaultStatus: 204}, h.CancelAuditJournal)
	huma.Register(journalsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/journal/{journalId}/review", Summary: "Review journal", DefaultStatus: 204}, h.ReviewJournal)
	huma.Register(journalsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/journal/{journalId}/cancel-review", Summary: "Cancel review journal", DefaultStatus: 204}, h.CancelReviewJournal)
	huma.Register(journalsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/journal/{journalId}/post", Summary: "Post journal", DefaultStatus: 204}, h.PostJournal)
	huma.Register(journalsApi, huma.Operation{Method: "DELETE", Path: "/api/v1/sob/{sobId}/journal/{journalId}", Summary: "Delete system journal", DefaultStatus: 204}, h.DeleteSystemJournal)
	huma.Register(journalsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/journals/monthly-closing-journal", Summary: "Create monthly closing journal", DefaultStatus: 201}, h.CreateMonthlyClosingJournal)
	huma.Register(journalsApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/journals/year-end-closing-journal", Summary: "Create year-end closing journal", DefaultStatus: 201}, h.CreateYearEndClosingJournal)
	huma.Register(journalsApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/journals/closing-journal", Summary: "Get closing journal IDs"}, h.GetClosingJournal)
}
