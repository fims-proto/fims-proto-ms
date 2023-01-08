package account

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account_type"
	"github/fims-proto/fims-proto-ms/internal/account/domain/balance_direction"
)

type Account struct {
	id                uuid.UUID
	sobId             uuid.UUID
	superiorAccountId uuid.UUID
	title             string
	accountNumber     string
	numberHierarchy   []int
	level             int
	accountType       account_type.AccountType
	balanceDirection  balance_direction.BalanceDirection
}

func New(id, sobId, superiorAccountId uuid.UUID, title, accountNumber string, numberHierarchy []int, level int, accountType, direction string) (*Account, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil account id")
	}

	if sobId == uuid.Nil {
		return nil, errors.New("nil sob")
	}

	if superiorAccountId == uuid.Nil && len(numberHierarchy) > 1 {
		return nil, errors.New("nil superior account id")
	}

	if title == "" {
		return nil, errors.New("empty account title")
	}

	if accountNumber == "" {
		return nil, errors.New("empty account number")
	}

	if level < 1 {
		return nil, errors.Errorf("level %d must >= 1", level)
	}

	if level != len(numberHierarchy) {
		return nil, errors.Errorf("level %d not match to number hierarchy %v", level, numberHierarchy)
	}

	at, err := account_type.FromString(accountType)
	if err != nil {
		return nil, err
	}

	bd, err := balance_direction.FromString(direction)
	if err != nil {
		return nil, err
	}

	return &Account{
		id:                id,
		sobId:             sobId,
		superiorAccountId: superiorAccountId,
		title:             title,
		accountNumber:     accountNumber,
		numberHierarchy:   numberHierarchy,
		level:             level,
		accountType:       at,
		balanceDirection:  bd,
	}, nil
}

func (ac Account) Id() uuid.UUID {
	return ac.id
}

func (ac Account) SobId() uuid.UUID {
	return ac.sobId
}

func (ac Account) SuperiorAccountId() uuid.UUID {
	return ac.superiorAccountId
}

func (ac Account) AccountNumber() string {
	return ac.accountNumber
}

func (ac Account) NumberHierarchy() []int {
	return ac.numberHierarchy
}

func (ac Account) Title() string {
	return ac.title
}

func (ac Account) Level() int {
	return ac.level
}

func (ac Account) AccountType() account_type.AccountType {
	return ac.accountType
}

func (ac Account) BalanceDirection() balance_direction.BalanceDirection {
	return ac.balanceDirection
}
