package syncpool

import (
	"sync"
)

const MaxDatagram = 1 << 16

var gBytePool = &BytePool{pool: sync.Pool{New: func() interface{} { return make([]byte, MaxDatagram) }}}

type BytePool struct {
	pool sync.Pool
}

func GetBytePool() *BytePool {
	return gBytePool
}

func (b *BytePool) Get() []byte {
	return b.pool.Get().([]byte)
}

func (b *BytePool) Put(buf []byte) {
	b.pool.Put(buf)
}
