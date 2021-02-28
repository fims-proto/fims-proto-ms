package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccountTypeFromConstant(t *testing.T) {
	accountType, _ := NewAccountType(Assets)
	assert.Equal(t, Assets, accountType)
}

func TestAccountTypeToString(t *testing.T) {
	accountType, err := NewAccountType(Assets)
	require.NoError(t, err)
	assert.Equal(t, "Assets", accountType.String())
}
