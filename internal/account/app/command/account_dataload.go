package command

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github/fims-proto/fims-proto-ms/internal/account/app/service"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
	"github/fims-proto/fims-proto-ms/internal/common/log"

	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type accountDataLoadEntry struct {
	number           string
	level            int
	title            string
	superiorNumber   string
	accountType      string
	balanceDirection string
}

type AccountDataLoadHandler struct {
	repo       domain.Repository
	sobService service.SobService
}

func NewAccountDataLoadHandler(repo domain.Repository, sobService service.SobService) AccountDataLoadHandler {
	if repo == nil {
		panic("nil repo")
	}
	return AccountDataLoadHandler{
		repo:       repo,
		sobService: sobService,
	}
}

func (h AccountDataLoadHandler) Handle(ctx context.Context, sobId uuid.UUID) (err error) {
	log.Info(ctx, "handle accounts data load for sob %s", sobId.String())
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle accounts data load for sob %s failed", sobId.String())
		}
	}()

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
	log.Info(ctx, "loaded csv file, size: %d", len(accountEntries))

	// 2. prepare accounts
	preparedAccounts, err := h.prepareAccounts(sobId, accountEntries, sob.AccountsCodeLength)
	if err != nil {
		return err
	}

	return h.repo.DataLoad(ctx, preparedAccounts)
}

func (h AccountDataLoadHandler) readFromCSV() ([]accountDataLoadEntry, error) {
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

	var entries []accountDataLoadEntry
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
			balanceDirection = commonAccount.UndefinedDirection.String()
		}
		entries = append(entries, accountDataLoadEntry{
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

func (h AccountDataLoadHandler) prepareAccounts(sobId uuid.UUID, accountEntries []accountDataLoadEntry, codeLengthLimits []int) ([]*domain.Account, error) {
	preparedAccounts := make(map[string]*domain.Account)
	for i := 0; i < len(codeLengthLimits); i++ {
		for _, entry := range accountEntries {
			if entry.level == i+1 {

				var accountNumber int
				var superiorAccountId uuid.UUID
				var superiorNumbers []int
				if entry.level == 1 {
					accountNumber, _ = strconv.Atoi(entry.number)
					superiorAccountId = uuid.Nil
					superiorNumbers = []int{}
				} else {
					accountNumber, _ = strconv.Atoi(strings.TrimPrefix(entry.number, entry.superiorNumber))
					superiorAccount, ok := preparedAccounts[entry.superiorNumber]
					if !ok {
						return nil, errors.Errorf("cannot find prepared superior account %s", entry.superiorNumber)
					}
					superiorAccountId = superiorAccount.Id()
					superiorNumbers = append(superiorAccount.SuperiorNumbers(), superiorAccount.LevelNumber())
				}
				account, err := domain.NewAccount(uuid.New(), sobId, superiorAccountId, superiorNumbers, entry.title, accountNumber, i+1, entry.accountType, entry.balanceDirection, codeLengthLimits[i])
				if err != nil {
					return nil, errors.Wrapf(err, "dataload failed on account %s", entry.number)
				}
				preparedAccounts[entry.number] = account

			}
		}
	}

	// to slice
	accounts := make([]*domain.Account, len(preparedAccounts))
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
