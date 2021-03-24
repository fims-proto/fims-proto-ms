package command

import (
	"context"
	"encoding/csv"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"io"
	"os"

	"github.com/pkg/errors"
)

type AccountDataloadCmd struct {
	Number         string
	Title          string
	SuperiorNumber string
	AccountType    commonAccount.Type
}

type AccountDataloadHandler struct {
	repo domain.Repository
}

func NewAccountDataloadHandler(repo domain.Repository) AccountDataloadHandler {
	if repo == nil {
		panic("nil repo")
	}
	return AccountDataloadHandler{repo: repo}
}

func (h AccountDataloadHandler) Handle(ctx context.Context) error {
	accountCmds, err := readFromCSV()
	if err != nil {
		return err
	}

	var accounts []*domain.Account
	for _, cmd := range accountCmds {
		account, err := domain.NewAccount(cmd.Number, cmd.Title, cmd.SuperiorNumber, cmd.AccountType)
		if err != nil {
			return errors.Wrapf(err, "dataload failed on account %s", cmd.Number)
		}
		accounts = append(accounts, account)
	}

	return h.repo.AddAccounts(ctx, accounts)
}

func readFromCSV() ([]AccountDataloadCmd, error) {
	csvFile, err := os.Open("dataload/accounts.csv")
	if err != nil {
		return nil, errors.Wrap(err, "could not load dataload/accounts.csv file")
	}

	csvReader := csv.NewReader(csvFile)

	cmds := []AccountDataloadCmd{}
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "could not read dataload/accounts.csv file")
		}
		accountType, err := commonAccount.NewAccountTypeFromString(line[0])
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert account type")
		}
		cmds = append(cmds, AccountDataloadCmd{
			Number:         line[1],
			Title:          line[2],
			SuperiorNumber: "",
			AccountType:    accountType,
		})
	}

	return cmds, nil
}
