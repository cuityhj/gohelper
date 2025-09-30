package lru

type Key string
type Value interface{}

var EmptyKey = Key("")

type Record interface {
	IsReachLimit(uint32) bool
	IsExpired() bool
	GetKey() Key
	GetValue() Value
}

type RecordStore interface {
	Save(Record) (Key, bool)
	Remove(Key) bool
	Get(Key) (Record, bool)
}
