package syncpool

import (
	"bytes"
	"sync"
)

var gBufferPool = &BufferPool{pool: sync.Pool{New: func() interface{} { return bytes.NewBuffer(make([]byte, MaxDatagram)) }}}

type BufferPool struct {
	pool sync.Pool
}

func GetBufferPool() *BufferPool {
	return gBufferPool
}

func (b *BufferPool) Get() *bytes.Buffer {
	buf := b.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (b *BufferPool) Put(buf *bytes.Buffer) {
	b.pool.Put(buf)
}
