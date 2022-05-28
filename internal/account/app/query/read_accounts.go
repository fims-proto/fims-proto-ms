package query

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AccountsReadModel interface {
	ReadAllAccounts(ctx context.Context, sobId uuid.UUID) ([]Account, error)
	ReadById(ctx context.Context, accountId uuid.UUID) (Account, error)
	ReadByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]*Account, error)
	ReadByAccountNumber(ctx context.Context, sobId uuid.UUID, levelNumber int, superiorNumbers []int) (Account, error)
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

func (h ReadAccountsHandler) HandleReadAll(ctx context.Context, sobId uuid.UUID) ([]Account, error) {
	accounts, err := h.readModel.ReadAllAccounts(ctx, sobId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read all accounts")
	}

	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read sob")
	}
	for i := range accounts {
		accountNumber, err := concatenateAccountNumber(accounts[i].LevelNumber, accounts[i].SuperiorNumbers, sob.AccountsCodeLength)
		if err != nil {
			return nil, errors.Wrap(err, "failed on concatenate account number")
		}
		accounts[i].AccountNumber = accountNumber
	}
	return accounts, nil
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
	accountNumber, err := concatenateAccountNumber(account.LevelNumber, account.SuperiorNumbers, sob.AccountsCodeLength)
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
		accountNumber, err := concatenateAccountNumber(account.LevelNumber, account.SuperiorNumbers, sob.AccountsCodeLength)
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
		levelNumber, superiorNumbers, err := cutAccountNumber(accountNumber, sob.AccountsCodeLength)
		if err != nil {
			return nil, errors.Wrapf(err, "validate existence of account %s failed", accountNumber)
		}
		account, err := h.readModel.ReadByAccountNumber(ctx, sobId, levelNumber, superiorNumbers)
		if err != nil {
			return nil, errors.Wrapf(err, "validate existence of account %s failed", accountNumber)
		}
		account.AccountNumber = accountNumber
		accounts[accountNumber] = account
	}
	return accounts, nil
}

// cutAccountNumber cuts given account number in string format in to levelNumber and superiorNumbers as per Sob setting
func cutAccountNumber(accountNumber string, accountCodeLengths []int) (int, []int, error) {
	if accountNumber == "" {
		return 0, nil, errors.New("empty account number")
	}
	if len(accountCodeLengths) == 0 {
		return 0, nil, errors.New("empty account code length array")
	}
	if len(accountNumber) < accountCodeLengths[0] {
		return 0, nil, errors.New("account number too short")
	}
	var totalLen int
	var levelNumbers []int
	for i := 0; i < len(accountCodeLengths); i++ {
		totalLen += accountCodeLengths[i]

		if len(accountNumber) < accountCodeLengths[i] {
			return 0, nil, errors.Errorf("invalid account number %s", accountNumber)
		}

		accountNumberString := accountNumber[0:accountCodeLengths[i]]
		levelNumber, err := strconv.Atoi(accountNumberString)
		if levelNumber == 0 || err != nil {
			return 0, nil, errors.Errorf("invalid account number %s", accountNumber)
		}
		levelNumbers = append(levelNumbers, levelNumber)

		accountNumber = strings.TrimPrefix(accountNumber, accountNumberString)
		if len(accountNumber) <= 0 {
			break
		}
	}
	if len(accountNumber) > 0 {
		return 0, nil, errors.New("account number too long")
	}
	if len(levelNumbers) == 0 {
		return 0, nil, errors.New("failed to cut account number array")
	}

	levelNumber := levelNumbers[len(levelNumbers)-1]

	return levelNumber, levelNumbers[:len(levelNumbers)-1], nil
}

// concatenateAccountNumber concatenates levelNumber and superiorNumbers into one accountNumber string as per Sob setting
func concatenateAccountNumber(levelNumber int, superiorNumbers, accountCodeLengths []int) (string, error) {
	if len(accountCodeLengths) < len(superiorNumbers)+1 {
		return "", errors.Errorf("accountCodeLengths too short: %d, len of superiorNumbers: %d", len(accountCodeLengths), len(superiorNumbers))
	}

	var builder strings.Builder
	for i, superiorNumber := range superiorNumbers {
		builder.WriteString(fmt.Sprintf("%0*d", accountCodeLengths[i], superiorNumber))
	}
	builder.WriteString(fmt.Sprintf("%0*d", accountCodeLengths[len(superiorNumbers)], levelNumber))

	return builder.String(), nil
}
