package http

import (
	"github/fims-proto/fims-proto-ms/internal/common/datav3"
	"github/fims-proto/fims-proto-ms/internal/journal/app/query"
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/journal/app"
	"github/fims-proto/fims-proto-ms/internal/journal/app/command"

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

// ReadAllJournalEntries godoc
// @Text List all journal entries by sob
// @Description List all journals by sob with pagination
// @Tags journals
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $filter query string false "filter on field(s)"
// @Success 200 {array} JournalEntryResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/journal-entries/ [get]
func (h Handler) ReadAllJournalEntries(c *gin.Context) {
	datav3.PagingResponseProcessor(
		c,
		func(pageRequest datav3.PageRequest) (datav3.Page[query.JournalEntry], error) {
			return h.app.Queries.PagingJournalEntries.Handle(c, uuid.MustParse(c.Param("sobId")), pageRequest)
		},
		JournalEntryDTOToVO,
	)
}

// ReadJournalEntryById godoc
// @Text Show journal entry by sob and entry id
// @Description Show journal entry by sob and entry id
// @Tags journals
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param entryId path string true "Entry ID"
// @Success 200 {object} JournalEntryResponse
// @Failure 404
// @Failure 500 {object} Error
// @Router /sob/{sobId}/journal-entry/{entryId} [get]
func (h Handler) ReadJournalEntryById(c *gin.Context) {
	entry, err := h.app.Queries.JournalEntryById.Handle(c, uuid.MustParse(c.Param("entryId")))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if entry.EntryId == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, JournalEntryDTOToVO(entry))
}

// Audit godoc
// @Text Audit journal entry
// @Description Audit journal entry
// @Tags journals
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param entryId path string true "Entry ID"
// @Param AuditJournalEntryRequest body AuditJournalEntryRequest true "Audit journal entry request, auditor user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/journal-entry/{entryId}/audit [post]
func (h Handler) Audit(c *gin.Context) {
	var req AuditJournalEntryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.AuditJournalEntryCmd{
		EntryId: uuid.MustParse(c.Param("entryId")),
		Auditor: req.Auditor,
	}
	if err := h.app.Commands.AuditJournalEntry.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// CancelAudit godoc
// @Text Cancel audit journal entry
// @Description Cancel audit journal entry
// @Tags journals
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param entryId path string true "Entry ID"
// @Param AuditJournalEntryRequest body AuditJournalEntryRequest true "Cancel audit journal entry request, auditor user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/journal-entry/{entryId}/cancel-audit [post]
func (h Handler) CancelAudit(c *gin.Context) {
	var req AuditJournalEntryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.CancelAuditJournalEntryCmd{
		EntryId: uuid.MustParse(c.Param("entryId")),
		Auditor: req.Auditor,
	}
	if err := h.app.Commands.CancelAuditJournalEntry.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Review godoc
// @Text Review journal entry
// @Description Review journal entry
// @Tags journals
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param entryId path string true "Entry ID"
// @Param ReviewJournalEntryRequest body ReviewJournalEntryRequest true "Review journal entry request, reviewer user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/journal-entry/{entryId}/review [post]
func (h Handler) Review(c *gin.Context) {
	var req ReviewJournalEntryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.ReviewJournalEntryCmd{
		EntryId:  uuid.MustParse(c.Param("entryId")),
		Reviewer: req.Reviewer,
	}
	if err := h.app.Commands.ReviewJournalEntry.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// CancelReview godoc
// @Text Cancel review journal entry
// @Description Cancel review journal entry
// @Tags journals
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param entryId path string true "Entry ID"
// @Param ReviewJournalEntryRequest body ReviewJournalEntryRequest true "Cancel review journal entry request, reviewer user ID"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/journal-entry/{entryId}/cancel-review [post]
func (h Handler) CancelReview(c *gin.Context) {
	var req ReviewJournalEntryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.CancelReviewJournalEntryCmd{
		EntryId:  uuid.MustParse(c.Param("entryId")),
		Reviewer: req.Reviewer,
	}
	if err := h.app.Commands.CancelReviewJournalEntry.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// Post godoc
// @Text Post journal entry
// @Description Post journal entry
// @Tags journals
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param entryId path string true "Entry ID"
// @Param PostJournalEntryRequest body PostJournalEntryRequest true "Post journal entry request, poster user ID"
// @Success 204
// @Failure 500 {object} Error
// @Router /sob/{sobId}/journal-entry/{entryId}/post [post]
func (h Handler) Post(c *gin.Context) {
	var req PostJournalEntryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	cmd := command.PostJournalEntryCmd{
		EntryId: uuid.MustParse(c.Param("entryId")),
		Poster:  req.Poster,
	}
	if err := h.app.Commands.PostJournalEntry.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateJournalEntry godoc
// @Text Update journal entry
// @Description Update journal entry
// @Tags journals
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param entryId path string true "Entry ID"
// @Param UpdateJournalEntryRequest body UpdateJournalEntryRequest true "Update journal entry request"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/journal-entry/{entryId} [patch]
func (h Handler) UpdateJournalEntry(c *gin.Context) {
	var req UpdateJournalEntryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var items []command.LineItemCmd
	for _, itemReq := range req.LineItems {
		item := itemReq.mapToCommand()
		items = append(items, item)
	}
	cmd := command.UpdateJournalEntryCmd{
		EntryId:         uuid.MustParse(c.Param("entryId")),
		LineItems:       items,
		TransactionTime: req.TransactionTime,
		Updater:         req.Updater,
	}
	if err := h.app.Commands.UpdateJournalEntry.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateJournalEntry godoc
// @Text Create journal entry
// @Description Create journal entry
// @Tags journals
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param CreateJournalEntryRequest body CreateJournalEntryRequest true "Create journal entry request"
// @Success 201 {object} JournalEntryResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/journal-entries/ [post]
func (h Handler) CreateJournalEntry(c *gin.Context) {
	var req CreateJournalEntryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := req.mapToCommand(uuid.MustParse(c.Param("sobId")))
	err := h.app.Commands.CreateJournalEntry.Handle(c, cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	createdEntry, err := h.app.Queries.JournalEntryById.Handle(c, cmd.EntryId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, JournalEntryDTOToVO(createdEntry))
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	r.GET("/sob/:sobId/journal-entries/", h.ReadAllJournalEntries)
	r.GET("/sob/:sobId/journal-entry/:entryId", h.ReadJournalEntryById)
	r.POST("/sob/:sobId/journal-entries/", h.CreateJournalEntry)
	r.PATCH("/sob/:sobId/journal-entry/:entryId", h.UpdateJournalEntry)
	r.POST("/sob/:sobId/journal-entry/:entryId/audit", h.Audit)
	r.POST("/sob/:sobId/journal-entry/:entryId/cancel-audit", h.CancelAudit)
	r.POST("/sob/:sobId/journal-entry/:entryId/review", h.Review)
	r.POST("/sob/:sobId/journal-entry/:entryId/cancel-review", h.CancelReview)
	r.POST("/sob/:sobId/journal-entry/:entryId/post", h.Post)
}
