package account

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account_type"
	"testing"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		number         string
		title          string
		superiorNumber string
		accType        accounttype.Type
		verify         func(t *testing.T, account *Account, err error)
	}{
		{
			"general_account_success", "1001", "库存现金", "", accounttype.Assets,
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, err)
				assert.Equal(t, "1001", account.Number())
				assert.Equal(t, "库存现金", account.Title())
				assert.Empty(t, account.SuperiorNumber())
				assert.Equal(t, accounttype.Assets, account.Type())
			},
		},
		{
			"subsidiary_account_success", "1001001", "库存现金某子项", "1001", accounttype.Assets,
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, err)
				assert.Equal(t, "1001001", account.Number())
				assert.Equal(t, "库存现金某子项", account.Title())
				assert.Equal(t, "1001", account.SuperiorNumber())
				assert.Equal(t, accounttype.Assets, account.Type())
			},
		},
		{
			"empty_number_error", "", "库存现金", "", accounttype.Assets,
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
		{
			"empty_title_error", "1001", "", "", accounttype.Assets,
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
		{
			"invalid_superior_account_number_error", "1001001", "库存现金某子项", "1002", accounttype.Assets,
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
			account, err := NewAccount(test.number, test.title, test.superiorNumber, test.accType)
			test.verify(t, account, err)
		})
	}
}
