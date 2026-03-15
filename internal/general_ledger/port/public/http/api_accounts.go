package http

import (
	"net/http"
	"strconv"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/class"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReadAllAccounts godoc
//
//	@Text			List all accounts
//	@Description	List all accounts
//	@Tags			accounts
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Success		200		{array}		AccountSlimResponse
//	@Failure		500		{object}	Error
//	@Router			/sob/{sobId}/accounts [get]
func (h Handler) ReadAllAccounts(c *gin.Context) {
	accounts, err := h.app.Queries.AllAccounts.Handle(c, uuid.MustParse(c.Param("sobId")))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, converter.DTOsToVOs(accounts, accountDTOToSlimVO))
}

// ReadAccountClasses godoc
//
//	@Text			List allowed account classes and their allowed groups
//	@Description	List allowed account classes and their allowed groups
//	@Tags			accounts
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path	string	true	"Sob ID"
//	@Success		200		{array}	AccountClass
//	@Router			/sob/{sobId}/account-classes [get]
func (h Handler) ReadAccountClasses(c *gin.Context) {
	var resp []AccountClass
	for _, c := range class.Classes {
		var groups []string
		for _, g := range c.Groups {
			groups = append(groups, strconv.Itoa(int(g)))
		}
		resp = append(resp, AccountClass{
			Class:  strconv.Itoa(int(c.Class)),
			Groups: groups,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// ReadAccountById godoc
//
//	@Text			Get an account by id
//	@Description	Get an account by id
//	@Tags			accounts
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId		path		string	true	"Sob ID"
//	@Param			accountId	path		string	true	"Account ID"
//	@Success		200			{object}	AccountDetailResponse
//	@Failure		404
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/account/{accountId} [get]
func (h Handler) ReadAccountById(c *gin.Context) {
	v, err := h.app.Queries.AccountById.Handle(c, uuid.MustParse(c.Param("accountId")))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if v.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, accountDTOToDetailVO(v))
}

// CreateAccount godoc
//
//	@Text			Create account
//	@Description	Create account
//	@Tags			accounts
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId					path		string					true	"Sob ID"
//	@Param			CreateAccountRequest	body		CreateAccountRequest	true	"Create account request"
//	@Success		201						{object}	AccountDetailResponse
//	@Failure		500						{object}	Error
//	@Router			/sob/{sobId}/accounts [post]
func (h Handler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	classReq, err := strconv.Atoi(req.Class)
	if err != nil {
		_ = c.Error(err)
		return
	}
	group, err := strconv.Atoi(req.Group)
	if err != nil {
		_ = c.Error(err)
		return
	}
	cmd := command.CreateAccountCmd{
		AccountId:             uuid.New(),
		SobId:                 uuid.MustParse(c.Param("sobId")),
		Title:                 req.Title,
		LevelNumber:           req.LevelNumber,
		BalanceDirection:      req.BalanceDirection,
		Class:                 classReq,
		Group:                 group,
		SuperiorAccountNumber: req.SuperiorAccountNumber,
		DimensionCategoryIds:  req.DimensionCategoryIds,
	}

	if err := h.app.Commands.CreateAccount.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	createdAccount, err := h.app.Queries.AccountById.Handle(c, cmd.AccountId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, accountDTOToDetailVO(createdAccount))
}

// UpdateAccount godoc
//
//	@Text			Update account
//	@Description	Update account
//	@Tags			accounts
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId					path	string					true	"Sob ID"
//	@Param			accountId				path	string					true	"Account ID"
//	@Param			UpdateAccountRequest	body	UpdateAccountRequest	true	"Update account request"
//	@Success		204
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/account/{accountId} [patch]
func (h Handler) UpdateAccount(c *gin.Context) {
	var req UpdateAccountRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	group, err := strconv.Atoi(req.Group)
	if err != nil {
		_ = c.Error(err)
		return
	}
	cmd := command.UpdateAccountCmd{
		AccountId:            uuid.MustParse(c.Param("accountId")),
		SobId:                uuid.MustParse(c.Param("sobId")),
		Title:                req.Title,
		LevelNumber:          req.LevelNumber,
		BalanceDirection:     req.BalanceDirection,
		Group:                group,
		DimensionCategoryIds: req.DimensionCategoryIds,
	}
	if err = h.app.Commands.UpdateAccount.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
