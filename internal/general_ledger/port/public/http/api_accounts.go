package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
)

// ReadPagingAccounts godoc
// @Text List all accounts
// @Description List all accounts
// @Tags accounts
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $filter query string false "filter on field(s)" example(title eq 'something' and amount lt 10)
// @Success 200 {array} AccountResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/accounts [get]
func (h Handler) ReadPagingAccounts(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.Account], error) {
			return h.app.Queries.PagingAccounts.Handle(c, uuid.MustParse(c.Param("sobId")), pageRequest)
		},
		accountDTOToVO,
	)
}

// ReadAccountById godoc
// @Text Get an account by id
// @Description Get an account by id
// @Tags accounts
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param accountId path string true "Account ID"
// @Success 200 {object} AccountResponse
// @Failure 404
// @Failure 500 {object} Error
// @Router /sob/{sobId}/account/{accountId} [get]
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
	c.JSON(http.StatusOK, accountDTOToVO(v))
}

// AssignAuxiliaryCategoriesToAccount godoc
// @Text Assign auxiliary categories to account
// @Description Assign auxiliary categories to account
// @Tags accounts
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param accountId path string true "Account ID"
// @Param AssignAuxiliaryCategoriesToAccountRequest body AssignAuxiliaryCategoriesToAccountRequest true "Account id and category ids"
// @Success 204
// @Failure 500 {object} Error
// @Router /sob/{sobId}/account/{accountId}/assign-auxiliaries [post]
func (h Handler) AssignAuxiliaryCategoriesToAccount(c *gin.Context) {
	var req AssignAuxiliaryCategoriesToAccountRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.AssignAuxiliaryCategoryCmd{
		AccountId:    uuid.MustParse(c.Param("accountId")),
		CategoryKeys: req.CategoryKeys,
	}
	if err := h.app.Commands.AssignAuxiliaryCategory.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
