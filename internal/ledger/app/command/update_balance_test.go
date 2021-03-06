package command

import (
	"context"
	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain/ledger"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp_UpdateLedgerBalance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		cmdConstructor func() UpdateLedgerBalanceCmd
		verify         func(t *testing.T, err error, repo ledgerRepoMock)
	}{
		{
			name: "update_success",
			cmdConstructor: func() UpdateLedgerBalanceCmd {
				return UpdateLedgerBalanceCmd{
					VoucherUUID: uuid.New(),
					LineItems: []LineItemCmd{
						{
							AccountNumber: "10000101",
							Debit:         decimal.RequireFromString("100"),
							Credit:        decimal.Zero,
						},
						{
							AccountNumber: "20000202",
							Debit:         decimal.RequireFromString("100"),
							Credit:        decimal.Zero,
						},
					},
				}
			},
			verify: func(t *testing.T, err error, repo ledgerRepoMock) {
				require.NoError(t, err)
				assert.Equal(t, decimal.RequireFromString("100"), repo.ledgers["1000"].Balance())
				assert.Equal(t, decimal.RequireFromString("100"), repo.ledgers["100001"].Balance())
				assert.Equal(t, decimal.RequireFromString("100"), repo.ledgers["10000101"].Balance())

				assert.Equal(t, decimal.RequireFromString("-100"), repo.ledgers["2000"].Balance())
				assert.Equal(t, decimal.RequireFromString("-100"), repo.ledgers["200002"].Balance())
				assert.Equal(t, decimal.RequireFromString("-100"), repo.ledgers["20000202"].Balance())
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repoMock := newLedgerRepoMock()
			repoMock.initData()
			h := NewUpdateLedgerBalanceHandler(repoMock, newAccountServiceMock(), newVoucherServiceMock())
			err := h.Handle(context.Background(), tt.cmdConstructor())
			tt.verify(t, err, repoMock)
		})
	}
}

type ledgerRepoMock struct {
	ledgers map[string]*ledger.Ledger
}

func (r ledgerRepoMock) UpdateLedgers(
	ctx context.Context,
	ledgerNumbers []string,
	updateFn func(ledgers []*ledger.Ledger) ([]*ledger.Ledger, error),
) error {
	// fetch entries from db
	var ledgers []*ledger.Ledger
	for dbNum, l := range r.ledgers {
		for _, inputNum := range ledgerNumbers {
			if inputNum == dbNum {
				ledgers = append(ledgers, l)
				break
			}
		}
	}
	// call updateFn
	afterUpdateLedgers, err := updateFn(ledgers)
	if err != nil {
		return errors.Wrap(err, "ledger list update failed")
	}
	// write to db
	for _, l := range afterUpdateLedgers {
		r.ledgers[l.Number()] = l
	}
	return nil
}

func (r ledgerRepoMock) initData() {
	r.ledgers["1000"], _ = ledger.NewLedger("1000", "1000 ledger", "", commonAccount.Assets)
	r.ledgers["100001"], _ = ledger.NewLedger("100001", "100001 ledger", "1000", commonAccount.Assets)
	r.ledgers["10000101"], _ = ledger.NewLedger("10000101", "10000101 ledger", "100001", commonAccount.Assets)

	r.ledgers["2000"], _ = ledger.NewLedger("2000", "2000 ledger", "", commonAccount.Liabilities)
	r.ledgers["200002"], _ = ledger.NewLedger("200002", "200002 ledger", "2000", commonAccount.Liabilities)
	r.ledgers["20000202"], _ = ledger.NewLedger("20000202", "20000202 ledger", "200002", commonAccount.Liabilities)
}

func newLedgerRepoMock() ledgerRepoMock {
	return ledgerRepoMock{ledgers: make(map[string]*ledger.Ledger)}
}

type accountServiceMock struct{}

func (s accountServiceMock) readSuperiorNumbers(ctx context.Context, accountNumber string) ([]string, error) {
	if len(accountNumber) != 8 {
		return nil, errors.New("let's test with 8 length account number")
	}
	return []string{accountNumber, accountNumber[:6], accountNumber[:4]}, nil
}

func newAccountServiceMock() accountServiceMock {
	return accountServiceMock{}
}

type voucherServiceMock struct{}

func (s voucherServiceMock) checkVoucherPosted(ctx context.Context, voucherUUID uuid.UUID) (bool, error) {
	return false, nil
}

func newVoucherServiceMock() voucherServiceMock {
	return voucherServiceMock{}
}
