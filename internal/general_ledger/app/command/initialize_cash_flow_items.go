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
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/cash_flow_item"

	"github.com/google/uuid"
)

func initializeCashFlowItems(ctx context.Context, sobId uuid.UUID, repo domain.Repository) error {
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %w", err)
	}

	csvFile, err := os.Open(filepath.Join(workDir, "dataload", "xqykjzz", "cash_flow_items.csv"))
	if err != nil {
		return fmt.Errorf("could not open cash_flow_items.csv: %w", err)
	}
	defer func() { _ = csvFile.Close() }()

	csvReader := csv.NewReader(csvFile)

	// skip header
	if _, err = csvReader.Read(); err != nil {
		return fmt.Errorf("could not read csv header: %w", err)
	}

	var items []*cash_flow_item.CashFlowItem
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("could not read cash_flow_items.csv: %w", err)
		}

		seq, err := strconv.Atoi(line[4])
		if err != nil {
			return fmt.Errorf("invalid sequence in cash flow item %s: %w", line[0], err)
		}

		item, err := cash_flow_item.New(uuid.New(), sobId, line[0], line[1], line[2], line[3], seq)
		if err != nil {
			return fmt.Errorf("invalid cash flow item %s: %w", line[0], err)
		}

		items = append(items, item)
	}

	return repo.InitializeCashFlowItems(ctx, items)
}
