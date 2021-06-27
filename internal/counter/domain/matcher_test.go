package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain_NewMatcher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		action func() (*Matcher, error)
		verify func(t *testing.T, m *Matcher, err error)
	}{
		{
			name: "normal",
			action: func() (*Matcher, error) {
				return NewMatcher("-", "obj1", "obj2")
			},
			verify: func(t *testing.T, m *Matcher, err error) {
				require.NoError(t, err)
				assert.Equal(t, "obj1-obj2", m.String())
			},
		},
		{
			name: "skip separator",
			action: func() (*Matcher, error) {
				return NewMatcher("", "obj1", "obj2")
			},
			verify: func(t *testing.T, m *Matcher, err error) {
				require.NoError(t, err)
				assert.Equal(t, "obj1-obj2", m.String())
			},
		},
		{
			name: "no objects",
			action: func() (*Matcher, error) {
				return NewMatcher("-")
			},
			verify: func(t *testing.T, m *Matcher, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "one object",
			action: func() (*Matcher, error) {
				return NewMatcher("-", "obj")
			},
			verify: func(t *testing.T, m *Matcher, err error) {
				require.NoError(t, err)
				assert.Equal(t, "obj", m.String())
			},
		},
		{
			name: "one empty object",
			action: func() (*Matcher, error) {
				return NewMatcher("-", "obj", "")
			},
			verify: func(t *testing.T, m *Matcher, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m, err := tt.action()
			tt.verify(t, m, err)
		})
	}
}
