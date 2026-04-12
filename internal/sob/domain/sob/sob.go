package sob

import (
	"fmt"
	"unicode/utf8"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

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
		return nil, commonErrors.NewInternalError(commonErrors.SlugSobEmptyId)
	}

	if name == "" {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugSobEmptyName)
	}

	if utf8.RuneCountInString(name) > 50 {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugSobNameTooLong)
	}

	if utf8.RuneCountInString(description) > 500 {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugSobDescriptionTooLong)
	}

	if baseCurrency == "" {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugSobEmptyBaseCurrency)
	}

	if startingPeriodYear < 2000 || startingPeriodYear > 3000 {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugSobInvalidStartingYear)
	}

	if startingPeriodMonth < 1 || startingPeriodMonth > 12 {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugSobInvalidStartingMonth)
	}

	if len(accountsCodeLength) < 2 || len(accountsCodeLength) > 10 {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugSobInvalidAccountLevel)
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
