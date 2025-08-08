package counter

import "sync"

type Counter struct {
	value uint16
	lock  sync.Mutex
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) GetNext() uint16 {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value += 1
	return c.value
}

func (c *Counter) Reset() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value = 0
}
