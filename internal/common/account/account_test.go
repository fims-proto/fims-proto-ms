package account

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestAccountType(t *testing.T) {
	accountType, err := NewAccountType("assets")
	require.NoError(t, err)
	assert.Equal(t, Assets, accountType)
	assert.Equal(t, "assets", accountType.String())
}

func TestNewDirectionFromString(t *testing.T) {
	direction, err := NewDirection("debit")
	require.NoError(t, err)
	assert.Equal(t, Debit, direction)
	assert.Equal(t, "debit", direction.String())
}
