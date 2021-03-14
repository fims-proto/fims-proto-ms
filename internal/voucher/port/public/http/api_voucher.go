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
	return Handler{app: app}
}

func (h Handler) AllVouchers(c *gin.Context) {
	vouchers, err := h.app.Queries.ReadVouchers.HandleReadAll(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	res := []VoucherResponse{}
	for _, voucher := range vouchers {
		res = append(res, mapFromVoucherQuery(voucher))
	}
	c.JSON(http.StatusOK, res)
}

func (h Handler) VoucherByUUID(c *gin.Context) {
	voucher, err := h.app.Queries.ReadVouchers.HandleReadByUUID(c.Request.Context(), uuid.MustParse(c.Param("uuid")))
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
		VoucherUUID: uuid.MustParse(c.Param("uuid")),
		Auditor:     req.Auditor,
	}
	if err := h.app.Commands.AuditVoucher.Handle(c.Request.Context(), cmd); err != nil {
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
		VoucherUUID: uuid.MustParse(c.Param("uuid")),
		Auditor:     req.Auditor,
	}
	if err := h.app.Commands.AuditVoucher.HandleCancel(c.Request.Context(), cmd); err != nil {
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
		VoucherUUID: uuid.MustParse(c.Param("uuid")),
		Reviewer:    req.Reviewer,
	}
	if err := h.app.Commands.ReviewVoucher.Handle(c.Request.Context(), cmd); err != nil {
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
		VoucherUUID: uuid.MustParse(c.Param("uuid")),
		Reviewer:    req.Reviewer,
	}
	if err := h.app.Commands.ReviewVoucher.HandleCancel(c.Request.Context(), cmd); err != nil {
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
		items = append(items, item)
	}
	cmd := command.UpdateVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("uuid")),
		LineItems:   items,
	}
	if err := h.app.Commands.UpdateVoucher.Handle(c.Request.Context(), cmd); err != nil {
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
	newUUID, err := h.app.Commands.RecordVoucher.Handle(c.Request.Context(), req.mapToCommand())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
	// TODO figure out a way to avoid hard coded string of url
	c.Writer.Header().Set("Content-Location", "/vouchers/"+newUUID.String())
}

func (h Handler) Post(c *gin.Context) {
	cmd := command.PostVoucherCmd{
		VoucherUUID: uuid.MustParse(c.Param("uuid")),
	}
	if err := h.app.Commands.PostVoucher.Handler(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func InitRouter(h Handler, r *gin.Engine) {
	g := r.Group("/vouchers")
	{
		g.GET("/", h.AllVouchers)
		g.GET("/:uuid", h.VoucherByUUID)
		g.POST("/", h.Record)
		g.PATCH("/:uuid", h.Update)
		g.POST("/:uuid/audit", h.Audit)
		g.POST("/:uuid/cancel-audit", h.CancelAudit)
		g.POST("/:uuid/review", h.Review)
		g.POST("/:uuid/cancel-review", h.CancelReview)
		g.POST("/:uuid/post", h.Post)
	}
}
