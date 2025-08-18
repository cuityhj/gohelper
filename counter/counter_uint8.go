package counter

import "sync"

type CounterUint8 struct {
	value uint8
	lock  sync.Mutex
}

func NewCounterUint8() *CounterUint8 {
	return &CounterUint8{}
}

func (c *CounterUint8) GetNext() uint8 {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value += 1
	return c.value
}

func (c *CounterUint8) Reset() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value = 0
}
