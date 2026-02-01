package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/report/app/command"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReadAllReports godoc
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
func (h Handler) ReadAllReports(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.Report], error) {
			return h.app.Queries.PagingReports.Handle(c, uuid.MustParse(c.Param("sobId")), pageRequest)
		},
		reportDTOToVO,
	)
}

// ReadReportById godoc
//
//	@Text			Show report by sob and id
//	@Description	Show report by sob and id
//	@Tags			reports
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId		path		string	true	"Sob ID"
//	@Param			reportId	path		string	true	"Report ID"
//	@Success		200			{object}	ReportResponse
//	@Failure		404
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/report/{reportId} [get]
func (h Handler) ReadReportById(c *gin.Context) {
	r, err := h.app.Queries.ReportById.Handle(c, uuid.MustParse(c.Param("reportId")))
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
//	@Text			Generate report based on given template
//	@Description	Generate report
//	@Tags			reports
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId					path		string					true	"Sob ID"
//	@Param			GenerateReportRequest	body		GenerateReportRequest	true	"Generate report request"
//	@Success		201						{object}	ReportResponse
//	@Failure		400						{object}	Error
//	@Failure		500						{object}	Error
//	@Router			/sob/{sobId}/report/{reportId}/generate [post]
func (h Handler) GenerateReport(c *gin.Context) {
	var req GenerateReportRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	newReportId := uuid.New()
	cmd := command.GenerateReportCmd{
		TemplateId:       uuid.MustParse(c.Param("reportId")),
		ReportId:         newReportId,
		SobId:            uuid.MustParse(c.Param("sobId")),
		Title:            req.Title,
		AmountTypes:      req.AmountTypes,
		PeriodFiscalYear: req.PeriodFiscalYear,
		PeriodNumber:     req.PeriodNumber,
	}
	if err := h.app.Commands.Generate.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	generatedReport, err := h.app.Queries.ReportById.Handle(c, newReportId)
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
//	@Param			sobId	path	string	true	"Sob ID"
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
//	@Tags        reports
//	@Summary     Update entire report structure
//	@Description Updates report metadata, sections, and items. Supports add, update, delete, and reorder operations in a single atomic transaction.
//	@Accept      application/json
//	@Produce     application/json
//	@Param       sobId               path     string                  true  "Sob ID"
//	@Param       reportId            path     string                  true  "Report ID"
//	@Param       UpdateReportRequest body     UpdateReportRequest     true  "Complete report structure"
//	@Success     200                 {object} UpdateReportResponse
//	@Failure     400                 {object} Error
//	@Failure     500                 {object} Error
//	@Router      /sob/{sobId}/report/{reportId} [patch]
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

	createdItemIds, err := h.app.Commands.UpdateReport.Handle(c, cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, UpdateReportResponse{CreatedItemIds: createdItemIds})
}
