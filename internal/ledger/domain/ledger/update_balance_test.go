package ledger

import (
	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"testing"

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
		accType commonAccount.Type
		args    args
		verify  func(t *testing.T, l Ledger, err error)
	}{
		{
			name:    "update assets account ledger",
			accType: commonAccount.Assets,
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
			accType: commonAccount.Liabilities,
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

func prepareLedger(accType commonAccount.Type) Ledger {
	l, _ := NewLedger("0000", "test", "", accType)
	return *l
}
