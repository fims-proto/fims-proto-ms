package data

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewPageableFromRequest(c *gin.Context) (Pageable, error) {
	page, err := strconv.ParseInt(c.DefaultQuery("$page", "1"), 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse page query parameter")
	}

	size, err := strconv.ParseInt(c.DefaultQuery("$size", "40"), 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse page query parameter")
	}

	sort := c.Query("$sort")
	sorts, err := newSortsFromQuery(sort)
	if err != nil {
		return nil, err
	}

	choose := c.Query("$choose")
	chooses, err := newChoosesFromQuery(choose)
	if err != nil {
		return nil, err
	}

	filter := c.Query("$filter")
	filters, err := newFiltersFromQuery(filter)
	if err != nil {
		return nil, err
	}

	return newPageRequest(int(page), int(size), sorts, chooses, filters)
}

func PagingResponseProcessor[DTO any, VO any](
	c *gin.Context,
	provider func(Pageable) (Page[DTO], error),
	converter func(DTO) VO,
) {
	pageable, err := NewPageableFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	dataPage, err := provider(pageable)
	if err != nil {
		_ = c.Error(err)
		return
	}
	vos := make([]VO, len(dataPage.Content()))
	for i, vo := range dataPage.Content() {
		vos[i] = converter(vo)
	}
	resp, _ := NewPage(vos, pageable, dataPage.NumberOfElements())
	c.JSON(http.StatusOK, resp)
}
