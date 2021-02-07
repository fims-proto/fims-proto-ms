package counter

import "github.com/pkg/errors"

type Counter struct {
	UUID    		string // this UUID should bind uniquely to the user
	current 		uint
	formater       Formater
}

func NewCounter(UUID string, len uint) (*Counter, error) {
	if UUID == "" {
		return nil, errors.New("empty Numbering service UUID provided")
	}
	return &Counter{
		UUID: UUID,
		current: 0,
		formater: Formater{length: len},
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
