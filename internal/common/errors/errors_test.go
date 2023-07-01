package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewErrNoWhereUsed(t *testing.T) {
	err := ErrNoWhereUsed("hello", "world")
	assert.True(t, errors.Is(err, ErrNoWhereUsed()))
	assert.True(t, errors.Is(fmt.Errorf("error occurred: %w", err), ErrNoWhereUsed()))
	assert.False(t, errors.Is(err, ErrRecordNotFound()))
	assert.False(t, errors.Is(err, gorm.ErrRecordNotFound))
}
