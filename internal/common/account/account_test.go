package account

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestAccountType(t *testing.T) {
	accountType, err := NewAccountType("ASSETS")
	require.NoError(t, err)
	assert.Equal(t, Assets, accountType)
	assert.Equal(t, "ASSETS", accountType.String())
}

func TestNewDirectionFromString(t *testing.T) {
	direction, err := NewDirection("Debit")
	require.NoError(t, err)
	assert.Equal(t, Debit, direction)
	assert.Equal(t, "Debit", direction.String())
}
