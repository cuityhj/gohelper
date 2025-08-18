package counter

import "sync"

type CounterUint16 struct {
	value uint16
	lock  sync.Mutex
}

func NewCounterUint16() *CounterUint16 {
	return &CounterUint16{}
}

func (c *CounterUint16) GetNext() uint16 {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value += 1
	return c.value
}

func (c *CounterUint16) Reset() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value = 0
}
