package counter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain_NewCounter(t * testing.T){
	t.Parallel()
	tests := []struct {
		UUID string
		len uint
		prefix string
		sufix string
		verify func(t * testing.T, counter *Counter, err error)
	}{
		{
			UUID: "testUUID",
			len: 6,
			prefix: "天字",
			sufix: "号",
			verify: func(t *testing.T, counter *Counter,err error){
				require.NoError(t,err)
				assert.Equal(t, "testUUID", counter.UUID)
				next,err1 := counter.Next()
				require.NoError(t,err1)
				assert.Equal(t, next,"天字000001号")
			},
		},
		{
			UUID: "testUUID",
			len: 4,
			prefix: "地煞",
			sufix: "位",
			verify: func(t *testing.T, counter *Counter,err error){
				require.NoError(t,err)
				assert.Equal(t, "testUUID", counter.UUID)
				counter.Next()
				nn,err2 := counter.Next()
				require.NoError(t,err2)
				assert.Equal(t, nn,"地煞0002位")
			},
		},
	}
	for _, test := range tests {
		test := test 
		t.Run(test.UUID,func(t *testing.T){
			t.Parallel()
			counter, err := NewCounter(test.UUID, test.len, test.prefix,test.sufix)
			test.verify(t,counter,err)
		})
	}
}

