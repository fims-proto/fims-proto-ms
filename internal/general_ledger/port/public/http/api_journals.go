package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SearchJournals godoc
//
//	@Text			List all journals by sob
//	@Description	List all journals by sob with pagination
//	@Tags			journals
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Param			$page	query		int		false	"page number"		default(1)
//	@Param			$size	query		int		false	"page size"			default(40)
//	@Param			$sort	query		string	false	"sort on field(s)"	example(updatedAt desc,createdAt)
//	@Param			$filter	query		string	false	"filter on field(s)"
//	@Success		200		{object}	data.PageResponse[JournalSlimResponse]
//	@Failure		500		{object}	Error
//	@Router			/sob/{sobId}/journals [get]
func (h Handler) SearchJournals(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.Journal], error) {
			return h.app.Queries.PagingJournals.Handle(c, uuid.MustParse(c.Param("sobId")), pageRequest)
		},
		journalDTOToSlimVO,
	)
}

// ReadJournalById godoc
//
//	@Text			Show journal by sob and id
//	@Description	Show journal by sob and id
//	@Tags			journals
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId		path		string	true	"Sob ID"
//	@Param			journalId	path		string	true	"Journal ID"
//	@Success		200			{object}	JournalDetailResponse
//	@Failure		404
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/journal/{journalId} [get]
func (h Handler) ReadJournalById(c *gin.Context) {
	j, err := h.app.Queries.JournalById.Handle(c, uuid.MustParse(c.Param("journalId")))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if j.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, journalDTOToDetailVO(j))
}

// CreateJournal godoc
//
//	@Text			Create journal
//	@Description	Create journal
//	@Tags			journals
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId					path		string					true	"Sob ID"
//	@Param			CreateJournalRequest	body		CreateJournalRequest	true	"Create journal request"
//	@Success		201						{object}	JournalDetailResponse
//	@Failure		400						{object}	Error
//	@Failure		500						{object}	Error
//	@Router			/sob/{sobId}/journals [post]
func (h Handler) CreateJournal(c *gin.Context) {
	var req CreateJournalRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := req.mapToCommand(uuid.MustParse(c.Param("sobId")))
	if err := h.app.Commands.CreateJournal.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	createdJournal, err := h.app.Queries.JournalById.Handle(c, cmd.JournalId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, journalDTOToDetailVO(createdJournal))
}

// UpdateJournal godoc
//
//	@Text			Update journal
//	@Description	Update journal
//	@Tags			journals
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId					path	string					true	"Sob ID"
//	@Param			journalId				path	string					true	"Journal ID"
//	@Param			UpdateJournalRequest	body	UpdateJournalRequest	true	"Update journal request"
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/journal/{journalId} [patch]
func (h Handler) UpdateJournal(c *gin.Context) {
	var req UpdateJournalRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var items []command.JournalLineCmd
	for _, itemReq := range req.JournalLines {
		item := itemReq.mapToCommand()
		items = append(items, item)
	}
	cmd := command.UpdateJournalCmd{
		JournalId:       uuid.MustParse(c.Param("journalId")),
		HeaderText:      req.HeaderText,
		JournalLines:    items,
		TransactionDate: req.TransactionDate,
		Updater:         req.Updater,
	}
	if err := h.app.Commands.UpdateJournal.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// AuditJournal godoc
//
//	@Text			AuditJournal journal
//	@Description	AuditJournal journal
//	@Tags			journals
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId				path	string				true	"Sob ID"
//	@Param			journalId			path	string				true	"Journal ID"
//	@Param			AuditJournalRequest	body	AuditJournalRequest	true	"AuditJournal journal request, auditor user ID"
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/journal/{journalId}/audit [post]
func (h Handler) AuditJournal(c *gin.Context) {
	var req AuditJournalRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.AuditJournalCmd{
		JournalId: uuid.MustParse(c.Param("journalId")),
		Auditor:   req.Auditor,
	}
	if err := h.app.Commands.AuditJournal.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// CancelAuditJournal godoc
//
//	@Text			Cancel audit journal
//	@Description	Cancel audit journal
//	@Tags			journals
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId				path	string				true	"Sob ID"
//	@Param			journalId			path	string				true	"Journal ID"
//	@Param			AuditJournalRequest	body	AuditJournalRequest	true	"Cancel audit journal request, auditor user ID"
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/journal/{journalId}/cancel-audit [post]
func (h Handler) CancelAuditJournal(c *gin.Context) {
	var req AuditJournalRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.CancelAuditJournalCmd{
		JournalId: uuid.MustParse(c.Param("journalId")),
		Auditor:   req.Auditor,
	}
	if err := h.app.Commands.CancelAuditJournal.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ReviewJournal godoc
//
//	@Text			ReviewJournal journal
//	@Description	ReviewJournal journal
//	@Tags			journals
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId					path	string					true	"Sob ID"
//	@Param			journalId				path	string					true	"Journal ID"
//	@Param			ReviewJournalRequest	body	ReviewJournalRequest	true	"ReviewJournal journal request, reviewer user ID"
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/journal/{journalId}/review [post]
func (h Handler) ReviewJournal(c *gin.Context) {
	var req ReviewJournalRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.ReviewJournalCmd{
		JournalId: uuid.MustParse(c.Param("journalId")),
		Reviewer:  req.Reviewer,
	}
	if err := h.app.Commands.ReviewJournal.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// CancelReviewJournal godoc
//
//	@Text			Cancel review journal
//	@Description	Cancel review journal
//	@Tags			journals
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId					path	string					true	"Sob ID"
//	@Param			journalId				path	string					true	"Journal ID"
//	@Param			ReviewJournalRequest	body	ReviewJournalRequest	true	"Cancel review journal request, reviewer user ID"
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/journal/{journalId}/cancel-review [post]
func (h Handler) CancelReviewJournal(c *gin.Context) {
	var req ReviewJournalRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.CancelReviewJournalCmd{
		JournalId: uuid.MustParse(c.Param("journalId")),
		Reviewer:  req.Reviewer,
	}
	if err := h.app.Commands.CancelReviewJournal.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

// PostJournal godoc
//
//	@Text			PostJournal journal
//	@Description	PostJournal journal
//	@Tags			journals
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId				path	string				true	"Sob ID"
//	@Param			journalId			path	string				true	"Journal ID"
//	@Param			PostJournalRequest	body	PostJournalRequest	true	"PostJournal journal request, poster user ID"
//	@Success		204
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/journal/{journalId}/post [post]
func (h Handler) PostJournal(c *gin.Context) {
	var req PostJournalRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	cmd := command.PostJournalCmd{
		JournalId: uuid.MustParse(c.Param("journalId")),
		Poster:    req.Poster,
	}
	if err := h.app.Commands.PostJournal.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
