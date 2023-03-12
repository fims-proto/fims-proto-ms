package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldType(t *testing.T) {
	testfield, err := New("test")
	require.NoError(t, err)
	assert.Equal(t, testfield.Name(), "test")
}
