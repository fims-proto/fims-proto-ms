package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

/*
 * currently we only use business object field to match counter.
 * in the future we may introduce another object - counter configuration to match business object to counter object
 */
type Counter struct {
	uuid           uuid.UUID
	businessObject string
	current        uint
	formatter      Formatter
	lastResetDate  time.Time
}

func NewCounter(counterUUID uuid.UUID, businessObject string, prefix string, sufix string) (*Counter, error) {
	if counterUUID == uuid.Nil {
		return nil, errors.New("empty numbering service UUID provided")
	}
	if businessObject == "" {
		return nil, errors.New("empty target business object")
	}

	return &Counter{
		uuid:           counterUUID,
		current:        0,
		businessObject: businessObject,
		formatter:      NewFormatter(prefix, sufix),
		lastResetDate:  time.Now(),
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
	c.lastResetDate = time.Now()
	return nil
}

func (c *Counter) UUID() uuid.UUID {
	return c.uuid
}

func (c Counter) BusinessObject() string {
	return c.businessObject
}
