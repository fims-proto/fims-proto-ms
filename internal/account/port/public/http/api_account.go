package http

import (
	"github/fims-proto/fims-proto-ms/internal/account/app"
	"net/http"

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
// @Success 200 {array} AccountResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/accounts/ [get]
func (h Handler) ReadAllAccounts(c *gin.Context) {
	accounts, err := h.app.Queries.ReadAccounts.HandleReadAll(c, uuid.MustParse(c.Param("sobId")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	resp := make([]AccountResponse, len(accounts))
	for i, account := range accounts {
		resp[i] = mapFromAccountQuery(account)
	}
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
