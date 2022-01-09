package domain

import (
	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		title             string
		levelNumber       int
		superiorAccountId uuid.UUID
		superiorNumbers   []int
		level             int
		levelCodeLength   int
		accountType       string
		balanceDirection  string
		verify            func(t *testing.T, account *Account, err error)
	}{
		{
			"general_account_success", "库存现金", 1001, uuid.Nil,
			[]int{},
			1, 4, "ASSETS", "DEBIT",
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, err)
				assert.Equal(t, 1001, account.LevelNumber())
				assert.Equal(t, "库存现金", account.Title())
				assert.Equal(t, uuid.Nil, account.SuperiorAccountId())
				assert.Empty(t, account.SuperiorNumbers())
				assert.Equal(t, commonAccount.Assets, account.Type())
				assert.Equal(t, commonAccount.Debit, account.BalanceDirection())
			},
		},
		{
			"subsidiary_account_success", "库存现金某子项", 1, uuid.New(),
			[]int{1001},
			2, 3, "ASSETS", "DEBIT",
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, err)
				assert.Equal(t, "1001001", account.LevelNumber())
				assert.Equal(t, "库存现金某子项", account.Title())
				assert.NotNil(t, account.SuperiorAccountId())
				assert.Equal(t, 1001, account.SuperiorNumbers()[0])
				assert.Equal(t, commonAccount.Assets, account.Type())
				assert.Equal(t, commonAccount.Debit, account.BalanceDirection())
			},
		},
		{
			"zero_number_error", "库存现金", 0, uuid.Nil,
			[]int{},
			1, 4, "ASSETS", "DEBIT",
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
		{
			"zero_number_error", "", 1001, uuid.Nil,
			[]int{},
			1, 4, "ASSETS", "DEBIT",
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
		{
			"empty_superior_account", "库存现金某子项", 1, uuid.Nil,
			[]int{},
			2, 3, "ASSETS", "DEBIT",
			func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
		{
			"incorrect_superior_number", "库存现金某子项", 1, uuid.New(),
			[]int{},
			2, 3, "ASSETS", "DEBIT",
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
			account, err := NewAccount(uuid.New(), uuid.New(), test.superiorAccountId, test.superiorNumbers, test.title, test.levelNumber, test.level, test.accountType, test.balanceDirection, test.levelCodeLength)
			test.verify(t, account, err)
		})
	}
}
