package sob

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Sob struct {
	id                  uuid.UUID
	name                string
	description         string
	baseCurrency        string
	startingPeriodYear  int
	startingPeriodMonth int
	accountsCodeLength  []int
}

func New(id uuid.UUID, name, description, baseCurrency string, startingPeriodYear, startingPeriodMonth int, accountsCodeLength []int) (*Sob, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil sob id")
	}

	if name == "" {
		return nil, errors.New("empty sob name")
	}

	if utf8.RuneCountInString(name) > 50 {
		return nil, errors.New("sob name too long")
	}

	if utf8.RuneCountInString(description) > 500 {
		return nil, errors.New("sob description too long")
	}

	if baseCurrency == "" {
		return nil, errors.New("empty base currency")
	}

	if startingPeriodYear < 2000 || startingPeriodYear > 3000 {
		return nil, errors.New("invalid starting period year")
	}

	if startingPeriodMonth < 1 || startingPeriodMonth > 12 {
		return nil, errors.New("invalid starting period month")
	}

	if len(accountsCodeLength) < 2 || len(accountsCodeLength) > 10 {
		return nil, errors.New("invalid account level")
	}

	for i, accountCodeLength := range accountsCodeLength {
		if accountCodeLength < 1 || accountCodeLength > 6 {
			return nil, fmt.Errorf("invalid account code length at level %d", i)
		}
	}

	return &Sob{
		id:                  id,
		name:                name,
		description:         description,
		baseCurrency:        baseCurrency,
		startingPeriodYear:  startingPeriodYear,
		startingPeriodMonth: startingPeriodMonth,
		accountsCodeLength:  accountsCodeLength,
	}, nil
}

func (s *Sob) Id() uuid.UUID {
	return s.id
}

func (s *Sob) Name() string {
	return s.name
}

func (s *Sob) Description() string {
	return s.description
}

func (s *Sob) BaseCurrency() string {
	return s.baseCurrency
}

func (s *Sob) StartingPeriodYear() int {
	return s.startingPeriodYear
}

func (s *Sob) StartingPeriodMonth() int {
	return s.startingPeriodMonth
}

func (s *Sob) AccountsCodeLength() []int {
	return s.accountsCodeLength
}
