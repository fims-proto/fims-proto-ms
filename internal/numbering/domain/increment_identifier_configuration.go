package domain

import (
	"sync"

	"github.com/google/uuid"
)

func (c *IdentifierConfiguration) IncrementCounter() {
	unlock := counterLock.lock(c.Id())
	defer unlock()

	c.counter++
}

type counterLockType struct {
	configurations sync.Map
}

func (c *counterLockType) lock(key uuid.UUID) func() {
	val, _ := c.configurations.LoadOrStore(key, &sync.Mutex{})
	configurationLock := val.(*sync.Mutex)
	configurationLock.Lock()

	return func() { configurationLock.Unlock() }
}

var counterLock = counterLockType{}
