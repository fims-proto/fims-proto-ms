package identifier_configuration

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIdentifierConfiguration_IncrementCounter(t *testing.T) {
	matcher, _ := NewPropertyMatcher("test", "dummy")
	config1, _ := New(uuid.New(), "dummy", []PropertyMatcher{*matcher}, 0, "", "")
	config2, _ := New(uuid.New(), "dummy", []PropertyMatcher{*matcher}, 0, "", "")

	var wg sync.WaitGroup
	wg.Add(190)

	start := make(chan struct{})

	for i := 0; i < 100; i++ {
		go func() {
			<-start
			defer wg.Done()
			config1.IncrementCounter()
		}()
	}

	for i := 0; i < 90; i++ {
		go func() {
			<-start
			defer wg.Done()
			config2.IncrementCounter()
		}()
	}

	close(start)
	wg.Wait()

	assert.Equal(t, 100, config1.Counter())
	assert.Equal(t, 90, config2.Counter())
}
