package counter

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Counter struct {
	uuid          uuid.UUID // this UUID should bind uniquely to the user
	Current       uint
	formatter     Formatter
	LastResetDate time.Time
}

func NewCounter(counterUUID uuid.UUID, prefix string, sufix string) (*Counter, error) {
	if counterUUID == uuid.Nil {
		return nil, errors.New("empty Numbering service UUID provided")
	}

	return &Counter{
		uuid:          counterUUID,
		Current:       0,
		formatter:     NewFormatter(prefix, sufix),
		LastResetDate: time.Now(),
	}, nil
}

func (c *Counter) Next() {
	c.Current++
}

func (c *Counter) Identifier() string {
	return c.formatter.format(c.Current)
}

func (c *Counter) Reset() error {
	c.Current = 0
	c.LastResetDate = time.Now()
	return nil
}

func (c *Counter) UUID() uuid.UUID {
	return c.uuid
}
