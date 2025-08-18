package counter

import "sync"

type Uint16Counter struct {
	value uint16
	lock  sync.Mutex
}

func NewUint16Counter() *Uint16Counter {
	return &Uint16Counter{}
}

func (c *Uint16Counter) GetNext() uint16 {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value += 1
	return c.value
}

func (c *Uint16Counter) Reset() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value = 0
}
