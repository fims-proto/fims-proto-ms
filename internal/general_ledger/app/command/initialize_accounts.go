package command

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/balance_direction"
	sobQuery "github/fims-proto/fims-proto-ms/internal/sob/app/query"

	"github.com/google/uuid"
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
	preparedAccounts, err := prepareAccounts(sob.Id, accountEntries)
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

func prepareAccounts(sobId uuid.UUID, accountEntries []accountEntry) ([]*account.Account, error) {
	// Build a set of superior account numbers for quick lookup (O(1) instead of binary search)
	superiorNumbers := make(map[string]bool)
	for _, entry := range accountEntries {
		if entry.superiorNumber != "" {
			superiorNumbers[entry.superiorNumber] = true
		}
	}

	// Map from raw account number to domain account object
	preparedAccounts := make(map[string]*account.Account) // keyed by raw account number

	// Process accounts level by level to ensure superiors are created first
	maxLevel := 0
	for _, entry := range accountEntries {
		if entry.level > maxLevel {
			maxLevel = entry.level
		}
	}

	for level := 1; level <= maxLevel; level++ {
		for _, entry := range accountEntries {
			if entry.level == level {
				var levelNumber int
				var superiorAccountId uuid.UUID
				var superiorRaw string

				if entry.level == 1 {
					// Level 1: extract the single 6-digit segment directly
					superiorAccountId = uuid.Nil
					levelNumberStr := entry.number[:6]
					var err error
					levelNumber, err = strconv.Atoi(levelNumberStr)
					if err != nil {
						return nil, fmt.Errorf("invalid level number in account %s: %w", entry.number, err)
					}
					superiorRaw = ""
				} else {
					// Level 2+: get superior from already-prepared accounts
					superiorAccount, ok := preparedAccounts[entry.superiorNumber]
					if !ok {
						return nil, fmt.Errorf("cannot find prepared superior account %s for %s", entry.superiorNumber, entry.number)
					}
					superiorAccountId = superiorAccount.Id()
					superiorRaw = superiorAccount.RawAccountNumber()

					// Extract just the last 6-digit segment (the level number)
					lastSegmentStr := entry.number[len(entry.number)-6:]
					var err error
					levelNumber, err = strconv.Atoi(lastSegmentStr)
					if err != nil {
						return nil, fmt.Errorf("invalid level number in account %s: %w", entry.number, err)
					}
				}

				// Check if this account is a superior for any other account (O(1) lookup)
				isLeaf := !superiorNumbers[entry.number]

				domainAccount, err := account.New(
					uuid.New(),
					sobId,
					superiorAccountId,
					entry.title,
					superiorRaw,
					levelNumber,
					entry.level,
					isLeaf,
					entry.class,
					entry.group,
					entry.balanceDirection,
					nil,
				)
				if err != nil {
					return nil, fmt.Errorf("dataload failed on account %s: %w", entry.number, err)
				}

				// Store using entry.number (the raw account number from CSV)
				preparedAccounts[entry.number] = domainAccount
			}
		}
	}

	// Convert map to slice
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
