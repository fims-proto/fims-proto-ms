package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
)

// ReadAllVouchers godoc
// @Text List all voucher by sob
// @Description List all vouchers by sob with pagination
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $filter query string false "filter on field(s)"
// @Success 200 {array} VoucherResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/vouchers [get]
func (h Handler) ReadAllVouchers(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.Voucher], error) {
			return h.app.Queries.PagingVouchers.Handle(c, uuid.MustParse(c.Param("sobId")), pageRequest)
		},
		VoucherDTOToVO,
	)
}

// ReadVoucherById godoc
// @Text Show voucher by sob and id
// @Description Show voucher by sob and id
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param id path string true "Voucher ID"
// @Success 200 {object} VoucherResponse
// @Failure 404
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{id} [get]
func (h Handler) ReadVoucherById(c *gin.Context) {
	v, err := h.app.Queries.VoucherById.Handle(c, uuid.MustParse(c.Param("id")))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if v.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, VoucherDTOToVO(v))
}

// AuditVoucher godoc
// @Text AuditVoucher voucher
// @Description AuditVoucher voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param id path string true "Voucher ID"
// @Param AuditVoucherRequest body AuditVoucherRequest true "AuditVoucher voucher request, auditor user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{id}/audit [post]
func (h Handler) AuditVoucher(c *gin.Context) {
	var req AuditVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.AuditVoucherCmd{
		VoucherId: uuid.MustParse(c.Param("id")),
		Auditor:   req.Auditor,
	}
	if err := h.app.Commands.AuditVoucher.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// CancelAuditVoucher godoc
// @Text Cancel audit voucher
// @Description Cancel audit voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param id path string true "Voucher ID"
// @Param AuditVoucherRequest body AuditVoucherRequest true "Cancel audit voucher request, auditor user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{id}/cancel-audit [post]
func (h Handler) CancelAuditVoucher(c *gin.Context) {
	var req AuditVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.CancelAuditVoucherCmd{
		VoucherId: uuid.MustParse(c.Param("id")),
		Auditor:   req.Auditor,
	}
	if err := h.app.Commands.CancelAuditVoucher.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ReviewVoucher godoc
// @Text ReviewVoucher voucher
// @Description ReviewVoucher voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param id path string true "Voucher ID"
// @Param ReviewVoucherRequest body ReviewVoucherRequest true "ReviewVoucher voucher request, reviewer user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{id}/review [post]
func (h Handler) ReviewVoucher(c *gin.Context) {
	var req ReviewVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.ReviewVoucherCmd{
		VoucherId: uuid.MustParse(c.Param("id")),
		Reviewer:  req.Reviewer,
	}
	if err := h.app.Commands.ReviewVoucher.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// CancelReviewVoucher godoc
// @Text Cancel review voucher
// @Description Cancel review voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param id path string true "Voucher ID"
// @Param ReviewVoucherRequest body ReviewVoucherRequest true "Cancel review voucher request, reviewer user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{id}/cancel-review [post]
func (h Handler) CancelReviewVoucher(c *gin.Context) {
	var req ReviewVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.CancelReviewVoucherCmd{
		VoucherId: uuid.MustParse(c.Param("id")),
		Reviewer:  req.Reviewer,
	}
	if err := h.app.Commands.CancelReviewVoucher.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// PostVoucher godoc
// @Text PostVoucher voucher
// @Description PostVoucher voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param id path string true "Voucher ID"
// @Param PostVoucherRequest body PostVoucherRequest true "PostVoucher voucher request, poster user ID"
// @Success 204
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{id}/post [post]
func (h Handler) PostVoucher(c *gin.Context) {
	var req PostVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	cmd := command.PostVoucherCmd{
		VoucherId: uuid.MustParse(c.Param("id")),
		Poster:    req.Poster,
	}
	if err := h.app.Commands.PostVoucher.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateVoucher godoc
// @Text Update voucher
// @Description Update voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param id path string true "Voucher ID"
// @Param UpdateVoucherRequest body UpdateVoucherRequest true "Update voucher request"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{id} [patch]
func (h Handler) UpdateVoucher(c *gin.Context) {
	var req UpdateVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var items []command.LineItemCmd
	for _, itemReq := range req.LineItems {
		item := itemReq.mapToCommand()
		items = append(items, item)
	}
	cmd := command.UpdateVoucherCmd{
		VoucherId:       uuid.MustParse(c.Param("id")),
		HeaderText:      req.HeaderText,
		LineItems:       items,
		TransactionTime: req.TransactionTime,
		Updater:         req.Updater,
	}
	if err := h.app.Commands.UpdateVoucher.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateVoucher godoc
// @Text Create voucher
// @Description Create voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param CreateVoucherRequest body CreateVoucherRequest true "Create voucher request"
// @Success 201 {object} VoucherResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/vouchers [post]
func (h Handler) CreateVoucher(c *gin.Context) {
	var req CreateVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := req.mapToCommand(uuid.MustParse(c.Param("sobId")))
	err := h.app.Commands.CreateVoucher.Handle(c, cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	createdVoucher, err := h.app.Queries.VoucherById.Handle(c, cmd.VoucherId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, VoucherDTOToVO(createdVoucher))
}
