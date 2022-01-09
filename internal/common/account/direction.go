package account

import "github.com/pkg/errors"

type Direction uint

const (
	InvalidDirection   = Direction(iota) // invalid direction
	Debit                                // 借
	Credit                               // 贷
	UndefinedDirection                   // 未指名
)

var availableDirections = map[Direction]string{
	Debit:              "DEBIT",
	Credit:             "CREDIT",
	UndefinedDirection: "NOT_DEFINED",
}

func NewDirection(s string) (Direction, error) {
	for i, v := range availableDirections {
		if v == s {
			return i, nil
		}
	}

	return InvalidDirection, errors.Errorf("invalid direction name: '%s'", s)
}

func (d Direction) String() string {
	for k, v := range availableDirections {
		if k == d {
			return v
		}
	}
	panic("direction string error, should not happen")
}
