package http

import (
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	app *app.Application
}

func NewHandler(app *app.Application) Handler {
	if app == nil {
		panic("nil application")
	}
	return Handler{app: app}
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	r.GET("/sob/:sobId/account-classes", h.ReadAccountClasses)
	r.GET("/sob/:sobId/accounts", h.ReadAllAccounts)
	r.GET("/sob/:sobId/account/:accountId", h.ReadAccountById)
	r.POST("/sob/:sobId/accounts", h.CreateAccount)
	r.PATCH("/sob/:sobId/account/:accountId", h.UpdateAccount)
	r.DELETE("/sob/:sobId/account/:accountId", h.DeleteAccount)

	r.GET("/sob/:sobId/first-period/ledgers", h.ReadFirstPeriodLedgers)
	r.POST("/sob/:sobId/ledgers/initialize", h.InitializeLedgers)
	r.GET("/sob/:sobId/ledger/:accountId", h.ReadLedgerSummary)
	r.GET("/sob/:sobId/ledger/:accountId/entries", h.ReadLedgerEntries)
	r.GET("/sob/:sobId/ledgers/:accountId/dimension/:dimensionCategoryId", h.ReadLedgerDimensionSummary)
	r.GET("/sob/:sobId/periods", h.ReadPeriods)
	r.GET("/sob/:sobId/ledgers", h.ReadLedgersByPeriodRange)
	r.GET("/sob/:sobId/period/:periodId/pre-close-check", h.PreCloseCheck)
	r.POST("/sob/:sobId/period/:periodId/close", h.ClosePeriod)

	r.GET("/sob/:sobId/journals", h.SearchJournals)
	r.GET("/sob/:sobId/journal/:journalId", h.ReadJournalById)
	r.POST("/sob/:sobId/journals", h.CreateJournal)
	r.PATCH("/sob/:sobId/journal/:journalId", h.UpdateJournal)
	r.POST("/sob/:sobId/journal/:journalId/audit", h.AuditJournal)
	r.POST("/sob/:sobId/journal/:journalId/cancel-audit", h.CancelAuditJournal)
	r.POST("/sob/:sobId/journal/:journalId/review", h.ReviewJournal)
	r.POST("/sob/:sobId/journal/:journalId/cancel-review", h.CancelReviewJournal)
	r.POST("/sob/:sobId/journal/:journalId/post", h.PostJournal)
}
