package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var accountId = uuid.New()

func TestDomain_NewLineItem(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		summary   string
		accountId uuid.UUID
		debit     decimal.Decimal
		credit    decimal.Decimal
		verify    func(t *testing.T, lineItem *LineItem, err error)
	}{
		{
			"debit_success", "Test Summary", accountId, decimal.RequireFromString("200.00"), decimal.Zero,
			func(t *testing.T, lineItem *LineItem, err error) {
				debit, _ := decimal.NewFromString("200.00")
				require.NoError(t, err)
				assert.Equal(t, "Test Summary", lineItem.Summary())
				assert.Equal(t, accountId, lineItem.AccountId())
				assert.Equal(t, debit, lineItem.Debit())
				assert.True(t, lineItem.Credit().IsZero())
			},
		},
		{
			"credit_success", "Test Summary", accountId, decimal.Zero, decimal.RequireFromString("200.00"),
			func(t *testing.T, lineItem *LineItem, err error) {
				credit, _ := decimal.NewFromString("200.00")
				require.NoError(t, err)
				assert.Equal(t, "Test Summary", lineItem.Summary())
				assert.Equal(t, accountId, lineItem.AccountId())
				assert.Equal(t, credit, lineItem.Credit())
				assert.True(t, lineItem.Debit().IsZero())
			},
		},
		{
			"debit_and_credit_error", "Test Summary", accountId, decimal.RequireFromString("200.00"), decimal.RequireFromString("200.00"),
			func(t *testing.T, lineItem *LineItem, err error) {
				require.Nil(t, lineItem)
				assert.Error(t, err)
			},
		},
		{
			"empty_debit_credit_error", "Test Summary", accountId, decimal.Zero, decimal.Zero,
			func(t *testing.T, lineItem *LineItem, err error) {
				require.Nil(t, lineItem)
				assert.Error(t, err)
			},
		},
		{
			"empty_summary_error", "", accountId, decimal.RequireFromString("200.00"), decimal.Zero,
			func(t *testing.T, lineItem *LineItem, err error) {
				require.Nil(t, lineItem)
				assert.Error(t, err)
			},
		},
		{
			"nil_account_id_error", "Test Summary", uuid.Nil, decimal.RequireFromString("200.00"), decimal.Zero,
			func(t *testing.T, lineItem *LineItem, err error) {
				require.Nil(t, lineItem)
				assert.Error(t, err)
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			item, err := NewLineItem(uuid.New(), test.accountId, test.summary, test.debit, test.credit)
			test.verify(t, item, err)
		})
	}
}
