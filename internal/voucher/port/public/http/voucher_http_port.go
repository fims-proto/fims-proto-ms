package http

import (
	"github.com/gin-gonic/gin"
	"github/fims-proto/fims-proto-ms/internal/voucher/app"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"net/http"
)

type Handler struct {
	app app.Application
}

func NewHandler(app app.Application) Handler {
	return Handler{app: app}
}

func (h Handler) allVouchers(c *gin.Context) {
	// TODO we haven't use openAPI yet. With openAPI we can have a dedicated view struct, as for now, use struct from application layer
	vouchers, err := h.app.Queries.ReadVouchers.HandleReadAll(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, vouchers)
}

func (h Handler) voucherForUUID(c *gin.Context) {
	// TODO openAPI
	v, err := h.app.Queries.ReadVouchers.HandleReadForUUID(c.Param("uuid"), c.Request.Context())
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h Handler) record(c *gin.Context) {
	// TODO openAPI
	var cmd command.RecordVoucherCmd
	if err := c.ShouldBind(&cmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := h.app.Commands.RecordVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
	// TODO figure out a way to avoid hard coded string of url
	c.Writer.Header().Set("Content-Location", "/vouchers/" + cmd.UUID)
}

func (h Handler) update(c *gin.Context){
	var cmd command.UpdateVoucherCmd
	if err := c.ShouldBind(&cmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return 
	}
	cmd.UUID = c.Param("uuid")
	if err := h.app.Commands.UpdateVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
}

func (h Handler) audit(c *gin.Context) {
	// TODO openAPI
	var cmd command.AuditVoucherCmd
	if err := c.ShouldBind(&cmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmd.VoucherUUID = c.Param("uuid")
	if err := h.app.Commands.AuditVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusAccepted)
}

func (h Handler) review(c *gin.Context) {
	// TODO openAPI
	var cmd command.ReviewVoucherCmd
	if err := c.ShouldBind(&cmd); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	cmd.VoucherUUID = c.Param("uuid")
	if err := h.app.Commands.ReviewVoucher.Handle(c.Request.Context(), cmd); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusAccepted)
}

func InitRouter(h Handler, r *gin.Engine) {
	g := r.Group("/vouchers")
	{
		g.GET("/", h.allVouchers)
		g.GET("/:uuid", h.voucherForUUID)
		g.POST("/", h.record)
		g.PATCH("/:uuid", h.update)
		g.POST("/:uuid/audit", h.audit)
		// TODO cancel audit
		g.POST("/:uuid/review", h.review)
		// TODO cancel review
	}
}