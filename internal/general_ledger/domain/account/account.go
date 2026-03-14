package account

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	accountNumber        string
	numberHierarchy      []int
	level                int
	isLeaf               bool
	class                class.Class
	group                class.Group
	balanceDirection     balance_direction.BalanceDirection
	dimensionCategoryIds []uuid.UUID
}

// New takes all fields except:
// - accountNumber: it's calculated from numberHierarchy
// - superiorAccount: this method cannot create an entity with such nested structure
func New(
	id uuid.UUID,
	sobId uuid.UUID,
	superiorAccountId uuid.UUID,
	title string,
	numberHierarchy []int,
	codeLengths []int,
	level int,
	isLeaf bool,
	classId int,
	groupId int,
	balanceDirection string,
	dimensionCategoryIds []uuid.UUID,
) (*Account, error) {
	accountNumber, err := composeAccountNumber(numberHierarchy, codeLengths)
	if err != nil {
		return nil, err
	}

	return NewByAllFields(id, sobId, superiorAccountId, nil, title, accountNumber, numberHierarchy, level, isLeaf, classId, groupId, balanceDirection, dimensionCategoryIds)
}

// NewByAllFields takes all attributes of Account, and doesn't validate accountNumber field
// Typically used in persistence level
func NewByAllFields(
	id uuid.UUID,
	sobId uuid.UUID,
	superiorAccountId uuid.UUID,
	superiorAccount *Account,
	title string,
	accountNumber string,
	numberHierarchy []int,
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

	if superiorAccountId == uuid.Nil && len(numberHierarchy) > 1 {
		return nil, errors.New("nil superior account id")
	}

	if title == "" {
		return nil, errors.New("empty account title")
	}

	if utf8.RuneCountInString(title) > 50 {
		return nil, errors.New("account title exceeds max length (50)")
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

	return &Account{
		id:                   id,
		sobId:                sobId,
		superiorAccountId:    superiorAccountId,
		superiorAccount:      superiorAccount,
		title:                title,
		accountNumber:        accountNumber,
		numberHierarchy:      numberHierarchy,
		level:                level,
		isLeaf:               isLeaf,
		class:                c,
		group:                g,
		balanceDirection:     bd,
		dimensionCategoryIds: dimensionCategoryIds,
	}, nil
}

func composeAccountNumber(numberHierarchy, codeLengths []int) (string, error) {
	if len(numberHierarchy) > len(codeLengths) {
		return "", fmt.Errorf("account number hierarchy %d exceeds max depth %d", len(numberHierarchy), len(codeLengths))
	}

	for i := 0; i < len(numberHierarchy); i++ {
		if numberHierarchy[i] < 1 {
			return "", fmt.Errorf("account number %d at level %d cannot be smaller than 1", numberHierarchy[i], i)
		}
		if len(strconv.Itoa(numberHierarchy[i])) > codeLengths[i] {
			return "", fmt.Errorf("account number %d at level %d exceeds max length (%d)", numberHierarchy[i], i, codeLengths[i])
		}
	}

	var builder strings.Builder
	for i, number := range numberHierarchy {
		if _, err := fmt.Fprintf(&builder, "%0*d", codeLengths[i], number); err != nil {
			return "", fmt.Errorf("failed to compose account number: %w", err)
		}
	}

	return builder.String(), nil
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

func (a *Account) AccountNumber() string {
	return a.accountNumber
}

func (a *Account) NumberHierarchy() []int {
	return a.numberHierarchy
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
