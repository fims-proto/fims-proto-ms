package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReadPagingAuxiliaryCategories godoc
//
//	@Text			List all auxiliary categories
//	@Description	List all auxiliary categories
//	@Tags			auxiliary accounts
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Param			$page	query		int		false	"page number"			default(1)
//	@Param			$size	query		int		false	"page size"				default(40)
//	@Param			$sort	query		string	false	"sort on field(s)"		example(updatedAt desc,createdAt)
//	@Param			$filter	query		string	false	"filter on field(s)"	example(title eq 'something' and amount lt 10)
//	@Success		200		{array}		AuxiliaryCategoryResponse
//	@Failure		500		{object}	Error
//	@Router			/sob/{sobId}/auxiliaries [get]
func (h Handler) ReadPagingAuxiliaryCategories(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.AuxiliaryCategory], error) {
			return h.app.Queries.PagingAuxiliaryCategories.Handle(c, uuid.MustParse(c.Param("sobId")), pageRequest)
		},
		auxiliaryCategoryDTOToVO,
	)
}

// ReadPagingAuxiliaryAccounts godoc
//
//	@Text			List all auxiliary accounts
//	@Description	List all auxiliary accounts
//	@Tags			auxiliary accounts
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId		path		string	true	"Sob ID"
//	@Param			categoryKey	path		string	true	"Category Key"
//	@Param			$page		query		int		false	"page number"			default(1)
//	@Param			$size		query		int		false	"page size"				default(40)
//	@Param			$sort		query		string	false	"sort on field(s)"		example(updatedAt desc,createdAt)
//	@Param			$filter		query		string	false	"filter on field(s)"	example(title eq 'something' and amount lt 10)
//	@Success		200			{object}	data.PageResponse[AuxiliaryAccountResponse]
//	@Failure		500			{object}	Error
//	@Router			/sob/{sobId}/auxiliary/{categoryKey}/accounts [get]
func (h Handler) ReadPagingAuxiliaryAccounts(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.AuxiliaryAccount], error) {
			return h.app.Queries.PagingAuxiliaryAccounts.Handle(c, uuid.MustParse(c.Param("sobId")), c.Param("categoryKey"), pageRequest)
		},
		auxiliaryAccountDTOToVO,
	)
}

// CreateAuxiliaryCategory godoc
//
//	@Text			Create auxiliary category
//	@Description	Create auxiliary category
//	@Tags			auxiliary accounts
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId							path	string							true	"Sob ID"
//	@Param			CreateAuxiliaryCategoryRequest	body	CreateAuxiliaryCategoryRequest	true	"Create auxiliary category request"
//	@Success		201
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/auxiliaries [post]
func (h Handler) CreateAuxiliaryCategory(c *gin.Context) {
	var req CreateAuxiliaryCategoryRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.CreateAuxiliaryCategoryCmd{
		SobId:      uuid.MustParse(c.Param("sobId")),
		CategoryId: uuid.New(),
		Key:        req.Key,
		Title:      req.Title,
		IsStandard: false,
	}
	if err := h.app.Commands.CreateAuxiliaryCategory.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusCreated)
}

// CreateAuxiliaryAccount godoc
//
//	@Text			Create auxiliary account
//	@Description	Create auxiliary account
//	@Tags			auxiliary accounts
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId							path	string							true	"Sob ID"
//	@Param			categoryKey						path	string							true	"Category Key"
//	@Param			CreateAuxiliaryAccountRequest	body	CreateAuxiliaryAccountRequest	true	"Create auxiliary account request"
//	@Success		201
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/auxiliary/{categoryKey}/accounts [post]
func (h Handler) CreateAuxiliaryAccount(c *gin.Context) {
	var req CreateAuxiliaryAccountRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.CreateAuxiliaryAccountCmd{
		AccountId:   uuid.New(),
		SobId:       uuid.MustParse(c.Param("sobId")),
		CategoryKey: c.Param("categoryKey"),
		Key:         req.Key,
		Title:       req.Title,
		Description: req.Description,
	}
	if err := h.app.Commands.CreateAuxiliaryAccount.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusCreated)
}
