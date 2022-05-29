package query

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AccountsReadModel interface {
	ReadAllAccounts(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Account], error)
	ReadById(ctx context.Context, accountId uuid.UUID) (Account, error)
	ReadByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]*Account, error)
	ReadByAccountNumber(ctx context.Context, sobId uuid.UUID, numberHierarchy []int) (Account, error)
}

type ReadAccountsHandler struct {
	readModel  AccountsReadModel
	sobService SobService
}

func NewReadAccountsHandler(readModel AccountsReadModel, sobService SobService) ReadAccountsHandler {
	if readModel == nil {
		panic("nil read model")
	}
	if sobService == nil {
		panic("nil sob service")
	}
	return ReadAccountsHandler{
		readModel:  readModel,
		sobService: sobService,
	}
}

func (h ReadAccountsHandler) HandleReadAll(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Account], error) {
	accountsPage, err := h.readModel.ReadAllAccounts(ctx, sobId, pageable)
	if err != nil {
		return data.Page[Account]{}, errors.Wrap(err, "failed to read all accounts")
	}

	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return data.Page[Account]{}, errors.Wrap(err, "failed to read sob")
	}
	for i := range accountsPage.Content {
		accountNumber, err := concatenateAccountNumber(accountsPage.Content[i].NumberHierarchy, sob.AccountsCodeLength)
		if err != nil {
			return data.Page[Account]{}, errors.Wrap(err, "failed on concatenate account number")
		}
		accountsPage.Content[i].AccountNumber = accountNumber
	}
	return accountsPage, nil
}

func (h ReadAccountsHandler) HandleReadById(ctx context.Context, accountId uuid.UUID) (Account, error) {
	account, err := h.readModel.ReadById(ctx, accountId)
	if err != nil {
		return Account{}, errors.Wrap(err, "failed to read account")
	}

	sob, err := h.sobService.ReadById(ctx, account.SobId)
	if err != nil {
		return Account{}, errors.Wrap(err, "failed to read sob")
	}
	accountNumber, err := concatenateAccountNumber(account.NumberHierarchy, sob.AccountsCodeLength)
	if err != nil {
		return Account{}, errors.Wrap(err, "failed on concatenate account number")
	}
	account.AccountNumber = accountNumber
	return account, nil
}

func (h ReadAccountsHandler) HandleReadByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]Account, error) {
	result := make(map[uuid.UUID]Account)
	if len(accountIds) == 0 {
		return result, nil
	}

	accounts, err := h.readModel.ReadByIds(ctx, accountIds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read accounts by ids")
	}

	// assume all accounts are from same sob
	sobId := accounts[accountIds[0]].SobId
	for _, account := range accounts {
		if sobId.String() != account.SobId.String() {
			return nil, errors.New("only support reading accounts in same sob for now")
		}
		sobId = account.SobId
	}

	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read sob")
	}

	for key, account := range accounts {
		accountNumber, err := concatenateAccountNumber(account.NumberHierarchy, sob.AccountsCodeLength)
		if err != nil {
			return nil, errors.Wrap(err, "failed on concatenate account number")
		}
		account.AccountNumber = accountNumber
		result[key] = *account
	}
	return result, nil
}

func (h ReadAccountsHandler) HandleReadByAccountNumber(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]Account, error) {
	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return nil, errors.Wrap(err, "validate existence failed to read sob")
	}

	accounts := make(map[string]Account)

	for _, accountNumber := range accountNumbers {
		numberHierarchy, err := cutAccountNumber(accountNumber, sob.AccountsCodeLength)
		if err != nil {
			return nil, errors.Wrapf(err, "validate existence of account %s failed", accountNumber)
		}
		account, err := h.readModel.ReadByAccountNumber(ctx, sobId, numberHierarchy)
		if err != nil {
			return nil, errors.Wrapf(err, "validate existence of account %s failed", accountNumber)
		}
		account.AccountNumber = accountNumber
		accounts[accountNumber] = account
	}
	return accounts, nil
}

// cutAccountNumber cuts given account number in string format in to levelNumber and numberHierarchy as per Sob setting
func cutAccountNumber(accountNumber string, accountCodeLengths []int) ([]int, error) {
	if accountNumber == "" {
		return nil, errors.New("empty account number")
	}
	if len(accountCodeLengths) == 0 {
		return nil, errors.New("empty account code length array")
	}
	if len(accountNumber) < accountCodeLengths[0] {
		return nil, errors.New("account number too short")
	}
	var totalLen int
	var numberHierarchy []int
	for i := 0; i < len(accountCodeLengths); i++ {
		totalLen += accountCodeLengths[i]

		if len(accountNumber) < accountCodeLengths[i] {
			return nil, errors.Errorf("invalid account number %s", accountNumber)
		}

		accountNumberString := accountNumber[0:accountCodeLengths[i]]
		levelNumber, err := strconv.Atoi(accountNumberString)
		if levelNumber == 0 || err != nil {
			return nil, errors.Errorf("invalid account number %s", accountNumber)
		}
		numberHierarchy = append(numberHierarchy, levelNumber)

		accountNumber = strings.TrimPrefix(accountNumber, accountNumberString)
		if len(accountNumber) <= 0 {
			break
		}
	}
	if len(accountNumber) > 0 {
		return nil, errors.New("account number too long")
	}
	if len(numberHierarchy) == 0 {
		return nil, errors.New("failed to cut account number array")
	}

	return numberHierarchy, nil
}

// concatenateAccountNumber concatenates levelNumber and numberHierarchy into one accountNumber string as per Sob setting
func concatenateAccountNumber(numberHierarchy, accountCodeLengths []int) (string, error) {
	if len(numberHierarchy) > len(accountCodeLengths) {
		return "", errors.Errorf("account depth %d exceeds max depth %d", len(numberHierarchy), len(accountCodeLengths))
	}

	var builder strings.Builder
	for i, number := range numberHierarchy {
		builder.WriteString(fmt.Sprintf("%0*d", accountCodeLengths[i], number))
	}

	return builder.String(), nil
}
