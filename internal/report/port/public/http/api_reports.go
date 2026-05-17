package http

import (
	"net/http"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/report/app/command"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/class"

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
//	@Description	Returns the report template for a given SoB and class. Returns 404 if not found.
//	@Tags			reports
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Param			class	query		string	true	"Report class"	Enums(balance_sheet, income_statement)
//	@Success		200		{object}	ReportResponse
//	@Failure		400		{object}	Error
//	@Failure		404
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/report/template [get]
func (h Handler) ReadReportTemplateByClass(c *gin.Context) {
	reportClass, err := class.FromString(c.Query("class"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r, err := h.app.Queries.ReportTemplateByClass.Handle(
		c,
		uuid.MustParse(c.Param("sobId")),
		reportClass,
	)
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
//	@Description	Returns the report instance for a given SoB, class, and period. Returns 404 if not yet generated.
//	@Tags			reports
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Param			class	query		string	true	"Report class"	Enums(balance_sheet, income_statement)
//	@Param			period	query		string	true	"Period (YYYY-MM)"
//	@Success		200		{object}	ReportResponse
//	@Failure		400		{object}	Error
//	@Failure		404
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/report [get]
func (h Handler) ReadReportByClassAndPeriod(c *gin.Context) {
	reportClass, err := class.FromString(c.Query("class"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := time.Parse("2006-01", c.Query("period"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period must be YYYY-MM"})
		return
	}

	r, err := h.app.Queries.ReportByClassAndPeriod.Handle(
		c,
		uuid.MustParse(c.Param("sobId")),
		reportClass,
		t.Year(),
		int(t.Month()),
	)
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
//	@Summary		Generate or regenerate a report instance
//	@Description	Generates a report instance for the given SoB, class, and period. If an instance already exists it is regenerated (amounts recalculated). Returns the resulting report.
//	@Tags			reports
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Param			class	query		string	true	"Report class"	Enums(balance_sheet, income_statement)
//	@Param			period	query		string	true	"Period (YYYY-MM)"
//	@Success		200		{object}	ReportResponse
//	@Failure		400		{object}	Error
//	@Failure		500		{object}	Error
//	@Router			/sob/{sobId}/report/generate [post]
func (h Handler) GenerateReport(c *gin.Context) {
	reportClass, err := class.FromString(c.Query("class"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := time.Parse("2006-01", c.Query("period"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period must be YYYY-MM"})
		return
	}

	cmd := command.GenerateReportCmd{
		SobId:        uuid.MustParse(c.Param("sobId")),
		Class:        reportClass,
		FiscalYear:   t.Year(),
		PeriodNumber: int(t.Month()),
	}
	actualId, err := h.app.Commands.Generate.Handle(c, cmd)
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

// RegenerateReport godoc
//
//	@Text			Regenerate report amounts
//	@Description	Regenerate report
//	@Tags			reports
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId		path	string	true	"Sob ID"
//	@Param			reportId	path	string	true	"Report ID"
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/report/{reportId}/regenerate [post]
func (h Handler) RegenerateReport(c *gin.Context) {
	cmd := command.RegenerateReportCmd{ReportId: uuid.MustParse(c.Param("reportId"))}
	if err := h.app.Commands.Regenerate.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateReport godoc
//
//	@Tags			reports
//	@Summary		Update entire report structure
//	@Description	Updates report metadata, sections, and items. Supports add, update, delete, and reorder operations in a single atomic transaction. Sections can contain nested sections. Items are sequenced by their position in the array.
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId				path	string				true	"Sob ID"
//	@Param			reportId			path	string				true	"Report ID"
//	@Param			UpdateReportRequest	body	UpdateReportRequest	true	"Complete report structure"
//	@Success		204					"No Content"
//	@Failure		400					{object}	Error
//	@Failure		500					{object}	Error
//	@Router			/sob/{sobId}/report/{reportId} [patch]
func (h Handler) UpdateReport(c *gin.Context) {
	var req UpdateReportRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	cmd, err := req.mapToCommand(
		uuid.MustParse(c.Param("reportId")),
		uuid.MustParse(c.Param("sobId")),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.app.Commands.UpdateReport.Handle(c, cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
