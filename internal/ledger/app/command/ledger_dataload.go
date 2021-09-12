package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type LedgerDataloadCmd struct {
	Number         string
	Title          string
	SuperiorNumber string
	AccountType    string
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

func (h LedgerDataloadHandler) Handle(ctx context.Context, sob string, cmds []LedgerDataloadCmd) (err error) {
	log.Info(ctx, "handle ledger dataload for sob %s", sob)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle ledger dataload for sob %s failed", sob)
		}
	}()

	decimalZero := decimal.RequireFromString("0")

	var ledgers []*domain.Ledger
	for _, cmd := range cmds {
		ledger, err := domain.NewLedger(uuid.New(), sob, cmd.Number, cmd.Title, cmd.SuperiorNumber, cmd.AccountType, decimalZero, decimalZero, decimalZero)
		if err != nil {
			return errors.Wrapf(err, "dataload failed on ledger number %s", cmd.Number)
		}
		ledgers = append(ledgers, ledger)
	}

	return h.repo.Dataload(ctx, ledgers)
}
