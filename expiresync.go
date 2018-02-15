package expiresync

import (
	"sync"
	"time"
)

type Map struct {
	storage *sync.Map
}

type expireValue struct {
	value  interface{}
	expire time.Time
}

func NewMap() *Map {
	return &Map{
		storage: &sync.Map{},
	}
}

func (s *Map) Get(key string) (value interface{}, ok bool) {
	if val, ok := s.storage.Load(key); ok {
		return val.(*expireValue).value, ok
	}
	return nil, false
}

func (s *Map) Set(key string, value interface{}, expire time.Duration) {
	s.storage.Store(key, &expireValue{
		value:  value,
		expire: time.Now().Add(expire),
	})
}

func (s *Map) DeleteExpired() {
	s.storage.Range(func(key interface{}, value interface{}) bool {
		session := (value.(expireValue))
		if session.expire.Before(time.Now()) {
			// is it safe?
			s.storage.Delete(key)
		}
		return true
	})
}

func (s *Map) Delete(key string) {
	s.storage.Delete(key)
}
