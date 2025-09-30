package lru

import (
	"container/list"
	"sync"
)

type safeRecordStore struct {
	maxRecordCount uint32
	lru            *list.List
	records        map[Key]*list.Element
	lock           sync.Mutex
}

func NewSafeRecordStore(maxRecordCount uint32) *safeRecordStore {
	return &safeRecordStore{
		maxRecordCount: maxRecordCount,
		lru:            list.New(),
		records:        make(map[Key]*list.Element, maxRecordCount),
	}
}

func (rs *safeRecordStore) Save(record Record) (Key, bool) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	if element, ok := rs.records[record.GetKey()]; ok {
		rs.lru.MoveToFront(element)
		element.Value = record
	} else {
		element := rs.lru.PushFront(record)
		rs.records[record.GetKey()] = element
		if uint32(rs.lru.Len()) >= rs.maxRecordCount {
			return rs.removeOldest()
		}
	}

	return EmptyKey, false
}

func (rs *safeRecordStore) Remove(key Key) bool {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	if element, ok := rs.records[key]; ok {
		delete(rs.records, key)
		rs.lru.Remove(element)
		return true
	}

	return false
}

func (rs *safeRecordStore) Get(key Key) (Record, bool) {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	if element, ok := rs.records[key]; ok {
		rs.lru.MoveToFront(element)
		return element.Value.(Record), true
	}

	return nil, false
}

func (rs *safeRecordStore) removeOldest() (Key, bool) {
	if element := rs.lru.Back(); element != nil {
		rs.lru.Remove(element)
		record := element.Value.(Record)
		delete(rs.records, record.GetKey())
		return record.GetKey(), true
	}

	return EmptyKey, false
}
