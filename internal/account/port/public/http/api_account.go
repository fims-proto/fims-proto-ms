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

// ReadAllAccountConfigurations godoc
// @Text List all account configurations
// @Description List all account configurations
// @Tags accounts
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $choose query string false "choose only field(s)"
// @Param $filter query string false "filter on field(s)" example(title eq 'something' and amount lt 10)
// @Success 200 {array} AccountConfigurationResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/account-configurations/ [get]
func (h Handler) ReadAllAccountConfigurations(c *gin.Context) {
	pageable, err := data.NewPageableFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	accountConfigurationsPage, err := h.app.Queries.PagingAccountConfigurations.Handle(c, uuid.MustParse(c.Param("sobId")), pageable)
	if err != nil {
		_ = c.Error(err)
		return
	}
	accountConfigurations := make([]AccountConfigurationResponse, len(accountConfigurationsPage.Content()))
	for i, account := range accountConfigurationsPage.Content() {
		accountConfigurations[i] = accountConfigurationDTOToVO(account)
	}
	resp, _ := data.NewPage(accountConfigurations, pageable, accountConfigurationsPage.NumberOfElements())
	c.JSON(http.StatusOK, resp)
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	r.GET("/sob/:sobId/account-configurations/", h.ReadAllAccountConfigurations)
}
