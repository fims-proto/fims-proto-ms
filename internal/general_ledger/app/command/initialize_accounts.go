package command

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/balance_direction"
	sobQuery "github/fims-proto/fims-proto-ms/internal/sob/app/query"
)

type accountEntry struct {
	number           string
	level            int
	title            string
	superiorNumber   string
	class            int
	group            int
	balanceDirection string
}

func initializeAccounts(ctx context.Context, sob sobQuery.Sob, repo domain.Repository) error {
	// 1. read CSV
	accountEntries, err := readFromCSV()
	if err != nil {
		return err
	}

	// 2. prepare accounts
	preparedAccounts, err := prepareAccounts(sob.Id, accountEntries, sob.AccountsCodeLength)
	if err != nil {
		return err
	}

	return repo.InitialAccounts(ctx, preparedAccounts)
}

func readFromCSV() ([]accountEntry, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get working directory: %w", err)
	}

	csvFile, err := os.Open(filepath.Join(workDir, "dataload", "xqykjzz", "accounts.csv"))
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	csvReader := csv.NewReader(csvFile)

	// skip first line
	_, err = csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	var entries []accountEntry
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not read file: %w", err)
		}
		level, err := strconv.Atoi(line[3])
		if err != nil {
			return nil, fmt.Errorf("failed to convert level to number: %w", err)
		}
		classId, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("failed to convert class id to number: %w", err)
		}
		groupId, err := strconv.Atoi(line[1])
		if err != nil {
			return nil, fmt.Errorf("failed to convert group id to number: %w", err)
		}
		balanceDirection := line[6]
		if balanceDirection == "" {
			balanceDirection = balance_direction.NotDefined.String()
		}
		entries = append(entries, accountEntry{
			number:           line[2],
			level:            level,
			title:            line[4],
			superiorNumber:   line[5],
			class:            classId,
			group:            groupId,
			balanceDirection: balanceDirection,
		})
	}

	return entries, nil
}

func prepareAccounts(sobId uuid.UUID, accountEntries []accountEntry, codeLengthLimits []int) ([]*account.Account, error) {
	var superiorNumbers []string
	for _, entry := range accountEntries {
		if entry.superiorNumber != "" {
			superiorNumbers = append(superiorNumbers, entry.superiorNumber)
		}
	}
	slices.Sort(superiorNumbers)

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
						return nil, fmt.Errorf("cannot find prepared superior account %s", entry.superiorNumber)
					}
					superiorAccountId = superiorAccount.Id()
					numberHierarchy = append(superiorAccount.NumberHierarchy(), levelNumber)
				}

				// when an account is not superior for all other accounts, it's a leaf
				_, found := slices.BinarySearch(superiorNumbers, entry.number)

				domainAccount, err := account.New(
					uuid.New(),
					sobId,
					superiorAccountId,
					entry.title,
					numberHierarchy,
					codeLengthLimits,
					entry.level,
					!found,
					entry.class,
					entry.group,
					entry.balanceDirection,
					nil,
				)
				if err != nil {
					return nil, fmt.Errorf("dataload failed on account %s: %w", entry.number, err)
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
		return nil, fmt.Errorf("prepared accounts size (%d) doesn't equal to CSV entries size (%d)", len(accounts), len(accountEntries))
	}
	return accounts, nil
}
