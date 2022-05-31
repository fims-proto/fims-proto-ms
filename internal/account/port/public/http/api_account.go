package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github/fims-proto/fims-proto-ms/internal/account/app"

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

// ReadAllAccounts godoc
// @Summary List all accounts
// @Description List all accounts
// @Tags accounts
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $choose query string false "choose only field(s)" example(accountNumber,title)
// @Param $filter query string false "filter on field(s)" example(title eq 'some thing' and amount lt 10)
// @Success 200 {array} AccountResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/accounts/ [get]
func (h Handler) ReadAllAccounts(c *gin.Context) {
	pageable, err := data.NewPageable(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	accountsPage, err := h.app.Queries.ReadAccounts.HandleReadAll(c, uuid.MustParse(c.Param("sobId")), pageable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	accountResponses := make([]AccountResponse, len(accountsPage.Content))
	for i, account := range accountsPage.Content {
		accountResponses[i] = mapFromAccountQuery(account)
	}
	resp, _ := data.NewPage(accountResponses, accountsPage.Page, accountsPage.Size, accountsPage.NumberOfElements)
	c.JSON(http.StatusOK, resp)
}

func wrapErr(e error) Error {
	var slug string
	se, ok := e.(slugErr)
	if ok {
		slug = se.Slug()
	} else {
		slug = "unknown-error"
	}
	return Error{
		Slug:    slug,
		Message: e.Error(),
	}
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	r.GET("/sob/:sobId/accounts/", h.ReadAllAccounts)
}
