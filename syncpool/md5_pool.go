package syncpool

import (
	"crypto/md5"
	"hash"
	"sync"
)

var gMD5Pool = &MD5Pool{pool: sync.Pool{New: func() interface{} { return md5.New() }}}

type MD5Pool struct {
	pool sync.Pool
}

func GetMD5Pool() *MD5Pool {
	return gMD5Pool
}

func (m *MD5Pool) Get() hash.Hash {
	h := m.pool.Get().(hash.Hash)
	h.Reset()
	return h
}

func (m *MD5Pool) Put(h hash.Hash) {
	m.pool.Put(h)
}
