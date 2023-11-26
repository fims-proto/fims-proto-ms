package account

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/balance_direction"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/class"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_category"
)

type Account struct {
	id                  uuid.UUID
	sobId               uuid.UUID
	superiorAccountId   uuid.UUID
	title               string
	accountNumber       string
	numberHierarchy     []int
	level               int
	class               class.Class
	group               class.Group
	balanceDirection    balance_direction.BalanceDirection
	auxiliaryCategories []*auxiliary_category.AuxiliaryCategory
}

// New takes all fields except accountNumber. It's calculated from numberHierarchy
func New(
	id uuid.UUID,
	sobId uuid.UUID,
	superiorAccountId uuid.UUID,
	title string,
	numberHierarchy []int,
	codeLengths []int,
	level int,
	classId int,
	groupId int,
	balanceDirection string,
	auxiliaryCategories []*auxiliary_category.AuxiliaryCategory,
) (*Account, error) {
	accountNumber, err := composeAccountNumber(numberHierarchy, codeLengths)
	if err != nil {
		return nil, err
	}

	return NewByAllFields(id, sobId, superiorAccountId, title, accountNumber, numberHierarchy, level, classId, groupId, balanceDirection, auxiliaryCategories)
}

// NewByAllFields only difference from New function, is NewByAllFields takes accountNumber, and doesn't validate it.
// Typically used in persistence level
func NewByAllFields(
	id uuid.UUID,
	sobId uuid.UUID,
	superiorAccountId uuid.UUID,
	title string,
	accountNumber string,
	numberHierarchy []int,
	level int,
	classId int,
	groupId int,
	balanceDirection string,
	auxiliaryCategories []*auxiliary_category.AuxiliaryCategory,
) (*Account, error) {
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
		return nil, fmt.Errorf("level %d must >= 1", level)
	}

	if level != len(numberHierarchy) {
		return nil, fmt.Errorf("level %d not match to number hierarchy %v", level, numberHierarchy)
	}

	c := class.Class(classId)
	g := class.Group(groupId)
	if err := class.Validate(c, g); err != nil {
		return nil, err
	}

	bd, err := balance_direction.FromString(balanceDirection)
	if err != nil {
		return nil, err
	}

	for _, category := range auxiliaryCategories {
		if category == nil {
			return nil, errors.New("nil auxiliary category")
		}
	}

	return &Account{
		id:                  id,
		sobId:               sobId,
		superiorAccountId:   superiorAccountId,
		title:               title,
		accountNumber:       accountNumber,
		numberHierarchy:     numberHierarchy,
		level:               level,
		class:               c,
		group:               g,
		balanceDirection:    bd,
		auxiliaryCategories: auxiliaryCategories,
	}, nil
}

func (a *Account) Id() uuid.UUID {
	return a.id
}

func (a *Account) SobId() uuid.UUID {
	return a.sobId
}

func (a *Account) SuperiorAccountId() uuid.UUID {
	return a.superiorAccountId
}

func (a *Account) Title() string {
	return a.title
}

func (a *Account) AccountNumber() string {
	return a.accountNumber
}

func (a *Account) NumberHierarchy() []int {
	return a.numberHierarchy
}

func (a *Account) Level() int {
	return a.level
}

func (a *Account) Class() class.Class {
	return a.class
}

func (a *Account) Group() class.Group {
	return a.group
}

func (a *Account) BalanceDirection() balance_direction.BalanceDirection {
	return a.balanceDirection
}

func (a *Account) AuxiliaryCategories() []*auxiliary_category.AuxiliaryCategory {
	return a.auxiliaryCategories
}
