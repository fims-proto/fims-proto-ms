package http

import (
	"net/http"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/report/app/command"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SearchReports godoc
//
//	@Text			List all reports by sob
//	@Description	List all reports by sob with pagination
//	@Tags			reports
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Param			$page	query		int		false	"page number"		default(1)
//	@Param			$size	query		int		false	"page size"			default(40)
//	@Param			$sort	query		string	false	"sort on field(s)"	example(updatedAt desc,createdAt)
//	@Param			$filter	query		string	false	"filter on field(s)"
//	@Success		200		{object}	data.PageResponse[ReportResponse]
//	@Failure		500		{object}	Error
//	@Router			/sob/{sobId}/reports [get]
func (h Handler) SearchReports(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.Report], error) {
			return h.app.Queries.PagingReports.Handle(c, uuid.MustParse(c.Param("sobId")), pageRequest)
		},
		reportDTOToVO,
	)
}

// ReadReportTemplateByClass godoc
//
//	@Summary		Get report template by class
//	@Description	Returns the report template for a given SoB and class.
//	@Tags			reports
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Param			class	query		string	true	"Report class"
//	@Success		200		{object}	ReportResponse
//	@Failure		404
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/report/template [get]
func (h Handler) ReadReportTemplateByClass(c *gin.Context) {
	r, err := h.app.Queries.ReportTemplateByClass.Handle(c, uuid.MustParse(c.Param("sobId")), c.Query("class"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if r.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, reportDTOToVO(r))
}

// ReadReportByClassAndPeriod godoc
//
//	@Summary		Get report instance by class and period
//	@Description	Returns the report instance for a given SoB, class, and period.
//	@Tags			reports
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Param			class	query		string	true	"Report class"
//	@Param			period	query		string	true	"Period (YYYY-MM)"
//	@Success		200		{object}	ReportResponse
//	@Failure		404
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/report [get]
func (h Handler) ReadReportByClassAndPeriod(c *gin.Context) {
	t, err := time.Parse("2006-01", c.Query("period"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period must be YYYY-MM"})
		return
	}
	r, err := h.app.Queries.ReportByClassAndPeriod.Handle(c, uuid.MustParse(c.Param("sobId")), c.Query("class"), t.Year(), int(t.Month()))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if r.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, reportDTOToVO(r))
}

// GenerateReport godoc
//
//	@Summary		Generate missing report instance
//	@Description	Creates a report instance from the latest template. If one already exists, returns it unchanged.
//	@Tags			reports
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Param			class	query		string	true	"Report class"
//	@Param			period	query		string	true	"Period (YYYY-MM)"
//	@Success		200		{object}	ReportResponse
//	@Failure		500		{object}	Error
//	@Router			/sob/{sobId}/report/generate [post]
func (h Handler) GenerateReport(c *gin.Context) {
	t, err := time.Parse("2006-01", c.Query("period"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period must be YYYY-MM"})
		return
	}
	actualId, err := h.app.Commands.Generate.Handle(c, command.GenerateReportCmd{
		SobId:        uuid.MustParse(c.Param("sobId")),
		Class:        c.Query("class"),
		FiscalYear:   t.Year(),
		PeriodNumber: int(t.Month()),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	generatedReport, err := h.app.Queries.ReportById.Handle(c, actualId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, reportDTOToVO(generatedReport))
}

// RecalculateReport godoc
//
//	@Summary		Recalculate report amounts
//	@Description	Recalculates amounts while preserving the existing report structure.
//	@Tags			reports
//	@Produce		application/json
//	@Param			sobId		path	string	true	"Sob ID"
//	@Param			reportId	path	string	true	"Report ID"
//	@Success		204
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/report/{reportId}/recalculate [post]
func (h Handler) RecalculateReport(c *gin.Context) {
	if err := h.app.Commands.Recalculate.Handle(c, command.RecalculateReportCmd{ReportId: uuid.MustParse(c.Param("reportId"))}); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// RegenerateReport godoc
//
//	@Summary		Regenerate report instance
//	@Description	Rebuilds an existing report instance from the latest template and recalculates amounts.
//	@Tags			reports
//	@Produce		application/json
//	@Param			sobId		path	string	true	"Sob ID"
//	@Param			reportId	path	string	true	"Report ID"
//	@Success		204
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/report/{reportId}/regenerate [post]
func (h Handler) RegenerateReport(c *gin.Context) {
	if err := h.app.Commands.Regenerate.Handle(c, command.RegenerateReportCmd{ReportId: uuid.MustParse(c.Param("reportId"))}); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateReport godoc
//
//	@Tags			reports
//	@Summary		Update report structure
//	@Description	Updates report columns, row tree, and expressions.
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId				path	string				true	"Sob ID"
//	@Param			reportId			path	string				true	"Report ID"
//	@Param			UpdateReportRequest	body	UpdateReportRequest	true	"Complete report structure"
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/report/{reportId} [patch]
func (h Handler) UpdateReport(c *gin.Context) {
	var req UpdateReportRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd, err := req.mapToCommand(uuid.MustParse(c.Param("reportId")), uuid.MustParse(c.Param("sobId")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = h.app.Commands.UpdateReport.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
