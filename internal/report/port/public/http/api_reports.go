package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/report/app/command"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"
)

// ReadAllReports godoc
// @Text List all reports by sob
// @Description List all reports by sob with pagination
// @Tags reports
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $filter query string false "filter on field(s)"
// @Success 200 {array} ReportResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/reports [get]
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
// @Text Show report by sob and id
// @Description Show report by sob and id
// @Tags reports
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param reportId path string true "Report ID"
// @Success 200 {object} ReportResponse
// @Failure 404
// @Failure 500 {object} Error
// @Router /sob/{sobId}/report/{reportId} [get]
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
// @Text Generate report based on given template
// @Description Generate report
// @Tags reports
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param GenerateReportRequest body GenerateReportRequest true "Generate report request"
// @Success 201 {object} ReportResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/report/{reportId}/generate [post]
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
// @Text Regenerate report amounts
// @Description Regenerate report
// @Tags reports
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/report/{reportId}/regenerate [post]
func (h Handler) RegenerateReport(c *gin.Context) {
	cmd := command.RegenerateReportCmd{ReportId: uuid.MustParse(c.Param("reportId"))}
	if err := h.app.Commands.Regenerate.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateItem godoc
// @Text Update a report item
// @Description Update a report item
// @Tags reports
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param reportId path string true "Report ID"
// @Param itemId path string true "Item ID"
// @Param UpdateItemRequest body UpdateItemRequest true "Update report item request"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/report/{reportId}/item/{itemId} [patch]
func (h Handler) UpdateItem(c *gin.Context) {
	var req UpdateItemRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	cmd := req.mapToCommand(uuid.MustParse(c.Param("sobId")), uuid.MustParse(c.Param("itemId")))
	if err := h.app.Commands.UpdateItem.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
