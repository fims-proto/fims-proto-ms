package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain_NewCounter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		prefix string
		sufix  string
		verify func(t *testing.T, counter *Counter, err error)
	}{
		{
			name:   "normal next",
			prefix: "天字",
			sufix:  "号",
			verify: func(t *testing.T, counter *Counter, err error) {
				require.NoError(t, err)
				assert.Equal(t, "天字0号", counter.Identifier())
				counter.Next()
				assert.Equal(t, "天字1号", counter.Identifier())
				counter.Next()
				assert.Equal(t, "天字2号", counter.Identifier())
			},
		},
		{
			name:   "normal reset and next",
			prefix: "地煞",
			sufix:  "位",
			verify: func(t *testing.T, counter *Counter, err error) {
				require.NoError(t, err)
				assert.Equal(t, "地煞0位", counter.Identifier())
				counter.Next()
				assert.Equal(t, "地煞1位", counter.Identifier())
				_ = counter.Reset()
				assert.Equal(t, "地煞0位", counter.Identifier())
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			counter, err := NewCounter(uuid.New(), "DUMMY", test.prefix, test.sufix)
			test.verify(t, counter, err)
		})
	}
}
