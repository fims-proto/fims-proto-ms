package http

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/app"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	app app.Application
}

func NewHandler(app app.Application) Handler {
	return Handler{app: app}
}

func (h Handler) AllVouchers(c *gin.Context) {
	vouchers, err := h.app.Queries.ReadVouchers.HandleReadAll(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	httpVouchers := []VoucherQry{}
	for _, voucher := range vouchers {
		httpItems := []LineItemQry{}
		for _, item := range voucher.LineItems {
			httpItem := LineItemQry{
				Summary:       item.Summary,
				AccountNumber: item.AccountNumber,
				Debit:         item.Debit,
				Credit:        item.Credit,
			}
			httpItems = append(httpItems, httpItem)
		}
		httpVoucher := VoucherQry{
			UUID:               voucher.UUID,
			Number:             string(voucher.Number),
			CreatedAt:          voucher.CreatedAt,
			AttachmentQuantity: int(voucher.AttachmentQuantity),
			LineItems:          httpItems,
			Debit:              voucher.Debit,
			Credit:             voucher.Credit,
			Creator:            voucher.Creator,
			Reviewer:           voucher.Reviewer,
			Auditor:            voucher.Auditor,
			IsReviewed:         voucher.IsReviewed,
			IsAudited:          voucher.IsAudited,
		}
		httpVouchers = append(httpVouchers, httpVoucher)
	}
	c.JSON(http.StatusOK, httpVouchers)
}

func (h Handler) Audit(c *gin.Context) {
	var httpCmd AuditVoucherCmd
	if err := c.ShouldBind(&httpCmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmd := command.AuditVoucherCmd{
		VoucherUUID: c.Param("uuid"),
		AuditorUUID: httpCmd.Auditor,
	}
	if err := h.app.Commands.AuditVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusAccepted)
}

func (h Handler) Review(c *gin.Context) {
	var httpCmd ReviewVoucherCmd
	if err := c.ShouldBind(&httpCmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmd := command.ReviewVoucherCmd{
		VoucherUUID:  c.Param("uuid"),
		ReviewerUUID: httpCmd.Reviewer,
	}
	if err := h.app.Commands.ReviewVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusAccepted)
}

func (h Handler) Update(c *gin.Context) {
	var httpCmd UpdateVoucherCmd
	if err := c.ShouldBind(&httpCmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	items := []command.LineItemCmd{}
	for _, httpItem := range httpCmd.LineItems {
		item := command.LineItemCmd{
			Summary:       httpItem.Summary,
			AccountNumber: httpItem.AccountNumber,
			Debit:         httpItem.Debit,
			Credit:        httpItem.Credit,
		}
		items = append(items, item)
	}
	cmd := command.UpdateVoucherCmd{
		VoucherUUID: c.Param("uuid"),
		LineItems:   items,
	}
	if err := h.app.Commands.UpdateVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
}

func (h Handler) Record(c *gin.Context) {
	var httpCmd RecordVoucherCmd
	if err := c.ShouldBind(&httpCmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	items := []command.LineItemCmd{}
	for _, httpItem := range httpCmd.LineItems {
		item := command.LineItemCmd{
			Summary:       httpItem.Summary,
			AccountNumber: httpItem.AccountNumber,
			Debit:         httpItem.Debit,
			Credit:        httpItem.Credit,
		}
		items = append(items, item)
	}
	cmd := command.RecordVoucherCmd{
		UUID:               httpCmd.UUID,
		Number:             httpCmd.Number,
		CreatedAt:          httpCmd.CreatedAt,
		AttachmentQuantity: uint(httpCmd.AttachmentQuantity),
		LineItems:          items,
		CreatorUUID:        httpCmd.Creator,
	}
	if err := h.app.Commands.RecordVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
	// TODO figure out a way to avoid hard coded string of url
	c.Writer.Header().Set("Content-Location", "/vouchers/"+cmd.UUID)
}

func (h Handler) VoucherForUUID(c *gin.Context) {
	voucher, err := h.app.Queries.ReadVouchers.HandleReadForUUID(c.Param("uuid"), c.Request.Context())
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	httpItems := []LineItemQry{}
	for _, item := range voucher.LineItems {
		httpItem := LineItemQry{
			Summary:       item.Summary,
			AccountNumber: item.AccountNumber,
			Debit:         item.Debit,
			Credit:        item.Credit,
		}
		httpItems = append(httpItems, httpItem)
	}

	httpVoucher := VoucherQry{
		UUID:               voucher.UUID,
		Number:             voucher.Number,
		CreatedAt:          voucher.CreatedAt,
		AttachmentQuantity: int(voucher.AttachmentQuantity),
		LineItems:          httpItems,
		Debit:              voucher.Debit,
		Credit:             voucher.Credit,
		Creator:            voucher.Creator,
		Reviewer:           voucher.Reviewer,
		Auditor:            voucher.Auditor,
		IsReviewed:         voucher.IsReviewed,
		IsAudited:          voucher.IsAudited,
	}
	c.JSON(http.StatusOK, httpVoucher)
}

func InitRouter(h Handler, r *gin.Engine) {
	g := r.Group("/vouchers")
	{
		g.GET("/", h.AllVouchers)
		g.GET("/:uuid", h.VoucherForUUID)
		g.POST("/", h.Record)
		g.PATCH("/:uuid", h.Update)
		g.POST("/:uuid/audit", h.Audit)
		// TODO cancel audit
		g.POST("/:uuid/review", h.Review)
		// TODO cancel review
	}
}
