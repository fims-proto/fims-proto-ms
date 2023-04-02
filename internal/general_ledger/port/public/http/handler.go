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

	r.GET("/sob/:sobId/periods", h.ReadPagingPeriods)
	r.GET("/sob/:sobId/periods/current", h.ReadSobCurrentPeriod)
	r.GET("/sob/:sobId/period/:periodId/ledgers", h.ReadPagingLedgersByPeriod)

	r.GET("/sob/:sobId/vouchers", h.ReadAllVouchers)
	r.GET("/sob/:sobId/voucher/:id", h.ReadVoucherById)
	r.POST("/sob/:sobId/vouchers", h.CreateVoucher)
	r.PATCH("/sob/:sobId/voucher/:id", h.UpdateVoucher)
	r.POST("/sob/:sobId/voucher/:id/audit", h.AuditVoucher)
	r.POST("/sob/:sobId/voucher/:id/cancel-audit", h.CancelAuditVoucher)
	r.POST("/sob/:sobId/voucher/:id/review", h.ReviewVoucher)
	r.POST("/sob/:sobId/voucher/:id/cancel-review", h.CancelReviewVoucher)
	r.POST("/sob/:sobId/voucher/:id/post", h.PostVoucher)
}
