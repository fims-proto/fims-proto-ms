package http

import (
	"fmt"
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

// AllVouchers godoc
// @Summary List all vouchers by sob
// @Description List all vouchers by sob with paginagion
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sob path string true "Sob Id"
// @Success 200 {array} VoucherResponse
// @Failure 500 {object} Error
// @Router /vouchers/{sob}/ [get]
func (h Handler) AllVouchers(c *gin.Context) {
	vouchers, err := h.app.Queries.ReadVouchers.HandleReadAll(c, c.Param("sob"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	res := []VoucherResponse{}
	for _, voucher := range vouchers {
		res = append(res, mapFromVoucherQuery(voucher))
	}
	c.JSON(http.StatusOK, res)
}

// VoucherByUUID godoc
// @Summary Show voucher by sob and voucher id
// @Description Show voucher by sob and voucher id
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sob path string true "Sob Id"
// @Param voucherId path string true "Voucher Id"
// @Success 200 {object} VoucherResponse
// @Failure 404
// @Failure 500 {object} Error
// @Router /vouchers/{sob}/{voucherId} [get]
func (h Handler) VoucherByUUID(c *gin.Context) {
	voucher, err := h.app.Queries.ReadVouchers.HandleReadByUUID(c, uuid.MustParse(c.Param("voucherId")))
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
// @Param sob path string true "Sob Id"
// @Param voucherId path string true "Voucher Id"
// @Param AuditVoucherRequest body AuditVoucherRequest true "Audit voucher request, auditor user Id"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /vouchers/{sob}/{voucherId}/audit [post]
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
// @Param sob path string true "Sob Id"
// @Param voucherId path string true "Voucher Id"
// @Param AuditVoucherRequest body AuditVoucherRequest true "Cancel audit voucher request, , auditor user Id"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /vouchers/{sob}/{voucherId}/cancel-audit [post]
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
// @Param sob path string true "Sob Id"
// @Param voucherId path string true "Voucher Id"
// @Param ReviewVoucherRequest body ReviewVoucherRequest true "Review voucher request, reviewer user Id"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /vouchers/{sob}/{voucherId}/review [post]
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
// @Param sob path string true "Sob Id"
// @Param voucherId path string true "Voucher Id"
// @Param ReviewVoucherRequest body ReviewVoucherRequest true "Cancel review voucher request, reviewer user Id"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /vouchers/{sob}/{voucherId}/cancel-review [post]
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

// Update godoc
// @Summary Update voucher
// @Description Update voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sob path string true "Sob Id"
// @Param voucherId path string true "Voucher Id"
// @Param UpdateVoucherRequest body UpdateVoucherRequest true "Update voucher request"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /vouchers/{sob}/{voucherId} [patch]
func (h Handler) Update(c *gin.Context) {
	var req UpdateVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	items := []command.LineItemCmd{}
	for _, itemReq := range req.LineItems {
		item := itemReq.mapToCommand()
		item.Id = uuid.MustParse(itemReq.Id)
		items = append(items, item)
	}
	cmd := command.UpdateVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("voucherId")),
		LineItems:   items,
	}
	if err := h.app.Commands.UpdateVoucher.Handle(c, cmd); err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.Status(http.StatusNoContent)
}

// Record godoc
// @Summary Create voucher
// @Description Create voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sob path string true "Sob Id"
// @Param RecordVoucherRequest body RecordVoucherRequest true "Create voucher request"
// @Success 201
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /vouchers/{sob}/ [post]
func (h Handler) Record(c *gin.Context) {
	var req RecordVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	cmd := req.mapToCommand()
	cmd.Sob = c.Param("sob")
	newUUID, err := h.app.Commands.RecordVoucher.Handle(c, cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.Status(http.StatusCreated)
	// TODO figure out a way to avoid hard coded string of url
	c.Writer.Header().Set("Content-Location", fmt.Sprintf("/vouchers/%s/%s", c.Param("sob"), newUUID.String()))
}

// Post godoc
// @Summary Post voucher
// @Description Post voucher
// @Tags vouchers
// @Accept application/json
// @Produce application/json
// @Param sob path string true "Sob Id"
// @Param voucherId path string true "Voucher Id"
// @Success 204
// @Failure 500 {object} Error
// @Router /vouchers/{sob}/{voucherId}/post [post]
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
	se, ok := e.(sluggableErr)
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
	g := r.Group("/vouchers/:sob/")
	{
		g.GET("", h.AllVouchers)
		g.GET(":voucherId", h.VoucherByUUID)
		g.POST("", h.Record)
		g.PATCH(":voucherId", h.Update)
		g.POST(":voucherId/audit", h.Audit)
		g.POST(":voucherId/cancel-audit", h.CancelAudit)
		g.POST(":voucherId/review", h.Review)
		g.POST(":voucherId/cancel-review", h.CancelReview)
		g.POST(":voucherId/post", h.Post)
	}
}
