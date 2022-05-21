package http

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/app"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// ReadAllVouchers godoc
// @Summary List all vouchers by sob
// @Description List all vouchers by sob with pagination
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Success 200 {array} VoucherResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/vouchers/ [get]
func (h Handler) ReadAllVouchers(c *gin.Context) {
	vouchers, err := h.app.Queries.ReadVouchers.HandleReadAll(c, uuid.MustParse(c.Param("sobId")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	res := make([]VoucherResponse, len(vouchers))
	for i, voucher := range vouchers {
		res[i] = mapFromVoucherQuery(voucher)
	}
	c.JSON(http.StatusOK, res)
}

// ReadVoucherById godoc
// @Summary Show voucher by sob and voucher id
// @Description Show voucher by sob and voucher id
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param voucherId path string true "Voucher ID"
// @Success 200 {object} VoucherResponse
// @Failure 404
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{voucherId} [get]
func (h Handler) ReadVoucherById(c *gin.Context) {
	voucher, err := h.app.Queries.ReadVouchers.HandleReadById(c, uuid.MustParse(c.Param("voucherId")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	if voucher.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, mapFromVoucherQuery(voucher))
}

// Audit godoc
// @Summary Audit voucher
// @Description Audit voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param voucherId path string true "Voucher ID"
// @Param AuditVoucherRequest body AuditVoucherRequest true "Audit voucher request, auditor user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{voucherId}/audit [post]
func (h Handler) Audit(c *gin.Context) {
	var req AuditVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	cmd := command.AuditVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("voucherId")),
		Auditor:     req.Auditor,
	}
	if err := h.app.Commands.AuditVoucher.Handle(c, cmd); err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.Status(http.StatusNoContent)
}

// CancelAudit godoc
// @Summary Cancel audit voucher
// @Description Cancel audit voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param voucherId path string true "Voucher ID"
// @Param AuditVoucherRequest body AuditVoucherRequest true "Cancel audit voucher request, auditor user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{voucherId}/cancel-audit [post]
func (h Handler) CancelAudit(c *gin.Context) {
	var req AuditVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	cmd := command.AuditVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("voucherId")),
		Auditor:     req.Auditor,
	}
	if err := h.app.Commands.AuditVoucher.HandleCancel(c, cmd); err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.Status(http.StatusNoContent)
}

// Review godoc
// @Summary Review voucher
// @Description Review voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param voucherId path string true "Voucher ID"
// @Param ReviewVoucherRequest body ReviewVoucherRequest true "Review voucher request, reviewer user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{voucherId}/review [post]
func (h Handler) Review(c *gin.Context) {
	var req ReviewVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	cmd := command.ReviewVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("voucherId")),
		Reviewer:    req.Reviewer,
	}
	if err := h.app.Commands.ReviewVoucher.Handle(c, cmd); err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.Status(http.StatusNoContent)
}

// CancelReview godoc
// @Summary Cancel review voucher
// @Description Cancel review voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param voucherId path string true "Voucher ID"
// @Param ReviewVoucherRequest body ReviewVoucherRequest true "Cancel review voucher request, reviewer user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{voucherId}/cancel-review [post]
func (h Handler) CancelReview(c *gin.Context) {
	var req ReviewVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	cmd := command.ReviewVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("voucherId")),
		Reviewer:    req.Reviewer,
	}
	if err := h.app.Commands.ReviewVoucher.HandleCancel(c, cmd); err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateVoucher godoc
// @Summary Update voucher
// @Description Update voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param voucherId path string true "Voucher ID"
// @Param UpdateVoucherRequest body UpdateVoucherRequest true "Update voucher request"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{voucherId} [patch]
func (h Handler) UpdateVoucher(c *gin.Context) {
	var req UpdateVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	var items []command.LineItemCmd
	for _, itemReq := range req.LineItems {
		item := itemReq.mapToCommand()
		items = append(items, item)
	}
	cmd := command.UpdateVoucherCmd{
		VoucherUUID:     uuid.MustParse(c.Param("voucherId")),
		LineItems:       items,
		TransactionTime: req.TransactionTime,
	}
	if err := h.app.Commands.UpdateVoucher.Handle(c, cmd); err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateVoucher godoc
// @Summary Create voucher
// @Description Create voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param CreateVoucherRequest body CreateVoucherRequest true "Create voucher request"
// @Success 201 {object} VoucherResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/vouchers/ [post]
func (h Handler) CreateVoucher(c *gin.Context) {
	var req CreateVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	cmd := req.mapToCommand()
	cmd.SobId = uuid.MustParse(c.Param("sobId"))
	createdId, err := h.app.Commands.CreateVoucher.Handle(c, cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	createdVoucher, err := h.app.Queries.ReadVouchers.HandleReadById(c, createdId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.JSON(http.StatusCreated, mapFromVoucherQuery(createdVoucher))
}

// Post godoc
// @Summary Post voucher
// @Description Post voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param voucherId path string true "Voucher ID"
// @Success 204
// @Failure 500 {object} Error
// @Router /sob/{sobId}/voucher/{voucherId}/post [post]
func (h Handler) Post(c *gin.Context) {
	cmd := command.PostVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("voucherId")),
	}
	if err := h.app.Commands.PostVoucher.Handle(c, cmd); err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.Status(http.StatusNoContent)
}

func wrapErr(e error) Error {
	var slug string
	se, ok := e.(slugErr)
	if ok {
		slug = se.Slug()
	} else {
		slug = "unknown-error"
	}
	return Error{
		Slug:    slug,
		Message: e.Error(),
	}
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	r.GET("/sob/:sobId/vouchers/", h.ReadAllVouchers)
	r.GET("/sob/:sobId/voucher/:voucherId", h.ReadVoucherById)
	r.POST("/sob/:sobId/vouchers/", h.CreateVoucher)
	r.PATCH("/sob/:sobId/voucher/:voucherId", h.UpdateVoucher)
	r.POST("/sob/:sobId/voucher/:voucherId/audit", h.Audit)
	r.POST("/sob/:sobId/voucher/:voucherId/cancel-audit", h.CancelAudit)
	r.POST("/sob/:sobId/voucher/:voucherId/review", h.Review)
	r.POST("/sob/:sobId/voucher/:voucherId/cancel-review", h.CancelReview)
	r.POST("/sob/:sobId/voucher/:voucherId/post", h.Post)
}
