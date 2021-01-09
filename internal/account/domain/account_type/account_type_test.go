package accounttype

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccountTypeFromString(t *testing.T) {
	accountType, _ := NewAccountTypeFromString("assets")
	assert.Equal(t, "assets", accountType.String())
}

func TestNewAccountTypeFromString_UnknownType(t *testing.T) {
	_, err := NewAccountTypeFromString("unknown")
	assert.Error(t, err)
}
