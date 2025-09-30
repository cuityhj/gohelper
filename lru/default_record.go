package lru

import "time"

type defaultRecord struct {
	key        Key
	count      uint32
	validTime  time.Duration
	expireTime time.Time
}

func NewDefaultRecord(key Key, validTime time.Duration) Record {
	r := &defaultRecord{
		key:       key,
		validTime: validTime,
	}
	r.reset()
	return r
}

func (r *defaultRecord) reset() {
	r.count = 1
	r.expireTime = time.Now().Add(r.validTime)
}

func (r *defaultRecord) IsReachLimit(limit uint32) bool {
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

func (r *defaultRecord) IsExpired() bool {
	return r.expireTime.Before(time.Now())
}

func (r *defaultRecord) GetKey() Key {
	return r.key
}

func (r *defaultRecord) GetValue() Value {
	return r.key
}
