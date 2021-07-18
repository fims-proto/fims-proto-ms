package command

import (
	"context"
	"encoding/csv"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
	commonaccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type AccountDataloadCmd struct {
	Number         string
	Title          string
	SuperiorNumber string
	AccountType    commonaccount.Type
}

type AccountDataloadHandler struct {
	repo          domain.Repository
	ledgerService LedgerService
}

func NewAccountDataloadHandler(repo domain.Repository, ledgerService LedgerService) AccountDataloadHandler {
	if repo == nil {
		panic("nil repo")
	}
	if ledgerService == nil {
		panic("nil ledger service")
	}
	return AccountDataloadHandler{
		repo:          repo,
		ledgerService: ledgerService,
	}
}

func (h AccountDataloadHandler) Handle(ctx context.Context, sob string) error {
	accountCmds, err := readFromCSV()
	if err != nil {
		return err
	}

	var accounts []*domain.Account
	for _, cmd := range accountCmds {
		account, err := domain.NewAccount(sob, cmd.Number, cmd.Title, cmd.SuperiorNumber, cmd.AccountType)
		if err != nil {
			return errors.Wrapf(err, "dataload failed on account %s", cmd.Number)
		}
		accounts = append(accounts, account)
	}

	var immutableAccounts []domain.Account
	for _, account := range accounts {
		immutableAccounts = append(immutableAccounts, *account)
	}
	if err := h.ledgerService.LoadLedgers(ctx, sob, immutableAccounts); err != nil {
		return errors.Wrap(err, "ledger service failed to load data")
	}

	return h.repo.Dataload(ctx, accounts)
}

func readFromCSV() ([]AccountDataloadCmd, error) {
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

	cmds := []AccountDataloadCmd{}
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "could not read file")
		}
		accountType, err := commonaccount.NewAccountTypeFromString(line[0])
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert account type")
		}
		cmds = append(cmds, AccountDataloadCmd{
			Number:         line[1],
			Title:          line[2],
			SuperiorNumber: line[3],
			AccountType:    accountType,
		})
	}

	return cmds, nil
}
