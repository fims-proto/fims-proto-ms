package command

import (
	"context"
	commonaccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
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
				assert.Equal(t, "100", repo.ledgers["1000"].Balance().String())
				assert.Equal(t, "100", repo.ledgers["100001"].Balance().String())
				assert.Equal(t, "100", repo.ledgers["10000101"].Balance().String())

				assert.Equal(t, "100", repo.ledgers["1000"].Debit().String())
				assert.Equal(t, "100", repo.ledgers["100001"].Debit().String())
				assert.Equal(t, "100", repo.ledgers["10000101"].Debit().String())

				assert.Equal(t, "-100", repo.ledgers["2000"].Balance().String())
				assert.Equal(t, "-100", repo.ledgers["200002"].Balance().String())
				assert.Equal(t, "-100", repo.ledgers["20000202"].Balance().String())

				assert.Equal(t, "100", repo.ledgers["2000"].Debit().String())
				assert.Equal(t, "100", repo.ledgers["200002"].Debit().String())
				assert.Equal(t, "100", repo.ledgers["20000202"].Debit().String())
			},
		},
		{
			name: "update_success_2",
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
							AccountNumber: "30000101",
							Debit:         decimal.Zero,
							Credit:        decimal.RequireFromString("100"),
						},
					},
				}
			},
			verify: func(t *testing.T, err error, repo ledgerRepoMock) {
				require.NoError(t, err)
				assert.Equal(t, "100", repo.ledgers["1000"].Balance().String())
				assert.Equal(t, "100", repo.ledgers["100001"].Balance().String())
				assert.Equal(t, "100", repo.ledgers["10000101"].Balance().String())

				assert.Equal(t, "100", repo.ledgers["1000"].Debit().String())
				assert.Equal(t, "100", repo.ledgers["100001"].Debit().String())
				assert.Equal(t, "100", repo.ledgers["10000101"].Debit().String())

				assert.Equal(t, "-100", repo.ledgers["3000"].Balance().String())
				assert.Equal(t, "-100", repo.ledgers["300001"].Balance().String())
				assert.Equal(t, "-100", repo.ledgers["30000101"].Balance().String())

				assert.Equal(t, "100", repo.ledgers["3000"].Credit().String())
				assert.Equal(t, "100", repo.ledgers["300001"].Credit().String())
				assert.Equal(t, "100", repo.ledgers["30000101"].Credit().String())
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
	ledgers map[string]*domain.Ledger
}

func (r ledgerRepoMock) AddLedger(ctx context.Context, l *domain.Ledger) error {
	panic("not implemented")
}

func (r ledgerRepoMock) Dataload(ctx context.Context, ls []*domain.Ledger) error {
	panic("not implemented")
}

func (r ledgerRepoMock) UpdateLedgers(
	ctx context.Context,
	ledgerNumbers []string,
	updateFn func(ledgers []*domain.Ledger) ([]*domain.Ledger, error),
) error {
	// fetch entries from db
	var ledgers []*domain.Ledger
	for _, inputNum := range ledgerNumbers {
		l, ok := r.ledgers[inputNum]
		if !ok {
			return errors.Errorf("ledger number %s not exists", inputNum)
		}
		ledgers = append(ledgers, l)
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
	r.ledgers["1000"], _ = domain.NewLedger("1000", "1000 ledger", "", commonaccount.Assets)
	r.ledgers["100001"], _ = domain.NewLedger("100001", "100001 ledger", "1000", commonaccount.Assets)
	r.ledgers["10000101"], _ = domain.NewLedger("10000101", "10000101 ledger", "100001", commonaccount.Assets)

	r.ledgers["2000"], _ = domain.NewLedger("2000", "2000 ledger", "", commonaccount.Liabilities)
	r.ledgers["200002"], _ = domain.NewLedger("200002", "200002 ledger", "2000", commonaccount.Liabilities)
	r.ledgers["20000202"], _ = domain.NewLedger("20000202", "20000202 ledger", "200002", commonaccount.Liabilities)

	r.ledgers["3000"], _ = domain.NewLedger("3000", "3000 ledger", "", commonaccount.Assets)
	r.ledgers["300001"], _ = domain.NewLedger("300001", "300001 ledger", "3000", commonaccount.Assets)
	r.ledgers["30000101"], _ = domain.NewLedger("30000101", "30000101 ledger", "300001", commonaccount.Assets)
}

func newLedgerRepoMock() ledgerRepoMock {
	return ledgerRepoMock{ledgers: make(map[string]*domain.Ledger)}
}

type accountServiceMock struct{}

func (s accountServiceMock) ReadSuperiorNumbers(ctx context.Context, accountNumber string) ([]string, error) {
	if len(accountNumber) != 8 {
		return nil, errors.New("let's test with 8 length account number")
	}
	return []string{accountNumber, accountNumber[:6], accountNumber[:4]}, nil
}

func newAccountServiceMock() accountServiceMock {
	return accountServiceMock{}
}

type voucherServiceMock struct{}

func (s voucherServiceMock) CheckVoucherPosted(ctx context.Context, voucherUUID uuid.UUID) (bool, error) {
	return false, nil
}

func newVoucherServiceMock() voucherServiceMock {
	return voucherServiceMock{}
}
