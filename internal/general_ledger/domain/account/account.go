package account

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/balance_direction"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/class"

	"github.com/google/uuid"
)

type Account struct {
	id                   uuid.UUID
	sobId                uuid.UUID
	superiorAccountId    uuid.UUID
	superiorAccount      *Account
	title                string
	rawAccountNumber     string
	level                int
	isLeaf               bool
	class                class.Class
	group                class.Group
	balanceDirection     balance_direction.BalanceDirection
	dimensionCategoryIds []uuid.UUID
}

// New takes all fields except:
// - rawAccountNumber: it's derived from numberHierarchy via PadRawAccountNumber
// - superiorAccount: this method cannot create an entity with such nested structure
func New(
	id uuid.UUID,
	sobId uuid.UUID,
	superiorAccountId uuid.UUID,
	title string,
	numberHierarchy []int,
	level int,
	isLeaf bool,
	classId int,
	groupId int,
	balanceDirection string,
	dimensionCategoryIds []uuid.UUID,
) (*Account, error) {
	rawAccountNumber, err := PadRawAccountNumber(numberHierarchy)
	if err != nil {
		return nil, err
	}

	return NewByAllFields(
		id,
		sobId,
		superiorAccountId,
		nil,
		title,
		rawAccountNumber,
		level,
		isLeaf,
		classId,
		groupId,
		balanceDirection,
		dimensionCategoryIds,
	)
}

// NewByAllFields takes all attributes of Account, and doesn't validate rawAccountNumber field
// Typically used in persistence level
func NewByAllFields(
	id uuid.UUID,
	sobId uuid.UUID,
	superiorAccountId uuid.UUID,
	superiorAccount *Account,
	title string,
	rawAccountNumber string,
	level int,
	isLeaf bool,
	classId int,
	groupId int,
	balanceDirection string,
	dimensionCategoryIds []uuid.UUID,
) (*Account, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil account id")
	}

	if sobId == uuid.Nil {
		return nil, errors.New("nil sob")
	}

	if superiorAccountId == uuid.Nil && level > 1 {
		return nil, errors.New("nil superior account id")
	}

	if title == "" {
		return nil, errors.New("empty account title")
	}

	if utf8.RuneCountInString(title) > 50 {
		return nil, errors.New("account title exceeds max length (50)")
	}

	if rawAccountNumber == "" {
		return nil, errors.New("empty raw account number")
	}

	// Validate rawAccountNumber format
	if _, err := HierarchyFromRaw(rawAccountNumber); err != nil {
		return nil, fmt.Errorf("invalid raw account number: %w", err)
	}

	// Verify level matches raw account number
	derivedLevel := LevelFromRaw(rawAccountNumber)
	if level != derivedLevel {
		return nil, fmt.Errorf("level %d does not match raw account number level %d", level, derivedLevel)
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

	return &Account{
		id:                   id,
		sobId:                sobId,
		superiorAccountId:    superiorAccountId,
		superiorAccount:      superiorAccount,
		title:                title,
		rawAccountNumber:     rawAccountNumber,
		level:                level,
		isLeaf:               isLeaf,
		class:                c,
		group:                g,
		balanceDirection:     bd,
		dimensionCategoryIds: dimensionCategoryIds,
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

func (a *Account) SuperiorAccount() *Account {
	return a.superiorAccount
}

func (a *Account) Title() string {
	return a.title
}

func (a *Account) RawAccountNumber() string {
	return a.rawAccountNumber
}

func (a *Account) Level() int {
	return a.level
}

func (a *Account) IsLeaf() bool {
	return a.isLeaf
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

func (a *Account) DimensionCategoryIds() []uuid.UUID {
	return a.dimensionCategoryIds
}
