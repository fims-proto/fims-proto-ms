package counter

import (
	"time"

	"github.com/pkg/errors"
)

type Counter struct {
	UUID    		string // this UUID should bind uniquely to the user
	current 		uint
	formater       Formater
	LastResetDate   time.Time
}

func NewCounter(UUID string, len uint,prefix string, sufix string) (*Counter, error) {
	if UUID == "" {
		return nil, errors.New("empty Numbering service UUID provided")
	}
	return &Counter{
		UUID: UUID,
		current: 0,
		formater: Formater{
			length: len,
			prefix: prefix,
			sufix: sufix,
		},
		LastResetDate: time.Now(),
	},nil
}

func (c *Counter) Next() (string, error) {
	c.current += 1
	s, err:= c.formater.format(c.current)
	if err != nil {
		return "", err
	}
	return s, nil

}

func (c *Counter) Reset() error {
	c.current = 0;
	c.LastResetDate = time.Now()
	return nil
}
