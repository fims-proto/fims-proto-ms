package command

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"github/fims-proto/fims-proto-ms/internal/account/domain/balance_direction"

	"github/fims-proto/fims-proto-ms/internal/account/app/service"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
)

type accountEntry struct {
	number           string
	level            int
	title            string
	superiorNumber   string
	accountType      string
	balanceDirection string
}

type InitialAccountsHandler struct {
	repo       domain.Repository
	sobService service.SobService
}

func NewInitialAccountHandler(repo domain.Repository, sobService service.SobService) InitialAccountsHandler {
	if repo == nil {
		panic("nil repo")
	}
	return InitialAccountsHandler{
		repo:       repo,
		sobService: sobService,
	}
}

func (h InitialAccountsHandler) Handle(ctx context.Context, sobId uuid.UUID) error {
	// 0. read sob
	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return errors.Wrap(err, "read sob failed")
	}

	// 1. read CSV
	accountEntries, err := h.readFromCSV()
	if err != nil {
		return err
	}

	// 2. prepare accounts
	preparedAccounts, err := h.prepareAccounts(sobId, accountEntries, sob.AccountsCodeLength)
	if err != nil {
		return err
	}

	return h.repo.InitialAccounts(ctx, preparedAccounts)
}

func (h InitialAccountsHandler) readFromCSV() ([]accountEntry, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "could not get working directory")
	}

	csvFile, err := os.Open(filepath.Join(workDir, "dataload", "account", "accounts.csv"))
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}

	csvReader := csv.NewReader(csvFile)

	// skip first line
	_, err = csvReader.Read()
	if err != nil {
		return nil, errors.Wrap(err, "could not read file")
	}

	var entries []accountEntry
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "could not read file")
		}
		level, err := strconv.Atoi(line[2])
		if err != nil {
			return nil, errors.Wrap(err, "convert level to number failed")
		}
		balanceDirection := line[5]
		if balanceDirection == "" {
			balanceDirection = balance_direction.NotDefined.String()
		}
		entries = append(entries, accountEntry{
			number:           line[1],
			level:            level,
			title:            line[3],
			superiorNumber:   line[4],
			accountType:      line[0],
			balanceDirection: balanceDirection,
		})
	}

	return entries, nil
}

func (h InitialAccountsHandler) prepareAccounts(sobId uuid.UUID, accountEntries []accountEntry, codeLengthLimits []int) ([]*account.Account, error) {
	preparedAccounts := make(map[string]*account.Account)
	for i := 0; i < len(codeLengthLimits); i++ {
		for _, entry := range accountEntries {
			if entry.level == i+1 {
				var levelNumber int
				var superiorAccountId uuid.UUID
				var numberHierarchy []int
				if entry.level == 1 {
					superiorAccountId = uuid.Nil
					levelNumber, _ = strconv.Atoi(entry.number)
					numberHierarchy = []int{levelNumber}
				} else {
					levelNumber, _ = strconv.Atoi(strings.TrimPrefix(entry.number, entry.superiorNumber))
					superiorAccount, ok := preparedAccounts[entry.superiorNumber]
					if !ok {
						return nil, errors.Errorf("cannot find prepared superior account %s", entry.superiorNumber)
					}
					superiorAccountId = superiorAccount.Id()
					numberHierarchy = append(superiorAccount.NumberHierarchy(), levelNumber)
				}
				domainAccount, err := account.New(uuid.New(), sobId, superiorAccountId, entry.title, entry.number, numberHierarchy, entry.level, entry.accountType, entry.balanceDirection)
				if err != nil {
					return nil, errors.Wrapf(err, "dataload failed on account %s", entry.number)
				}
				preparedAccounts[entry.number] = domainAccount
			}
		}
	}

	// to slice
	accounts := make([]*account.Account, len(preparedAccounts))
	i := 0
	for _, v := range preparedAccounts {
		accounts[i] = v
		i++
	}
	if len(accounts) != len(accountEntries) {
		return nil, errors.Errorf("prepared accounts size (%d) doesn't equal to CSV entries size (%d)", len(accounts), len(accountEntries))
	}
	return accounts, nil
}
