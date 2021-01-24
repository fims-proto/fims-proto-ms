package accounttype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountTypeFromString(t *testing.T) {
	accountType, _ := NewAccountTypeFromString("assets")
	assert.Equal(t, "assets", accountType.String())
}

func TestNewAccountTypeFromString_UnknownType(t *testing.T) {
	_, err := NewAccountTypeFromString("unknown")
	assert.Error(t, err)
}
