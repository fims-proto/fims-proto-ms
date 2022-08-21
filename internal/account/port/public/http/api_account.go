package http

import (
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
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

// ReadPagingAccountConfigurations godoc
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
func (h Handler) ReadPagingAccountConfigurations(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageable data.Pageable) (data.Page[query.AccountConfiguration], error) {
			return h.app.Queries.PagingAccountConfigurations.Handle(c, uuid.MustParse(c.Param("sobId")), pageable)
		},
		accountConfigurationDTOToVO,
	)
}

// ReadPagingAccountsByPeriod godoc
// @Text List accounts in period
// @Description List accounts in period
// @Tags accounts
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param periodId path string true "Period ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $choose query string false "choose only field(s)"
// @Param $filter query string false "filter on field(s)" example(title eq 'something' and amount lt 10)
// @Success 200 {array} AccountResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/period/{periodId}/accounts/ [get]
func (h Handler) ReadPagingAccountsByPeriod(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageable data.Pageable) (data.Page[query.Account], error) {
			return h.app.Queries.PagingAccountsByPeriod.Handle(c, uuid.MustParse(c.Param("sobId")), uuid.MustParse(c.Param("periodId")), pageable)
		},
		accountDTOToVO,
	)
}

// ReadPagingPeriods godoc
// @Text List periods
// @Description List periods
// @Tags accounts
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $choose query string false "choose only field(s)"
// @Param $filter query string false "filter on field(s)" example(title eq 'something' and amount lt 10)
// @Success 200 {array} PeriodResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/periods/ [get]
func (h Handler) ReadPagingPeriods(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageable data.Pageable) (data.Page[query.Period], error) {
			return h.app.Queries.PagingPeriods.Handle(c, uuid.MustParse(c.Param("sobId")), pageable)
		},
		periodDTOToVO,
	)
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	r.GET("/sob/:sobId/account-configurations/", h.ReadPagingAccountConfigurations)
	r.GET("/sob/:sobId/periods/", h.ReadPagingPeriods)
	r.GET("/sob/:sobId/period/:periodId/accounts/", h.ReadPagingAccountsByPeriod)
}
