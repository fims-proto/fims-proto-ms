package domain

import (
	"strconv"

	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Account struct {
	id                uuid.UUID
	sobId             uuid.UUID
	superiorAccountId uuid.UUID
	numberHierarchy   []int
	title             string
	accountType       commonAccount.Type
	balanceDirection  commonAccount.Direction
}

func NewAccount(id, sobId, superiorAccountId uuid.UUID, numberHierarchy []int, title string, accountType, balanceDirection string, codeLength []int) (*Account, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil uuid")
	}
	if sobId == uuid.Nil {
		return nil, errors.New("nil sob")
	}
	if superiorAccountId == uuid.Nil && len(numberHierarchy) > 1 {
		return nil, errors.New("nil superior account")
	}
	if title == "" {
		return nil, errors.New("empty account title")
	}
	at, err := commonAccount.NewAccountType(accountType)
	if err != nil {
		return nil, errors.Wrap(err, "invalid account type")
	}
	bd, err := commonAccount.NewDirection(balanceDirection)
	if err != nil {
		return nil, errors.Wrap(err, "invalid balance direction")
	}
	if len(numberHierarchy) > len(codeLength) {
		return nil, errors.Errorf("account depth %d exceeds max depth %d", len(numberHierarchy), len(codeLength))
	}
	for i := 0; i < len(numberHierarchy); i++ {
		if numberHierarchy[i] < 1 {
			return nil, errors.Errorf("account number %d at level %d cannot be smaller than 1", numberHierarchy[i], i)
		}
		if len(strconv.Itoa(numberHierarchy[i])) > codeLength[i] {
			return nil, errors.Errorf("account number %d at level %d exceeds max length (%d)", numberHierarchy[i], i, codeLength[i])
		}
	}

	return &Account{
		id:                id,
		sobId:             sobId,
		superiorAccountId: superiorAccountId,
		numberHierarchy:   numberHierarchy,
		title:             title,
		accountType:       at,
		balanceDirection:  bd,
	}, nil
}

func (acc Account) Id() uuid.UUID {
	return acc.id
}

func (acc Account) SobId() uuid.UUID {
	return acc.sobId
}

func (acc Account) SuperiorAccountId() uuid.UUID {
	return acc.superiorAccountId
}

func (acc Account) NumberHierarchy() []int {
	return acc.numberHierarchy
}

func (acc Account) Title() string {
	return acc.title
}

func (acc Account) Type() commonAccount.Type {
	return acc.accountType
}

func (acc Account) BalanceDirection() commonAccount.Direction {
	return acc.balanceDirection
}
