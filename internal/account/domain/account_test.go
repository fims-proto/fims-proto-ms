package domain

import (
	"testing"

	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		title             string
		superiorAccountId uuid.UUID
		accountNumber     string
		level             int
		numberHierarchy   []int
		codeLength        []int
		accountType       string
		balanceDirection  string
		verify            func(t *testing.T, account *Account, err error)
	}{
		{
			name:              "general_account_success",
			title:             "库存现金",
			superiorAccountId: uuid.Nil,
			accountNumber:     "1001",
			level:             1,
			numberHierarchy:   []int{1001},
			codeLength:        []int{4, 3, 3},
			accountType:       "assets",
			balanceDirection:  "DEBIT",
			verify: func(t *testing.T, account *Account, err error) {
				require.Nil(t, err)
				assert.Equal(t, []int{1001}, account.NumberHierarchy())
				assert.Equal(t, "库存现金", account.Title())
				assert.Equal(t, "1001", account.AccountNumber())
				assert.Equal(t, 1, account.Level())
				assert.Equal(t, uuid.Nil, account.SuperiorAccountId())
				assert.Equal(t, commonAccount.Assets, account.Type())
				assert.Equal(t, commonAccount.Debit, account.BalanceDirection())
			},
		},
		{
			name:              "subsidiary_account_success",
			title:             "库存现金某子项",
			superiorAccountId: uuid.New(),
			accountNumber:     "1001001",
			level:             2,
			numberHierarchy:   []int{1001, 1},
			codeLength:        []int{4, 3, 3},
			accountType:       "assets",
			balanceDirection:  "DEBIT",
			verify: func(t *testing.T, account *Account, err error) {
				require.Nil(t, err)
				assert.Equal(t, []int{1001, 1}, account.NumberHierarchy())
				assert.Equal(t, "库存现金某子项", account.Title())
				assert.Equal(t, "1001001", account.AccountNumber())
				assert.Equal(t, 2, account.Level())
				assert.NotNil(t, account.SuperiorAccountId())
				assert.Equal(t, commonAccount.Assets, account.Type())
				assert.Equal(t, commonAccount.Debit, account.BalanceDirection())
			},
		},
		{
			name:              "zero_number_error",
			title:             "库存现金",
			superiorAccountId: uuid.Nil,
			accountNumber:     "0000",
			level:             1,
			numberHierarchy:   []int{0},
			codeLength:        []int{4, 3, 3},
			accountType:       "assets",
			balanceDirection:  "DEBIT",
			verify: func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
		{
			name:              "empty_superior_account",
			title:             "库存现金",
			superiorAccountId: uuid.Nil,
			accountNumber:     "1001001",
			level:             2,
			numberHierarchy:   []int{1001, 1},
			codeLength:        []int{4, 3, 3},
			accountType:       "assets",
			balanceDirection:  "DEBIT",
			verify: func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
		{
			name:              "account_length_too_long",
			title:             "库存现金",
			superiorAccountId: uuid.Nil,
			accountNumber:     "10011111",
			level:             2,
			numberHierarchy:   []int{1001, 1111},
			codeLength:        []int{4, 3, 3},
			accountType:       "assets",
			balanceDirection:  "DEBIT",
			verify: func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
		{
			name:              "account_depth_too_long",
			title:             "库存现金",
			superiorAccountId: uuid.New(),
			accountNumber:     "1001001001001",
			level:             1,
			numberHierarchy:   []int{1001, 1, 1, 1},
			codeLength:        []int{4, 3, 3},
			accountType:       "assets",
			balanceDirection:  "DEBIT",
			verify: func(t *testing.T, account *Account, err error) {
				require.Nil(t, account)
				assert.Error(t, err)
			},
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			account, err := NewAccount(uuid.New(), uuid.New(), tt.superiorAccountId, tt.title, tt.accountNumber, tt.accountType, tt.balanceDirection, tt.level, tt.numberHierarchy, tt.codeLength)
			tt.verify(t, account, err)
		})
	}
}
