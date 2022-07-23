package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIdentifierConfiguration_IsMatchProperties(t *testing.T) {
	matcher1, _ := NewPropertyMatcher("name1", "dummy1")
	matcher2, _ := NewPropertyMatcher("name2", "dummy2")
	config, _ := NewIdentifierConfiguration(uuid.New(), "dummy", []PropertyMatcher{*matcher1, *matcher2}, 0, "", "")

	assert.True(t, config.IsMatchProperties(map[string]string{
		"name1": "dummy1",
		"name2": "dummy2",
	}))

	assert.False(t, config.IsMatchProperties(map[string]string{
		"name1": "value1",
		"name2": "value2",
	}))
}
