package data

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewPageable(c *gin.Context) (Pageable, error) {
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

	return NewPageRequest(int(page), int(size), sorts, chooses, filters)
}
