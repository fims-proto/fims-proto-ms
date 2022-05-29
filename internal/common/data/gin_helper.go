package data

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewPageable(c *gin.Context) (Pageable, error) {
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse page query parameter")
	}

	size, err := strconv.ParseInt(c.DefaultQuery("size", "40"), 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse page query parameter")
	}

	sort := c.Query("sort")
	sortFields := make(map[string]string)
	if sort != "" {
		sortSegments := strings.Split(sort, ",")
		for _, segment := range sortSegments {
			elements := strings.Split(strings.TrimSpace(segment), " ")
			if len(elements) == 1 {
				sortFields[elements[0]] = "asc"
			} else if len(elements) == 2 {
				sortFields[elements[0]] = elements[1]
			} else {
				return nil, errors.Errorf("invalid sort query parameter %s", sort)
			}
		}
	}

	choose := c.Query("choose")
	var chooseFields []string
	if choose != "" {
		chooseSegments := strings.Split(choose, ",")
		for _, segment := range chooseSegments {
			chooseFields = append(chooseFields, strings.TrimSpace(segment))
		}
	}

	return NewPageRequest(int(page), int(size), sortFields, chooseFields)
}
