package domain

import (
	"fmt"
	"strconv"
	"strings"

	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Account struct {
	id                uuid.UUID
	sobId             uuid.UUID
	title             string
	accountNumber     string
	numberHierarchy   []int
	superiorAccountId uuid.UUID
	level             int
	accountType       commonAccount.Type
	balanceDirection  commonAccount.Direction
}

func NewAccount(id, sobId, superiorAccountId uuid.UUID, title, accountNumber, accountType, balanceDirection string, level int, numberHierarchy, codeLength []int) (*Account, error) {
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
	if level != len(numberHierarchy) {
		return nil, errors.Errorf("level %d not match to number hierarchy %v", level, numberHierarchy)
	}
	concatenatedNumber, err := concatenateAccountNumber(numberHierarchy, codeLength)
	if err != nil {
		return nil, err
	}
	if concatenatedNumber != accountNumber {
		return nil, errors.Errorf("account number %s not match to number hierarchy %v", accountNumber, numberHierarchy)
	}

	return &Account{
		id:                id,
		sobId:             sobId,
		title:             title,
		accountNumber:     accountNumber,
		numberHierarchy:   numberHierarchy,
		superiorAccountId: superiorAccountId,
		level:             level,
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

func (acc Account) AccountNumber() string {
	return acc.accountNumber
}

func (acc Account) NumberHierarchy() []int {
	return acc.numberHierarchy
}

func (acc Account) Title() string {
	return acc.title
}

func (acc Account) Level() int {
	return acc.level
}

func (acc Account) Type() commonAccount.Type {
	return acc.accountType
}

func (acc Account) BalanceDirection() commonAccount.Direction {
	return acc.balanceDirection
}

func concatenateAccountNumber(numberHierarchy, codeLengths []int) (string, error) {
	if len(numberHierarchy) > len(codeLengths) {
		return "", errors.Errorf("account depth %d exceeds max depth %d", len(numberHierarchy), len(codeLengths))
	}

	for i := 0; i < len(numberHierarchy); i++ {
		if numberHierarchy[i] < 1 {
			return "", errors.Errorf("account number %d at level %d cannot be smaller than 1", numberHierarchy[i], i)
		}
		if len(strconv.Itoa(numberHierarchy[i])) > codeLengths[i] {
			return "", errors.Errorf("account number %d at level %d exceeds max length (%d)", numberHierarchy[i], i, codeLengths[i])
		}
	}

	var builder strings.Builder
	for i, number := range numberHierarchy {
		builder.WriteString(fmt.Sprintf("%0*d", codeLengths[i], number))
	}

	return builder.String(), nil
}
