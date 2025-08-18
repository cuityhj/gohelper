package counter

import "sync"

type Uint8Counter struct {
	value uint8
	lock  sync.Mutex
}

func NewUint8Counter() *Uint8Counter {
	return &Uint8Counter{}
}

func (c *Uint8Counter) GetNext() uint8 {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value += 1
	return c.value
}

func (c *Uint8Counter) Reset() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value = 0
}
