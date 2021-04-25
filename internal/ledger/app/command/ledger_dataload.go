package command

import (
	"context"
	commonaccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"

	"github.com/pkg/errors"
)

type LedgerDataloadCmd struct {
	Number         string
	Title          string
	SuperiorNumber string
	AccountType    commonaccount.Type
}

type LedgerDataloadHandler struct {
	repo domain.Repository
}

func NewLedgerDataloadHandler(repo domain.Repository) LedgerDataloadHandler {
	if repo == nil {
		panic("nil repo")
	}
	return LedgerDataloadHandler{repo: repo}
}

func (h LedgerDataloadHandler) Handle(ctx context.Context, cmds []LedgerDataloadCmd) error {
	var ledgers []*domain.Ledger
	for _, cmd := range cmds {
		ledger, err := domain.NewLedger(cmd.Number, cmd.Title, cmd.SuperiorNumber, cmd.AccountType)
		if err != nil {
			return errors.Wrapf(err, "dataload failed on ledger number %s", cmd.Number)
		}
		ledgers = append(ledgers, ledger)
	}

	return h.repo.Dataload(ctx, ledgers)
}
