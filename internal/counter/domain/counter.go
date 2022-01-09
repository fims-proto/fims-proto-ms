package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Counter struct {
	id             uuid.UUID
	businessObject string
	current        uint
	formatter      Formatter
	lastResetAt    time.Time
}

func NewCounter(counterUUID uuid.UUID, current uint, prefix, suffix string, lastResetAt time.Time, matcherSep string, objs ...string) (*Counter, error) {
	if counterUUID == uuid.Nil {
		return nil, errors.New("nil uuid")
	}

	m, err := NewMatcher(matcherSep, objs...)
	if err != nil {
		return nil, errors.Wrap(err, "counter business match create failed")
	}

	return &Counter{
		id:             counterUUID,
		current:        current,
		businessObject: m.String(),
		formatter:      NewFormatter(prefix, suffix),
		lastResetAt:    lastResetAt,
	}, nil
}

func NewCounterFromDB(counterUUID uuid.UUID, current uint, businessObject, prefix, suffix string, lastResetAt time.Time) (*Counter, error) {
	if counterUUID == uuid.Nil {
		return nil, errors.New("nil uuid")
	}

	if businessObject == "" {
		return nil, errors.New("empty business object")
	}

	return &Counter{
		id:             counterUUID,
		current:        current,
		businessObject: businessObject,
		formatter:      NewFormatter(prefix, suffix),
		lastResetAt:    lastResetAt,
	}, nil
}

func (c *Counter) Next() {
	c.current++
}

func (c *Counter) Identifier() string {
	return c.formatter.format(c.current)
}

func (c *Counter) Reset() error {
	c.current = 0
	c.lastResetAt = time.Now()
	return nil
}

func (c Counter) Id() uuid.UUID {
	return c.id
}

func (c Counter) CurrentIndex() uint {
	return c.current
}

func (c Counter) Prefix() string {
	return c.formatter.prefix
}

func (c Counter) Suffix() string {
	return c.formatter.suffix
}

func (c Counter) BusinessObject() string {
	return c.businessObject
}

func (c Counter) LastResetAt() time.Time {
	return c.lastResetAt
}
