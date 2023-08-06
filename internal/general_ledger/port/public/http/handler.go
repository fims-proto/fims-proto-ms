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
	r.GET("/sob/:sobId/accounts", h.ReadPagingAccounts)
	r.GET("/sob/:sobId/account/:accountId", h.ReadAccountById)
	r.POST("/sob/:sobId/account/:accountId/assign-auxiliaries", h.AssignAuxiliaryCategoriesToAccount)

	r.GET("/sob/:sobId/auxiliaries", h.ReadPagingAuxiliaryCategories)
	r.POST("/sob/:sobId/auxiliaries", h.CreateAuxiliaryCategory)
	r.GET("/sob/:sobId/auxiliary/:categoryKey/accounts", h.ReadPagingAuxiliaryAccounts)
	r.POST("/sob/:sobId/auxiliary/:categoryKey/accounts", h.CreateAuxiliaryAccount)

	r.GET("/sob/:sobId/periods", h.ReadPagingPeriods)
	r.GET("/sob/:sobId/periods/current", h.ReadSobCurrentPeriod)
	r.GET("/sob/:sobId/period/:periodId/ledgers", h.ReadPagingLedgersByPeriod)
	r.GET("/sob/:sobId/period/:periodId/auxiliary-ledgers", h.ReadPagingAuxiliaryLedgers)
	r.POST("/sob/:sobId/period/:periodId/close", h.ClosePeriod)

	r.GET("/sob/:sobId/vouchers", h.ReadAllVouchers)
	r.GET("/sob/:sobId/voucher/:voucherId", h.ReadVoucherById)
	r.POST("/sob/:sobId/vouchers", h.CreateVoucher)
	r.PATCH("/sob/:sobId/voucher/:voucherId", h.UpdateVoucher)
	r.POST("/sob/:sobId/voucher/:voucherId/audit", h.AuditVoucher)
	r.POST("/sob/:sobId/voucher/:voucherId/cancel-audit", h.CancelAuditVoucher)
	r.POST("/sob/:sobId/voucher/:voucherId/review", h.ReviewVoucher)
	r.POST("/sob/:sobId/voucher/:voucherId/cancel-review", h.CancelReviewVoucher)
	r.POST("/sob/:sobId/voucher/:voucherId/post", h.PostVoucher)
}
