package lru

import (
	"container/list"
)

type unsafeRecordStore struct {
	maxRecordCount uint32
	lru            *list.List
	records        map[Key]*list.Element
}

func NewUnsafeRecordStore(maxRecordCount uint32) *unsafeRecordStore {
	return &unsafeRecordStore{
		maxRecordCount: maxRecordCount,
		lru:            list.New(),
		records:        make(map[Key]*list.Element, maxRecordCount),
	}
}

func (rs *unsafeRecordStore) Save(record Record) (Key, bool) {
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

func (rs *unsafeRecordStore) Remove(key Key) bool {
	if element, ok := rs.records[key]; ok {
		delete(rs.records, key)
		rs.lru.Remove(element)
		return true
	}

	return false
}

func (rs *unsafeRecordStore) Get(key Key) (Record, bool) {
	if element, ok := rs.records[key]; ok {
		rs.lru.MoveToFront(element)
		return element.Value.(Record), true
	}

	return nil, false
}

func (rs *unsafeRecordStore) removeOldest() (Key, bool) {
	if element := rs.lru.Back(); element != nil {
		rs.lru.Remove(element)
		record := element.Value.(Record)
		delete(rs.records, record.GetKey())
		return record.GetKey(), true
	}

	return EmptyKey, false
}
