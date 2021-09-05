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

func (h Handler) AllVouchers(c *gin.Context) {
	vouchers, err := h.app.Queries.ReadVouchers.HandleReadAll(c, c.Param("sob"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	res := VouchersResponse{}
	for _, voucher := range vouchers {
		res = append(res, mapFromVoucherQuery(voucher))
	}
	c.JSON(http.StatusOK, res)
}

func (h Handler) VoucherByUUID(c *gin.Context) {
	voucher, err := h.app.Queries.ReadVouchers.HandleReadByUUID(c, uuid.MustParse(c.Param("voucherId")))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, mapFromVoucherQuery(voucher))
}

func (h Handler) Audit(c *gin.Context) {
	var req AuditVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmd := command.AuditVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("voucherId")),
		Auditor:     req.Auditor,
	}
	if err := h.app.Commands.AuditVoucher.Handle(c, cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) CancelAudit(c *gin.Context) {
	var req AuditVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmd := command.AuditVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("voucherId")),
		Auditor:     req.Auditor,
	}
	if err := h.app.Commands.AuditVoucher.HandleCancel(c, cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) Review(c *gin.Context) {
	var req ReviewVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmd := command.ReviewVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("voucherId")),
		Reviewer:    req.Reviewer,
	}
	if err := h.app.Commands.ReviewVoucher.Handle(c, cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) CancelReview(c *gin.Context) {
	var req ReviewVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmd := command.ReviewVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("voucherId")),
		Reviewer:    req.Reviewer,
	}
	if err := h.app.Commands.ReviewVoucher.HandleCancel(c, cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) Update(c *gin.Context) {
	var req UpdateVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
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
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) Record(c *gin.Context) {
	var req RecordVoucherRequest
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmd := req.mapToCommand()
	cmd.Sob = c.Param("sob")
	newUUID, err := h.app.Commands.RecordVoucher.Handle(c, cmd)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
	// TODO figure out a way to avoid hard coded string of url
	c.Writer.Header().Set("Content-Location", fmt.Sprintf("/vouchers/%s/%s", c.Param("sob"), newUUID.String()))
}

func (h Handler) Post(c *gin.Context) {
	cmd := command.PostVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("voucherId")),
	}
	if err := h.app.Commands.PostVoucher.Handle(c, cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func InitRouter(h Handler, r *gin.Engine) {
	g := r.Group("/vouchers/:sob")
	{
		g.GET("/", h.AllVouchers)
		g.GET("/:voucherId", h.VoucherByUUID)
		g.POST("/", h.Record)
		g.PATCH("/:voucherId", h.Update)
		g.POST("/:voucherId/audit", h.Audit)
		g.POST("/:voucherId/cancel-audit", h.CancelAudit)
		g.POST("/:voucherId/review", h.Review)
		g.POST("/:voucherId/cancel-review", h.CancelReview)
		g.POST("/:voucherId/post", h.Post)
	}
}
