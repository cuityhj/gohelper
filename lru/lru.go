package lru

import (
	"container/list"
	"sync"
	"time"
)

type Record interface {
	IsReachLimit(uint32) bool
	IsExpired() bool
	GetKey() string
}

type DefaultRecord struct {
	key        string
	count      uint32
	validTime  time.Duration
	expireTime time.Time
}

func NewDefaultRecord(key string, validTime time.Duration) Record {
	r := &DefaultRecord{
		key:       key,
		validTime: validTime,
	}
	r.reset()
	return r
}

func (r *DefaultRecord) reset() {
	r.count = 1
	r.expireTime = time.Now().Add(r.validTime)
}

func (r *DefaultRecord) IsReachLimit(limit uint32) bool {
	if r.IsExpired() {
		r.reset()
		return false
	}

	if r.count+1 > limit {
		return true
	}

	r.count++
	return false
}

func (r *DefaultRecord) IsExpired() bool {
	return r.expireTime.Before(time.Now())
}

func (r *DefaultRecord) GetKey() string {
	return r.key
}

type ConcurrencyRecordStore struct {
	maxRecordCount uint32
	lru            *list.List
	records        map[string]*list.Element
	lock           sync.Mutex
}

func NewConcurrencyRecordStore(maxRecordCount uint32) *ConcurrencyRecordStore {
	return &ConcurrencyRecordStore{
		maxRecordCount: maxRecordCount,
		lru:            list.New(),
		records:        make(map[string]*list.Element),
	}
}

func (rs *ConcurrencyRecordStore) Add(record Record) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	if element, ok := rs.records[record.GetKey()]; ok {
		rs.lru.MoveToFront(element)
		element.Value = record
	} else {
		element := rs.lru.PushFront(record)
		rs.records[record.GetKey()] = element
		if uint32(rs.lru.Len()) > rs.maxRecordCount {
			rs.removeOldest()
		}
	}
}

func (rs *ConcurrencyRecordStore) Remove(key string) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	if element, ok := rs.records[key]; ok {
		delete(rs.records, key)
		rs.lru.Remove(element)
	}
}

func (rs *ConcurrencyRecordStore) Get(key string) (Record, bool) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	if element, ok := rs.records[key]; ok {
		rs.lru.MoveToFront(element)
		return element.Value.(Record), true
	} else {
		return nil, false
	}
}

func (rs *ConcurrencyRecordStore) removeOldest() {
	if element := rs.lru.Back(); element != nil {
		rs.lru.Remove(element)
		record := element.Value.(Record)
		delete(rs.records, record.GetKey())
	}
}
