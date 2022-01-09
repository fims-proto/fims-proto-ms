package query

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/app/command"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AccountsReadModel interface {
	ReadAllAccounts(ctx context.Context, sobId uuid.UUID) ([]Account, error)
	ReadById(ctx context.Context, accountId uuid.UUID) (Account, error)
	ReadByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]Account, error)
	ReadByAccountNumber(ctx context.Context, sobId uuid.UUID, levelNumber int, superiorNumbers []int) (Account, error)
}

type ReadAccountsHandler struct {
	readModel  AccountsReadModel
	sobService command.SobService
}

func NewReadAccountsHandler(readModel AccountsReadModel, sobService command.SobService) ReadAccountsHandler {
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
	return h.readModel.ReadAllAccounts(ctx, sobId)
}

func (h ReadAccountsHandler) HandleReadById(ctx context.Context, accountId uuid.UUID) (Account, error) {
	return h.readModel.ReadById(ctx, accountId)
}

func (h ReadAccountsHandler) HandleReadByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]Account, error) {
	return h.readModel.ReadByIds(ctx, accountIds)
}

func (h ReadAccountsHandler) HandleReadByAccountNumber(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]Account, error) {
	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return nil, errors.Wrap(err, "validate existence failed on reading sob")
	}

	accounts := make(map[string]Account)

	for _, accountNumber := range accountNumbers {
		levelNumber, superiorNumbers, err := h.cutAccountNumber(accountNumber, sob.AccountsCodeLength)
		if err != nil {
			return nil, errors.Wrapf(err, "validate existence of account %s failed", accountNumber)
		}
		account, err := h.readModel.ReadByAccountNumber(ctx, sobId, levelNumber, superiorNumbers)
		if err != nil {
			return nil, errors.Wrapf(err, "validate existence of account %s failed", accountNumber)
		}
		accounts[accountNumber] = account
	}
	return accounts, nil
}

func (h ReadAccountsHandler) cutAccountNumber(accountNumber string, accountCodeLengths []int) (int, []int, error) {
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
