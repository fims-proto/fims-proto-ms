package domain

import (
	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"strconv"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Account struct {
	id                uuid.UUID
	sobId             uuid.UUID
	superiorAccountId uuid.UUID
	superiorNumbers   []int
	title             string
	levelNumber       int
	level             int
	accountType       commonAccount.Type
	balanceDirection  commonAccount.Direction
}

func NewAccount(id, sobId, superiorAccountId uuid.UUID, superiorNumbers []int, title string, levelNumber, level int, accountType, balanceDirection string, levelCodeLength int) (*Account, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil uuid")
	}
	if sobId == uuid.Nil {
		return nil, errors.New("nil sob")
	}
	if superiorAccountId == uuid.Nil && level > 1 {
		return nil, errors.New("nil superior account")
	}
	if len(superiorNumbers) != level-1 {
		return nil, errors.New("account level doesn't match superior numbers")
	}
	if title == "" {
		return nil, errors.New("empty account title")
	}
	if level < 1 {
		return nil, errors.New("invalid account level")
	}
	if levelNumber == 0 {
		return nil, errors.New("invalid account levelNumber")
	}
	if len(strconv.Itoa(levelNumber)) > levelCodeLength {
		return nil, errors.New("account levelNumber exceeds max length")
	}
	at, err := commonAccount.NewAccountType(accountType)
	if err != nil {
		return nil, errors.Wrap(err, "invalid account type")
	}
	bd, err := commonAccount.NewDirection(balanceDirection)
	if err != nil {
		return nil, errors.Wrap(err, "invalid balance direction")
	}

	return &Account{
		id:                id,
		sobId:             sobId,
		superiorAccountId: superiorAccountId,
		superiorNumbers:   superiorNumbers,
		title:             title,
		levelNumber:       levelNumber,
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

func (acc Account) SuperiorNumbers() []int {
	return acc.superiorNumbers
}

func (acc Account) Title() string {
	return acc.title
}

func (acc Account) LevelNumber() int {
	return acc.levelNumber
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
