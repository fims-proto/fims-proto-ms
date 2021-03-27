package adapter

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
	commonaccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdapter_MemoryRepository_ReadOneByNumber(t *testing.T) {
	t.Parallel()

	repo := prepareMemoryRepo(t)
	account, err := repo.AccountByNumber(context.Background(), "10000101")
	require.NoError(t, err)

	assert.Equal(t, "10000101", account.Number)
	assert.Equal(t, "100001", account.SuperiorAccount.Number)
	assert.Equal(t, "1000", account.SuperiorAccount.SuperiorAccount.Number)
}

func prepareMemoryRepo(t *testing.T) AccountMemoryRepository {
	repo := NewAccountMemoryRepository()
	var accounts []*domain.Account

	a, err := domain.NewAccount("1000", "1000 title", "", commonaccount.Assets)
	require.NoError(t, err)
	accounts = append(accounts, a)

	a, err = domain.NewAccount("100001", "100001 title", "1000", commonaccount.Assets)
	require.NoError(t, err)
	accounts = append(accounts, a)

	a, err = domain.NewAccount("10000101", "10000101 title", "100001", commonaccount.Assets)
	require.NoError(t, err)
	accounts = append(accounts, a)

	err = repo.Dataload(context.Background(), accounts)
	require.NoError(t, err)

	return repo
}
