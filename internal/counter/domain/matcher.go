package domain

import (
	"strings"

	"github.com/pkg/errors"
)

type Matcher struct {
	busiObjects []string
	sep         string
}

func NewMatcher(sep string, objs ...string) (*Matcher, error) {
	if sep == "" {
		sep = "-"
	}
	if len(objs) == 0 {
		return nil, errors.New("empty business objects")
	}

	for i, v := range objs {
		if v == "" {
			return nil, errors.Errorf("empty business object at index %d", i)
		}
	}

	return &Matcher{
		busiObjects: objs,
		sep:         sep,
	}, nil
}

func (m Matcher) String() string {
	return strings.Join(m.busiObjects, m.sep)
}
