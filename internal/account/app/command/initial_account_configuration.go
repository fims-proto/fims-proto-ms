package command

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github/fims-proto/fims-proto-ms/internal/account/domain/account_configuration"
	"github/fims-proto/fims-proto-ms/internal/account/domain/balance_direction"

	"github/fims-proto/fims-proto-ms/internal/account/app/service"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
)

type accountConfigurationEntry struct {
	number           string
	level            int
	title            string
	superiorNumber   string
	accountType      string
	balanceDirection string
}

type InitialAccountConfigurationHandler struct {
	repo       domain.Repository
	sobService service.SobService
}

func NewInitialAccountConfigurationHandler(repo domain.Repository, sobService service.SobService) InitialAccountConfigurationHandler {
	if repo == nil {
		panic("nil repo")
	}
	return InitialAccountConfigurationHandler{
		repo:       repo,
		sobService: sobService,
	}
}

func (h InitialAccountConfigurationHandler) Handle(ctx context.Context, sobId uuid.UUID) error {
	// 0. read sob
	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return errors.Wrap(err, "read sob failed")
	}

	// 1. read CSV
	accountConfigurationEntries, err := h.readFromCSV()
	if err != nil {
		return err
	}

	// 2. prepare accounts
	preparedAccounts, err := h.prepareAccountConfigurations(sobId, accountConfigurationEntries, sob.AccountsCodeLength)
	if err != nil {
		return err
	}

	return h.repo.InitialAccountConfiguration(ctx, preparedAccounts)
}

func (h InitialAccountConfigurationHandler) readFromCSV() ([]accountConfigurationEntry, error) {
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

	var entries []accountConfigurationEntry
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
			balanceDirection = balance_direction.Unknown.String()
		}
		entries = append(entries, accountConfigurationEntry{
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

func (h InitialAccountConfigurationHandler) prepareAccountConfigurations(sobId uuid.UUID, accountEntries []accountConfigurationEntry, codeLengthLimits []int) ([]*account_configuration.AccountConfiguration, error) {
	preparedAccounts := make(map[string]*account_configuration.AccountConfiguration)
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
					superiorAccountId = superiorAccount.AccountId()
					numberHierarchy = append(superiorAccount.NumberHierarchy(), levelNumber)
				}
				account, err := account_configuration.New(sobId, uuid.New(), superiorAccountId, entry.title, entry.number, numberHierarchy, entry.level, entry.accountType, entry.balanceDirection)
				if err != nil {
					return nil, errors.Wrapf(err, "dataload failed on account %s", entry.number)
				}
				preparedAccounts[entry.number] = account
			}
		}
	}

	// to slice
	accounts := make([]*account_configuration.AccountConfiguration, len(preparedAccounts))
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
