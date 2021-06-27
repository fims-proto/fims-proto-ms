package domain

import (
	commonaccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		number         string
		title          string
		superiorNumber string
		accType        commonaccount.Type
		verify         func(t *testing.T, account *Account, err error)
	}{
		{
			"general_account_success", "1001", "库存现金", "", commonaccount.Assets,
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, err)
				assert.Equal(t, "1001", account.Number())
				assert.Equal(t, "库存现金", account.Title())
				assert.Empty(t, account.SuperiorNumber())
				assert.Equal(t, commonaccount.Assets, account.Type())
			},
		},
		{
			"subsidiary_account_success", "1001001", "库存现金某子项", "1001", commonaccount.Assets,
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, err)
				assert.Equal(t, "1001001", account.Number())
				assert.Equal(t, "库存现金某子项", account.Title())
				assert.Equal(t, "1001", account.SuperiorNumber())
				assert.Equal(t, commonaccount.Assets, account.Type())
			},
		},
		{
			"empty_number_error", "", "库存现金", "", commonaccount.Assets,
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
		{
			"empty_title_error", "1001", "", "", commonaccount.Assets,
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
		{
			"invalid_superior_account_number_error", "1001001", "库存现金某子项", "1002", commonaccount.Assets,
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			account, err := NewAccount("test_sob", test.number, test.title, test.superiorNumber, test.accType)
			test.verify(t, account, err)
		})
	}
}
