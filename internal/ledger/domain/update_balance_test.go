package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain_LedgerUpdateBalance(t *testing.T) {
	t.Parallel()
	type args struct {
		debit  decimal.Decimal
		credit decimal.Decimal
	}
	tests := []struct {
		name    string
		accType string
		args    args
		verify  func(t *testing.T, l Ledger, err error)
	}{
		{
			name:    "update assets account ledger",
			accType: "Assets",
			args: args{
				debit:  decimal.RequireFromString("100"),
				credit: decimal.RequireFromString("50"),
			},
			verify: func(t *testing.T, l Ledger, err error) {
				require.NoError(t, err)
				assert.True(t, decimal.RequireFromString("50").Equal(l.Balance()))
			},
		},
		{
			name:    "update liabilities account ledger",
			accType: "Liabilities",
			args: args{
				debit:  decimal.RequireFromString("100"),
				credit: decimal.RequireFromString("50"),
			},
			verify: func(t *testing.T, l Ledger, err error) {
				require.NoError(t, err)
				assert.True(t, decimal.RequireFromString("-50").Equal(l.Balance()))
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := prepareLedger(tt.accType)
			err := l.UpdateBalance(tt.args.debit, tt.args.credit)
			tt.verify(t, l, err)
		})
	}
}

func prepareLedger(accType string) Ledger {
	zero := decimal.RequireFromString("0")
	l, _ := NewLedger(uuid.New(), "test_sob", "0000", "test", "", accType, zero, zero, zero)
	return *l
}
