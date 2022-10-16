package data

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"
)

func NewPageRequestFromQuery(c *gin.Context) (PageRequest, error) {
	page, err := strconv.ParseInt(c.DefaultQuery("$page", "1"), 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse page query parameter")
	}

	size, err := strconv.ParseInt(c.DefaultQuery("$size", "40"), 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse page query parameter")
	}

	p, err := pageable.NewPageableFromQuery(int(page), int(size))
	if err != nil {
		return nil, err
	}

	sort := c.Query("$sort")
	sorts, err := sortable.NewSortableFromQuery(sort)
	if err != nil {
		return nil, err
	}

	filter := c.Query("$filter")
	filters, err := filterable.NewFilterableFromQuery(filter)
	if err != nil {
		return nil, err
	}

	return NewPageRequest(p, sorts, filters), nil
}

func PagingResponseProcessor[DTO any, VO any](
	c *gin.Context,
	provider func(request PageRequest) (Page[DTO], error),
	converter func(DTO) VO,
) {
	r, err := NewPageRequestFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	dataPage, err := provider(r)
	if err != nil {
		_ = c.Error(err)
		return
	}
	vos := make([]VO, len(dataPage.Content()))
	for i, vo := range dataPage.Content() {
		vos[i] = converter(vo)
	}
	resp, _ := NewPage(vos, r, dataPage.NumberOfElements())
	c.JSON(http.StatusOK, resp)
}
