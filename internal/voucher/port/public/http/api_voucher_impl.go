//  implementation of generated openapi
package http

import (
	"github.com/gin-gonic/gin"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"net/http"
)

func (h Handler) _AllVouchers(c *gin.Context) {
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
			Number:             int32(voucher.Number),
			CreatedAt:          voucher.CreatedAt,
			AttachmentQuantity: int32(voucher.AttachmentQuantity),
			LineItems:          httpItems,
			Debit:              voucher.Debit,
			Credit:             voucher.Credit,
			CreatorUUID:        voucher.CreatorUUID,
			ReviewerUUID:       voucher.ReviewerUUID,
			AuditorUUID:        voucher.AuditorUUID,
			IsReviewed:         voucher.IsReviewed,
			IsAudited:          voucher.IsAudited,
		}
		httpVouchers = append(httpVouchers, httpVoucher)
	}
	c.JSON(http.StatusOK, httpVouchers)
}

func (h Handler) _Audit(c *gin.Context) {
	var httpCmd AuditVoucherCmd
	if err := c.ShouldBind(&httpCmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	httpCmd.VoucherUUID = c.Param("uuid")
	cmd := command.AuditVoucherCmd{
		AuditorUUID: httpCmd.AuditorUUID,
		VoucherUUID: httpCmd.VoucherUUID,
	}
	if err := h.app.Commands.AuditVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusAccepted)
}

func (h Handler) _Review(c *gin.Context) {
	var httpCmd ReviewVoucherCmd
	if err := c.ShouldBind(&httpCmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	httpCmd.VoucherUUID = c.Param("uuid")
	cmd := command.ReviewVoucherCmd{
		VoucherUUID:  httpCmd.VoucherUUID,
		ReviewerUUID: httpCmd.ReviewerUUID,
	}
	if err := h.app.Commands.ReviewVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusAccepted)
}

func (h Handler) _Update(c *gin.Context) {
	var httpCmd UpdateVoucherCmd
	if err := c.ShouldBind(&httpCmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	httpCmd.VoucherUUID = c.Param("uuid")
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
		VoucherUUID: httpCmd.VoucherUUID,
		LineItems:   items,
	}
	if err := h.app.Commands.UpdateVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
}

func (h Handler) _Record(c *gin.Context) {
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
		Number:             uint(httpCmd.Number),
		CreatedAt:          httpCmd.CreatedAt,
		AttachmentQuantity: uint(httpCmd.AttachmentQuantity),
		LineItems:          items,
		Debit:              httpCmd.Debit,
		Credit:             httpCmd.Credit,
		CreatorUUID:        httpCmd.CreatorUUID,
	}
	if err := h.app.Commands.RecordVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
	// TODO figure out a way to avoid hard coded string of url
	c.Writer.Header().Set("Content-Location", "/vouchers/"+cmd.UUID)
}

func (h Handler) _VoucherForUUID(c *gin.Context) {
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
		Number:             int32(voucher.Number),
		CreatedAt:          voucher.CreatedAt,
		AttachmentQuantity: int32(voucher.AttachmentQuantity),
		LineItems:          httpItems,
		Debit:              voucher.Debit,
		Credit:             voucher.Credit,
		CreatorUUID:        voucher.CreatorUUID,
		ReviewerUUID:       voucher.ReviewerUUID,
		AuditorUUID:        voucher.AuditorUUID,
		IsReviewed:         voucher.IsReviewed,
		IsAudited:          voucher.IsAudited,
	}
	c.JSON(http.StatusOK, httpVoucher)
}
